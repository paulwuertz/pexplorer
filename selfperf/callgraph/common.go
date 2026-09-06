package callgraph

import (
	"encoding/binary"
	"fmt"
	"log"
	"sort"
	"strconv"
	"strings"

	"github.com/go-delve/delve/pkg/dwarf/frame"
	"github.com/go-delve/delve/pkg/dwarf/godwarf"
	"github.com/paulwuertz/pexplorer/selfperf/config"
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
	if l < 2 || l > 6 {
		return false
	}
	// catch non bl(x)'s, but b + condition le is a if or loop jump
	if l == 3 && instr != "blx" || string([]byte(instr)[:2]) != "bl" {
		return false
	}
	return true
}

func IsDynamicCallInstr(instr string) bool {
	if instr == "blx" {
		return true
	}
	return false
}

// TODO - are there more like these?
func IsForwardedCall(instr symbolextraction.DisAsm, f *symbolextraction.FunctionSymbol, s *symbolextraction.SElfReport) (bool, uint64) {
	var calladdr uint64 = 0
	noKnownForwardedInstruction := instr.Instruction != "b.w"
	foundNumber, err := fmt.Sscanf(instr.Opstr, "#0x%X", &calladdr)
	if f.Name == "log_backend_enable" {
		f.StackQualifiers = "estimated+experimental"
	}
	if err != nil || foundNumber != 1 || noKnownForwardedInstruction {
		return false, 0
	}
	startAddr := f.Address
	endAddr := f.Address + f.FlashSize
	outsiteCallrange := (calladdr < startAddr) || (calladdr > endAddr)
	if !outsiteCallrange {
		return false, 0
	}
	_, ok := s.Addr2FnMap[calladdr]
	if !ok {
		return false, 0
	} else {
		return true, calladdr
	}
}

func GetFunctionCallList(s *symbolextraction.SElfReport) symbolextraction.FunctionCallList {
	var l symbolextraction.FunctionCallList = make(symbolextraction.FunctionCallList, 0)
	for _, f := range s.Functions {
		var callsTo []string = make([]string, 0)
		for _, to := range f.Callees {
			// removing add constprops,isra,... and other qualifiers
			fn_name, _, _ := strings.Cut(to.CallToFunctionName, ".")
			callsTo = append(callsTo, fn_name)
		}
		fn_name, _, _ := strings.Cut(f.Name, ".")
		fun_calls := symbolextraction.FunctionCallEntry{
			From: fn_name,
			To:   callsTo,
		}
		l = append(l, fun_calls)
	}
	return l
}

func DynamicCallResolutionToMap(dynamicCalls []config.DynamicCallResolution) (callMap config.DynamicCallResolutionMap) {
	callMap = make(config.DynamicCallResolutionMap, len(dynamicCalls))
	for _, d := range dynamicCalls {
		callerName := d.Caller
		_, isNew := callMap[callerName]
		if !isNew {
			fmt.Println("Warning: ", callerName, "configured as a dynamic callback is not unique -> TODO fix symbol ident not by name")
		}
		// assumem symbol by name is unique for now, even if we know better...
		callMap[callerName] = d
	}
	return callMap
}

