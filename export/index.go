package export

import (
	"github.com/henryhale/depgraph/lang"
	"slices"
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

func Format(f *string, deps *lang.DependencyGraph) string {
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
