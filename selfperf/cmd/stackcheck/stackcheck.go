package main

import (
	"crypto/sha256"
	"debug/elf"
	"flag"
	"fmt"
	"log"

	"github.com/paulwuertz/pexplorer/selfperf/callgraph"
	"github.com/paulwuertz/pexplorer/selfperf/config"
	"github.com/paulwuertz/pexplorer/selfperf/rtos"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func main() {
	infile := flag.String("i", "", "input ELF file - obligatory")
	conffile := flag.String("c", "", "config file - containing dynamic threads and calls")

	flag.Parse()

	fw_hash := sha256.Sum256([]byte(*infile))
	fw_hash_str := fmt.Sprintf("%x", fw_hash)
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

	elfReport := symbolextraction.GetFWReport(elfFile, fw_hash_str)
	callgraph.EnhanceByDisasm(&elfReport, p.DynamicCalls)
	callgraph.TraverseCallGraph(&elfReport)
	threads := rtos.GetAllThreads(&elfReport, p)
	rtos.PrintStackStats(threads)
}
