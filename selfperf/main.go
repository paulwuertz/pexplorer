package main

import (
	"debug/dwarf"
	"debug/elf"
	"fmt"
	"io"
	"log"
)

type Symbol struct {
	Name              string
	UnmangledName     string
	Address           uint64
	FlashSize         uint64
	FunctionStackSize uint64
	SourceFilePath    string
	SourceFileLine    uint64
	// SymbolType
	// (asm code)
	// (variable type)
	// calls
}

func main() {
	// elfFile, err := elf.Open("/home/paul/git/Prusa-Firmware-Buddy/build/mini_release_noboot/firmware")
	elfFile, err := elf.Open("/home/paul/git/ztest/build_hello_world_frdm_k64f_42/zephyr/zephyr.elf")

	if err != nil {
		log.Fatal(err)
	}

	dwarfData, err := elfFile.DWARF()
	if err != nil {
		log.Fatal(err)
	}

	entryReader := dwarfData.Reader()

	for {
		// Read all entries in sequence
		entry, err := entryReader.Next()
		if entry == nil || err == io.EOF {
			// We've reached the end of DWARF entries
			break
		}

		// Check if this entry is a function
		if entry.Tag == dwarf.TagSubprogram {

			// Go through fields
			for _, field := range entry.Field {

				if field.Attr == dwarf.AttrName {
					fmt.Println(field.Val.(string))
				}
			}
		}
	}

	symData, _ := elfFile.Symbols()

	for _, sec := range elfFile.Sections {
		sh := sec.SectionHeader
		fmt.Println(sh.Name, sh.Size)
	}

	i := 0
	symbols := make([]Symbol, 0)
	// dwarfData.

	for _, sym := range symData {
		if sym.Size == 0 {
			continue
		}
		i += 1
		symbols = append(
			symbols,
			Symbol{
				Name:              sym.Name,
				UnmangledName:     "",
				Address:           sym.Value,
				FlashSize:         sym.Size,
				FunctionStackSize: 0,
				SourceFilePath:    "",
				SourceFileLine:    0,
			},
		)
		fmt.Println(sym, sym.Value, fmt.Sprintf("%x", sym.Value))

	}
	// fmt.Println(dwarfData)
	fmt.Println(i, "symbols found")
}
