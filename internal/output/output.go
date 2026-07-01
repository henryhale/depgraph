package output

import (
	"slices"
	"sort"

	"github.com/henryhale/depgraph/internal/graph"
)

// sortedFiles returns the dependency graph's file paths in lexical order so
// that formatters emit reproducible output regardless of map iteration order.
func sortedFiles(deps *graph.DependencyGraph) []string {
	keys := make([]string, 0, len(*deps))
	for k := range *deps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// sortedImports returns a file's imported paths in lexical order.
func sortedImports(imports map[string][]string) []string {
	keys := make([]string, 0, len(imports))
	for k := range imports {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

func FormatSupported(f *string) bool {
	formats := []string{
		"json",
		"mermaid",
		"jsoncanvas",
		"dot",
	}
	return slices.Contains(formats, *f)
}

func Format(f *string, deps *graph.DependencyGraph) string {
	var output string
	switch *f {
	case "json":
		output = JSON(deps)
	case "mermaid":
		output = Mermaid(deps)
	case "jsoncanvas":
		output = JSONCanvas(deps)
	case "dot":
		output = DOT(deps)

	// case "html":

	default:
		output = ""
	}

	return output
}
