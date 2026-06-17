package main

import (
	"crypto/sha256"
	"debug/elf"
	"flag"
	"fmt"
	"log"

	"github.com/paulwuertz/pexplorer/selfperf/callgraph"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func main() {
	infile := flag.String("i", "", "input ELF file - obligatory")
	// unitFilterName := flag.String("f", "", "filter to show only the dwarf tree of compilation unita containing this string in its pathname")

	fw_hash := sha256.Sum256([]byte(*infile))
	fw_hash_str := fmt.Sprintf("%x", fw_hash)
	flag.Parse()
	if *infile == "" {
		log.Fatal("Please add an ELF file to generate a report for.")
	}
	elfFile, err := elf.Open(*infile)

	elfReport := symbolextraction.GetFWReport(elfFile, fw_hash_str)
	callgraph.EnhanceByDisasm(&elfReport)
	// if *unitFilterName == "" || strings.Contains(cname, *unitFilterName) {
	callgraph.GetStackUseDetails(&elfReport)
	callgraph.TraverseCallGraph(&elfReport)
	// }

	if err != nil {
		log.Fatal(err)
	}

}
