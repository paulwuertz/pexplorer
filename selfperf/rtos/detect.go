package rtos

import (
	"fmt"
	"maps"
	"slices"
	"strings"

	"github.com/paulwuertz/pexplorer/selfperf/config"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func GetVarByAddr(addr uint64, s *symbolextraction.SElfReport) *symbolextraction.VariableSymbol {
	for i := range s.Variables {
		v := &s.Variables[i]
		if v.Address-addr < 2 || addr-v.Address < 2 { // of by 1 is ok
			return v
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
	var val uint64
	for i, b := range data {
		val |= uint64(b) << (8 * i)
	}
	return val
}

func PrintStackStats(threads []config.RTOSThread, stats symbolextraction.UnresolvedCallStats) error {
	var has_overflows bool = false
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
			has_overflows = true
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

	if has_overflows {
		return fmt.Errorf("stack overflow detected")
	}
	return nil
}

type ThreadMap map[uint64]config.RTOSThread

// TODO only zephyr for now... how to definitly detect it though...
func FindStaticZephyrRtosThreads(s *symbolextraction.SElfReport) (tm ThreadMap) {
	tm = make(ThreadMap)
	idx := slices.IndexFunc(s.Types, func(c symbolextraction.Typedef) bool { return c.Name == "_static_thread_data" })
	if idx < 0 || idx >= len(s.Types) {
		s.Errors = append(s.Errors, "_static_thread_data type not found in DWARF metadata")
		return tm
	}
	static_thread_data_struct := s.Types[idx]
	for _, v := range s.Variables {
		// static threads created by macro
		if IsStaticZephyrThread(v, s) {
			thread_struct := mapMemberBytes(static_thread_data_struct, v.Data)
			_, nameOk := thread_struct["init_name"]
			stackAddrArr, stackOk := thread_struct["init_stack"]
			stackEntryAddr, entryOk := thread_struct["init_entry"]
			if !nameOk || !stackOk || !entryOk {
				s.Errors = append(s.Errors, fmt.Sprintf("Static thread without valid init_{name|stack|entry} %v", v.Name))
				continue
			}
			stackAddr := arrayToUint64(stackAddrArr)
			threadEntryVarAddr := arrayToUint64(stackEntryAddr)
			threadEntryFn, fnFound := s.Addr2FnMap[threadEntryVarAddr]
			if !fnFound {
				s.Errors = append(s.Errors, fmt.Sprintf("Static thread init_entry not found at address: %v", threadEntryVarAddr))
				continue
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
			s.Errors = append(s.Errors, fmt.Sprintf("Configured thread %s not found in ELF functions", tName))
			continue
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
				s.Errors = append(s.Errors, fmt.Sprintf("Configured thread - associated thread %s not found in ELF functions", sName))
				continue
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

type DeviveDriverCall struct {
	// assume a call when address in the range of the device API struct
	// is refereneced before an unresolved call
	DeviceStructStart    uint64
	DeviceStructEnd      uint64
	DeviceAPIStructStart uint64
	DeviceAPIStructEnd   uint64
	DeviceVar            *symbolextraction.VariableSymbol
	DeviceAPIVar         *symbolextraction.VariableSymbol
	Callbacks            []*symbolextraction.FunctionSymbol
}
type DeviveDriverCalls []DeviveDriverCall

func GetDeviceAPICalls(s *symbolextraction.SElfReport) {
	idx := slices.IndexFunc(s.Types, func(c symbolextraction.Typedef) bool { return c.Name == "device" })
	if idx == -1 {
		fmt.Println("No device struct found for resolving zephyr device API calls.")
		return
	}
	device_data_struct := s.Types[idx]

	var deviceVars []*symbolextraction.VariableSymbol
	var deviceAPIVars []*symbolextraction.VariableSymbol
	for _, v := range s.Variables {
		if v.VariableType == "device" {
			if v.FlashSize != uint64(device_data_struct.Size) {
				fmt.Println("Size mismatch for device struct, skipping.")
				continue
			}
			deviceStruct := mapMemberBytes(device_data_struct, v.Data)
			for name, m := range deviceStruct {
				if len(m) == 4 {
					address := arrayToUint64(m)
					// fmt.Printf("\t%s 0x%x\n", name, address)
					idx := slices.IndexFunc(s.Variables, func(v symbolextraction.VariableSymbol) bool { return v.Address == address })

					//just care for API calls now...?!
					if name == "api" && address != 0 && idx != -1 {
						apiVar := &s.Variables[idx]
						// fmt.Printf("device var: 0x%x %s %s %d\n", v.Address, v.Name, apiVar.Name, apiVar.FlashSize)
						deviceVars = append(deviceVars, &v)
						deviceAPIVars = append(deviceAPIVars, apiVar)
					}
				} else {
					fmt.Printf("\t%s len: %d %v\n", name, len(m), m)
				}
			}
		}
	}
	fmt.Println("device vars:", len(deviceVars))

	apicallcnt := 0
	var DeviceApiCallMap DeviveDriverCalls = make(DeviveDriverCalls, 0)
	for i, apiVar := range deviceAPIVars {
		idx := slices.IndexFunc(s.Types, func(c symbolextraction.Typedef) bool { return c.Name == apiVar.VariableType })
		if idx == -1 {
			fmt.Println("No type for device API struct, skipping...")
			continue
		}
		device_api_struct := s.Types[idx]
		apiStruct := mapMemberBytes(device_api_struct, apiVar.Data)
		devDriverCalls := DeviveDriverCall{
			DeviceStructStart:    deviceVars[i].Address,
			DeviceStructEnd:      deviceVars[i].Address + deviceVars[i].FlashSize,
			DeviceAPIStructStart: apiVar.Address,
			DeviceAPIStructEnd:   apiVar.Address + apiVar.FlashSize,
			DeviceVar:            deviceVars[i],
			DeviceAPIVar:         apiVar,
			Callbacks:            []*symbolextraction.FunctionSymbol{},
		}
		for name, fnCallAddrArr := range apiStruct {
			if len(fnCallAddrArr) == 4 {
				fnCallAddr := arrayToUint64(fnCallAddrArr)
				// fmt.Printf("\t%s 0x%x\n", name, address)
				f, isCallbackFound := s.Addr2FnMap[fnCallAddr]

				//just care for API calls now...?!
				if isCallbackFound && fnCallAddr != 0 {
					fmt.Printf("\tapicall found %s\n", f.Name)
					apicallcnt += 1
					devDriverCalls.Callbacks = append(devDriverCalls.Callbacks, f)
				} else {
					fmt.Printf("\t%s !!!nocallbackfound: %d\n", apiVar.Name, fnCallAddr)
				}
			} else {
				fmt.Printf("\t%s !!!len: %d %v\n", name, len(fnCallAddrArr), fnCallAddrArr)
			}
		}
		if len(devDriverCalls.Callbacks) != 0 {
			DeviceApiCallMap = append(DeviceApiCallMap, devDriverCalls)
		} else {
			fmt.Printf("!!!no Dev Driver Calls found for: %s\n", apiVar.Name)
		}
	}
	fmt.Println("API calls", DeviceApiCallMap)
	fmt.Println("len API calls", apicallcnt, " over ", len(DeviceApiCallMap), " apis")
	fmt.Println("len API calls", apicallcnt, " over ", len(DeviceApiCallMap), " apis")
}
