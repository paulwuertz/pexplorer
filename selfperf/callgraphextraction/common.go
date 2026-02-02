package callgraphextraction

import (
	"fmt"

	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

// ref:
// * https://developer.arm.com/documentation/dui0489/i/arm-and-thumb-instructions/b
// * https://developer.arm.com/documentation/dui0489/i/arm-and-thumb-instructions/bl
// * https://developer.arm.com/documentation/dui0489/i/arm-and-thumb-instructions/blx
// * https://developer.arm.com/documentation/dui0489/i/arm-and-thumb-instructions/condition-codes
func IsFnCallInstr(instr string) bool {
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
	for i := 0; i < len(s.Functions); i++ {
		f := &s.Functions[i]
		if len(f.Asm) == 0 {
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
