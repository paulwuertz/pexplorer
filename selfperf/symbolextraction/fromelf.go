package symbolextraction

// import "symbolextraction/types"

import (
	"debug/elf"
	"fmt"
	"io"

	"github.com/ianlancetaylor/demangle"
)

type FlashSectionMap map[string]*elf.Section

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
			SectionIndex:      uint8(sym.Section),
			// section
			// asm
		})
		// fmt.Println(fmt.Sprintf("fun %s at %x", sym.Name, address), sym)
	}
	return functions
}

func ExtractSections(elfFile elf.File) (sections []ElfSection, secRefs FlashSectionMap) {
	secData := elfFile.Sections
	secRefs = make(FlashSectionMap)

	// build a symbol map from the symbol section
	// this should be always present...
	for i, section := range secData {
		sections = append(sections, ElfSection{
			Name:    section.Name,
			Address: section.Addr,
			Size:    section.Size,
			Index:   uint8(i),
		})
		if section.Addr != 0 {
			secRefs[section.Name] = section
		}
		// fmt.Println(fmt.Sprintf("fun %s at %x", sym.Name, address), sym)
	}
	return
}

func (fm FlashSectionMap) getSectionByName(name string) (sec *elf.Section) {
	sec, ok := fm[name]
	if !ok {
		sec, ok = fm["."+name]
	}
	return // might be nil anyway if section is not present...
}

func AddASMToFunctions(syms []FunctionSymbol, fm FlashSectionMap) {
	sec := fm.getSectionByName("text")
	sr := sec.Open()
	textStartAddr := sec.Addr // todo addr or offset?
	textEndAddr := sec.Addr + sec.Size
	fmt.Println("sec ", sec)
	fmt.Println("sec Size", sec.Size)
	for i, _ := range syms {
		sym := &syms[i]
		addr, size := sym.Address, sym.FlashSize
		if addr >= textStartAddr && addr <= textEndAddr {
			symSecOffset := int64(addr - textStartAddr)
			sr.Seek(symSecOffset, io.SeekStart)
			// fmt.Println(sym, "@:", addr, "newPos ", newPos)
			sym.Asm = make([]byte, size)
			nb, err := sr.Read(sym.Asm)
			if err != nil {
				fmt.Println("error reading asm-bytes", nb, "/", size, err, sym)
			}
			// fmt.Println("bytes read", nb, "/", size, err)
			// fmt.Println(fmt.Sprintf("%X", b))
		} else {
			fmt.Println(sym, " asm outside text section")
		}
	}
}
