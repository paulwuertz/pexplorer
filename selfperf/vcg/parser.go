package vcg

import (
	"fmt"
	"regexp"
	"strings"
)

// Graph is a compact representation of a VCG graph that is easy to marshal to JSON.
type Graph struct {
	Title string `json:"title,omitempty"`
	Nodes []Node `json:"nodes,omitempty"`
	Edges []Edge `json:"edges,omitempty"`
}

// Node is a single VCG node.
type Node struct {
	Title string `json:"title,omitempty"`
	Label string `json:"label,omitempty"`
}

// Edge is a single VCG edge.
type Edge struct {
	Label      string `json:"label,omitempty"`
	Sourcename string `json:"sourcename,omitempty"`
	Targetname string `json:"targetname,omitempty"`
}

// ParseString parses a VCG-formatted string and converts it into a JSON-friendly Graph.
func ParseString(input string) (Graph, error) {
	g := Graph{}
	trimmed := strings.TrimSpace(input)
	if trimmed == "" {
		return g, fmt.Errorf("empty VCG input")
	}

	// accept a minimal wrapper like: graph: { ... }
	body, err := extractGraphBody(trimmed)
	if err != nil {
		return g, err
	}

	for _, entry := range splitTopLevelEntries(body) {
		token := strings.TrimSpace(entry)
		if token == "" {
			continue
		}
		switch {
		case strings.HasPrefix(token, "title:"):
			g.Title = parseStringValue(token[len("title:"):])
		case strings.HasPrefix(token, "node:"):
			n, err := parseNode(token[len("node:"):])
			if err != nil {
				return g, err
			}
			g.Nodes = append(g.Nodes, n)
		case strings.HasPrefix(token, "edge:"):
			e, err := parseEdge(token[len("edge:"):])
			if err != nil {
				return g, err
			}
			g.Edges = append(g.Edges, e)
		}
	}
	return g, nil
}

func extractGraphBody(input string) (string, error) {
	start := strings.Index(input, "graph")
	if start == -1 {
		return "", fmt.Errorf("missing graph declaration")
	}
	bodyStart := strings.Index(input[start:], "{")
	if bodyStart == -1 {
		return "", fmt.Errorf("missing graph body")
	}
	body := input[start+bodyStart+1:]
	depth := 1
	for i, r := range body {
		if r == '{' {
			depth++
		}
		if r == '}' {
			depth--
			if depth == 0 {
				return body[:i], nil
			}
		}
	}
	return "", fmt.Errorf("unterminated graph body")
}

func splitTopLevelEntries(body string) []string {
	parts := make([]string, 0)
	current := strings.Builder{}
	depth := 0
	inQuote := false
	quoteChar := byte(0)

	flush := func() {
		if piece := strings.TrimSpace(current.String()); piece != "" {
			parts = append(parts, piece)
		}
		current.Reset()
	}

	for i := 0; i < len(body); i++ {
		r := body[i]
		switch {
		case inQuote:
			current.WriteByte(r)
			if r == '\\' && i+1 < len(body) {
				current.WriteByte(body[i+1])
				i++
				continue
			}
			if r == quoteChar {
				inQuote = false
				quoteChar = 0
			}
		case r == '"' || r == '\'':
			inQuote = true
			quoteChar = r
			current.WriteByte(r)
		case r == '{' || r == '[' || r == '(':
			depth++
			current.WriteByte(r)
		case r == '}' || r == ']' || r == ')':
			if depth > 0 {
				depth--
			}
			current.WriteByte(r)
		case depth == 0 && (r == ';' || r == '\n' || r == '\r'):
			flush()
		default:
			current.WriteByte(r)
		}
	}
	flush()
	return parts
}

func loadAttrsFromBraces(input string) (attrs map[string]string, err error) {
	input = strings.TrimSpace(input)
	attrs = make(map[string]string)
	if input == "" {
		return attrs, fmt.Errorf("empty edge/node declaration")
	}
	containsBody := strings.Contains(input, "{") && strings.Contains(input, "}")
	if !containsBody {
		return attrs, fmt.Errorf("node without attributes: %q", input)
	}
	inner := input[strings.Index(input, "{")+1 : strings.LastIndex(input, "}")]
	for _, kv := range splitPairs(inner) {
		key, val := kv[0], kv[1]
		attrs[key] = val
	}
	return attrs, nil
}

func parseNode(input string) (Node, error) {
	attrs, err := loadAttrsFromBraces(input)
	n := Node{}
	if v, ok := attrs["title"]; ok {
		n.Title = v
	}
	if v, ok := attrs["label"]; ok {
		n.Label = v
	}
	return n, err
}

func parseEdge(input string) (Edge, error) {
	attrs, err := loadAttrsFromBraces(input)
	e := Edge{}
	if v, ok := attrs["sourcename"]; ok {
		e.Sourcename = v
	}
	if v, ok := attrs["targetname"]; ok {
		e.Targetname = v
	}
	if v, ok := attrs["label"]; ok {
		e.Label = v
	}
	return e, err
}

func splitPairs(text string) [][2]string {
	pairs := make([][2]string, 0)
	for i := 0; i < len(text); {
		for i < len(text) && isWhitespace(text[i]) {
			i++
		}
		if i >= len(text) {
			break
		}
		keyStart := i
		for i < len(text) && text[i] != ':' && !isWhitespace(text[i]) {
			i++
		}
		if i >= len(text) || text[i] != ':' {
			break
		}
		key := strings.TrimSpace(text[keyStart:i])
		value := ""
		i++
		for i < len(text) && isWhitespace(text[i]) {
			i++
		}
		if i >= len(text) {
			pairs = append(pairs, [2]string{key, value})
			break
		}
		valueStart := i
		if text[i] == '"' || text[i] == '\'' {
			quote := text[i]
			i++
			for i < len(text) {
				if text[i] == '\\' && i+1 < len(text) {
					i++
				}
				if text[i] == quote {
					i++
					break
				}
				i++
			}
		} else {
			for i < len(text) && !isWhitespace(text[i]) {
				i++
			}
		}
		pairs = append(pairs, [2]string{key, parseStringValue(strings.TrimSpace(text[valueStart:i]))})
	}
	return pairs
}

func isWhitespace(b byte) bool {
	return b == ' ' || b == '\t' || b == '\n' || b == '\r'
}

func parseStringValue(raw string) string {
	trimmed := strings.TrimSpace(raw)
	trimmed = strings.TrimSuffix(trimmed, ";")
	if len(trimmed) >= 2 && trimmed[0] == '"' && trimmed[len(trimmed)-1] == '"' {
		unquoted := trimmed[1 : len(trimmed)-1]
		unquoted = strings.ReplaceAll(unquoted, "\\\"", "\"")
		unquoted = strings.ReplaceAll(unquoted, "\\n", "\n")
		return unquoted
	}
	if len(trimmed) >= 2 && trimmed[0] == '\'' && trimmed[len(trimmed)-1] == '\'' {
		unquoted := trimmed[1 : len(trimmed)-1]
		unquoted = strings.ReplaceAll(unquoted, "\\'", "'")
		return unquoted
	}
	if strings.HasPrefix(trimmed, "\"") && strings.HasSuffix(trimmed, "\"") {
		return parseStringValue(trimmed[1 : len(trimmed)-1])
	}
	if ok, _ := regexp.MatchString(`^[-+]?\d+$`, trimmed); ok {
		return trimmed
	}
	return strings.TrimSpace(trimmed)
}
