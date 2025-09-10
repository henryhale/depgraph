package output

import (
	"slices"

	"github.com/henryhale/depgraph/internal/graph"
)

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
