package symbolextraction

import (
	"debug/elf"
	"encoding/json"
	"fmt"
	"log"
	"testing"
)

func Test_extractFunctions(t *testing.T) {

	elfFile1, err := elf.Open("/home/paul/git/ztest/build_hello_world_frdm_k64f_42/zephyr/zephyr.elf")
	elfFile2, err := elf.Open("/home/paul/git/Prusa-Firmware-Buddy/build/mini_release_noboot/firmware")

	if err != nil {
		log.Fatal(err)
	}
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		elfFile elf.File
		want    []FunctionSymbol
	}{
		{"", *elfFile1, []FunctionSymbol{}},
		{"", *elfFile2, []FunctionSymbol{}},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := extractFunctions(tt.elfFile)
			b, err := json.Marshal(got)
			fmt.Println(string(b), err)
			// TODO: update the condition below to compare got with tt.want.
			if true {
				t.Errorf("extractFunctions() = %v, want %v", got, tt.want)
			}
		})
	}
}
