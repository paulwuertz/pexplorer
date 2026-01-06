//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"debug/dwarf"
	"debug/elf"
	"fmt"
	"log"
	"syscall/js"
)

type ElfSections struct {
	Name    string
	Address uint64
	Size    uint64
}

type FunctionSymbol struct {
	Name              string
	UnmangledName     string
	Address           uint64
	FlashSize         uint64
	FunctionStackSize uint64
	SourceFilePath    string
	SourceFileLine    uint64
	Variables         []VariableSymbol
	// (asm code)
	// calls
}

type VariableSymbol struct {
	Name           string
	UnmangledName  string
	Address        uint64
	FlashSize      uint64
	SourceFilePath string
	SourceFileLine uint64
	VariableType   string
}

// CompileUnit represents a compilation unit,
// including a series of source files and function definitions
type CompileUnit struct {
	Source    []string
	functions []FunctionSymbol
	variables []VariableSymbol
}

type SElfReport struct {
	srcFiles  []CompileUnit
	sections  []ElfSections
	functions []FunctionSymbol
	variables []VariableSymbol
}

func read_sections_as_json(elfFile elf.File) {
	// elf, err := NewFile(r io.ReaderAt) (*File, error)

}

func extractSections(elfFile elf.File) any {
	var sections = make([]any, len(elfFile.Sections))
	for _, sec := range elfFile.Sections {
		sh := sec.SectionHeader
		sections = append(sections, sh.Name)
		sections = append(sections, sh.Size)
		sections = append(sections, sh.Addr)
		// append(sections, ElfSections{
		// 	sh.Name, sh.Size, sh.Addr,
		// })
	}
	return sections
}

func wasm_read_sections_as_json() js.Func {
	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		elf_bin_js := args[0]
		elf_size := elf_bin_js.Length()
		elf_binary := make([]byte, elf_size)

		fmt.Printf("args # %d\n", len(args))
		fmt.Printf("input %d\n", elf_size)
		bytes_copied := js.CopyBytesToGo(elf_binary, elf_bin_js)
		fmt.Printf("n %d bytes\n", bytes_copied)

		r := bytes.NewReader(elf_binary)
		elf, _ := elf.NewFile(r)
		return extractSections(*elf)
	})
	return jsonFunc
}

func extractFunctions(elfFile elf.File) []FunctionSymbol {
	functions := []FunctionSymbol{}
	fnSymByAddr := map[uint64]FunctionSymbol{}
	symData, _ := elfFile.Symbols()

	// build a symbol map from the symbol section
	// this should be always present...
	for _, sym := range symData {
		symType := elf.SymType(sym.Info & 0xf)
		isVariable := (elf.STT_OBJECT == symType)
		isFunc := (elf.STT_FUNC == symType)
		if (!isVariable && !isFunc) || sym.Size == 0 {
			continue
		}

		address := sym.Value
		if isVariable {
			fmt.Println(fmt.Sprintf("var %s at %x", sym.Name, address), sym)
		} else if isFunc {
			fnSymByAddr[address] = FunctionSymbol{
				Name:              sym.Name,
				UnmangledName:     "",
				Address:           address,
				FlashSize:         sym.Size,
				FunctionStackSize: 0,
				SourceFilePath:    "",
				SourceFileLine:    0,
			}
			fmt.Println(fmt.Sprintf("fun %s at %x", sym.Name, address), sym)
		}
	}
	return functions
}

func extractVariables(elfFile elf.File) []VariableSymbol {
	variables := []VariableSymbol{}
	// varSymByAddr := map[uint64]VariableSymbol{}
	return variables
}

func enhanceByDwarfDebugInfo(elfFile elf.File, functions []FunctionSymbol, variables []VariableSymbol) {
	var compileUnits = []*CompileUnit{}

	dw, _ := elfFile.DWARF()
	rd := dw.Reader()

	var curCompileUnit *CompileUnit
	var curFunction *FunctionSymbol

	for idx := 0; ; idx++ {
		entry, err := rd.Next()
		if err != nil {
			fmt.Println("next symbols found")
			// return fmt.Errorf("iterate entry error: %v", err)
		}
		if entry == nil {
			fmt.Println("last entry")
			// return nil
		}

		// parse compilation unit
		if entry.Tag == dwarf.TagCompileUnit {
			lrd, err := dw.LineReader(entry)
			if err != nil {
				// return err
			}

			cu := &CompileUnit{}
			curCompileUnit = cu

			// record the files contained in this compilation unit
			for _, v := range lrd.Files() {
				if v == nil {
					continue
				}
				cu.Source = append(cu.Source, v.Name)
			}
			compileUnits = append(compileUnits, cu)
		}

		// pare subprogram
		if entry.Tag == dwarf.TagSubprogram {
			fn := FunctionSymbol{
				Name:           entry.Val(dwarf.AttrName).(string),
				SourceFilePath: curCompileUnit.Source[entry.Val(dwarf.AttrDeclFile).(int64)-1],
			}
			curFunction = &fn
			curCompileUnit.functions = append(curCompileUnit.functions, fn)
		}

		// parse variable
		if entry.Tag == dwarf.TagVariable {
			variable := VariableSymbol{
				Name: entry.Val(dwarf.AttrName).(string),
			}
			curFunction.Variables = append(curFunction.Variables, variable)
		}
	}
	// return SElfReport{
	// 	srcFiles: compileUnits,
	// 	// sections: ,
	// 	functions: functions,
	// 	variables: variables,
	// }
}

func dissassembleASM() {

}

func extractStaticCallsFromASM() {

}

func extractStaticCallsFromBuildDir() {

}

// func extractStaticCallsFromASM() {

// }

// func extractStaticCallsFromBuildDir() {

// }

func main() {

	fmt.Println("Go Web Assembly")
	js.Global().Set("read_sections_as_json", wasm_read_sections_as_json())
	<-make(chan struct{})
	// elfFile, err := elf.Open("/home/paul/git/Prusa-Firmware-Buddy/build/mini_release_noboot/firmware")
	elfFile, err := elf.Open("/home/paul/git/ztest/build_hello_world_frdm_k64f_42/zephyr/zephyr.elf")

	if err != nil {
		log.Fatal(err)
	}

	// sections :=
	extractSections(*elfFile)
	functions := extractFunctions(*elfFile)
	variables := extractVariables(*elfFile)
	// srcFiles := enhanceByDwarfDebugInfo(*elfFile, *functions, *variables)
	// elfReport := SElfReport(srcFiles, sections, functions, variables)

	fmt.Println("functions found:", functions)
	fmt.Println("variables found:", variables)
	// fmt.Println(i, "symbols found", len(varSymByAddr))
	// fmt.Println(i, "symbols found", len(fnSymByAddr))
	// for key, value := range symByAddr {
	// 	fmt.Println("Key:", fmt.Sprintf("%x", key), "Value:", value)
	// }

}
