// package config
package main

import (
	"encoding/json"
	"fmt"
)

type RTOSThread struct {
	ThreadEntryName   string `json:"thread_entry_name"`
	StackVariableName string `json:"stack_variable_name,omitempty"`
	Size              uint64 `json:"size,omitempty,omitzero"`
}

type DynamicCallResolution struct {
	CallFrom string `json:"call_from"`
	CallTo   string `json:"call_to"`
}

type PexplorerConfig struct {
	Threads      []RTOSThread            `json:"threads,omitempty"`
	DynamicCalls []DynamicCallResolution `json:"dynamic_calls,omitempty"`
}

func main() {
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
		CallFrom: "work_queue_main",
		CallTo:   "led_event_triggered_work_handler",
	})
	p.DynamicCalls = append(p.DynamicCalls, DynamicCallResolution{
		CallFrom: "work_queue_main",
		CallTo:   "cannectivity_usb_reboot",
	})
	p.DynamicCalls = append(p.DynamicCalls, DynamicCallResolution{
		CallFrom: "work_queue_main",
		CallTo:   "dfu_button_poll",
	})
	datajson, _ := json.Marshal(p)
	// datajson, _ := json.MarshalIndent(p, "", "    ")
	fmt.Println(string(datajson))
}
