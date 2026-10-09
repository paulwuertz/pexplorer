package vcg

import "testing"

func TestParseString(t *testing.T) {
	input := `graph: {
 title: "callgraph"
 node: { title: "n0" label: "main" shape: box }
 node: { title: "n1" label: "helper" }
 edge: { sourcename: "n0" targetname: "n1" label: "calls" }
}`

	g, err := ParseString(input)
	if err != nil {
		t.Fatalf("ParseString returned error: %v", err)
	}
	if g.Title != "callgraph" {
		t.Fatalf("unexpected title: got %q want %q", g.Title, "callgraph")
	}
	if len(g.Nodes) != 2 {
		t.Fatalf("unexpected node count: got %d want 2", len(g.Nodes))
	}
	if g.Nodes[0].Title != "n0" {
		t.Fatalf("unexpected first node title: got %q want %q", g.Nodes[0].Title, "n0")
	}
	if len(g.Edges) != 1 {
		t.Fatalf("unexpected edge count: got %d want 1", len(g.Edges))
	}
	if g.Edges[0].Source != "n0" || g.Edges[0].Target != "n1" {
		t.Fatalf("unexpected edge endpoints: got %q -> %q", g.Edges[0].Source, g.Edges[0].Target)
	}
	if g.Edges[0].Label != "calls" {
		t.Fatalf("unexpected edge label: got %q want %q", g.Edges[0].Label, "calls")
	}
}
