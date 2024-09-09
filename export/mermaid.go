package export

import (
	"strings"
)

func Mermaid(deps *AnalysisResultMap) string {
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
		fileId := strings.ReplaceAll(file, "/", "_")
		fileId = strings.ReplaceAll(file, "@", "_")

		if idExists(fileId) {
			continue
		}
		ids[fileId] = struct{}{}

		subgraph := indent1 + "subgraph " + fileId + nl

		// exports as nodes in subgraph
		for _, export := range analysis.Exports {
			id := fileId + "_" + export

			if idExists(id) {
				continue
			}
			ids[id] = struct{}{}

			text := "[" + export + "]"
			subgraph += indent2 + id + text + nl
		}

		subgraph += indent1 + "end" + nl

		graph += subgraph

		// resolve edges
		for importedFile, items := range analysis.Imports {
			iFileId := strings.ReplaceAll(importedFile, "/", "_")
			iFileId = strings.ReplaceAll(importedFile, "@", "_")
			for _, item := range items {
				itemnode := iFileId
				if !strings.Contains(item, "*") {
					itemnode += "_" + item
				}

				if idExists("edge" + itemnode) {
					continue
				}
				ids["edge" + itemnode] = struct{}{}

				edges += indent1 + fileId + "-->|uses|" + itemnode + nl
			}
		}
	}

	return graph + edges
}
