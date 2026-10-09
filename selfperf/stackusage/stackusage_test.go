package stackusage

import (
	"os"
	"testing"
)

func TestParseString(t *testing.T) {
	src := `/home/paul/git/cannecti/zephyr/lib/utils/ring_buffer.c:12:10:ring_buf_area_claim	16	static
/home/paul/git/cannecti/zephyr/include/zephyr/sys/ring_buffer.h:304:24:ring_buf_put_claim	12	static
`

	entries, err := ParseString(src)
	if err != nil {
		t.Fatalf("ParseString returned error: %v", err)
	}
	if len(entries) != 2 {
		t.Fatalf("expected 2 entries, got %d", len(entries))
	}

	if entries[0].File != "/home/paul/git/cannecti/zephyr/lib/utils/ring_buffer.c" {
		t.Fatalf("unexpected file path: %q", entries[0].File)
	}
	if entries[0].Function != "ring_buf_area_claim" {
		t.Fatalf("unexpected function: %q", entries[0].Function)
	}
	if entries[0].Size != 16 {
		t.Fatalf("unexpected size: %d", entries[0].Size)
	}
	if entries[0].Qualifier != "static" {
		t.Fatalf("unexpected qualifier: %q", entries[0].Qualifier)
	}
	if entries[1].Line != 304 || entries[1].Column != 24 {
		t.Fatalf("unexpected source position: %d:%d", entries[1].Line, entries[1].Column)
	}
}

func TestParseDir(t *testing.T) {
	dir := t.TempDir()
	p := dir + "/example.su"
	payload := "/tmp/example.c:10:2:demo_fn	8	dynamic\n"
	if err := os.WriteFile(p, []byte(payload), 0o600); err != nil {
		t.Fatalf("write test file: %v", err)
	}

	entries, err := ParseDir(dir)
	if err != nil {
		t.Fatalf("ParseDir returned error: %v", err)
	}
	if len(entries) != 1 {
		t.Fatalf("expected 1 entry, got %d", len(entries))
	}
	if entries[0].Function != "demo_fn" {
		t.Fatalf("unexpected function: %q", entries[0].Function)
	}
}
