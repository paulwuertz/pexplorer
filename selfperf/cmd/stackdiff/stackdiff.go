package main

import (
	"debug/elf"
	"flag"
	"log"

	"github.com/paulwuertz/pexplorer/selfperf/config"
	"github.com/paulwuertz/pexplorer/selfperf/diff"
	"github.com/paulwuertz/pexplorer/selfperf/report"
	"github.com/paulwuertz/pexplorer/selfperf/rtos"
)

func main() {
	infile := flag.String("i", "", "input ELF file - obligatory")
	reffile := flag.String("ref", "", "reference ELF file - obligatory")
	conffile := flag.String("c", "", "config file - containing dynamic threads and calls")

	settings := diff.DiffSettings{
		ShowFilePath:      true,
		ShowStackDiff:     true,
		ShowAllSymbolDiff: true,
	}

	flag.Parse()

	if *infile == "" {
		log.Fatal("Please add an ELF file to generate a diff report for.")
	}
	if *reffile == "" {
		log.Fatal("Please add a reference ELF file to generate a diff report for.")
	}
	newElfFile, err := elf.Open(*infile)
	refElfFile, err := elf.Open(*reffile)

	var p config.PexplorerConfig
	if *conffile != "" {
		p, err = config.Import_config_from_file(*conffile)
	}

	if err != nil {
		log.Fatal(err)
	}

	newReport, newTreadStats := report.GetReportAndRTOSStats(newElfFile, p)
	refReport, refTreadStats := report.GetReportAndRTOSStats(refElfFile, p)
	rtos.PrintStackStats(newTreadStats)
	rtos.PrintStackStats(refTreadStats)
	log.Printf("%d %d", &newReport, &refReport)
}
