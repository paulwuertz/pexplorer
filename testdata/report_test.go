package main

// package testdata
import (
	"debug/elf"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"testing"

	"github.com/paulwuertz/pexplorer/selfperf/callgraph"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func ReportFromJsonFile(jsonfile string) symbolextraction.SElfReport {
	var report symbolextraction.SElfReport
	content, err := os.ReadFile(jsonfile)
	if err != nil {
		log.Fatal("Error when opening file: ", err)
	}

	// Now let's unmarshall the data into `payload`
	err = json.Unmarshal(content, &report)
	return report
}

func ReportFromElfFile(elffile string) symbolextraction.SElfReport {
	elfFile, err := elf.Open(elffile)

	if err != nil {
		log.Fatal(err)
	}

	elfReport := symbolextraction.GetFWReport(elfFile)
	callgraph.EnhanceByDisasm(&elfReport)
	return elfReport
}

func same_functions_extracted(json_rep, elf_rep symbolextraction.SElfReport) bool {
	return len(json_rep.Functions) != len(elf_rep.Functions)
}

func TestExtractReport(t *testing.T) {
	tests := []struct {
		name    string
		refFile string
		elfFile string
	}{
		// TODO: Add test cases.
		{
			"cannectivity_13_gcc_lpc55s16",
			"./ref_puncover_gcc/zephyr_cannectivity_13_gcc_lpc55s16.json",
			"./elf_testdata/zephyr_cannectivity_13_gcc_lpc55s16.elf",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			json_rep := ReportFromJsonFile(tt.refFile)
			elf_rep := ReportFromElfFile(tt.elfFile)

			nr_fn_matches := len(json_rep.Functions) == len(elf_rep.Functions)
			if !nr_fn_matches {
				t.Logf("#fns in json (%d) vs elf report(%d) not matching\n", len(json_rep.Functions), len(elf_rep.Functions))
			}

			matched := make(map[uint64]*symbolextraction.FunctionSymbol, 0)
			matchedSameAddr := make(map[uint64]*symbolextraction.FunctionSymbol, 0)
			nrSymsMatchedOneOff := 0

			t.Logf("functions in json but not in elf report =>\n")
			t.Logf("====================================================\n")
			for _, f := range json_rep.Functions {
				addr := f.Address
				addrOff1 := f.Address + 1
				fMatch, isInOtherReport := elf_rep.Addr2FnMap[addr]
				fMatch, isInOtherReportOff := elf_rep.Addr2FnMap[addrOff1]
				if isInOtherReport {
					matched[addr] = fMatch
					matchedSameAddr[addr] = fMatch
				} else if isInOtherReportOff {
					matched[addrOff1] = fMatch
					matchedSameAddr[addr] = fMatch
					nrSymsMatchedOneOff += 1
					//t.Logf("%s off by +1 %d in json compared to elf report\n", f.Name, addr)
				} else {
					t.Logf("%s at %d (0x%X)\n", f.Name, addr, addr)
				}
			}
			t.Logf("====================================================\n")
			t.Logf("%d syms are off by +1 in json compared to elf report\n", nrSymsMatchedOneOff)
			t.Logf("====================================================\n")
			t.Logf("functions only in elf report            =>\n")
			t.Logf("====================================================\n")
			for addr, f := range elf_rep.Addr2FnMap {
				_, inElfReport := matched[addr]
				if !inElfReport {
					t.Logf("%s at %d (0x%X)\n", f.Name, addr, addr)
				}
			}
			t.Logf("====================================================\n")
			t.Logf("functions different props   elf/json         =>\n")
			t.Logf("====================================================\n")
			nrSameFns := 0
			errs := map[string]int{
				"Address":          0,
				"FlashSize":        0,
				"StackSize":        0,
				"StackSizeToSmall": 0,
				"StackSizeToBig":   0,
				"StackSizeMatch":   0,
				"NoStackSize":      0,
				"Callees":          0,
				"Callers":          0,
				"Name":             0,
				"Src":              0,
			}
			for _, f := range json_rep.Functions {
				fMatch, isInOtherReport := matchedSameAddr[f.Address]
				if isInOtherReport {
					diffs := ""
					if f.Name == "main" {
						f.StackQualifiers = "estimated+experimental"

					}
					if fMatch.Address != f.Address {
						errs["Address"] += 1
						diffs += fmt.Sprintf("Address %d != %d, ", fMatch.Address, f.Address)
					}
					if fMatch.FlashSize != f.FlashSize {
						errs["FlashSize"] += 1
						diffs += fmt.Sprintf("FlashSize %d != %d, ", fMatch.FlashSize, f.FlashSize)
					}
					if fMatch.StackSize != 0 && f.StackSize != 0 && fMatch.StackSize != f.StackSize {
						errs["StackSize"] += 1
						if fMatch.StackSize > f.StackSize {
							errs["StackSizeToBig"] += 1
						} else {
							errs["StackSizeToSmall"] += 1
						}
						diffs += fmt.Sprintf("StackSize %d != %d, ", fMatch.StackSize, f.StackSize)
					} else if fMatch.StackSize == 0 && fMatch.StackSize != f.StackSize {
						errs["NoStackSize"] += 1
					} else if fMatch.StackSize != 0 && fMatch.StackSize == f.StackSize {
						errs["StackSizeMatch"] += 1
					}
					if len(fMatch.Callees) != len(f.Callees) {
						errs["Callees"] += 1
						diffs += fmt.Sprintf("#Callees %d != %d, ", len(fMatch.Callees), len(f.Callees))
					}
					if len(fMatch.Callers) != len(f.Callers) {
						errs["Callers"] += 1
						diffs += fmt.Sprintf("#Callers %d != %d, ", len(fMatch.Callers), len(f.Callers))
					}
					if fMatch.Name != f.Name {
						errs["Name"] += 1
						diffs += fmt.Sprintf("Name %s != %s, ", fMatch.Name, f.Name)
					}
					if fMatch.SourceFilePath != "/"+f.SourceFilePath {
						errs["Src"] += 1
						diffs += fmt.Sprintf("Src %s != %s, ", fMatch.SourceFilePath, f.SourceFilePath)
					}
					if diffs == "" {
						nrSameFns += 1
					} else {
						t.Logf("%s - %s\n", f.Name, diffs)
					}
				}
			}
			t.Logf("====================================================\n")
			t.Logf("# fns with same props %d/%d\n", nrSameFns, len(json_rep.Functions))
			v, _ := json.Marshal(errs)
			t.Logf("%s", string(v))
		})
	}
}

func main() {
	// first iteration
	// report_test.go:148: # fns with same props 289/583
	// report_test.go:153: {"Callees":0,"Callers":0,"FlashSize":4,"Name":24,"NoStackSize":371,"Src":267,"StackSize":0}

	// json_rep := ReportFromJsonFile("./ref_puncover_gcc/ref_puncover_gcc_zephyr_cannectivity_13_gcc_lpc55s16.json")
	// elf_rep := ReportFromElfFile("./ref_puncover_gcc/ref_puncover_gcc_zephyr_cannectivity_13_gcc_lpc55s16.json")
}
