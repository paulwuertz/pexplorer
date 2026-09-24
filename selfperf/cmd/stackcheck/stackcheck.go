package main

import (
	"debug/elf"
	"flag"
	"log"

	"github.com/paulwuertz/pexplorer/selfperf/config"
	"github.com/paulwuertz/pexplorer/selfperf/report"
	"github.com/paulwuertz/pexplorer/selfperf/rtos"
)

func main() {
	infile := flag.String("i", "", "input ELF file - obligatory")
	conffile := flag.String("c", "", "config file - containing dynamic threads and calls")
	flag.Parse()

	if *infile == "" {
		log.Fatal("Please add an ELF file to generate a report for.")
	}
	elfFile, err := elf.Open(*infile)

	var p config.PexplorerConfig
	if *conffile != "" {
		p, err = config.Import_config_from_file(*conffile)
	}

	if err != nil {
		log.Fatal(err)
	}

	elfReport, threads := report.GetReportAndRTOSStats(elfFile, p)
	rtos.PrintStackStats(threads, elfReport.UnresolvedStats)
}
