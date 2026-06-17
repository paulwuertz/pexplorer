package main

import (
	"debug/elf"
	"log"
	"testing"

	"github.com/paulwuertz/pexplorer/selfperf/callgraph"
	"github.com/paulwuertz/pexplorer/selfperf/rtos"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func BenchmarkReportGen(b *testing.B) {
	benchmarks := []struct {
		firmwareFile string
	}{
		{"../../../testdata/elf_testdata/zswatch_nrf5340_07.elf"},
		{"../../../testdata/elf_testdata/Pinecilv2_EN_v2_23.elf"},
		{"../../../testdata/elf_testdata/Pinecilv1_EN_v2_23.elf"},
		{"../../../testdata/elf_testdata/zephyr_cannectivity_12_llvm_lpc55s16.elf"},
		{"../../../testdata/elf_testdata/zephyr_cannectivity_13_llvm_lpc55s16.elf"},
		{"../../../testdata/elf_testdata/zephyr_cannectivity_13_gcc_lpc55s16.elf"},
		{"../../../testdata/elf_testdata/zephyr_cannectivity_12_gcc_lpc55s16.elf"},
		{"../../../testdata/elf_testdata/prusa_buddy_boot_64.elf"},
	}

	for _, bm := range benchmarks {
		b.Run(bm.firmwareFile,
			func(b *testing.B) {
				elfFile, err := elf.Open(bm.firmwareFile)
				if err != nil {
					log.Fatal(err)
				}
				elfReport := symbolextraction.GetFWReport(elfFile, "")
				callgraph.EnhanceByDisasm(&elfReport)
				rtos.ScanForRtosFeatures(&elfReport)
			},
		)
	}
}
