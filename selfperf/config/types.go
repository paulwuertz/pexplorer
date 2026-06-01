package config

type RTOSThread struct {
	ThreadEntryName   string `json:"thread_entry_name"`
	StackVariableName string `json:"stack_variable_name"`
	Size              uint64 `json:"size"`
}

type DynamicCallResolution struct {
	CallFrom string `json:"call_from"`
	CallTo   string `json:"call_to"`
}

type PexplorerConfig struct {
	Threads        []RTOSThread            `json:"threads"`
	DynamicThreads []DynamicCallResolution `json:"dynamic_threads"`
}
