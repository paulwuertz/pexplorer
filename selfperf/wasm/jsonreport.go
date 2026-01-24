//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"debug/elf"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

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

		elfReport := symbolextraction.GetFWReport(elfFile)
		datajson, _ := json.Marshal(elfReport)
		return string(datajson)
	})
	return jsonFunc
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("get_elf_report", wasm_get_elf_report())
	<-make(chan struct{})
}
