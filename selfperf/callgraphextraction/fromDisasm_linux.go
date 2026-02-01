//go:build linux
// +build linux

package callgraphextraction

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

func AddCallGraph(s *symbolextraction.SElfReport) {
	for i := 0; i < len(s.Functions); i++ {
		f := &s.Functions[i]
		if len(f.DisAsm) == 0 {
			msg := fmt.Sprintf("no call data for %s at %d function with no disasm data", f.Name, f.Address)
			s.Info = append(s.Info, msg)
			continue
		}

		for _, insn := range f.DisAsm {
			if IsFnCallInstr(insn.Instruction) {
				//stackoverflow.com/questions/75285743/arm-gcc-cortex-m4-calling-address-as-function-generates-blx-instead-of-bl
				var calladdr uint64
				n, err := fmt.Sscanf(insn.Opstr, "#0x%X", &calladdr)
				if err == nil && n == 1 {
					f.Callees = append(f.Callees, symbolextraction.FunctionCall{f.Address, uint64(calladdr), false})
				} else {
					f.Callees = append(f.Callees, symbolextraction.FunctionCall{CallFrom: f.Address, DynamicCall: true})
					continue
				}
				fnCalled, fnFound := s.Addr2FnMap[uint64(calladdr)]
				if fnFound {
					fnCalled.Callers = append(fnCalled.Callers, symbolextraction.FunctionCall{f.Address, uint64(calladdr), false})
				} else {
					msg := fmt.Sprintf("static call from %s at %d to unknown function", f.Name, f.Address)
					s.Info = append(s.Info, msg)
				}

			}
		}
	}
}

func EnhanceByDisasm(s *symbolextraction.SElfReport) {
	// mv to fromelf :)
	for i := 0; i < len(s.Functions); i++ {
		f := &s.Functions[i]
		s.Addr2FnMap[f.Address] = f
	}
	// get calls from disasm
	AddDisAsmFromAsm(s)
	AddCallGraph(s)
}
