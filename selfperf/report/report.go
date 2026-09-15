package report

import (
	"crypto/sha256"
	"debug/elf"
	"fmt"

	"github.com/paulwuertz/pexplorer/selfperf/callgraph"
	"github.com/paulwuertz/pexplorer/selfperf/config"
	"github.com/paulwuertz/pexplorer/selfperf/rtos"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func GetReportAndRTOSStats(elfFile *elf.File, p config.PexplorerConfig) (symbolextraction.SElfReport, []config.RTOSThread) {
	fw_hash := sha256.Sum256([]byte(elfFile.Data.GoString()))
	fw_hash_str := fmt.Sprintf("%x", fw_hash)

	elfReport := symbolextraction.GetFWReport(elfFile, fw_hash_str)
	callgraph.EnhanceByDisasm(&elfReport, p.DynamicCalls)
	callgraph.TraverseCallGraph(&elfReport)
	threads := rtos.GetAllThreads(&elfReport, p)
	return elfReport, threads
}
