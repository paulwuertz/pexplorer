package symbolextraction

import "debug/elf"

func GetFWReport(elfFile *elf.File) SElfReport {

	sectionJsonInfo, sectionsRef := ExtractSections(*elfFile)
	functions := ExtractFunctions(*elfFile)
	variables := ExtractVariables(*elfFile)

	AddASMToFunctions(functions, sectionsRef)
	AddDataToVar(variables, sectionsRef)
	elfReport := SElfReport{
		Elf:       elfFile,
		Sections:  sectionJsonInfo,
		Functions: functions,
		Variables: variables,
	}
	EnhanceByDwarfDebugInfo(&elfReport)
	return elfReport
}
