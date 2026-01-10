package symbolextraction

type ElfSection struct {
	Name    string `json:"name"`
	Address uint64 `json:"address"`
	Size    uint64 `json:"size"`
	Index   uint8  `json:"index"`
}

type FunctionSymbol struct {
	Name              string           `json:"name"`
	Address           uint64           `json:"address,omitempty,omitzero"`
	FlashSize         uint64           `json:"size"`
	FunctionStackSize uint64           `json:"stack_size,omitempty,omitzero"`
	SectionIndex      uint8            `json:"secidx"`
	SourceFilePath    string           `json:"file,omitempty,omitzero"`
	SourceFileLine    uint64           `json:"line,omitempty,omitzero"`
	Variables         []VariableSymbol `json:"vars,omitempty,omitzero"`
	Asm               []byte
	// calls
}

type VariableSymbol struct {
	Name           string `json:"name"`
	Address        uint64 `json:"address,omitempty,omitzero"`
	FlashSize      uint64 `json:"size"`
	SectionIndex   uint8  `json:"secidx"`
	SourceFilePath string `json:"file,omitempty,omitzero"`
	SourceFileLine uint64 `json:"line,omitempty,omitzero"`
	VariableType   string `json:"type,omitempty,omitzero"`
	Data           []byte `json:"byte,omitempty,omitzero"`
}

// CompileUnit represents a compilation unit,
// including a series of source files and function definitions
type CompileUnit struct {
	Source    []string
	functions []FunctionSymbol
	variables []VariableSymbol
}

type SElfReport struct {
	srcFiles  []CompileUnit
	sections  []ElfSection
	functions []FunctionSymbol
	variables []VariableSymbol
}
