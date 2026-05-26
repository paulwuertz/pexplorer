package symbolextraction

import (
	"debug/dwarf"
	"debug/elf"
	"fmt"
	"log"
	"maps"
	"reflect"
	"slices"

	"github.com/go-delve/delve/pkg/dwarf/frame"
	"github.com/go-delve/delve/pkg/dwarf/godwarf"
	"github.com/go-delve/delve/pkg/dwarf/line"
	"github.com/go-delve/delve/pkg/dwarf/reader"
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
	linedata_str, _ := godwarf.GetDebugSectionElf(elf, "line_str")
	debugLines := line.ParseAll(linedata, linedata_str, nil, 0, true, 4)

	// map all possible src breakpoint adresses to its dbgline compilation unit
	var addr2fileMap map[uint64]*line.DebugLineInfo = make(map[uint64]*line.DebugLineInfo)
	for _, dbl := range debugLines {
		// assume 32-bi addresses - hack for getting all addresses
		addresses, err := dbl.AllPCsBetween(0, 0xFFFF_FFFF, "", 0)
		if err != nil {
			log.Fatal("err it all pcs", err)
		}
		for _, addr := range addresses {
			_, ok := addr2fileMap[addr]
			if ok { // TODO understand and handle duplicates
				// fmt.Println("duplicates address in", addr, "already", "\n\t", al.FileNames[0], "and\n\t", dbl.FileNames[0])
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

	// rangesdata, _ := godwarf.GetDebugSectionElf(elf, "aranges")
	// // find the closest addr for each function and set it
	// for i := 0; i < len(vars); i++ {
	// 	function := &funcs[i]
	// 	pos, _ := slices.BinarySearch(keys, function.Address) // exact
	// 	if pos == len(keys) {
	// 		pos = len(funcs) - 1
	// 	}
	// 	closest_addr := keys[pos]
	// 	dbl := addr2fileMap[closest_addr]
	// 	fn, ln := dbl.PCToLine(function.Address, function.Address)
	// 	function.SourceFilePath = fn
	// 	function.SourceFileLine = uint64(ln)
	// 	// fmt.Println("\tfile ", exact, fn, ln, function.Address, closest_addr, "diff", int64(function.Address)-int64(closest_addr), function.Name)
	// }
}

func ExtractType(typeRef *godwarf.Type, cache map[string]Typedef, s *SElfReport) Typedef {
	name := (*typeRef).Common().Name
	switch (*typeRef).(type) {
	case *godwarf.StructType:
		structRef := (*typeRef).(*godwarf.StructType)
		typestr := structRef.Kind + " " + structRef.StructName
		typedef := Typedef{Name: name, Type: typestr, Size: structRef.Common().ByteSize, Members: make([]Typedef, 0)}
		cachedef, ok := cache[typestr]
		if ok {
			return cachedef
		} else {
			cache[typestr] = Typedef{}
		}
		for _, f := range structRef.Field {
			member := ExtractType(&f.Type, cache, s)
			member.ByteOffset = f.ByteOffset
			member.BitOffset = f.BitOffset
			member.BitSize = f.BitSize
			member.Name = f.Name
			typedef.Members = append(typedef.Members, member)
			// fmt.Println("\tStructType f", f)
		}
		cache[typestr] = typedef
		return typedef
	case godwarf.Type:
		if name == "" {
			switch (*typeRef).(type) {
			case *godwarf.PtrType:
				ptrRef := (*typeRef).(*godwarf.PtrType)
				ptrData := ExtractType(&ptrRef.Type, cache, s)
				ptrData.IsPointer = true
				ptrData.Size = (*ptrRef).Common().ByteSize
				return ptrData
			case *godwarf.QualType:
				ptrRef := (*typeRef).(*godwarf.QualType)
				return ExtractType(&ptrRef.Type, cache, s)
			case *godwarf.ArrayType:
				arrRef := (*typeRef).(*godwarf.ArrayType)
				return ExtractType(&arrRef.Type, cache, s)
			default:
				msg := fmt.Sprintf("uncatched unnammed Type: %s", reflect.TypeOf(typeRef).Name())
				s.Info = append(s.Info, msg)
				return Typedef{}
			}
		}
		return Typedef{Type: name, Size: (*typeRef).Common().ByteSize, Members: make([]Typedef, 0)}
	default:
		msg := fmt.Sprintf("uncatchedType: %s of %s", reflect.TypeOf(typeRef).Name(), name)
		s.Info = append(s.Info, msg)
		return Typedef{}
	}
}

func getVariableTypes(s *SElfReport) []Typedef {

	var dwarfData, _ = s.Elf.DWARF()
	var rd = dwarfData.Reader()
	var delveReader = reader.New(dwarfData)
	var funcs = s.Functions
	var vars = s.Variables
	var cache map[dwarf.Offset]godwarf.Type = make(map[dwarf.Offset]godwarf.Type, 0)
	var varMap map[uint64]*VariableSymbol = make(map[uint64]*VariableSymbol, len(vars))
	var fnMap map[uint64]*FunctionSymbol = make(map[uint64]*FunctionSymbol, len(vars))
	var typeMap = make(map[string]Typedef, 0)
	var curCompileUnit *CompileUnit
	var curFunction *FunctionSymbol

	// create maps by address
	for i := 0; i < len(vars); i++ {
		varMap[vars[i].Address] = &vars[i]
	}
	for i := 0; i < len(funcs); i++ {
		fnMap[funcs[i].Address] = &funcs[i]
	}
	// iterate over debug data
	for idx := 0; ; idx++ {
		entry, err := rd.Next()
		if err != nil {
			// return fmt.Errorf("iterate entry error: %v", err)
		}
		if entry == nil {
			break
		}
		// parse compilation unit
		if entry.Tag == dwarf.TagCompileUnit {
			lrd, err := dwarfData.LineReader(entry)
			if err != nil {
				log.Fatal("no lrd")
			}

			cu := CompileUnit{}
			curCompileUnit = &cu
			cu.Source = make([]string, len(lrd.Files()))

			// record the files contained in this compilation unit
			for i, v := range lrd.Files() {
				if v == nil {
					continue
				}
				cu.Source[i] = v.Name
			}
			s.CompileUnits = append(s.CompileUnits, cu)
		}

		var filename string = ""
		// TODO +AttrName and friends?!
		cuFileIndex, fileNameFound := entry.Val(dwarf.AttrDeclFile).(int64)
		cuFileLine, _ := entry.Val(dwarf.AttrDeclLine).(int64)
		name, isNameOk := entry.Val(dwarf.AttrName).(string)
		if fileNameFound {
			filename = curCompileUnit.Source[cuFileIndex-1]
		}

		// pare subprogram
		if entry.Tag == dwarf.TagSubprogram {
			addr, ok := entry.Val(dwarf.AttrLowpc).(uint64)
			fn, ok := fnMap[addr]
			if !ok {
				if addr == 0 || !isNameOk {
					// might be a virtual fn, ignore for now...
					// TODO any reason to export them as well?
					continue
				} else {
					// TODO sym in dwarf is offset by 1 in addr - why
					// and are there syms exclusivly in dwarf or elf?
					addrBegin, ok1 := entry.Val(dwarf.AttrLowpc).(uint64)
					addrEnd, ok2 := entry.Val(dwarf.AttrHighpc).(int64)
					if ok1 && ok2 {
						size := int64(addrBegin) + addrEnd + 1
						funcs = append(funcs, FunctionSymbol{
							Name:           name,
							Address:        addrBegin,
							FlashSize:      uint64(size),
							SourceFilePath: filename,
							SourceFileLine: uint64(cuFileLine),
							// SectionIndex:      uint8(sym.Section),
							// asm
						})
						continue
					}
					log.Fatal(name, "could not get mapped to fn:", entry, "at addr", addr)
				}
			}
			fn.SourceFilePath = filename
			curFunction = fn
			curCompileUnit.Functions = append(curCompileUnit.Functions, fn)
		} else {
			curFunction = nil
		}

		// parse variable
		if entry.Tag == dwarf.TagVariable {
			var varRef *VariableSymbol
			var found bool
			// get var addr by name and at file info
			if isNameOk {
				if name == "z_interrupt_stacks" {
					fmt.Println("entry var withot fn")
				}
				addr, err := delveReader.AddrFor(name, 0, 4)
				if err != nil {
					for _, v := range s.Variables {
						if v.Name == name {
							addr = v.Address
							break
						}
					}
				}
				varRef, found = varMap[addr]
				if found {
					// Variables do not come with an file attribute
					// so use the current compilation units main file
					varRef.SourceFilePath = curCompileUnit.Source[1]
					varRef.SourceFileLine = uint64(cuFileLine)
				}
			}
			// get type info
			atoff, ok := entry.Val(dwarf.AttrType).(dwarf.Offset)
			if ok {
				typeRef, err := godwarf.ReadType(dwarfData, 0, atoff, cache)
				if err != nil {
					continue
				}

				if varRef != nil {
					// TODO add test with samples
					// fmt.Println("entry var withot fn", entry)
				}
				typeStr := typeRef.Common().Name
				// try to get the type
				if typeStr != "" {
					// https://github.com/ARM-software/abi-aa/blob/main/aaelf32/aaelf32.rst
					// https://github.com/ARM-software/abi-aa/blob/main/aadwarf32/aadwarf32.rst
					// https://github.com/ARM-software/abi-aa/blob/main/aapcs32/aapcs32.rst#the-base-procedure-call-standard
					ExtractType(&typeRef, typeMap, s)
					if varRef != nil {
						varRef.VariableType = typeStr
						if curFunction != nil {
							curFunction.Variables = append(curFunction.Variables, varRef)
						} else {
							// cu
							// TODO delete or do smth?
							// msg := fmt.Sprintf("entry var withot fn: %s", entry)
							// s.Info = append(s.Info, msg)
						}
					}
				} else {
					t := ExtractType(&typeRef, typeMap, s)
					if t.Name != "" && varRef != nil {
						varRef.VariableType = t.Name
					} else {
						// ?
					}
				}
			} else {
				// fmt.Println("\t- no type info found", entry.Tag.String())
			}
		}
	}
	return slices.Collect(maps.Values(typeMap))
}

func EnhanceByDwarfDebugInfo(s *SElfReport) (srcFiles []string) {
	s.Types = getVariableTypes(s)
	setLineInfo(s.Elf, s.Functions, s.Variables)
	return
}
