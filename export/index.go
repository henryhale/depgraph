package export

import (
	"slices"
	"github.com/henryhale/depgraph/lang"
)

func FormatSupported(f *string) bool {
	formats := []string{
		"json",
		"mermaid",
		"jsoncanvas",
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

	// case "html":

	default:
		output = ""
	}

	return output
}
