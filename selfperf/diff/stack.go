package diff

import (
	"fmt"
	"math"
	"slices"

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

func getColoredDiffedProgessBar(now, before int, isOverflow bool) string {
	const colorRed = "\033[0;31m"
	const colorGreen = "\033[0;32m"
	const colorNone = "\033[0m"
	var color string

	// either freed some stack -> green, used more -> red or neutral -> white
	if now > before {
		color = colorRed
	} else if now < before {
		color = colorGreen
	} else {
		color = colorNone
	}

	var bar []byte = make([]byte, 0)
	if isOverflow {
		bar = append(bar, []byte(colorRed)...)
	}

	// color in starting from 1% in 5% steps
	// nr_changed_blocks := int(math.Abs(float64(now)-float64(before))+4.0) / 5
	max_used := (max(now, before) + 4) / 5
	color_starts_at := int(min(now, before)) / 5
	for i := 0; i < 20; i++ {
		// color start and stop
		if i == color_starts_at+1 && !isOverflow {
			bar = append(bar, []byte(color)...)
		}
		// chars
		if i >= max_used {
			bar = append(bar, '-')
		} else {
			bar = append(bar, []byte("█")...)
		}
		if i == max_used && !isOverflow {
			bar = append(bar, []byte(colorNone)...)
		}
	}

	if isOverflow {
		bar = append(bar, []byte(colorNone)...)
	}
	return string(bar)
}

func RTOSStackDiff(new []config.RTOSThread, ref []config.RTOSThread) {
	commonThreads := getCommonThreads(new, ref)
	fmt.Println("┌────────────────────────────────────────────────────────────────────────────────────────┐")
	for _, ts := range commonThreads {
		thread := ts[0]
		ref_thread := ts[1]
		var stackusage_diff int = int(thread.Used) - int(ref_thread.Used)
		stackdiff_percent := float64(stackusage_diff) / float64(thread.Size) * 100.0
		stackusage_percent := float64(thread.Used) / float64(thread.Size) * 100.0
		var is_overflow = thread.Used > thread.Size
		var bar = getColoredDiffedProgessBar(int(stackusage_percent), int(stackusage_percent+stackdiff_percent), is_overflow)
		fmt.Printf("│ %-25s stack-use changed %+5d / %5d (%+-3.1f%%) |%s│\n", thread.ThreadEntryName, stackusage_diff, thread.Size, stackdiff_percent, bar)
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
	fmt.Println("└────────────────────────────────────────────────────────────────────────────────────────┘")

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

type FunctionStackDiff struct {
	NewFunctionName string  `json:"function_name"`
	NewFunctionPath string  `json:"function_path"`
	DiffStackSize   int64   `json:"diff_stacksize"`
	NewStackSize    int64   `json:"new_stacksize"`
	OldStackSize    int64   `json:"old_stacksize"`
	DiffPercentage  float64 `json:"diff_percentage"`
}
type FunctionStackDiffReport []FunctionStackDiff

func getStackDiff(diff FunctionDiffReport) FunctionStackDiffReport {
	stackDiff := make(FunctionStackDiffReport, 0)
	for _, pair := range diff.functionsMatches {
		newF := pair[0]
		refF := pair[1]
		var diff int64 = newF.StackSize - refF.StackSize
		if diff != 0 {
			diff_perc := 100.0*float64(newF.StackSize)/float64(refF.StackSize) - 100.0
			if math.IsInf(diff_perc, 1) {
				diff_perc = 100.0
			} else if math.IsInf(diff_perc, -1) {
				diff_perc = -100.0
			}
			stackDiff = append(stackDiff, FunctionStackDiff{
				NewFunctionName: newF.Name,
				NewFunctionPath: newF.SourceFilePath,
				DiffStackSize:   diff,
				NewStackSize:    newF.StackSize,
				OldStackSize:    refF.StackSize,
				DiffPercentage:  diff_perc,
			})
			// fmt.Println("\tstackDiff", k, diff, newF.StackSize, refF.StackSize)
		}
	}
	slices.SortFunc(stackDiff, func(i, j FunctionStackDiff) int {
		return int(j.DiffStackSize) - int(i.DiffStackSize)
	})
	return stackDiff
}

// func getCommonVariableSymbols(new FunctionSymbolMap, ref FunctionSymbolMap) FunctionDiffReport {

// }

func SymbolDiff(new symbolextraction.SElfReport, ref symbolextraction.SElfReport) FunctionStackDiffReport {
	diffReport := getCommonFunctionSymbols(new.Name2FnMap, ref.Name2FnMap)
	stackDiffReport := getStackDiff(diffReport)
	// for _, sd := range stackDiffReport {
	// 	fmt.Printf("stack_diff %40s:%d -> %d (%d / %-3.1f%%)\n", sd.NewFunction.Name, sd.OldStackSize, sd.NewFunction.StackSize, sd.DiffStackSize, sd.DiffPercentage)
	// }
	return stackDiffReport
	// getCommonVariableSymbols(new.Name2FnMap, ref.Name2FnMap)
}
