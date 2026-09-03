package main

import (
	"crypto/sha256"
	"debug/elf"
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/paulwuertz/pexplorer/selfperf/callgraph"
	"github.com/paulwuertz/pexplorer/selfperf/config"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func main() {
	// arg input, output file, indentation

	infile := flag.String("i", "", "input ELF file - obligatory")
	conffile := flag.String("c", "", "config file - containing dynamic threads and calls")
	outfile := flag.String("o", "", "output report to this file - if omited print to stdout")
	function_call_list_file := flag.String("f", "", "output list of all function calls to a file in json format to compare /regression")
	pretty := flag.Bool("p", false, "pretty-print the report else it is compact")

	flag.Parse()

	if *infile == "" {
		log.Fatal("Please add an ELF file to generate a report for.")
	}
	var conf config.PexplorerConfig
	elfFile, err := elf.Open(*infile)
	fw_hash := sha256.Sum256([]byte(*infile))
	fw_hash_str := fmt.Sprintf("%x", fw_hash)

	if err != nil {
		log.Fatal(err)
	}

	if *conffile != "" {
		conf, err = config.Import_config_from_file(*conffile)
		if err != nil {
			fmt.Println("Error importing config file:", err)
			return
		}
	}

	elfReport := symbolextraction.GetFWReport(elfFile, fw_hash_str)
	callgraph.EnhanceByDisasm(&elfReport, conf.DynamicCalls)
	callgraph.TraverseCallGraph(&elfReport)
	elfReport.SingleFirmware = true
	elfReport.FirmwareIdentifier = "unspecified"

	var datajson []byte
	if *pretty {
		datajson, _ = json.MarshalIndent(elfReport, "", "    ")
	} else {
		datajson, _ = json.Marshal(elfReport)
	}

	if *function_call_list_file != "" {
		fn_calls := callgraph.GetFunctionCallList(&elfReport)
		fnjson, _ := json.MarshalIndent(fn_calls, "", "    ")
		if err := os.WriteFile(*function_call_list_file, fnjson, 0666); err != nil {
			log.Fatal(err)
		}
	}

	if *outfile == "" {
		fmt.Println(string(datajson))
	} else {
		if err := os.WriteFile(*outfile, datajson, 0666); err != nil {
			log.Fatal(err)
		}
	}
}
