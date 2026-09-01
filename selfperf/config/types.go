// package config
package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

type RTOSThread struct {
	ThreadEntryName   string `json:"thread_entry_name"`
	StackVariableName string `json:"stack_variable_name,omitempty"`
	Size              uint64 `json:"size,omitempty,omitzero"`
	Used              uint64 `json:"-"`
	NrUnresolvedCalls uint64 `json:"-"`
}

type DynamicCallResolution struct {
	Caller  string   `json:"caller"`
	Callees []string `json:"callees"`
}

type PexplorerConfig struct {
	Threads      []RTOSThread            `json:"threads,omitempty"`
	DynamicCalls []DynamicCallResolution `json:"dynamic_calls,omitempty"`
}

func test_export() {
	var p PexplorerConfig
	// 	logging_stack - size: 768 bytes - addr: 0x20002400
	// z_interrupt_stacks - size: 2048 bytes - addr: 0x20003f00
	// z_idle_stacks - size: 320 bytes - addr: 0x20004700
	// z_main_stack - size: 1024 bytes - addr: 0x20004840
	// sys_work_q_stack - size: 1024 bytes - addr: 0x20004c40
	p.Threads = make([]RTOSThread, 0)
	p.Threads = append(p.Threads, RTOSThread{
		ThreadEntryName:   "log_process_thread_func",
		StackVariableName: "logging_stack",
	})
	p.Threads = append(p.Threads, RTOSThread{
		ThreadEntryName:   "bg_thread_main",
		StackVariableName: "z_main_stack",
	})
	p.Threads = append(p.Threads, RTOSThread{
		ThreadEntryName: "gs_usb_rx_thread",
		Size:            1024,
	})
	p.Threads = append(p.Threads, RTOSThread{
		ThreadEntryName: "gs_usb_tx_thread",
		Size:            1024,
	})
	p.Threads = append(p.Threads, RTOSThread{
		ThreadEntryName:   "work_queue_main",
		StackVariableName: "sys_work_q_stack",
	})

	p.DynamicCalls = make([]DynamicCallResolution, 0)
	p.DynamicCalls = append(p.DynamicCalls, DynamicCallResolution{
		Caller:  "work_queue_main",
		Callees: []string{"led_event_triggered_work_handler"},
	})
	p.DynamicCalls = append(p.DynamicCalls, DynamicCallResolution{
		Caller:  "work_queue_main",
		Callees: []string{"cannectivity_usb_reboot"},
	})
	p.DynamicCalls = append(p.DynamicCalls, DynamicCallResolution{
		Caller:  "work_queue_main",
		Callees: []string{"dfu_button_poll"},
	})
	datajson, _ := json.Marshal(p)
	// datajson, _ := json.MarshalIndent(p, "", "    ")
	fmt.Println(string(datajson))
}

func Import_config_from_file(filename string) (p PexplorerConfig, err error) {
	// Open our jsonFile
	jsonFile, err := os.Open(filename)
	if err != nil {
		fmt.Println(err)
		return p, err
	}
	defer jsonFile.Close()

	byteValue, err := io.ReadAll(jsonFile)
	if err != nil {
		fmt.Println("Error reading file:", err)
		return p, err
	}

	err = json.Unmarshal([]byte(byteValue), &p)
	if err != nil {
		fmt.Println("Error unmarshalling JSON:", err)
		return p, err
	}
	return p, nil
}

func test_import() {
	p, _ := Import_config_from_file("/home/paul/git/pexplorer/selfperf/tmp/pexplorer-CANnectivity - NXP LPC55S16v1.4 - GCC(8).json")
	fmt.Println(p)
}

// func main() {
// 	fmt.Println("Test export:")
// 	test_export()
// 	fmt.Println("\n\nTest import:\n\n")
// 	test_import()
// }
