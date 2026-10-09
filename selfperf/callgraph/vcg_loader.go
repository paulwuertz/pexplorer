package callgraph

import (
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"strings"

	"github.com/paulwuertz/pexplorer/selfperf/vcg"
)

// ParseVCGDir recursively walks a directory tree and loads all .ci files into VCG graphs.
func ParseVCGDir(root string) ([]vcg.Graph, error) {
	graphs := make([]vcg.Graph, 0)
	if root == "" {
		return nil, fmt.Errorf("empty directory path")
	}

	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d == nil || d.IsDir() {
			return nil
		}
		if !strings.EqualFold(filepath.Ext(d.Name()), ".ci") {
			return nil
		}

		graph, err := ParseVCGFile(path)
		if err != nil {
			return fmt.Errorf("parse VCG file %q: %w", path, err)
		}
		graphs = append(graphs, graph)
		return nil
	})
	if err != nil {
		return nil, err
	}

	return graphs, nil
}

// ParseVCGFiles parses a list of VCG files or directories. Directories are expanded recursively.
func ParseVCGFiles(paths []string) ([]vcg.Graph, error) {
	graphs := make([]vcg.Graph, 0)
	for _, path := range paths {
		if path == "" {
			continue
		}

		info, err := os.Stat(path)
		if err != nil {
			graph, parseErr := ParseVCGFile(path)
			if parseErr != nil {
				return nil, fmt.Errorf("parse %q: %w", path, parseErr)
			}
			graphs = append(graphs, graph)
			continue
		}
		if info.IsDir() {
			more, err := ParseVCGDir(path)
			if err != nil {
				return nil, err
			}
			graphs = append(graphs, more...)
			continue
		}
		if strings.EqualFold(filepath.Ext(info.Name()), ".ci") {
			graph, err := ParseVCGFile(path)
			if err != nil {
				return nil, err
			}
			graphs = append(graphs, graph)
		}
	}
	return graphs, nil
}

// ParseVCGFile reads a single VCG file and parses it into a Graph.
func ParseVCGFile(path string) (vcg.Graph, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return vcg.Graph{}, fmt.Errorf("read file: %w", err)
	}
	return vcg.ParseString(string(data))
}
