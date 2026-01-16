package symbolextraction

import (
	"debug/elf"
	"fmt"
	"log"
	"slices"

	"github.com/go-delve/delve/pkg/dwarf/godwarf"
	"github.com/go-delve/delve/pkg/dwarf/line"
)

func setLineInfo(elf *elf.File, funcs []FunctionSymbol, vars []VariableSymbol) {
	linedata, _ := godwarf.GetDebugSectionElf(elf, "line")
	debugLines := line.ParseAll(linedata, nil, nil, 0, true, 4)

	// map all possible src breakpoint adresses to its dbgline compilation unit
	var addr2fileMap map[uint64]*line.DebugLineInfo = make(map[uint64]*line.DebugLineInfo)
	for _, dbl := range debugLines {
		// assume 32-bi addresses - hack for getting all addresses
		addresses, err := dbl.AllPCsBetween(0, 0xFFFF_FFFF, "", 0)
		if err != nil {
			log.Fatal("err it all pcs", err)
		}
		for _, addr := range addresses {
			al, ok := addr2fileMap[addr]
			if ok { // TODO understand and handle duplicates
				fmt.Println("duplicates address in", addr, "already", "\n\t", al.FileNames[0], "and\n\t", dbl.FileNames[0])
			}
			addr2fileMap[addr] = dbl
		}
		// fmt.Println("\tfile ", dbl.FileNames[0], len(addresses), addresses)
	}
	// create a searchable address map
	keys := make([]uint64, 0, len(addr2fileMap))
	for k := range addr2fileMap {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	// find the closest addr for each function and set it
	for i := 0; i < len(funcs); i++ {
		function := &funcs[i]
		pos, _ := slices.BinarySearch(keys, function.Address) // exact
		closest_addr := keys[pos]
		dbl := addr2fileMap[closest_addr]
		fn, ln := dbl.PCToLine(function.Address, function.Address)
		function.SourceFilePath = fn
		function.SourceFileLine = uint64(ln)
		// fmt.Println("\tfile ", exact, fn, ln, function.Address, closest_addr, "diff", int64(function.Address)-int64(closest_addr), function.Name)
	}
}

func EnhanceByDwarfDebugInfo(elf *elf.File, funcs []FunctionSymbol, vars []VariableSymbol) (srcFiles []string) {
	setLineInfo(elf, funcs, vars)
	return
}
