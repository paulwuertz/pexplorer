package symbolextraction

import "debug/elf"

func GetFWReport(elfFile *elf.File) SElfReport {

	sectionJsonInfo, sectionsRef := ExtractSections(*elfFile)
	functions := ExtractFunctions(*elfFile)
	variables := ExtractVariables(*elfFile)

	info := make([]string, 0)
	AddASMToFunctions(functions, sectionsRef, info)
	AddDataToVar(variables, sectionsRef, info)
	elfReport := SElfReport{
		Elf:          elfFile,
		Sections:     sectionJsonInfo,
		Functions:    functions,
		Variables:    variables,
		Info:         info,
		Addr2FnMap:   map[uint64]*FunctionSymbol{},
		SectionsMap:  map[uint64]*ElfSection{},
		Architecture: elfFile.FileHeader.Machine.String(),
	}
	for i := 0; i < len(elfReport.Functions); i++ {
		f := &elfReport.Functions[i]
		elfReport.Addr2FnMap[f.Address] = f
	}
	for i := 0; i < len(elfReport.Sections); i++ {
		s := &elfReport.Sections[i]
		elfReport.SectionsMap[s.Index] = s
	}
	EnhanceByDwarfDebugInfo(&elfReport)
	return elfReport
}
