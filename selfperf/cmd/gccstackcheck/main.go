package main

import (
	"flag"
	"log"
	"strings"

	"github.com/paulwuertz/pexplorer/selfperf/callgraph"
	"github.com/paulwuertz/pexplorer/selfperf/config"
	"github.com/paulwuertz/pexplorer/selfperf/rtos"
	"github.com/paulwuertz/pexplorer/selfperf/stackusage"
	"github.com/paulwuertz/pexplorer/selfperf/symbolextraction"
)

func main() {
	buildDir := flag.String("d", "", "build directory containing .ci and .su files")
	configFile := flag.String("c", "", "JSON file containing a PexplorerConfig with thread names and stack sizes")
	flag.Parse()
	if *buildDir == "" && flag.NArg() > 0 {
		*buildDir = flag.Arg(0)
	}
	if *configFile == "" && flag.NArg() > 1 {
		*configFile = flag.Arg(1)
	}
	if *buildDir == "" || *configFile == "" {
		log.Fatal("usage: gccstackcheck -d <build-dir> -c <pexplorer-config.json>")
	}

	cfg, err := config.Import_config_from_file(*configFile)
	if err != nil {
		log.Fatal(err)
	}

	report, threads, err := buildReportFromBuildDir(*buildDir, cfg)
	if err != nil {
		log.Fatal(err)
	}
	if err := rtos.PrintStackStats(threads, report.UnresolvedStats); err != nil {
		log.Fatal(err)
	}
}

func buildReportFromBuildDir(buildDir string, cfg config.PexplorerConfig) (*symbolextraction.SElfReport, []config.RTOSThread, error) {
	graphs, err := callgraph.ParseVCGDir(buildDir)
	if err != nil {
		return nil, nil, err
	}
	callList := callgraph.GraphsToFunctionCallList(graphs)
	stackEntries, err := stackusage.ParseDir(buildDir)
	if err != nil {
		return nil, nil, err
	}

	report := &symbolextraction.SElfReport{
		Functions:  []symbolextraction.FunctionSymbol{},
		Addr2FnMap: map[uint64]*symbolextraction.FunctionSymbol{},
		Name2FnMap: map[string]*symbolextraction.FunctionSymbol{},
		Info:       []string{},
		Errors:     []string{},
		UnresolvedStats: symbolextraction.UnresolvedCallStats{
			DynamicCalls:  map[string]int{},
			ResolvedCalls: map[string]int{},
		},
	}

	stackSizes := make(map[string]int64)
	for _, entry := range stackEntries {
		name := normalizeFuncName(entry.Function)
		if name == "" {
			continue
		}
		if cur, ok := stackSizes[name]; !ok || int64(entry.Size) > cur {
			stackSizes[name] = int64(entry.Size)
		}
	}

	nameMap := make(map[string]*symbolextraction.FunctionSymbol)
	for _, call := range callList {
		fromName := normalizeFuncName(call.From)
		if fromName != "" {
			if _, ok := nameMap[fromName]; !ok {
				fn := &symbolextraction.FunctionSymbol{Name: fromName}
				if size, ok := stackSizes[fromName]; ok {
					fn.StackSize = size
				}
				nameMap[fromName] = fn
			}
		}
		for _, toNameRaw := range call.To {
			toName := normalizeFuncName(toNameRaw)
			if toName == "" {
				continue
			}
			if _, ok := nameMap[toName]; !ok {
				fn := &symbolextraction.FunctionSymbol{Name: toName}
				if size, ok := stackSizes[toName]; ok {
					fn.StackSize = size
				}
				nameMap[toName] = fn
			}
		}
	}
	for _, entry := range stackEntries {
		name := normalizeFuncName(entry.Function)
		if name == "" {
			continue
		}
		if fn, ok := nameMap[name]; ok {
			fn.StackSize = int64(entry.Size)
		}
	}

	for _, fn := range nameMap {
		report.Functions = append(report.Functions, *fn)
	}
	for i := range report.Functions {
		fn := &report.Functions[i]
		report.Name2FnMap[normalizeFuncName(fn.Name)] = fn
		fn.Address = uint64(i + 1)
		report.Addr2FnMap[fn.Address] = fn
	}

	seenCalls := make(map[string]struct{})
	for _, call := range callList {
		fromName := normalizeFuncName(call.From)
		fromFn, ok := report.Name2FnMap[fromName]
		if !ok {
			continue
		}
		for _, toNameRaw := range call.To {
			toName := normalizeFuncName(toNameRaw)
			toFn, ok := report.Name2FnMap[toName]
			if !ok {
				continue
			}
			key := fromName + "->" + toName
			if _, seen := seenCalls[key]; seen {
				continue
			}
			seenCalls[key] = struct{}{}
			callAddr := toFn.Address
			fromFn.Callees = append(fromFn.Callees, symbolextraction.FunctionCall{
				CallFromFunctionName: fromName,
				CallFrom:             &fromFn.Address,
				CallTo:               &callAddr,
				CallToFunctionName:   toName,
			})
			toFn.Callers = append(toFn.Callers, symbolextraction.FunctionCall{
				CallFromFunctionName: fromName,
				CallFrom:             &fromFn.Address,
				CallTo:               &callAddr,
				CallToFunctionName:   toName,
			})
		}
	}

	threadsOut := make([]config.RTOSThread, 0, len(cfg.Threads))
	for _, thread := range cfg.Threads {
		threadName := normalizeFuncName(thread.ThreadEntryName)
		rootFn, ok := report.Name2FnMap[threadName]
		if !ok {
			log.Printf("thread %q not found in call graph, skipping", thread.ThreadEntryName)
			continue
		}
		stackSize := uint64(rootFn.StackSize)
		if thread.Size != 0 {
			stackSize = thread.Size
		}
		tree, overflowPaths := rootFn.GetCallTreeJson(report, stackSize, 0)
		used := uint64(rootFn.StackSize)
		worstBranch := symbolextraction.CallBranch{CallList: []symbolextraction.CallNode{rootFn.ToUnlinkedCallNode()}, StackSize: rootFn.StackSize}
		if len(tree.Branches) > 0 {
			used = uint64(tree.Branches[0].StackSize)
			worstBranch = tree.Branches[0]
			rootFn.MaxStackSizeCallees = tree.Branches[0].StackSize
			tree.Tree.MaxStackSizeCallees = tree.Branches[0].StackSize
		}
		threadsOut = append(threadsOut, config.RTOSThread{
			ThreadEntryName:   rootFn.Name,
			StackVariableName: thread.StackVariableName,
			Size:              stackSize,
			Used:              used,
			NrUnresolvedCalls: uint64(len(tree.UnresolvedCalls)),
			NrOverflowPaths:   overflowPaths,
			WorstStackBranch:  worstBranch,
			Calltree:          tree,
		})
	}

	report.UnresolvedStats = config.GetUnresolvedCallStats(report, nil)
	return report, threadsOut, nil
}

func normalizeFuncName(name string) string {
	name = strings.TrimSpace(name)
	name = strings.Trim(name, "\"'")
	if idx := strings.Index(name, "."); idx >= 0 {
		name = name[:idx]
	}
	return name
}
