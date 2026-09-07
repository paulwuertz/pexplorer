//go:build linux
// +build linux

package callgraph

import (
	"fmt"
	"log"

	"github.com/bpfsnoop/gapstone"
	"github.com/paulwuertz/pexplorer/selfperf/config"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func AddDisAsmFromAsm(s *symbolextraction.SElfReport) {
	g, err := gapstone.New(gapstone.CS_ARCH_ARM, gapstone.CS_MODE_THUMB+gapstone.CS_MODE_MCLASS)
	if err != nil {
		log.Fatalf("Failed to initialize engine: %v", err)
	}
	for i := 0; i < len(s.Functions); i++ {
		f := &s.Functions[i]
		if len(f.Asm) == 0 {
			msg := fmt.Sprintf("no call data for %s at %d function with no asm data", f.Name, f.Address)
			s.Info = append(s.Info, msg)
			continue
		}

		if f.Name == "net_buf_unref" {
			fmt.Println("p-p")
		}
		insns, err := g.Disasm(
			f.Asm,     // code buffer
			f.Address, // starting address
			0,         // insns to disassemble, 0 for all
		)

		// TODO gapstone.ErrOK (=> 0) means failed disasm...
		if err != nil && err != gapstone.ErrOK {
			fmt.Println("Disassembly error: ", err)
			continue
		}
		f.DisAsm = make([]symbolextraction.DisAsm, len(insns))
		for i, insn := range insns {
			f.DisAsm[i] = symbolextraction.DisAsm{
				Addr:        uint64(insn.InstructionHeader.Address),
				Instruction: insn.InstructionHeader.Mnemonic,
				Opstr:       insn.InstructionHeader.OpStr,
				InsBytes:    insn.InstructionHeader.Bytes,
			}
		}
	}
	g.Close()
}

func EnhanceByDisasm(s *symbolextraction.SElfReport, dynamicCalls []config.DynamicCallResolution) {
	// get calls from disasm
	AddDisAsmFromAsm(s)
	AddCallGraph(s, dynamicCalls)
	GetStackUseDetails(s)
}
