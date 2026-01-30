package callgraphextraction

import (
	"fmt"
	"log"

	"github.com/knightsc/gapstone"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

// ref:
// * https://developer.arm.com/documentation/dui0489/i/arm-and-thumb-instructions/b
// * https://developer.arm.com/documentation/dui0489/i/arm-and-thumb-instructions/bl
// * https://developer.arm.com/documentation/dui0489/i/arm-and-thumb-instructions/blx
// * https://developer.arm.com/documentation/dui0489/i/arm-and-thumb-instructions/condition-codes
func isFnCallInstr(instr string) bool {
	l := len(instr)
	// allow bl+blx with optional 2 letter condition -> 2, 3, 4 or 5 letters
	if l < 2 || l > 5 {
		return false
	}
	// catch non bl(x)'s, but b + condition le is a if or loop jump
	if l == 3 && instr != "blx" || string([]byte(instr)[:2]) != "bl" {
		return false
	}
	return true
}

func AddCallGraph(s *symbolextraction.SElfReport) {
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

		if err == nil || err == gapstone.ErrOK {
			for _, insn := range insns {
				if isFnCallInstr(insn.Mnemonic) {
					//stackoverflow.com/questions/75285743/arm-gcc-cortex-m4-calling-address-as-function-generates-blx-instead-of-bl
					var calladdr uint64
					n, err := fmt.Sscanf(insn.OpStr, "#0x%X", &calladdr)
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
			continue
		}
		log.Fatalf("Disassembly error: %v", err)
	}
	g.Close()
}

func EnhanceByDisasm(s *symbolextraction.SElfReport) {
	// mv to fromelf :)
	for i := 0; i < len(s.Functions); i++ {
		f := &s.Functions[i]
		s.Addr2FnMap[f.Address] = f
	}
	// get calls from disasm
	AddCallGraph(s)
}
