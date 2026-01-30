package symbolextraction

import (
	"debug/elf"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/ianlancetaylor/demangle"
)

type SectionMaps struct {
	byName  map[string]*elf.Section
	byIndex map[uint8]*elf.Section
}

func ExtractFunctions(elfFile elf.File) []FunctionSymbol {
	functions := []FunctionSymbol{}
	symData, _ := elfFile.Symbols()

	for _, sym := range symData {
		symType := elf.ST_TYPE(sym.Info)
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
	// sort lowest address first
	sort.Slice(functions, func(i, j int) bool {
		return functions[i].Address < functions[j].Address
	})
	return functions
}

func ExtractVariables(elfFile elf.File) []VariableSymbol {
	variables := []VariableSymbol{}
	symData, _ := elfFile.Symbols()

	for _, sym := range symData {
		symType := elf.ST_TYPE(sym.Info)
		isVar := (elf.STT_OBJECT == symType)
		if (!isVar) || sym.Size == 0 {
			continue
		}

		name, err := demangle.ToString(sym.Name) // try to demangle
		if err != nil {
			name = sym.Name
		}

		address := sym.Value
		variables = append(variables, VariableSymbol{
			Name:         name,
			Address:      address,
			FlashSize:    sym.Size,
			SectionIndex: uint8(sym.Section),
		})
		// fmt.Println(fmt.Sprintf("fun %s at %x", sym.Name, address), sym)
	}
	// sort lowest address first
	sort.Slice(variables, func(i, j int) bool {
		return variables[i].Address < variables[j].Address
	})
	return variables
}

func ExtractSections(elfFile elf.File) (sections []ElfSection, secRefs SectionMaps) {
	secData := elfFile.Sections
	secRefs.byName = make(map[string]*elf.Section)
	secRefs.byIndex = make(map[uint8]*elf.Section)

	// build a symbol map from the symbol section
	// this should be always present...
	for i, section := range secData {
		sections = append(sections, ElfSection{
			Name:    section.Name,
			Address: section.Addr,
			Size:    section.Size,
			Index:   uint8(i),
		})
		secRefs.byName[section.Name] = section
		secRefs.byIndex[uint8(i)] = section
		// fmt.Println(fmt.Sprintf("fun %s at %x", sym.Name, address), sym)
	}
	return
}

func (fm SectionMaps) getSectionByName(name string) (sec *elf.Section) {
	sec, ok := fm.byName[name]
	if !ok {
		sec, ok = fm.byName["."+name]
	}
	return // might be nil anyway if section is not present...
}

func (fm SectionMaps) getSectionByIndex(i uint8) (sec *elf.Section) {
	sec = fm.byIndex[i]
	return // might be nil anyway if section is not present...
}

func AddASMToFunctions(syms []FunctionSymbol, fm SectionMaps, info []string) {
	sec := fm.getSectionByName("text")
	sr := sec.Open()
	textStartAddr := sec.Addr // todo addr or offset?
	textEndAddr := sec.Addr + sec.Size
	for i, _ := range syms {
		sym := &syms[i]
		addr, size := sym.Address, sym.FlashSize
		if addr >= textStartAddr && addr <= textEndAddr {
			symSecOffset := int64(addr - textStartAddr - 1)
			sr.Seek(symSecOffset, io.SeekStart)
			// fmt.Println(sym, "@:", addr, "newPos ", newPos)
			sym.Asm = make([]byte, size)
			nb, err := sr.Read(sym.Asm)
			if err != nil || nb != int(size) {
				msg := fmt.Sprintf("error reading asm for '%s' at %X, reading %d/%d", sym.Name, sym.Address, nb, size)
				info = append(info, msg)
			}
		} else {
			msg := fmt.Sprintf("asm for '%s' outside text section at %X", sym.Name, sym.Address)
			info = append(info, msg)
		}
	}
}

func AddDataToVar(syms []VariableSymbol, fm SectionMaps, info []string) {
	for i, _ := range syms {
		sym := &syms[i]
		addr, size := sym.Address, sym.FlashSize
		sec := fm.getSectionByIndex(sym.SectionIndex)
		if strings.Contains(sec.Name, "bss") {
			continue // will be zero data anyway...
		}
		sr := sec.Open()
		// fmt.Println(addr, size, sym.SectionIndex)
		symSecOffset := int64(addr - sec.Addr)
		sr.Seek(symSecOffset, io.SeekStart)
		// fmt.Println(sym, "@:", addr, "newPos ", newPos)
		sym.Data = make([]byte, size)
		nb, err := sr.Read(sym.Data)
		if err != nil || nb != int(size) {
			// fmt.Println("error reading data-bytes", nb, "/", size, err, sym)
		}
	}
}
