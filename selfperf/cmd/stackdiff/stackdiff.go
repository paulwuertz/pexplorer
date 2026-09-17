package main

import (
	"debug/elf"
	"encoding/json"
	"flag"
	"fmt"
	"log"

	"github.com/paulwuertz/pexplorer/selfperf/config"
	"github.com/paulwuertz/pexplorer/selfperf/diff"
	"github.com/paulwuertz/pexplorer/selfperf/report"
)

func main() {
	infile := flag.String("i", "", "input ELF file - obligatory")
	ref_file := flag.String("ref", "", "reference ELF file - obligatory")
	conffile := flag.String("c", "", "config file - containing dynamic threads and calls")
	ref_conffile := flag.String("ref_config", "", "config file - containing dynamic threads and calls")

	// settings := diff.DiffSettings{
	// 	ShowFilePath:      true,
	// 	ShowStackDiff:     true,
	// 	ShowAllSymbolDiff: true,
	// }

	flag.Parse()

	if *infile == "" {
		log.Fatal("Please add an ELF file to generate a diff report for.")
	}
	if *ref_file == "" {
		log.Fatal("Please add a reference ELF file to generate a diff report for.")
	}
	newElfFile, err := elf.Open(*infile)
	refElfFile, err := elf.Open(*ref_file)

	var p config.PexplorerConfig
	if *conffile != "" {
		p, err = config.Import_config_from_file(*conffile)
	}

	var ref_p config.PexplorerConfig
	if *ref_conffile != "" {
		ref_p, err = config.Import_config_from_file(*ref_conffile)
	} else {
		ref_p, err = config.Import_config_from_file(*conffile)
	}

	if err != nil {
		log.Fatal(err)
	}

	newReport, newTreadStats := report.GetReportAndRTOSStats(newElfFile, p)
	refReport, refTreadStats := report.GetReportAndRTOSStats(refElfFile, ref_p)
	diff.RTOSStackDiff(newTreadStats, refTreadStats)
	stack_diff := diff.SymbolDiff(newReport, refReport)

	var datajson, _ = json.MarshalIndent(stack_diff, "", "    ")
	fmt.Println(string(datajson))
	// rtos.PrintStackStats(newTreadStats)
	// rtos.PrintStackStats(refTreadStats)
	// log.Printf("%d %d", &newReport, &refReport)
}
