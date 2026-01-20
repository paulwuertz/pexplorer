package main

import (
	"debug/elf"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func main() {
	// arg input, output file, indentation

	infile := flag.String("i", "", "input ELF file - obligatory")
	outfile := flag.String("o", "", "output report to this file - if omited print to stdout")
	pretty := flag.Bool("p", false, "pretty-print the report else it is compact")

	flag.Parse()

	if *infile == "" {
		log.Fatal("Please add an ELF file to generate a report for.")
	}
	elfFile, err := elf.Open(*infile)

	if err != nil {
		log.Fatal(err)
	}

	elfReport := symbolextraction.GetFWReport(elfFile)
	var datajson []byte
	if *pretty {
		datajson, _ = json.MarshalIndent(elfReport, "", "    ")
	} else {
		datajson, _ = json.Marshal(elfReport)
	}

	if *outfile == "" {
		fmt.Println(string(datajson))
	} else {
		if err := os.WriteFile(*outfile, datajson, 0666); err != nil {
			log.Fatal(err)
		}
	}
}
