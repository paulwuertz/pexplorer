//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	"crypto/sha256"
	"debug/elf"
	"encoding/json"
	"fmt"
	"syscall/js"

	"github.com/paulwuertz/pexplorer/selfperf/callgraph"
	"github.com/paulwuertz/pexplorer/selfperf/config"
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

		fw_hash := sha256.Sum256([]byte(elf_binary))
		fw_hash_str := fmt.Sprintf("%x", fw_hash)
		fmt.Printf("args # %d\n", len(args))
		fmt.Printf("input %d\n", elf_size)
		bytes_copied := js.CopyBytesToGo(elf_binary, elf_bin_js)
		fmt.Printf("n %d bytes\n", bytes_copied)

		r := bytes.NewReader(elf_binary)
		elfFile, _ := elf.NewFile(r)

		elfReport = symbolextraction.GetFWReport(elfFile, fw_hash_str)
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
		if len(args) != 2 {
			return "Invalid no of arguments passed"
		}
		disasm_js := args[0]
		config_js := args[1]
		disasm_size := disasm_js.Length()
		disasm_binary := make([]byte, disasm_size)
		bytes_copied := js.CopyBytesToGo(disasm_binary, disasm_js)
		config_size := config_js.Length()
		config_binary := make([]byte, config_size)
		js.CopyBytesToGo(config_binary, config_js)

		var addr2Disasm map[uint64][]symbolextraction.DisAsm
		json.Unmarshal(disasm_binary, &addr2Disasm)
		var pexplorer_config config.PexplorerConfig
		json.Unmarshal(config_binary, &pexplorer_config)

		fmt.Printf("args # %d\n", len(args))
		fmt.Printf("n %d bytes\n", bytes_copied)
		fmt.Printf("config %d bytes = '%s'\n", config_size, config_binary)

		for addr, disasm := range addr2Disasm {
			f, ok := elfReport.Addr2FnMap[addr]
			if !ok {
				fmt.Println("cannot lookup fn with disasm again from addr: ", addr)
			}
			f.DisAsm = disasm
		}
		callgraph.GetStackUseDetails(&elfReport)
		// TODO pass report from web localstorage
		callgraph.AddCallGraph(&elfReport, pexplorer_config.DynamicCalls)
		callgraph.TraverseCallGraph(&elfReport)

		datajson, _ := json.Marshal(elfReport)
		datajson2, _ := json.Marshal(addr2Disasm)
		fmt.Printf("afta %s bytes\n", string(datajson2))

		return string(datajson)
	})
	return jsonFunc
}

func wasm_get_fn_calltree() js.Func {

	jsonFunc := js.FuncOf(func(this js.Value, args []js.Value) any {
		if len(args) != 1 {
			return "Invalid no of arguments passed"
		}
		fn_addr_js := args[0]
		fn_addr := fn_addr_js.Int()

		fn := elfReport.Addr2FnMap[uint64(fn_addr)]
		tree := fn.GetCallTreeJson(&elfReport, 0)

		fmt.Printf("args # %d\n", len(args))
		fmt.Printf("fn addr %d\n", fn_addr)

		datajson, _ := json.Marshal(tree)
		fmt.Printf("afta %s bytes\n", string(datajson))

		return string(datajson)
	})
	return jsonFunc
}

func main() {
	fmt.Println("Go Web Assembly")
	js.Global().Set("get_elf_report", wasm_get_elf_report())
	js.Global().Set("add_fn_calls_from_disasm", wasm_get_fn_calls_from_disasm())
	js.Global().Set("get_fn_calltree", wasm_get_fn_calltree())
	<-make(chan struct{})
}
