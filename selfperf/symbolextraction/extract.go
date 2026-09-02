package symbolextraction

import "debug/elf"

func GetFWReport(elfFile *elf.File, fwhash string) SElfReport {

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
		Name2FnMap:   map[string]*FunctionSymbol{},
		SectionsMap:  map[uint64]*ElfSection{},
		Architecture: elfFile.FileHeader.Machine.String(),
		FirmwareHash: fwhash,
	}
	for i := 0; i < len(elfReport.Functions); i++ {
		f := &elfReport.Functions[i]
		elfReport.Addr2FnMap[f.Address] = f
		elfReport.Name2FnMap[f.Name] = f
	}
	for i := 0; i < len(elfReport.Sections); i++ {
		s := &elfReport.Sections[i]
		elfReport.SectionsMap[s.Index] = s
	}
	EnhanceByDwarfDebugInfo(&elfReport)
	return elfReport
}
