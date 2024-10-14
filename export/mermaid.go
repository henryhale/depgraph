package export

import (
	"regexp"
	"strings"

	"github.com/henryhale/depgraph/lang"
)

func Mermaid(deps *lang.DependencyGraph) string {
	indent1 := "  "
	indent2 := indent1 + indent1
	nl := "\n"

	graph := "graph LR" + nl
	edges := ""

	// to avoid duplicates, everything has a unique id
	ids := make(map[string]struct{})
	idExists := func(id string) bool {
		_, found := ids[id]
		return found
	}

	for file, analysis := range *deps {
		fileID := cleanNodeID(file)

		if idExists(fileID) {
			continue
		}
		ids[fileID] = struct{}{}

		subgraph := indent1 + "subgraph " + fileID + nl

		// exports as nodes in subgraph
		for _, export := range analysis.Exports {
			id := fileID + "_" + cleanNodeID(export)

			if idExists(file + id) {
				continue
			}
			ids[file+id] = struct{}{}

			text := `["` + cleanLabel(export) + `"]`
			subgraph += indent2 + id + text + nl
		}

		subgraph += indent1 + "end" + nl

		graph += subgraph

		// resolve edges
		for importedFile, items := range analysis.Imports {
			iFileID := cleanNodeID(importedFile)
			for _, item := range items {
				itemnode := iFileID
				if !strings.Contains(item, "*") {
					itemnode += "_" + item
				}

				if idExists("$$edge$$" + file + itemnode) {
					continue
				}
				ids["$$edge$$"+file+itemnode] = struct{}{}

				edges += indent1 + fileID + "-->|uses|" + itemnode + nl
			}
		}
	}

	return graph + edges
}

func cleanNodeID(s string) string {
	re := regexp.MustCompile(`[^a-zA-Z0-9\/\._]`)
	return re.ReplaceAllString(s, "_")
}

func cleanLabel(s string) string {
	t := strings.ReplaceAll(s, "*", "#42;")
	return t
}
