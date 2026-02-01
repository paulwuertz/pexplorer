package callgraphextraction

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
