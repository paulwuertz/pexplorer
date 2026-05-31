//go:build js && wasm
// +build js,wasm

package main

import (
	"debug/elf"
	"log"

	"github.com/paulwuertz/pexplorer/selfperf/callgraph"
	"github.com/paulwuertz/pexplorer/selfperf/rtos"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func main() {
	// elfFile, err := elf.Open("/home/paul/git/Prusa-Firmware-Buddy/build/mini_release_noboot/firmware")
	// elfFile, err := elf.Open("/home/paul/git/cannecti/cannectivity/build_1.2/zephyr/zephyr.elf")
	// elfFile, err := elf.Open("/home/paul/git/cannecti/build_13_gcc_lpc55s16/zephyr/zephyr.elf")
	elfFile, err := elf.Open("/home/paul/git/pexplorer/testdata/elf_testdata/zswatch_nrf5340_070.elf")
	// elfFile, err := elf.Open("/home/paul/git//ztest/build_mqtt_publisher_frdm_k64f_42/zephyr/zephyr.elf")

	if err != nil {
		log.Fatal(err)
	}

	elfReport := symbolextraction.GetFWReport(elfFile)
	callgraph.EnhanceByDisasm(&elfReport)
	callgraph.TraverseCallGraph(&elfReport)
	rtos.ScanForRtosFeatures(&elfReport)

	// datajson, _ := json.MarshalIndent(elfReport, "", "    ")
	// fmt.Println(string(datajson))
	// fmt.Println("functions found:", len(elfReport.Functions))
	// fmt.Println("variables found:", len(elfReport.Variables))
	// fmt.Println("sections found:", len(elfReport.Sections))
}
