package callgraph

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
		{"instr test bleq.w", "bleq.w", true},
		{"instr test blne.w", "blne.w", true},
		{"instr test blcs.w", "blcs.w", true},
		{"instr test blhs.w", "blhs.w", true},
		{"instr test blcc.w", "blcc.w", true},
		{"instr test bllo.w", "bllo.w", true},
		{"instr test blmi.w", "blmi.w", true},
		{"instr test blpl.w", "blpl.w", true},
		{"instr test blvs.w", "blvs.w", true},
		{"instr test blvc.w", "blvc.w", true},
		{"instr test blhi.w", "blhi.w", true},
		{"instr test blls.w", "blls.w", true},
		{"instr test blge.w", "blge.w", true},
		{"instr test bllt.w", "bllt.w", true},
		{"instr test blgt.w", "blgt.w", true},
		{"instr test blle.w", "blle.w", true},
		{"instr test blal.w", "blal.w", true},
		{"instr test bleq.n", "bleq.n", true},
		{"instr test blne.n", "blne.n", true},
		{"instr test blcs.n", "blcs.n", true},
		{"instr test blhs.n", "blhs.n", true},
		{"instr test blcc.n", "blcc.n", true},
		{"instr test bllo.n", "bllo.n", true},
		{"instr test blmi.n", "blmi.n", true},
		{"instr test blpl.n", "blpl.n", true},
		{"instr test blvs.n", "blvs.n", true},
		{"instr test blvc.n", "blvc.n", true},
		{"instr test blhi.n", "blhi.n", true},
		{"instr test blls.n", "blls.n", true},
		{"instr test blge.n", "blge.n", true},
		{"instr test bllt.n", "bllt.n", true},
		{"instr test blgt.n", "blgt.n", true},
		{"instr test blle.n", "blle.n", true},
		{"instr test blal.n", "blal.n", true},
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
