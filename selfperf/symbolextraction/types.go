package symbolextraction

import (
	"debug/dwarf"
	"debug/elf"
)

type ElfSection struct {
	Name    string `json:"name"`
	Address uint64 `json:"address"`
	Size    uint64 `json:"size"`
	RamSize uint64 `json:"ram_size"`
	RomSize uint64 `json:"rom_size"`
	Index   uint64 `json:"index"`
}

type FunctionSymbol struct {
	Name                string            `json:"name"`
	Address             uint64            `json:"address,omitempty,omitzero"`
	FlashSize           uint64            `json:"size"`
	SectionIndex        uint64            `json:"secidx"`
	SourceFilePath      string            `json:"file,omitempty,omitzero"`
	SourceFileLine      uint64            `json:"line,omitempty,omitzero"`
	Variables           []*VariableSymbol `json:"vars,omitempty,omitzero"`
	Asm                 []byte            `json:"asm,omitempty,omitzero"`
	Callees             []FunctionCall    `json:"callees,omitempty,omitzero"`
	Callers             []FunctionCall    `json:"callers,omitempty,omitzero"`
	StackSize           int64             `json:"stack_size"`
	MaxStackSizeCallees int64             `json:"max_stack_size_callees,omitempty,omitzero"`
	StackQualifiers     string            `json:"stack_qualifiers,omitempty,omitzero"`
	// calls + refs
	entry     *dwarf.Entry      `json:"-"`
	variables []*VariableSymbol `json:"-"`
	cu        *CompileUnit      `json:"-"`
	Visited   bool              `json:"-"`
	DisAsm    []DisAsm          `json:"-"` // smaller to store just asm
}

type FunctionCall struct {
	// *uint64 for empty json export, TODO maybe use a string address?
	CallFrom             *uint64 `json:"from,omitempty,omitzero"`
	CallFromFunctionName string  `json:"from_function_name,omitempty"`
	CallTo               *uint64 `json:"to,omitempty,omitzero"`
	CallToFunctionName   string  `json:"to_function_name,omitempty"`
	DynamicCall          bool    `json:"dynamic"`
}

// smaller version of FunctionSymbol to export json for flamegraph rendering
type CallNode struct {
	Name                string     `json:"name"`
	Calls               []CallNode `json:"calls,omitempty,omitzero"`
	MaxStackSizeCallees int64      `json:"max_stack_size_callees,omitempty,omitzero"`
	StackSize           int64      `json:"stack_size,omitempty,omitzero"`
	Address             uint64     `json:"address,omitempty,omitzero"`
	// for backtracking
	Caller    *CallNode `json:"-"`
	Recursion bool      `json:"-"`
}

type CallBranch struct {
	CallList  []CallNode `json:"call_list,omitempty,omitzero"`
	StackSize int64      `json:"stack_size,omitempty,omitzero"`
	//TODO are we interested in flash size? imagine only one call to a fn
	// ie a library and so it is clear how much cutting it out would free...
}

type CallTree struct {
	Tree            CallNode       `json:"tree,omitempty,omitzero"`
	Branches        []CallBranch   `json:"branches,omitempty,omitzero"`
	UnresolvedCalls []FunctionCall `json:"unresolved,omitempty,omitzero"`
	// temp
	CurrentBranch CallBranch `json:"-"`
}

type DisAsm struct {
	Addr        uint64 `json:"addr"`
	Instruction string `json:"instruction"`
	Opstr       string `json:"opstr"`
	InsBytes    []byte `json:"insBytes"`
}

type VariableSymbol struct {
	Name           string       `json:"name"`
	Address        uint64       `json:"address,omitempty,omitzero"`
	FlashSize      uint64       `json:"size"`
	SectionIndex   uint64       `json:"secidx"`
	SourceFilePath string       `json:"file,omitempty,omitzero"`
	SourceFileLine uint64       `json:"line,omitempty,omitzero"`
	VariableType   string       `json:"type,omitempty,omitzero"`
	Data           []byte       `json:"staticInitData,omitempty,omitzero"`
	cu             *CompileUnit `json:"-"`
}

type Typedef struct {
	Name       string       `json:"name,omitempty,omitzero"`
	Type       string       `json:"type"`
	Size       int64        `json:"size"`
	BitSize    int64        `json:"bitsize,omitzero"`
	ByteOffset int64        `json:"byte_offset"`
	BitOffset  int64        `json:"bit_offset,omitzero"`
	IsPointer  bool         `json:"is_pointer,omitempty,omitzero"`
	Members    []Typedef    `json:"members,omitempty,omitzero"`
	cu         *CompileUnit `json:"-"`
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
	Architecture       string           `json:"architecture"`
	Elf                *elf.File        `json:"-"`
	FirmwareHash       string           `json:"firmware_hash"`
	CompileUnits       []CompileUnit    `json:"compile_units"`
	Sections           []ElfSection     `json:"sections"`
	Functions          []FunctionSymbol `json:"functions"`
	Variables          []VariableSymbol `json:"variables"`
	Types              []Typedef        `json:"types,omitempty,omitzero"`
	Info               []string         `json:"info"`
	// lookup
	Addr2FnMap  map[uint64]*FunctionSymbol `json:"-"`
	Name2FnMap  map[string]*FunctionSymbol `json:"-"` // todo no more lookup by just symbol name...
	SectionsMap map[uint64]*ElfSection     `json:"-"`
}

type FunctionCallEntry struct {
	From string   `json:"from"`
	To   []string `json:"to"`
}

type FunctionCallList []FunctionCallEntry
