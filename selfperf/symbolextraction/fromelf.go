package symbolextraction

// import "symbolextraction/types"

import (
	"debug/elf"

	"github.com/ianlancetaylor/demangle"
)

func ExtractFunctions(elfFile elf.File) []FunctionSymbol {
	functions := []FunctionSymbol{}
	symData, _ := elfFile.Symbols()

	// build a symbol map from the symbol section
	// this should be always present...
	for _, sym := range symData {
		symType := elf.SymType(sym.Info & 0xf)
		isFunc := (elf.STT_FUNC == symType)
		if (!isFunc) || sym.Size == 0 {
			continue
		}

		name, err := demangle.ToString(sym.Name) // try to demangle
		if err != nil {
			name = sym.Name
		}

		address := sym.Value
		functions = append(functions, FunctionSymbol{
			Name:              name,
			Address:           address,
			FlashSize:         sym.Size,
			FunctionStackSize: 0,
			SourceFilePath:    "",
			SourceFileLine:    0,
			// section
			// asm
		})
		// fmt.Println(fmt.Sprintf("fun %s at %x", sym.Name, address), sym)
	}
	return functions
}

func ExtractSections(elfFile elf.File) []ElfSection {
	sections := []ElfSection{}
	secData := elfFile.Sections

	// build a symbol map from the symbol section
	// this should be always present...
	for _, section := range secData {
		sections = append(sections, ElfSection{
			Name:    section.Name,
			Address: section.Addr,
			Size:    section.Size,
			//
		})
		// fmt.Println(fmt.Sprintf("fun %s at %x", sym.Name, address), sym)
	}
	return sections
}
