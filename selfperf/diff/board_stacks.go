package diff

import (
	"fmt"
	"log"
)

func getAllThreadNames(reports ReportMap) map[string]uint64 {
	var threadNames map[string]uint64 = make(map[string]uint64)
	for _, r := range reports {
		for _, t := range r.Threads {
			threadNames[t.ThreadEntryName] = t.Size
		}
	}
	return threadNames
}

func PrintBoardsStackReportMarkdown(reports ReportMap, errors []string) {
	var has_overflows bool = false
	fmt.Print("## Boards Stack check\n\n")
	fmt.Print("### Boards Stack usage summary\n\n")
	for thread_name, stacksize := range getAllThreadNames(reports) {
		fmt.Print("\n### " + thread_name + " Stack use:\n\n")
		fmt.Printf("| Board name | Max Stack Use Found / %d bytes | Usage %% | Unresolved function calls in calltree|\n", stacksize)
		fmt.Println("| ----------- | ----------- | ----------- | ----------- |")

		for board_name, report := range reports {
			for _, thread := range report.Threads {
				if thread_name == thread.ThreadEntryName {
					stackusage_percent := float64(thread.Used) / float64(thread.Size) * 100.0
					fmt.Printf("| %s | %d | %3.1f %% | %d |\n", board_name, thread.Used, stackusage_percent, thread.NrUnresolvedCalls)
				}
			}
		}
		fmt.Println()
	}

	// fmt.Print("### Stack usage details - worst cases found\n\n")
	// for _, thread := range threads {
	// 	var is_overflow = thread.Used > thread.Size
	// 	var note_type = "caution"
	// 	if is_overflow {
	// 		note_type = "warning"
	// 		has_overflows = true
	// 	}
	// 	function_call_path := thread.WorstStackBranch
	// 	var stack_sum int64 = 0
	// 	var stackliststr = "Calltree\n\n"
	// 	for j := 0; j < len(function_call_path.CallList); j++ {
	// 		c := function_call_path.CallList[j]
	// 		stack_sum += c.StackSize
	// 		stackliststr = stackliststr + fmt.Sprintf("* %s StackSize: %d bytes (Sum: %d)\n", c.Name, c.StackSize, stack_sum)
	// 	}
	// 	threadWarningHeader := fmt.Sprintf("Worst Calltree for %s", thread.ThreadEntryName)
	// 	threadWarningSummary := fmt.Sprintf("Is %d calls deep, expand for details.", len(function_call_path.CallList))
	// 	warningDetails := printMdCollapsableMsgWithDetails(note_type, threadWarningHeader, threadWarningSummary, stackliststr)
	// 	fmt.Println(warningDetails)

	// 	if thread.NrOverflowPaths > 1 {
	// 		fmt.Printf("│ └ WARNING: %-14d paths are overflowing, check pexplorer for more details! │\n", thread.NrOverflowPaths)
	// 	}
	// }

	fmt.Println("### Unresolved call details\n")
	for board_name, report := range reports {
		stats := report.Report.UnresolvedStats
		fmt.Printf("* %s at least %d dynamic calls in %d functions left to resolve. %d dynamic calls are already resolved.\n", board_name, stats.TotalNrDynamicCalls, stats.NrFunctionsWithDynamicCalls, stats.TotalNrResolvedCalls)
	}

	if len(errors) > 0 {
		fmt.Println("\n### Errors during analysis\n")
		for _, e := range errors {
			fmt.Println("* " + e)
		}
	}

	if has_overflows {
		log.Fatalln("Stackoverflow detected!")
	}
}
