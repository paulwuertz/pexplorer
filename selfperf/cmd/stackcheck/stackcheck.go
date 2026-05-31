package main

import (
	"debug/elf"
	"flag"
	"log"

	"github.com/paulwuertz/pexplorer/selfperf/callgraph"
	"github.com/paulwuertz/pexplorer/selfperf/rtos"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func main() {
	infile := flag.String("i", "", "input ELF file - obligatory")

	flag.Parse()

	if *infile == "" {
		log.Fatal("Please add an ELF file to generate a report for.")
	}
	elfFile, err := elf.Open(*infile)

	if err != nil {
		log.Fatal(err)
	}

	elfReport := symbolextraction.GetFWReport(elfFile)
	callgraph.EnhanceByDisasm(&elfReport)
	callgraph.GetStackUseDetails(&elfReport)
	callgraph.TraverseCallGraph(&elfReport)
	rtos.ScanForRtosFeatures(&elfReport)
}
