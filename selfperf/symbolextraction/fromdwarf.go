package symbolextraction

import (
	"debug/dwarf"
	"debug/elf"
	"encoding/binary"
	"fmt"
	"log"
	"maps"
	"slices"
	"sort"

	"github.com/go-delve/delve/pkg/dwarf/frame"
	"github.com/go-delve/delve/pkg/dwarf/godwarf"
	"github.com/go-delve/delve/pkg/dwarf/line"
)

var RegRuleEnum2String = map[frame.Rule]string{
	frame.RuleUndefined:     "RuleUndefined",
	frame.RuleSameVal:       "RuleSameVal",
	frame.RuleOffset:        "RuleOffset",
	frame.RuleValOffset:     "RuleValOffset",
	frame.RuleRegister:      "RuleRegister",
	frame.RuleExpression:    "RuleExpression",
	frame.RuleValExpression: "RuleValExpression",
	frame.RuleArchitectural: "RuleArchitectural",
	frame.RuleCFA:           "RuleCFA",
	frame.RuleFramePointer:  "RuleFramePointer",
}

func setLineInfo(elf *elf.File, funcs []FunctionSymbol, vars []VariableSymbol) {
	// TODO try Lookup by Address using .debug_aranges.
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
	// create a searchable address to source map
	keys := make([]uint64, 0, len(addr2fileMap))
	for k := range addr2fileMap {
		keys = append(keys, k)
	}
	slices.Sort(keys)
	// find the closest addr for each function and set it
	for i := 0; i < len(funcs); i++ {
		function := &funcs[i]
		pos, _ := slices.BinarySearch(keys, function.Address) // exact
		if pos == len(keys) {
			pos = len(funcs) - 1
		}
		closest_addr := keys[pos]
		dbl := addr2fileMap[closest_addr]
		fn, ln := dbl.PCToLine(function.Address, function.Address)
		function.SourceFilePath = fn
		function.SourceFileLine = uint64(ln)
		// fmt.Println("\tfile ", exact, fn, ln, function.Address, closest_addr, "diff", int64(function.Address)-int64(closest_addr), function.Name)
	}
}

func getVariableTypes(elf *elf.File, vars []VariableSymbol) []Typedef {
	typeMap := make(map[string]Typedef, 0)
	dwarfData, _ := elf.DWARF()
	rd := dwarfData.Reader()
	var cache map[dwarf.Offset]godwarf.Type = make(map[dwarf.Offset]godwarf.Type, 0)
	var varMap map[string]*VariableSymbol = make(map[string]*VariableSymbol, len(vars))
	for i := 0; i < len(vars); i++ {
		v := &vars[i]
		varMap[v.Name] = v
	}
	for idx := 0; ; idx++ {
		entry, err := rd.Next()
		if err != nil {
			// return fmt.Errorf("iterate entry error: %v", err)
		}
		if entry == nil {
			break
		}
		// parse compilation unit
		if entry.Tag == dwarf.TagVariable {
			tree, err := godwarf.LoadTree(entry.Offset, dwarfData, 0)
			atoff, ok := entry.Val(dwarf.AttrType).(dwarf.Offset)
			// fmt.Println("tree -", tree.Val(dwarf.AttrName), tree.Val(dwarf.AttrType), err)
			if ok && err == nil {
				var varName string = ""
				ty, _ := godwarf.ReadType(dwarfData, 0, atoff, cache)
				varName, ok := tree.Val(dwarf.AttrName).(string)
				if !ok {
					continue
				}
				typeStr := ty.Common().Name
				typeMap[typeStr] = Typedef{Name: typeStr, Size: uint64(ty.Common().ByteSize)}
				v, ok := varMap[varName]
				if ok && typeStr != "" {
					v.VariableType = ty.Common().Name
					// fmt.Println("\t- type info found", v, entry.Tag.String())
				} else {
					// ?
				}
			} else {
				// fmt.Println("\t- no type info found", entry.Tag.String())
			}
		}
	}
	return slices.Collect(maps.Values(typeMap))
}

func getStackUseDetails(elf *elf.File, funcs []FunctionSymbol, vars []VariableSymbol) (srcFiles []string) {
	dwarfData, _ := elf.DWARF()
	rd := dwarfData.Reader()
	for idx := 0; ; idx++ {
		entry, err := rd.Next()
		if err != nil {
			// return fmt.Errorf("iterate entry error: %v", err)
		}
		if entry == nil {
			return nil
		}

		// parse compilation unit
		if entry.Tag == dwarf.TagTypedef {
			tree, err := godwarf.LoadTree(entry.Offset, dwarfData, 0)
			fmt.Println("tree -", tree.Val(dwarf.AttrName), tree.Val(dwarf.AttrType), err)
			ty, err := godwarf.ReadType(dwarfData, int(entry.Offset), 0, nil)
			fmt.Println("\t-", ty, err)
		}
	}

	framedata, _ := godwarf.GetDebugSectionElf(elf, "frame")
	fe, err := frame.Parse(framedata, binary.LittleEndian, 0, 4, 0)
	if err != nil {
		log.Fatal("could not parse frame data of elffile", err)
	}
	sort.Slice(fe, func(i, j int) bool {
		return fe[i].Begin() < fe[j].Begin()
	})

	for i := 0; i < len(funcs); i++ {
		function := &funcs[i]
		mainfde, err := fe.FDEForPC(function.Address)
		fmt.Println("\t\tfn", err)
		if err != nil {
			continue
		}
		fmt.Println("\t\tfn", function.Name, mainfde.Length)
		for i := mainfde.Begin(); i < mainfde.End(); i = i + 2 {
			// for ARM the return addr is saved in r13,
			// when CFA is in R13 and only there, then CFA.off seems to be pretty much the frame size...
			// if not and the CFA reg changes and the SP is pushed around we do not know from this table...
			s, err := mainfde.EstablishFrame(i)
			if err != nil {
				continue
			}
			fmt.Println(fmt.Sprintf("%x", i), "-> cfa/off/rule:", s.CFA.Reg, "/", s.CFA.Offset, "/", RegRuleEnum2String[s.CFA.Rule], " expr:", s.CFA.Expression, "regs:", s.Regs, " reta:", s.RetAddrReg)
		}
	}
	return
}

func EnhanceByDwarfDebugInfo(s *SElfReport) (srcFiles []string) {
	s.Types = getVariableTypes(s.Elf, s.Variables)
	getStackUseDetails(s.Elf, s.Functions, s.Variables)
	setLineInfo(s.Elf, s.Functions, s.Variables)
	return
}
