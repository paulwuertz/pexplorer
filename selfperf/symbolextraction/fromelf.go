package symbolextraction

// import "symbolextraction/types"

import (
	"debug/elf"
)

func extractFunctions(elfFile elf.File) []FunctionSymbol {
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

		address := sym.Value
		functions = append(functions, FunctionSymbol{
			Name:              sym.Name,
			UnmangledName:     "",
			Address:           address,
			FlashSize:         sym.Size,
			FunctionStackSize: 0,
			SourceFilePath:    "",
			SourceFileLine:    0,
		})
		// fmt.Println(fmt.Sprintf("fun %s at %x", sym.Name, address), sym)
	}
	return functions
}
