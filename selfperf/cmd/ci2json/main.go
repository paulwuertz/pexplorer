package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"os"

	"github.com/paulwuertz/pexplorer/selfperf/stackusage"
)

func main() {
	buildDir := flag.String("d", "", "build directory containing .su files")
	outfile := flag.String("o", "", "output file for the JSON stack-usage entries")
	flag.Parse()
	if *buildDir == "" && flag.NArg() > 0 {
		*buildDir = flag.Arg(0)
	}
	if *buildDir == "" {
		log.Fatal("Please provide a build directory with -d <dir> or as the first positional argument")
	}

	entries, err := stackusage.ParseDir(*buildDir)
	if err != nil {
		log.Fatal(err)
	}

	data, err := json.MarshalIndent(entries, "", "  ")
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
