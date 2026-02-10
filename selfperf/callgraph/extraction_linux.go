//go:build linux
// +build linux

package callgraph

import (
	"fmt"
	"log"

	"github.com/knightsc/gapstone"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func AddDisAsmFromAsm(s *symbolextraction.SElfReport) {
	g, err := gapstone.New(gapstone.CS_ARCH_ARM, gapstone.CS_MODE_THUMB)
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
				Addr:        uint64(insn.Address),
				Instruction: insn.Mnemonic,
				Opstr:       insn.OpStr,
				InsBytes:    insn.Bytes,
			}
		}
	}
	g.Close()
}

func EnhanceByDisasm(s *symbolextraction.SElfReport) {
	// get calls from disasm
	AddDisAsmFromAsm(s)
	AddCallGraph(s)
	GetStackUseDetails(s)
}
