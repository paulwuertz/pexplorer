package callgraph

import (
	"encoding/binary"
	"fmt"
	"log"
	"sort"

	"github.com/go-delve/delve/pkg/dwarf/frame"
	"github.com/go-delve/delve/pkg/dwarf/godwarf"
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
				var calladdr []uint64 = make([]uint64, 1) // wasteful hack to get a nullable int... TODO any better way?
				n, err := fmt.Sscanf(insn.Opstr, "#0x%X", &calladdr[0])
				if err == nil && n == 1 {
					f.Callees = append(f.Callees, symbolextraction.FunctionCall{&f.Address, &calladdr[0], false})
				} else {
					f.Callees = append(f.Callees, symbolextraction.FunctionCall{CallFrom: &f.Address, DynamicCall: true})
					continue
				}
				fnCalled, fnFound := s.Addr2FnMap[calladdr[0]]
				if fnFound {
					fnCalled.Callers = append(fnCalled.Callers, symbolextraction.FunctionCall{&f.Address, &calladdr[0], false})
				} else {
					msg := fmt.Sprintf("static call from %s at %d to unknown function", f.Name, f.Address)
					s.Info = append(s.Info, msg)
				}
			}
		}
	}
}

func GetFunctionStackUsage(f *symbolextraction.FunctionSymbol, frames frame.FrameDescriptionEntries) (uint64, error) {
	mainfde, err := frames.FDEForPC(f.Address)
	// fmt.Println("\t\tfn", err)
	if err != nil {
		return 0, err
	}
	// fmt.Println("\t\tfn", f.Name, mainfde.Length, f.SourceFilePath, f.SourceFileLine)
	var max uint64 = 0
	for _, d := range f.DisAsm {
		i := d.Addr
		// for ARM the return addr is saved in r13,
		// when CFA is in R13 and only there, then CFA.off seems to be pretty much the frame size...
		// if not and the CFA reg changes and the SP is pushed around we do not know from this table...
		s, err := mainfde.EstablishFrame(i)
		if err != nil {
			// fmt.Println(err, "skip frame at addr", i, "for fn:", f.Name)
			continue
		}
		if uint64(s.CFA.Offset) >= max {
			max = uint64(s.CFA.Offset)
		}
		// fmt.Println(fmt.Sprintf("%x", i), "off:", s.CFA.Offset, d.Instruction, d.Opstr, "-> cfa reg:", s.CFA.Reg, "rule:", symbolextraction.RegRuleEnum2String[s.CFA.Rule], "expr:", s.CFA.Expression, "regs:", s.Regs, "reta:", s.RetAddrReg)
	}
	f.StackSize = max
	f.StackQualifiers = "estimated+experimental"
	// for i := mainfde.Begin(); i < mainfde.End(); i = i + 2 {
	return 0, err
}

func GetStackUseDetails(s *symbolextraction.SElfReport) {
	framedata, _ := godwarf.GetDebugSectionElf(s.Elf, "frame")
	fe, err := frame.Parse(framedata, binary.LittleEndian, 0, 4, 0)
	if err != nil {
		log.Fatal("could not parse frame data of elffile", err)
	}
	sort.Slice(fe, func(i, j int) bool {
		return fe[i].Begin() < fe[j].Begin()
	})

	for i := 0; i < len(s.Functions); i++ {
		function := &s.Functions[i]
		GetFunctionStackUsage(function, fe)
	}
}

func TraverseCallSubGraph(s *symbolextraction.SElfReport, f *symbolextraction.FunctionSymbol, subgraphIndex uint, calldepth uint) uint64 {
	var biggestSubStackSize uint64 = 0
	if len(f.Callees) == 0 {
		// fmt.Println(strings.Repeat("\t", int(calldepth)), f.Name, " endtree stacksize:", f.StackSize)
		return f.StackSize
	}
	if f.Visited {
		// fmt.Println(strings.Repeat("\t", int(calldepth)), f.Name, "revisited stacksize:", f.StackSize, "biggest calletreesize:", f.MaxStackSizeCallees)
		return f.MaxStackSizeCallees
	}
	f.Visited = true
	for _, callees := range f.Callees {
		callAddr := callees.CallTo
		if callAddr == nil {
			continue
		}
		callee, ok := s.Addr2FnMap[*callAddr]
		if !ok {
			// log.Fatal("fn not found", f)
			continue
		}
		subStackSize := f.StackSize + TraverseCallSubGraph(s, callee, subgraphIndex, calldepth+1)
		if subStackSize > biggestSubStackSize {
			biggestSubStackSize = subStackSize
		}
	}
	// fmt.Println(strings.Repeat("\t", int(calldepth)), f.Name, "stacksize:", f.StackSize, "biggest calletreesize:", biggestSubStackSize)
	f.MaxStackSizeCallees = biggestSubStackSize
	return biggestSubStackSize
}

func TraverseCallGraph(s *symbolextraction.SElfReport) {
	var subgraphIndex uint = 0
	for i := 0; i < len(s.Functions); i++ {
		function := &s.Functions[i]
		isRootOfCalltree := len(function.Callers) == 0
		if !function.Visited && isRootOfCalltree {
			TraverseCallSubGraph(s, function, subgraphIndex, 0)
			subgraphIndex++
		}
	}
}
