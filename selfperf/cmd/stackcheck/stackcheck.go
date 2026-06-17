package main

import (
	"crypto/sha256"
	"debug/elf"
	"flag"
	"fmt"
	"log"

	"github.com/paulwuertz/pexplorer/selfperf/callgraph"
	"github.com/paulwuertz/pexplorer/selfperf/rtos"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func main() {
	infile := flag.String("i", "", "input ELF file - obligatory")

	flag.Parse()

	fw_hash := sha256.Sum256([]byte(*infile))
	fw_hash_str := fmt.Sprintf("%x", fw_hash)
	if *infile == "" {
		log.Fatal("Please add an ELF file to generate a report for.")
	}
	elfFile, err := elf.Open(*infile)

	if err != nil {
		log.Fatal(err)
	}

	elfReport := symbolextraction.GetFWReport(elfFile, fw_hash_str)
	callgraph.EnhanceByDisasm(&elfReport)
	callgraph.GetStackUseDetails(&elfReport)
	callgraph.TraverseCallGraph(&elfReport)
	rtos.ScanForRtosFeatures(&elfReport)
}
