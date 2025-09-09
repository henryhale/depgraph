package graph

import (
	"github.com/henryhale/depgraph/internal/lang"
)

type Node struct {
	ID     string `json:"id"`
	Label  string `json:"label"`
	Parent string `json:"parent"`
	Type   string `json:"type"`
}

type Edge struct {
	From  string `json:"from"`
	To    string `json:"to"`
	Label string `json:"label"`
}

type Graph struct {
	Nodes []Node `json:"nodes"`
	Edges []Edge `json:"edges"`
}

type DependencyGraph map[string]lang.SourceFile

func GenerateGraphData(deps *DependencyGraph) *Graph {
	nodes := []Node{}
	edges := []Edge{}

	ids := make(map[string]struct{})
	idExists := func(id string) bool {
		_, found := ids[id]
		return found
	}

	for file, result := range *deps {
		if idExists(file) {
			continue
		}
		ids[file] = struct{}{}

		// add file node
		nodes = append(nodes, Node{
			ID:     file,
			Label:  file,
			Parent: "",
			Type:   "group",
		})

		// add exports as child nodes to file
		for _, export := range result.Exports {
			id := file + "_" + export
			if len(export) == 0 || idExists(id) {
				continue
			}
			ids[id] = struct{}{}

			nodes = append(nodes, Node{
				ID:     id,
				Label:  export,
				Parent: file,
				Type:   "text",
			})
		}

		// add edges for imports
		for importedFile, items := range result.Imports {
			for _, item := range items {
				id := file + importedFile
				if len(item) == 0 || idExists(id) {
					continue
				}
				ids[id] = struct{}{}

				edges = append(edges, Edge{
					From:  file,
					To:    importedFile,
					Label: "imports " + item,
				})
			}
		}

	}

	graph := Graph{Nodes: nodes, Edges: edges}

	return &graph
}