func AddCallGraph(s *symbolextraction.SElfReport, dynamicCalls []config.DynamicCallResolution) {
	callMap := DynamicCallResolutionToMap(dynamicCalls)
	for i := 0; i < len(s.Functions); i++ {
		f := &s.Functions[i]

		if len(f.Asm) == 0 {
			msg := fmt.Sprintf("no static call data for %s at %d function with no disasm data", f.Name, f.Address)
			s.Info = append(s.Info, msg)
			continue
		}

		for _, insn := range f.DisAsm {
			isForwardedCall, forwardedAddr := IsForwardedCall(insn, f, s)
			// if f.Name == "z_log_msg_post_finalize" {
			// 	fmt.Println("p-p")
			// }
			if IsFnCallInstr(insn.Instruction) {
				//stackoverflow.com/questions/75285743/arm-gcc-cortex-m4-calling-address-as-function-generates-blx-instead-of-bl
				var calladdr []uint64 = make([]uint64, 1) // wasteful hack to get a nullable int... TODO any better way?
				foundNumber, err := fmt.Sscanf(insn.Opstr, "#0x%X", &calladdr[0])
				if err != nil || foundNumber != 1 {
					// branch instruction with ill formed address assume dynamic call
					f.Callees = append(f.Callees, symbolextraction.FunctionCall{
						CallFromFunctionName: f.Name,
						CallFrom:             &f.Address,
						DynamicCall:          true,
					})
					continue
				}
				call := symbolextraction.FunctionCall{
					CallFromFunctionName: f.Name,
					CallFrom:             &f.Address,
					CallTo:               &calladdr[0],
					DynamicCall:          false,
				}
				fnCalled, fnFound := s.Addr2FnMap[calladdr[0]]
				if fnFound {
					call.CallToFunctionName = fnCalled.Name
					fnCalled.Callers = append(fnCalled.Callers, call)
				} else {
					msg := fmt.Sprintf("static call from %s at %d to unknown function", f.Name, f.Address)
					s.Info = append(s.Info, msg)
				}
				f.Callees = append(f.Callees, call)
			} else if isForwardedCall {
				call := symbolextraction.FunctionCall{
					CallFromFunctionName: f.Name,
					CallFrom:             &f.Address,
					CallTo:               &forwardedAddr,
					DynamicCall:          false,
				}
				fnCalled, fnFound := s.Addr2FnMap[forwardedAddr]
				if fnFound {
					call.CallToFunctionName = fnCalled.Name
					fnCalled.Callers = append(fnCalled.Callers, call)
				} else {
					msg := fmt.Sprintf("static call from %s at %d to unknown function", f.Name, f.Address)
					s.Info = append(s.Info, msg)
				}
				f.Callees = append(f.Callees, call)
			}
		}
		// resolve dynamic calls from conffile
		dynamicCalls, hasResolvedCalls := callMap[f.Name]
		if hasResolvedCalls {
			for _, calleeName := range dynamicCalls.Callees {
				callee, found := s.Name2FnMap[calleeName]
				if !found {
					log.Fatal("Callee of '", f.Name, "' to '", calleeName, "' from conf file not found")
				}
				filled_unresolved := false
				// fill all unresolved calls first
				// some dynamic calls might call into more then one function...
				for i := 0; i < len(f.Callees); i++ {
					if f.Callees[i].DynamicCall && f.Callees[i].CallTo == nil {
						f.Callees[i].CallToFunctionName = calleeName
						f.Callees[i].CallTo = &callee.Address
						filled_unresolved = true
						break
					}
				}
				// ... in that case append further
				if !filled_unresolved {
					call := symbolextraction.FunctionCall{
						CallFromFunctionName: f.Name,
						CallFrom:             &f.Address,
						CallTo:               &callee.Address,
						CallToFunctionName:   callee.Name,
						DynamicCall:          true,
					}
					f.Callees = append(f.Callees, call)
				}
				fmt.Println(dynamicCalls.Caller, "p-p", calleeName)
			}
		}
	}
}

