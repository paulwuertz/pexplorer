package diff

import (
	"fmt"
	"slices"
	"strings"

	"github.com/paulwuertz/pexplorer/selfperf/config"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

type ThreadMap map[string][2]config.RTOSThread

func getCommonThreads(new []config.RTOSThread, ref []config.RTOSThread) ThreadMap {
	var commonThreads ThreadMap = make(ThreadMap)
	for _, n := range new {
		for _, r := range ref {
			if n.ThreadEntryName == r.ThreadEntryName {
				commonThreads[n.ThreadEntryName] = [2]config.RTOSThread{n, r}
			}
		}
	}
	return commonThreads
}

func RTOSStackDiff(new []config.RTOSThread, ref []config.RTOSThread) {
	commonThreads := getCommonThreads(new, ref)
	const colorRed = "\033[0;31m"
	const colorNone = "\033[0m"
	fmt.Println("┌────────────────────────────────────────────────────────────────────────────────────┐")
	for _, ts := range commonThreads {
		thread := ts[0]
		ref_thread := ts[1]
		var stackusage_diff int = int(thread.Used) - int(ref_thread.Used)
		stackdiff_percent := float64(stackusage_diff) / float64(thread.Size) * 100.0
		stackusage_percent := float64(thread.Used) / float64(thread.Size) * 100.0
		var stackusage_barfill int = min(int(stackusage_percent+4)/5, 20)
		var bar = strings.Repeat("█", stackusage_barfill) + strings.Repeat("-", max(0, 20-stackusage_barfill))
		var is_overflow = thread.Used > thread.Size
		var color_start = colorNone
		var color_end = colorNone
		if is_overflow {
			color_start = colorRed
		}
		fmt.Printf("│ %-26s uses at least %5d / %5d (%-3.1f%%) |%s%20s%s│\n", thread.ThreadEntryName, stackusage_diff, thread.Size, stackdiff_percent, color_start, bar, color_end)
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

	// any_unresolved_fn_in_thread := false
	// for _, thread := range threads {
	// 	if thread.NrUnresolvedCalls != 0 {
	// 		fmt.Printf("WARNING: %-25s - incomplete calltree: %-10d unresolved calls)\n", thread.ThreadEntryName, thread.NrUnresolvedCalls)
	// 		any_unresolved_fn_in_thread = true
	// 	}
	// }
	// if any_unresolved_fn_in_thread {
	// 	fmt.Println("        --> consider adding a config and add unresolved calls for better results")
	// }
}

type FunctionSymbolMap map[string]*symbolextraction.FunctionSymbol
type VariableSymbolMap map[string]*symbolextraction.VariableSymbol
type FunctionSymbolMatch map[string][2]*symbolextraction.FunctionSymbol
type VariableSymbolMatch map[string][2]*symbolextraction.VariableSymbol
type FunctionDiffReport struct {
	functionsMatches FunctionSymbolMatch
	functionsDeleted FunctionSymbolMap
	functionsAdded   FunctionSymbolMap
	functionsRenamed FunctionSymbolMatch // TODO
}

func getCommonFunctionSymbols(new FunctionSymbolMap, ref FunctionSymbolMap) FunctionDiffReport {
	var diff FunctionDiffReport
	diff.functionsMatches = make(FunctionSymbolMatch)
	diff.functionsDeleted = make(FunctionSymbolMap)
	diff.functionsAdded = make(FunctionSymbolMap)
	diff.functionsRenamed = make(FunctionSymbolMatch)

	for k, newF := range new {
		refF, isInRef := ref[k]
		if !isInRef {
			// TODO WARNING existing insert
			diff.functionsAdded[k] = newF
			continue
		}
		if slices.Compare(newF.Asm, refF.Asm) != 0 {
			diff.functionsMatches[k] = [2]*symbolextraction.FunctionSymbol{newF, refF}
		}
		// fmt.Println("kvloop", k, v, refF)
	}

	for k, fRef := range ref {
		_, isInNew := new[k]
		if !isInNew {
			// TODO WARNING existing insert
			for kAdded, fAdded := range diff.functionsAdded {
				if slices.Compare(fAdded.Asm, fRef.Asm) == 0 {
					// identical fn with different name
					// is renamed not added and deleted
					diff.functionsRenamed[kAdded] = [2]*symbolextraction.FunctionSymbol{
						fRef, fAdded,
					}
					continue
				}
			}
			diff.functionsDeleted[k] = fRef
		}
	}
	// make a new added slice with no remakes
	functionsAddedNoRenames := make(FunctionSymbolMap)
	for kRenamed, fAdded := range diff.functionsAdded {
		_, isRenamed := diff.functionsRenamed[kRenamed]
		if !isRenamed {
			functionsAddedNoRenames[kRenamed] = fAdded
		}
	}
	diff.functionsAdded = functionsAddedNoRenames

	fmt.Println("matches", len(diff.functionsMatches))
	fmt.Println("added", len(diff.functionsAdded))
	fmt.Println("deleted", len(diff.functionsDeleted))
	fmt.Println("renamed", len(diff.functionsRenamed))
	return diff
}

func getCommonVariableSymbols(new FunctionSymbolMap, ref FunctionSymbolMap) FunctionDiffReport {

}

func SymbolDiff(new symbolextraction.SElfReport, ref symbolextraction.SElfReport) {
	getCommonFunctionSymbols(new.Name2FnMap, ref.Name2FnMap)
	getCommonVariableSymbols(new.Name2FnMap, ref.Name2FnMap)
}
