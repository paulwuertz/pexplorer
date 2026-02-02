//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"debug/elf"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/paulwuertz/pexplorer/selfperf/callgraphextraction"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

var elfReport symbolextraction.SElfReport

func wasm_get_elf_report() js.Func {

	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		elf_bin_js := args[0]
		elf_size := elf_bin_js.Length()
		elf_binary := make([]byte, elf_size)

		fmt.Printf("args # %d\n", len(args))
		fmt.Printf("input %d\n", elf_size)
		bytes_copied := js.CopyBytesToGo(elf_binary, elf_bin_js)
		fmt.Printf("n %d bytes\n", bytes_copied)

		r := bytes.NewReader(elf_binary)
		elfFile, _ := elf.NewFile(r)

		elfReport = symbolextraction.GetFWReport(elfFile)
		elfReport.SingleFirmware = true
		elfReport.FirmwareIdentifier = "unspecified"
		elfReport.Timestamp = "just now"

		datajson, _ := json.Marshal(elfReport)
		return string(datajson)
	})
	return jsonFunc
}

func wasm_get_fn_calls_from_disasm() js.Func {

	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		disasm_js := args[0]
		disasm_size := disasm_js.Length()
		disasm_binary := make([]byte, disasm_size)
		bytes_copied := js.CopyBytesToGo(disasm_binary, disasm_js)

		var addr2Disasm map[uint64][]symbolextraction.DisAsm
		json.Unmarshal(disasm_binary, &addr2Disasm)

		fmt.Printf("args # %d\n", len(args))
		fmt.Printf("n %d bytes\n", bytes_copied)

		for addr, disasm := range addr2Disasm {
			f, ok := elfReport.Addr2FnMap[addr]
			if !ok {
				fmt.Println("cannot lookup fn with disasm again from addr: ", addr)
			}
			f.DisAsm = disasm
		}
		callgraphextraction.AddCallGraph(&elfReport)

		datajson, _ := json.Marshal(elfReport)
		datajson2, _ := json.Marshal(addr2Disasm)
		fmt.Printf("afta %s bytes\n", string(datajson2))

		return string(datajson)
	})
	return jsonFunc
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("get_elf_report", wasm_get_elf_report())
	js.Global().Set("add_fn_calls_from_disasm", wasm_get_fn_calls_from_disasm())
	<-make(chan struct{})
}
