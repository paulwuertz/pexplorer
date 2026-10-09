package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/paulwuertz/pexplorer/selfperf/callgraph"
)

func main() {
	buildDir := flag.String("d", "", "build directory containing .ci files")
	outfile := flag.String("o", "", "output file for the JSON call list")
	flag.Parse()
	if *buildDir == "" && flag.NArg() > 0 {
		*buildDir = flag.Arg(0)
	}
	if *buildDir == "" {
		log.Fatal("Please provide a build directory with -d <dir> or as the first positional argument")
	}

	graphs, err := callgraph.ParseVCGDir(*buildDir)
	if err != nil {
		log.Fatal(err)
	}
	callList := callgraph.GraphsToFunctionCallList(graphs)
	data, err := json.MarshalIndent(callList, "", "  ")
	if err != nil {
		log.Fatal(err)
	}

	if *outfile == "" {
		fmt.Println(string(data))
		return
	}
	if err := os.WriteFile(*outfile, data, 0o644); err != nil {
		log.Fatal(err)
	}
}
