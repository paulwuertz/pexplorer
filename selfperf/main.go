//go:build js && wasm
// +build js,wasm

package main

import (
	"debug/elf"
	"encoding/json"
	"fmt"
	"log"

	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

// func enhanceByDwarfDebugInfo(elfFile elf.File, functions []FunctionSymbol, variables []VariableSymbol) {
// 	var compileUnits = []*CompileUnit{}

// 	dw, _ := elfFile.DWARF()
// 	rd := dw.Reader()

// 	var curCompileUnit *CompileUnit
// 	var curFunction *FunctionSymbol

// 	for idx := 0; ; idx++ {
// 		entry, err := rd.Next()
// 		if err != nil {
// 			fmt.Println("next symbols found")
// 			// return fmt.Errorf("iterate entry error: %v", err)
// 		}
// 		if entry == nil {
// 			fmt.Println("last entry")
// 			// return nil
// 		}

// 		// parse compilation unit
// 		if entry.Tag == dwarf.TagCompileUnit {
// 			lrd, err := dw.LineReader(entry)
// 			if err != nil {
// 				// return err
// 			}

// 			cu := &CompileUnit{}
// 			curCompileUnit = cu

// 			// record the files contained in this compilation unit
// 			for _, v := range lrd.Files() {
// 				if v == nil {
// 					continue
// 				}
// 				cu.Source = append(cu.Source, v.Name)
// 			}
// 			compileUnits = append(compileUnits, cu)
// 		}

// 		// pare subprogram
// 		if entry.Tag == dwarf.TagSubprogram {
// 			fn := FunctionSymbol{
// 				Name:           entry.Val(dwarf.AttrName).(string),
// 				SourceFilePath: curCompileUnit.Source[entry.Val(dwarf.AttrDeclFile).(int64)-1],
// 			}
// 			curFunction = &fn
// 			curCompileUnit.functions = append(curCompileUnit.functions, fn)
// 		}

// 		// parse variable
// 		if entry.Tag == dwarf.TagVariable {
// 			variable := VariableSymbol{
// 				Name: entry.Val(dwarf.AttrName).(string),
// 			}
// 			curFunction.Variables = append(curFunction.Variables, variable)
// 		}
// 	}
// 	// return SElfReport{
// 	// 	srcFiles: compileUnits,
// 	// 	// sections: ,
// 	// 	functions: functions,
// 	// 	variables: variables,
// 	// }
// }

// func dissassembleASM() {}
// func extractStaticCallsFromASM() {}
// func extractStaticCallsFromBuildDir() {}
// func extractStaticCallsFromASM() {}
// func extractStaticCallsFromBuildDir() {}

func main() {
	// fmt.Println("Go Web Assembly")
	// js.Global().Set("read_sections_as_json", wasm_read_sections_as_json())
	// <-make(chan struct{})
	// elfFile, err := elf.Open("/home/paul/git/Prusa-Firmware-Buddy/build/mini_release_noboot/firmware")
	elfFile, err := elf.Open("/home/paul/git/cannecti/cannectivity/build_1.2/zephyr/zephyr.elf")

	if err != nil {
		log.Fatal(err)
	}

	sectionJsonInfo, sectionsRef := symbolextraction.ExtractSections(*elfFile)
	functions := symbolextraction.ExtractFunctions(*elfFile)
	variables := symbolextraction.ExtractVariables(*elfFile)

	symbolextraction.AddASMToFunctions(functions, sectionsRef)
	symbolextraction.AddDataToVar(variables, sectionsRef)
	elfReport := symbolextraction.SElfReport{
		Elf:       elfFile,
		Sections:  sectionJsonInfo,
		Functions: functions,
		Variables: variables,
	}
	symbolextraction.EnhanceByDwarfDebugInfo(&elfReport)
	datajson, _ := json.MarshalIndent(elfReport, "", "    ")

	fmt.Println("functions found:", string(datajson))
	fmt.Println("functions found:", len(functions))
	fmt.Println("variables found:", len(variables))
	fmt.Println("sections found:", len(sectionJsonInfo))
}
