package rtos

import (
	"fmt"
	"log"
	"maps"
	"math"
	"slices"

	"github.com/paulwuertz/pexplorer/selfperf/config"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func GetVarByAddr(addr uint64, s *symbolextraction.SElfReport) *symbolextraction.VariableSymbol {
	for _, v := range s.Variables {
		if v.Address-addr < 2 || addr-v.Address < 2 { // of by 1 is ok
			return &v
		}
	}
	return nil
}

func IsStaticZephyrThread(sym symbolextraction.VariableSymbol, s *symbolextraction.SElfReport) bool {
	secidx := sym.SectionIndex
	section, found := s.SectionsMap[secidx]
	if !found {
		// fmt.Printf("unknown section for sym: ", sym.Name)
		return false
	}
	section_name := section.Name
	return section_name == "_static_thread_data_area"
}

func mapMemberBytes(t symbolextraction.Typedef, data []byte) map[string][]byte {
	lookup := make(map[string]([]byte), 10)
	for i := 0; i < len(t.Members); i++ {
		field := t.Members[i]
		lookup[field.Name] = make([]byte, field.Size)
		for j := field.ByteOffset; j < field.ByteOffset+field.Size; j++ {
			byte_nr := j - field.ByteOffset
			lookup[field.Name][byte_nr] = data[j]
		}
	}
	return lookup
}

func arrayToUint64(data []byte) uint64 {
	var val uint64 = 0
	for i := 0; i < len(data); i++ {
		val += uint64(data[i]) * uint64(math.Pow(256, float64(i)))
	}
	return val
}

func PrintStackStats(threads []config.RTOSThread) {
	// ┌─────────────────────────────────────────────────────────────────────────┐
	// │ log_process_thread_func   -  42.3% -   352/  832b - ████████------------│
	// │ shell_thread              -   7.3% -   304/ 4160b - █-------------------│
	// │ mgmt_event_work_handler   -  31.2% -   280/  896b - ██████--------------│
	// │ bg_thread_main            -  58.3% -  1232/ 2112b - ███████████---------│
	// │ work_queue_main           -  25.0% -   272/ 1088b - █████---------------│
	// └─────────────────────────────────────────────────────────────────────────┘
	for _, thread := range threads {
		stackusage_percent := float64(thread.Used) / float64(thread.Size) * 100.0
		fmt.Printf("| %.25s uses at least %d / %d (%.2f%%)|", thread.ThreadEntryName, thread.Used, thread.Size, stackusage_percent)
		fmt.Println()
	}
	for _, thread := range threads {
		if thread.NrUnresolvedCalls != 0 {
			fmt.Print("WARNING: ", thread.NrUnresolvedCalls, " unresolved calls - incomplete calltree, the analysis needs to be resolved for better results)")
		}
		fmt.Println()
	}
}

type ThreadMap map[uint64]config.RTOSThread

// TODO only zephyr for now... how to definitly detect it though...
func FindStaticZephyrRtosThreads(s *symbolextraction.SElfReport) (tm ThreadMap) {
	idx := slices.IndexFunc(s.Types, func(c symbolextraction.Typedef) bool { return c.Name == "_static_thread_data" })
	static_thread_data_struct := s.Types[idx]
	tm = make(ThreadMap)
	for _, v := range s.Variables {
		// static threads created by macro
		if IsStaticZephyrThread(v, s) {
			thread_struct := mapMemberBytes(static_thread_data_struct, v.Data)
			_, nameOk := thread_struct["init_name"]
			stackAddrArr, stackOk := thread_struct["init_stack"]
			stackEntryAddr, entryOk := thread_struct["init_entry"]
			if !nameOk || !stackOk || !entryOk {
				log.Fatal("static thread without valid init_{name|stack|entry}")
			}
			stackAddr := arrayToUint64(stackAddrArr)
			threadEntryVarAddr := arrayToUint64(stackEntryAddr)
			threadEntryFn, fnFound := s.Addr2FnMap[threadEntryVarAddr]
			if !fnFound {
				log.Fatal("static thread init_entry not found at address:", threadEntryVarAddr)
			}
			stackVar := GetVarByAddr(stackAddr, s)
			stackSize := len(stackVar.Data)
			thread_fn_calltree := threadEntryFn.GetCallTreeJson(s, 0)
			tm[threadEntryVarAddr] = config.RTOSThread{
				ThreadEntryName:   threadEntryFn.Name,
				StackVariableName: stackVar.Name,
				Size:              uint64(stackSize),
				Used:              uint64(thread_fn_calltree.Tree.MaxStackSizeCallees),
				NrUnresolvedCalls: uint64(len(thread_fn_calltree.UnresolvedCalls)),
			}
			fmt.Println("Found Zephyr thread: ", tm[threadEntryVarAddr])
		}
	}
	return tm
}

func FindConfiguredZephyrRtosThreads(s *symbolextraction.SElfReport, conf config.PexplorerConfig, tm ThreadMap) ThreadMap {
	var threadEntryFn symbolextraction.FunctionSymbol
	for _, t := range conf.Threads {
		tName := t.ThreadEntryName
		sName := t.StackVariableName
		stackSize := 0
		threadFound := false
		stackFound := false
		for _, f := range s.Functions {
			if f.Name == tName {
				threadFound = true
				threadEntryFn = f
				break
			}
		}
		if !threadFound {
			log.Fatal("Configured thread ", tName, " not found in ELF functions")
		}

		if t.StackVariableName != "" {
			for _, v := range s.Variables {
				if v.Name == sName {
					stackFound = true
					stackSize = len(v.Data)
					break
				}
			}
		}

		if !stackFound {
			if t.Size != 0 {
				stackSize = int(t.Size)
			} else {
				log.Fatal("Configured thread - associated thread ", sName, "not found in ELF functions")
			}
		}
		thread_fn_calltree := threadEntryFn.GetCallTreeJson(s, 0)
		tm[threadEntryFn.Address] = config.RTOSThread{
			ThreadEntryName:   tName,
			StackVariableName: sName,
			Size:              uint64(stackSize),
			Used:              uint64(thread_fn_calltree.Tree.MaxStackSizeCallees),
			NrUnresolvedCalls: uint64(len(thread_fn_calltree.UnresolvedCalls)),
		}
	}
	return tm
}

func GetAllThreads(s *symbolextraction.SElfReport, conf config.PexplorerConfig) []config.RTOSThread {
	tm := FindStaticZephyrRtosThreads(s)
	tm = FindConfiguredZephyrRtosThreads(s, conf, tm)
	threads := slices.Collect(maps.Values(tm))
	return threads
}
