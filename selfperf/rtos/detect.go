package rtos

import (
	"fmt"
	"slices"
	"strings"

	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func ScanForRtosFeatures(s *symbolextraction.SElfReport) {
	// TODO only zephyr for now... how to definitly detect it though...
	for _, v := range s.Variables {
		// what to get from them...
		if strings.Contains(v.Name, "dts_ord") {
			fmt.Println("found device:", v.Name, v)
		}
		// static threads created by macro
		if v.VariableType == "k_thread" {
			fmt.Println("found static thread:", v)
		}
	}
	// static threads created by macro
	var thread_create_fn symbolextraction.FunctionSymbol
	if slices.ContainsFunc(s.Functions, func(e symbolextraction.FunctionSymbol) bool {
		if e.Name == "z_impl_k_thread_create" {
			thread_create_fn = e
			return true
		}
		return false
	}) {
		thread_creating_fns := thread_create_fn.Callers
		fmt.Println("found dynamic thread:", len(thread_creating_fns))
		for i, t_addr := range thread_creating_fns {
			t, known := s.Addr2FnMap[*t_addr.CallFrom]
			if known {
				fmt.Println("\t#", i, "initialized in", t.Name)
			} else {
				fmt.Println("\t#", i, "initialized by unknown fn from addr", t_addr)
			}
		}
	}
}
