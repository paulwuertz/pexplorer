package rtos

import (
	"fmt"
	"log"
	"maps"
	"math"
	"slices"
	"strings"

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

func PrintStackStats(threads []config.RTOSThread, stats symbolextraction.UnresolvedCallStats) {

	fmt.Println("┌────────────────────────────────────────────────────────────────────────────────────┐")
	//           │ gs_usb_tx_thread           uses at least   728 /  1024 ( 71%) |███████████████-----│ ...
	const colorRed = "\033[0;31m"
	const colorNone = "\033[0m"
	for _, thread := range threads {
		stackusage_percent := float64(thread.Used) / float64(thread.Size) * 100.0
		var stackusage_barfill int = min(int(stackusage_percent+4)/5, 20)
		var bar = strings.Repeat("█", stackusage_barfill) + strings.Repeat("-", max(0, 20-stackusage_barfill))
		var is_overflow = thread.Used > thread.Size
		var color_start = colorNone
		var color_end = colorNone
		if is_overflow {
			color_start = colorRed
		}
		fmt.Printf("│ %-26s uses at least %5d / %5d (%3.0f%%) |%s%20s%s│\n", thread.ThreadEntryName, thread.Used, thread.Size, stackusage_percent, color_start, bar, color_end)
		if is_overflow {
			function_call_path := thread.WorstStackBranch
			var stack_sum int64 = 0
			for j := 0; j < len(function_call_path.CallList); j++ {
				c := function_call_path.CallList[j]
				list_str := "├──"
				if j == len(function_call_path.CallList)-1 {
					list_str = "└──"
				}
				stack_sum += c.StackSize
				fmt.Printf("│    %-3s %-35s StackSize: %-5d bytes (Sum: %-8d)  │\n", list_str, c.Name, c.StackSize, stack_sum)
			}
		}
		if thread.NrOverflowPaths > 1 {
			fmt.Printf("│ └ WARNING: %-14d paths are overflowing, check pexplorer for more details! │\n", thread.NrOverflowPaths)
		}
	}
	fmt.Println("└────────────────────────────────────────────────────────────────────────────────────┘")

	any_unresolved_fn_in_thread := false
	for _, thread := range threads {
		if thread.NrUnresolvedCalls != 0 {
			fmt.Printf("WARNING: %-25s - incomplete calltree: %-10d unresolved calls)\n", thread.ThreadEntryName, thread.NrUnresolvedCalls)
			any_unresolved_fn_in_thread = true
		}
	}
	if any_unresolved_fn_in_thread {
		fmt.Println("        --> consider adding or extending a config and add unresolved calls for better results")
		fmt.Printf("            there are at least %d dynamic calls in %d functions left to resolve. %d dynamic calls are already resolved.\n \n", stats.TotalNrDynamicCalls, stats.NrFunctionsWithDynamicCalls, stats.TotalNrResolvedCalls)
		fmt.Println("            (Note: the actual number of dynamic calls might be higher then the number of dynamic branch instructions.")
		fmt.Println("             Like in a workqueue a single dynamic branch can execute more then on work tasks.)")
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
			thread_fn_calltree, ovb := threadEntryFn.GetCallTreeJson(s, uint64(stackSize), 0)
			tm[threadEntryVarAddr] = config.RTOSThread{
				ThreadEntryName:   threadEntryFn.Name,
				StackVariableName: stackVar.Name,
				Size:              uint64(stackSize),
				Used:              uint64(thread_fn_calltree.Tree.MaxStackSizeCallees),
				NrUnresolvedCalls: uint64(len(thread_fn_calltree.UnresolvedCalls)),
				NrOverflowPaths:   ovb,
				Calltree:          thread_fn_calltree,
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
		thread_fn_calltree, ovb := threadEntryFn.GetCallTreeJson(s, uint64(stackSize), 0)
		tm[threadEntryFn.Address] = config.RTOSThread{
			ThreadEntryName:   tName,
			StackVariableName: sName,
			Size:              uint64(stackSize),
			Used:              uint64(thread_fn_calltree.Tree.MaxStackSizeCallees),
			NrUnresolvedCalls: uint64(len(thread_fn_calltree.UnresolvedCalls)),
			WorstStackBranch:  thread_fn_calltree.Branches[0],
			NrOverflowPaths:   ovb,
			Calltree:          thread_fn_calltree,
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
