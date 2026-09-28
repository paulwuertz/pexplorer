package main

import (
	"debug/elf"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/paulwuertz/pexplorer/selfperf/config"
	"github.com/paulwuertz/pexplorer/selfperf/diff"
	"github.com/paulwuertz/pexplorer/selfperf/report"
)

type StringSlice []string

func (s *StringSlice) String() string {
	return fmt.Sprint(*s)
}

func (s *StringSlice) Set(value string) error {
	*s = append(*s, value)
	return nil
}

func main() {
	var board_elfs StringSlice
	flag.Var(&board_elfs, "elf", "input all your board ELF files - obligatory")
	conffile := flag.String("c", "", "config file - containing dynamic threads and calls")
	flag.Parse()

	if len(board_elfs) < 2 {
		log.Fatal("Please add at least two ELF file to generate a board stack report.")
	}
	var p, err = config.Import_config_from_file(*conffile)
	if err != nil {
		log.Fatal(err)
	}

	reports := make(diff.ReportMap, 0)
	for _, board_elf := range board_elfs {
		// TODO just a demo, make it something better :)
		board_elf_pre := strings.Replace(board_elf, "/home/runner/work/cannectivity/cannectivity/cannectivity/twister-out/", "", -1)
		board_name := strings.Replace(board_elf_pre, "/zephyr_gnu/cannectivity/app/app.cannectivity/zephyr/zephyr.elf", "", -1)
		newElfFile, err := elf.Open(board_elf)
		if err != nil {
			log.Fatal(err)
		}
		newReport, newTreadStats := report.GetReportAndRTOSStats(newElfFile, p)
		reports[board_name] = diff.ElfReport{
			Report:  newReport,
			Threads: newTreadStats,
		}
	}
	diff.PrintBoardsStackReportMarkdown(reports, []string{})
}
