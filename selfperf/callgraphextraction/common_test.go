package callgraphextraction

import "testing"

func Test_isFnCallInstr(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		instr string
		want  bool
	}{
		{"instr test bl", "bl", true},
		{"instr test bleq", "bleq", true},
		{"instr test blne", "blne", true},
		{"instr test blcs", "blcs", true},
		{"instr test blhs", "blhs", true},
		{"instr test blcc", "blcc", true},
		{"instr test bllo", "bllo", true},
		{"instr test blmi", "blmi", true},
		{"instr test blpl", "blpl", true},
		{"instr test blvs", "blvs", true},
		{"instr test blvc", "blvc", true},
		{"instr test blhi", "blhi", true},
		{"instr test blls", "blls", true},
		{"instr test blge", "blge", true},
		{"instr test bllt", "bllt", true},
		{"instr test blgt", "blgt", true},
		{"instr test blle", "blle", true},
		{"instr test blal", "blal", true},
		{"instr test b", "b", false},
		{"instr test beq", "beq", false},
		{"instr test bne", "bne", false},
		{"instr test bcs", "bcs", false},
		{"instr test bhs", "bhs", false},
		{"instr test bcc", "bcc", false},
		{"instr test blo", "blo", false},
		{"instr test bmi", "bmi", false},
		{"instr test bpl", "bpl", false},
		{"instr test bvs", "bvs", false},
		{"instr test bvc", "bvc", false},
		{"instr test bhi", "bhi", false},
		{"instr test bls", "bls", false},
		{"instr test bge", "bge", false},
		{"instr test blt", "blt", false},
		{"instr test bgt", "bgt", false},
		{"instr test ble", "ble", false},
		{"instr test bal", "bal", false},
		{"instr test beq", "beq", false},
		{"instr test bne", "bne", false},
		{"instr test bcs", "bcs", false},
		{"instr test bhs", "bhs", false},
		{"instr test bcc", "bcc", false},
		{"instr test blo", "blo", false},
		{"instr test bmi", "bmi", false},
		{"instr test bpl", "bpl", false},
		{"instr test bvs", "bvs", false},
		{"instr test bvc", "bvc", false},
		{"instr test bhi", "bhi", false},
		{"instr test bls", "bls", false},
		{"instr test bge", "bge", false},
		{"instr test blt", "blt", false},
		{"instr test bgt", "bgt", false},
		{"instr test ble", "ble", false},
		{"instr test bal", "bal", false},
		{"instr test add", "add", false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := IsFnCallInstr(tt.instr)
			// TODO: update the condition below to compare got with tt.want.
			if tt.want != got {
				t.Errorf("isFnCallInstr() = %v, want %v", got, tt.want)
			}
		})
	}
}
