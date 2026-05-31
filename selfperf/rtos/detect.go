package rtos

import (
	"fmt"
	"log"
	"math"
	"slices"

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

func ScanForRtosFeatures(s *symbolextraction.SElfReport) {
	idx := slices.IndexFunc(s.Types, func(c symbolextraction.Typedef) bool { return c.Name == "_static_thread_data" })
	static_thread_data_struct := s.Types[idx]
	// TODO only zephyr for now... how to definitly detect it though...
	fmt.Println("User thread static stack usage analysis - worst case found might be higher during runtime,")
	fmt.Println("due to unresolved calls, unknown to the static analysis or how interrupts are handled on your architecture")
	fmt.Println("Results:")
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
			threadName := threadEntryFn.Name
			stackSize := len(stackVar.Data)

			thread_fn_calltree := threadEntryFn.GetCallTreeJson(s)
			nr_unresolved_calls := len(thread_fn_calltree.UnresolvedCalls)
			stackusage_percent := float64(threadEntryFn.MaxStackSizeCallees) / float64(stackSize) * 100.0
			fmt.Print("\t - ", threadName, " uses at least ", threadEntryFn.MaxStackSizeCallees, "/")
			fmt.Print(stackSize, " (", stackusage_percent, "%) ")
			if nr_unresolved_calls != 0 {
				fmt.Print("(WARNING - ", nr_unresolved_calls, " unresolved calls - incomplete calltree, the analysis needs to be resolved for better results)")
			}
			fmt.Println()
		}
	}
}
