//go:build js && wasm
// +build js,wasm

package main

import (
	"bytes"
	// "debug/dwarf"
	"debug/elf"
	"fmt"
	"syscall/js"
)

func extractSections(elfFile elf.File) any {
	var sections = make([]any, len(elfFile.Sections))
	for _, sec := range elfFile.Sections {
		sh := sec.SectionHeader
		sections = append(sections, sh.Name)
		sections = append(sections, sh.Size)
		sections = append(sections, sh.Addr)
		// append(sections, ElfSections{
		// 	sh.Name, sh.Size, sh.Addr,
		// })
	}
	return sections
}

func wasm_read_sections_as_json() js.Func {
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
		elf, _ := elf.NewFile(r)
		return extractSections(*elf)
	})
	return jsonFunc
}

// func main() {
// 	// fmt.Println("Go Web Assembly")
// 	// js.Global().Set("read_sections_as_json", wasm_read_sections_as_json())
// 	// <-make(chan struct{})
// }
