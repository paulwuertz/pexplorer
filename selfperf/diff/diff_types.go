package diff

import (
	"github.com/paulwuertz/pexplorer/selfperf/config"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

// TODO implement :)?
type DiffSettings struct {
	ShowFilePath      bool
	ShowStackDiff     bool
	ShowAllSymbolDiff bool
}

type ElfReport struct {
	Report  symbolextraction.SElfReport
	Threads []config.RTOSThread
}

type ReportMap map[string]ElfReport
