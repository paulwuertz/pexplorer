package callgraph

import (
	"encoding/base64"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"testing"

	"github.com/bpfsnoop/gapstone"
	"github.com/paulwuertz/pexplorer/selfperf/vcg"
)

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

func TestParseVCGDir(t *testing.T) {
	tempDir := t.TempDir()

	rootFile := filepath.Join(tempDir, "root.ci")
	if err := os.WriteFile(rootFile, []byte(`graph: {
 title: "root"
 node: { title: "n0" label: "main" }
 edge: { sourcename: "n0" targetname: "n1" label: "calls" }
}`), 0o600); err != nil {
		t.Fatalf("write root .ci file: %v", err)
	}

	nestedDir := filepath.Join(tempDir, "nested")
	if err := os.MkdirAll(nestedDir, 0o755); err != nil {
		t.Fatalf("create nested dir: %v", err)
	}

	nestedFile := filepath.Join(nestedDir, "child.ci")
	if err := os.WriteFile(nestedFile, []byte(`graph: {
 title: "child"
 node: { title: "n1" label: "helper" }
}`), 0o600); err != nil {
		t.Fatalf("write child .ci file: %v", err)
	}

	if err := os.WriteFile(filepath.Join(tempDir, "ignore.txt"), []byte("nope"), 0o600); err != nil {
		t.Fatalf("write ignored file: %v", err)
	}

	graphs, err := ParseVCGDir(tempDir)
	if err != nil {
		t.Fatalf("ParseVCGDir returned error: %v", err)
	}
	if len(graphs) != 2 {
		t.Fatalf("unexpected graph count: got %d want 2", len(graphs))
	}
	if graphs[0].Title == "" && len(graphs) > 0 {
		// nothing to do here; titles are validated below
	}
	if graphs[0].Title == "root" && len(graphs[1].Nodes) == 1 && graphs[1].Title == "child" {
		return
	}
	if graphs[1].Title == "root" && len(graphs[0].Nodes) == 1 && graphs[0].Title == "child" {
		return
	}
	for _, g := range graphs {
		t.Logf("graph title=%q nodes=%d edges=%d", g.Title, len(g.Nodes), len(g.Edges))
	}
	t.Fatal("unexpected parsed graph set")
}

func TestGraphToFunctionCallList(t *testing.T) {
	g := vcg.Graph{
		Title: "callgraph",
		Nodes: []vcg.Node{
			{Title: "n0", Label: "main"},
			{Title: "n1", Label: "helper"},
			{Title: "n2", Label: "cleanup"},
		},
		Edges: []vcg.Edge{
			{Source: "n0", Target: "n1", Label: "calls"},
			{Source: "n0", Target: "n2", Label: "calls"},
			{Source: "n1", Target: "n2", Label: "calls"},
		},
	}

	list := GraphsToFunctionCallList([]vcg.Graph{g})
	if len(list) != 2 {
		t.Fatalf("unexpected entry count: got %d want 2", len(list))
	}

	lookup := make(map[string][]string, len(list))
	for _, entry := range list {
		lookup[entry.From] = entry.To
	}

	mainCalls := lookup["main"]
	if len(mainCalls) != 2 {
		t.Fatalf("unexpected calls from main: %#v", mainCalls)
	}
	if !contains(mainCalls, "helper") || !contains(mainCalls, "cleanup") {
		t.Fatalf("missing expected calls from main: %#v", mainCalls)
	}
	helperCalls := lookup["helper"]
	if len(helperCalls) != 1 || !contains(helperCalls, "cleanup") {
		t.Fatalf("unexpected calls from helper: %#v", helperCalls)
	}
}

func contains(values []string, wanted string) bool {
	for _, value := range values {
		if value == wanted {
			return true
		}
	}
	return false
}

func Test_capstoneUse(t *testing.T) {
	tests := []struct {
		name string // description of this test case
		// Named input parameters for target function.
		instr string
	}{
		{
			"Dissassemble some Arm instructions",
			"ELX/9/f/EkgSSRNK//f4/hJLG2gT8AEPDtFA9hgAAPAz+QRGR/D6/w1LHGABIgtLGmAMSLLwRPsJSxhoA2icaAlJASJP9BZzoEep8Jf9/ufcTgsI8E4LCPROCwhoBwAgZAcAIF0GAAgETwsI",
		},
	}
	engine, err := gapstone.New(gapstone.CS_ARCH_ARM, gapstone.CS_MODE_THUMB)
	assembly, err := base64.StdEncoding.DecodeString(tests[0].instr)
	fmt.Printf("%X", assembly)
	if err == nil {

		defer engine.Close()

		maj, min := engine.Version()
		log.Printf("Hello Capstone! Version: %v.%v\n", maj, min)
		insns, err := engine.Disasm(
			assembly, // code buffer
			0x10000,  // starting address
			0,        // insns to disassemble, 0 for all
		)

		if err == nil {
			log.Printf("Disasm:\n")
			for _, insn := range insns {
				log.Printf("0x%x:\t%s\t\t%s\n", insn.Address, insn.Mnemonic, insn.OpStr)
			}
			return
		} else {
			t.Error("Disassembly error: ", err)
		}
		log.Fatalf("Disassembly error: %v", err)
	}
	log.Fatalf("Failed to initialize engine: %v", err)
}
