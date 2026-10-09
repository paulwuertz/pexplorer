package stackusage

import (
	"bufio"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strconv"
	"strings"
)

// Entry describes one function from a GCC -fstack-usage .su file.
type Entry struct {
	File      string `json:"file"`
	Line      int    `json:"line"`
	Column    int    `json:"column"`
	Function  string `json:"function"`
	Size      int64  `json:"size"`
	Qualifier string `json:"qualifier,omitempty"`
}

// ParseString parses the contents of one .su file.
func ParseString(input string) ([]Entry, error) {
	var entries []Entry
	scanner := bufio.NewScanner(strings.NewReader(input))
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" {
			continue
		}
		entry, err := parseLine(line)
		if err != nil {
			return nil, err
		}
		entries = append(entries, entry)
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return entries, nil
}

// Parse is an alias for ParseString.
func Parse(input string) ([]Entry, error) {
	return ParseString(input)
}

// ParseFile parses one .su file.
func ParseFile(path string) ([]Entry, error) {
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("read %q: %w", path, err)
	}
	entries, err := ParseString(string(data))
	if err != nil {
		return nil, fmt.Errorf("parse %q: %w", path, err)
	}
	return entries, nil
}

// ParseDir recursively walks a directory tree and loads all .su files.
func ParseDir(root string) ([]Entry, error) {
	if root == "" {
		return nil, fmt.Errorf("empty directory path")
	}
	var files []string
	err := filepath.WalkDir(root, func(path string, d fs.DirEntry, walkErr error) error {
		if walkErr != nil {
			return walkErr
		}
		if d == nil || d.IsDir() {
			return nil
		}
		if strings.EqualFold(filepath.Ext(d.Name()), ".su") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	sort.Strings(files)

	entries := make([]Entry, 0)
	for _, file := range files {
		parsed, err := ParseFile(file)
		if err != nil {
			return nil, err
		}
		entries = append(entries, parsed...)
	}
	return entries, nil
}

// ParseFiles parses a list of .su files or directories recursively.
func ParseFiles(paths []string) ([]Entry, error) {
	entries := make([]Entry, 0)
	for _, path := range paths {
		if path == "" {
			continue
		}
		info, err := os.Stat(path)
		if err != nil {
			parsed, parseErr := ParseFile(path)
			if parseErr != nil {
				return nil, fmt.Errorf("parse %q: %w", path, parseErr)
			}
			entries = append(entries, parsed...)
			continue
		}
		if info.IsDir() {
			more, parseErr := ParseDir(path)
			if parseErr != nil {
				return nil, parseErr
			}
			entries = append(entries, more...)
			continue
		}
		if strings.EqualFold(filepath.Ext(info.Name()), ".su") {
			parsed, parseErr := ParseFile(path)
			if parseErr != nil {
				return nil, parseErr
			}
			entries = append(entries, parsed...)
		}
	}
	return entries, nil
}

func parseLine(line string) (Entry, error) {
	fields := strings.Fields(line)
	if len(fields) < 2 {
		return Entry{}, fmt.Errorf("invalid stack-usage line %q", line)
	}

	meta := fields[0]
	file, lineNo, column, functionName, err := splitMeta(meta)
	if err != nil {
		return Entry{}, err
	}

	size, err := strconv.ParseInt(fields[1], 10, 64)
	if err != nil {
		return Entry{}, fmt.Errorf("invalid stack size in %q: %w", line, err)
	}

	qualifier := ""
	if len(fields) > 2 {
		qualifier = fields[2]
	}

	return Entry{
		File:      file,
		Line:      lineNo,
		Column:    column,
		Function:  functionName,
		Size:      size,
		Qualifier: qualifier,
	}, nil
}

func splitMeta(meta string) (file string, line int, column int, function string, err error) {
	lastColon := strings.LastIndex(meta, ":")
	if lastColon < 0 {
		return "", 0, 0, "", fmt.Errorf("missing function name in %q", meta)
	}
	function = meta[lastColon+1:]
	if function == "" {
		return "", 0, 0, "", fmt.Errorf("missing function name in %q", meta)
	}

	beforeFunc := meta[:lastColon]
	prevColon := strings.LastIndex(beforeFunc, ":")
	if prevColon < 0 {
		return "", 0, 0, "", fmt.Errorf("missing column in %q", meta)
	}
	column, err = strconv.Atoi(beforeFunc[prevColon+1:])
	if err != nil {
		return "", 0, 0, "", fmt.Errorf("invalid column in %q: %w", meta, err)
	}

	beforeColumn := beforeFunc[:prevColon]
	lineColon := strings.LastIndex(beforeColumn, ":")
	if lineColon < 0 {
		return "", 0, 0, "", fmt.Errorf("missing line in %q", meta)
	}
	line, err = strconv.Atoi(beforeColumn[lineColon+1:])
	if err != nil {
		return "", 0, 0, "", fmt.Errorf("invalid line in %q: %w", meta, err)
	}
	file = beforeColumn[:lineColon]
	return
}
