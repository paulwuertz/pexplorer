package main

import (
	"debug/dwarf"
	"debug/elf"
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/go-delve/delve/pkg/dwarf/godwarf"
)

var basicDwarf2attr []dwarf.Attr = []dwarf.Attr{
	dwarf.AttrSibling,
	dwarf.AttrLocation,
	dwarf.AttrName,
	dwarf.AttrOrdering,
	dwarf.AttrByteSize,
	dwarf.AttrBitOffset,
	dwarf.AttrBitSize,
	dwarf.AttrStmtList,
	dwarf.AttrLowpc,
	dwarf.AttrHighpc,
	dwarf.AttrLanguage,
	dwarf.AttrDiscr,
	dwarf.AttrDiscrValue,
	dwarf.AttrVisibility,
	dwarf.AttrImport,
	dwarf.AttrStringLength,
	dwarf.AttrCommonRef,
	dwarf.AttrCompDir,
	dwarf.AttrConstValue,
	dwarf.AttrContainingType,
	dwarf.AttrDefaultValue,
	dwarf.AttrInline,
	dwarf.AttrIsOptional,
	dwarf.AttrLowerBound,
	dwarf.AttrProducer,
	dwarf.AttrPrototyped,
	dwarf.AttrReturnAddr,
	dwarf.AttrStartScope,
	dwarf.AttrStrideSize,
	dwarf.AttrUpperBound,
	dwarf.AttrAbstractOrigin,
	dwarf.AttrAccessibility,
	dwarf.AttrAddrClass,
	dwarf.AttrArtificial,
	dwarf.AttrBaseTypes,
	dwarf.AttrCalling,
	dwarf.AttrCount,
	dwarf.AttrDataMemberLoc,
	dwarf.AttrDeclColumn,
	dwarf.AttrDeclFile,
	dwarf.AttrDeclLine,
	dwarf.AttrDeclaration,
	dwarf.AttrDiscrList,
	dwarf.AttrEncoding,
	dwarf.AttrExternal,
	dwarf.AttrFrameBase,
	dwarf.AttrFriend,
	dwarf.AttrIdentifierCase,
	dwarf.AttrMacroInfo,
	dwarf.AttrNamelistItem,
	dwarf.AttrPriority,
	dwarf.AttrSegment,
	dwarf.AttrSpecification,
	dwarf.AttrStaticLink,
	dwarf.AttrType,
	dwarf.AttrUseLocation,
	dwarf.AttrVarParam,
	dwarf.AttrVirtuality,
	dwarf.AttrVtableElemLoc,
}

func indent(depth int) string {
	return strings.Repeat("\t", depth)
}

func dumpCompilationUnitTreeRecursive(roots []*godwarf.Tree, d int) {
	for i, c := range roots {
		cname, _ := c.Entry.Val(dwarf.AttrName).(string)
		fmt.Println(indent(d), i, c.Tag.GoString(), "name:", cname, "nr children:", len(c.Children))
		for _, a := range basicDwarf2attr {
			attr := c.Val(a)
			if attr != nil {
				fmt.Println(indent(d), "-", a, ":", attr)
			}
		}
		dumpCompilationUnitTreeRecursive(c.Children, d+1)
	}
}

func main() {
	// arg input, output file, indentation

	infile := flag.String("i", "", "input ELF file - obligatory")
	unitFilterName := flag.String("f", "", "filter to show only the dwarf tree of compilation unita containing this string in its pathname")

	flag.Parse()
	if *infile == "" {
		log.Fatal("Please add an ELF file to generate a report for.")
	}
	elfFile, err := elf.Open(*infile)

	var dwarfData, _ = elfFile.DWARF()
	var rd = dwarfData.Reader()

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
			tree, _ := godwarf.LoadTree(entry.Offset, dwarfData, 0)
			var root []*godwarf.Tree = []*godwarf.Tree{tree}
			var cname, ok = tree.Entry.Val(dwarf.AttrName).(string)
			if *unitFilterName == "" || ok && strings.Contains(cname, *unitFilterName) {
				dumpCompilationUnitTreeRecursive(root, 0)
			}
		}
	}

	if err != nil {
		log.Fatal(err)
	}

}
