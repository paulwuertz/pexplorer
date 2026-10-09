package rtos

import (
	"testing"

	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func TestGetVarByAddrReturnsActualStructElement(t *testing.T) {
	report := &symbolextraction.SElfReport{
		Variables: []symbolextraction.VariableSymbol{
			{Name: "first", Address: 0x100, Data: []byte{1}},
			{Name: "second", Address: 0x200, Data: []byte{2}},
		},
	}

	got := GetVarByAddr(0x100, report)
	if got == nil {
		t.Fatal("expected variable, got nil")
	}
	if got.Name != "first" {
		t.Fatalf("expected first variable, got %q", got.Name)
	}

	got = GetVarByAddr(0x200, report)
	if got == nil {
		t.Fatal("expected variable, got nil")
	}
	if got.Name != "second" {
		t.Fatalf("expected second variable, got %q", got.Name)
	}
}

func TestFindStaticZephyrRtosThreadsHandlesMissingTypeData(t *testing.T) {
	report := &symbolextraction.SElfReport{
		Types:     nil,
		Variables: nil,
	}

	threads := FindStaticZephyrRtosThreads(report)
	if len(threads) != 0 {
		t.Fatalf("expected no threads, got %d", len(threads))
	}
	if len(report.Errors) != 1 {
		t.Fatalf("expected one recovery error, got %d: %v", len(report.Errors), report.Errors)
	}
	if report.Errors[0] != "_static_thread_data type not found in DWARF metadata" {
		t.Fatalf("unexpected error: %v", report.Errors)
	}
}