// test only for ARM for a working web demo
func ExtractFunctionStackUsage(f *symbolextraction.FunctionSymbol) {
	// fmt.Println("\t\tfn", f.Name, mainfde.Length, f.SourceFilePath, f.SourceFileLine)
	var current_stacksize int64 = 0
	if f.Name == "gs_usb_rx_thread" {
		f.StackQualifiers = "experimental-estimate"
	}
	for _, d := range f.DisAsm {
		sub := strings.HasPrefix(d.Instruction, "sub")
		stack_pointer := strings.HasPrefix(d.Opstr, "sp, #0x")
		if strings.HasPrefix(d.Instruction, "push") {
			//push   {r0, r1, r2, r3, r4, lr} OR
			//push.w {r0, r1, r2, r3, r4, lr} with a suffix condition - maybe TODO distinguish cond.?
			nr_regs := strings.Count(d.Opstr, ",") + 1
			current_stacksize += int64(nr_regs) * 4
			// fmt.Println(f.Name, "push now", current_stacksize, "@", d.Addr)
		} else if sub && stack_pointer {
			// sub   sp, #0x10
			stackSubSize, err := strconv.ParseInt(d.Opstr[5:], 0, 64)
			if err != nil {
				log.Fatalf(f.Name, "sub now hexstr err", d.Opstr)
			}
			current_stacksize += stackSubSize
			// fmt.Println(f.Name, "sub now", current_stacksize, "@", d.Addr)
		}
		// fmt.Println(fmt.Sprintf("%x", i), "off:", s.CFA.Offset, d.Instruction, d.Opstr, "-> cfa reg:", s.CFA.Reg, "rule:", symbolextraction.RegRuleEnum2String[s.CFA.Rule], "expr:", s.CFA.Expression, "regs:", s.Regs, "reta:", s.RetAddrReg)
	}
	f.StackSize = current_stacksize
	f.StackQualifiers = "estimated+experimental"
	// fmt.Println(f.Name, "result:", current_stacksize)
}

// maybe it could be useful in the future... was my first intend, but was off to much...
func GetFunctionStackUsage_DebugFrameUnwinding(f *symbolextraction.FunctionSymbol, frames frame.FrameDescriptionEntries) (uint64, error) {
	mainfde, err := frames.FDEForPC(f.Address)
	// fmt.Println("\t\tfn", err)
	if err != nil {
		return 0, err
	}
	// fmt.Println("\t\tfn", f.Name, mainfde.Length, f.SourceFilePath, f.SourceFileLine)
	var max int64 = 0
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
		if f.Name == "led_state_identify_run" && d.Addr > 0x601 {
			f.StackQualifiers = "estimated+experimental"
		}
		var current_stacksize int64 = 0
		switch s.CFA.Rule {
		case frame.RuleOffset:
			current_stacksize = s.CFA.Offset * mainfde.CIE.DataAlignmentFactor
		// case frame.RuleCFA:
		case frame.RuleCFA:
			offset := s.CFA.Offset * mainfde.CIE.DataAlignmentFactor
			reg := s.CFA.Reg
			current_stacksize = offset + int64(reg)
		}

		if current_stacksize >= max {
			max = current_stacksize
		}
		fmt.Println(fmt.Sprintf("%x", i), "off:", s.CFA.Offset, d.Instruction, d.Opstr, "-> cfa reg:", s.CFA.Reg, "rule:", symbolextraction.RegRuleEnum2String[s.CFA.Rule], "expr:", s.CFA.Expression, "regs:", s.Regs, "reta:", s.RetAddrReg)
	}
	f.StackSize = max
	f.StackQualifiers = "estimated+experimental"
	// for i := mainfde.Begin(); i < mainfde.End(); i = i + 2 {
	return 0, err
}

func GetStackUseDetails_DebugFrameUnwinding(s *symbolextraction.SElfReport) {
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
		GetFunctionStackUsage_DebugFrameUnwinding(function, fe)
	}
}

func GetStackUseDetails(s *symbolextraction.SElfReport) {
	for i := 0; i < len(s.Functions); i++ {
		function := &s.Functions[i]
		ExtractFunctionStackUsage(function)
	}
}

func TraverseCallSubGraph(s *symbolextraction.SElfReport, f *symbolextraction.FunctionSymbol, subgraphIndex uint, calldepth uint) int64 {
	var biggestSubStackSize int64 = 0
	if len(f.Callees) == 0 {
		// fmt.Println(strings.Repeat("\t", int(calldepth)), f.Name, " endtree stacksize:", f.StackSize)
		f.MaxStackSizeCallees = f.StackSize
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
