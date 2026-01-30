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
		Elf:        elfFile,
		Sections:   sectionJsonInfo,
		Functions:  functions,
		Variables:  variables,
		Info:       info,
		Addr2FnMap: map[uint64]*FunctionSymbol{},
	}
	EnhanceByDwarfDebugInfo(&elfReport)
	return elfReport
}
