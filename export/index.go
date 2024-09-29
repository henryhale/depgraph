package export

import (
	"github.com/henryhale/depgraph/lang"
)

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
