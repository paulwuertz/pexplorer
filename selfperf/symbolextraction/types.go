package symbolextraction

import (
	"debug/dwarf"
	"debug/elf"
)

type ElfSection struct {
	Name    string `json:"name"`
	Address uint64 `json:"address"`
	Size    uint64 `json:"size"`
	Index   uint8  `json:"index"`
}

type FunctionSymbol struct {
	Name              string            `json:"name"`
	Address           uint64            `json:"address,omitempty,omitzero"`
	FlashSize         uint64            `json:"size"`
	FunctionStackSize uint64            `json:"stack_size,omitempty,omitzero"`
	SectionIndex      uint8             `json:"secidx"`
	SourceFilePath    string            `json:"file,omitempty,omitzero"`
	SourceFileLine    uint64            `json:"line,omitempty,omitzero"`
	Variables         []*VariableSymbol `json:"vars,omitempty,omitzero"`
	Asm               []byte            `json:"asm,omitempty,omitzero"`
	// calls
	entry     *dwarf.Entry      `json:"-"`
	variables []*VariableSymbol `json:"-"`
	cu        *CompileUnit      `json:"-"`
}

type VariableSymbol struct {
	Name           string       `json:"name"`
	Address        uint64       `json:"address,omitempty,omitzero"`
	FlashSize      uint64       `json:"size"`
	SectionIndex   uint8        `json:"secidx"`
	SourceFilePath string       `json:"file,omitempty,omitzero"`
	SourceFileLine uint64       `json:"line,omitempty,omitzero"`
	VariableType   string       `json:"type,omitempty,omitzero"`
	Data           []byte       `json:"staticInitData,omitempty,omitzero"`
	cu             *CompileUnit `json:"-"`
}

type Typedef struct {
	Name    string       `json:"name"`
	Size    uint64       `json:"size"`
	Members []Typedef    `json:"members,omitempty,omitzero"`
	cu      *CompileUnit `json:"-"`
}

// CompileUnit represents a compilation unit,
// including a series of source files and function definitions
type CompileUnit struct {
	Source    []string
	Functions []*FunctionSymbol
	Variables []*VariableSymbol
}

type SElfReport struct {
	SingleFirmware     bool             `json:"singlefirmware"`
	FirmwareIdentifier string           `json:"firmwareID"`
	Timestamp          string           `json:"timestamp"`
	Elf                *elf.File        `json:"-"`
	CompileUnits       []CompileUnit    `json:"compile_units"`
	Sections           []ElfSection     `json:"sections"`
	Functions          []FunctionSymbol `json:"functions"`
	Variables          []VariableSymbol `json:"variables"`
	Types              []Typedef        `json:"types,omitempty,omitzero"`
}
