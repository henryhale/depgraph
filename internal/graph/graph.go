package graph

import (
	"sort"

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

	nodeIDs := make(map[string]struct{})
	nodeExists := func(id string) bool {
		_, found := nodeIDs[id]
		return found
	}

	// edges are deduplicated on their own key space so that a file->file
	// dependency carrying several distinct imports is not collapsed into one.
	edgeIDs := make(map[string]struct{})

	// iterate files in a stable order so the generated graph is reproducible
	// across runs (Go map iteration order is randomized).
	for _, file := range sortedKeys(*deps) {
		result := (*deps)[file]
		if nodeExists(file) {
			continue
		}
		nodeIDs[file] = struct{}{}

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
			if len(export) == 0 || nodeExists(id) {
				continue
			}
			nodeIDs[id] = struct{}{}

			nodes = append(nodes, Node{
				ID:     id,
				Label:  export,
				Parent: file,
				Type:   "text",
			})
		}

		// add edges for imports, one per distinct imported item
		for _, importedFile := range sortedImportKeys(result.Imports) {
			for _, item := range result.Imports[importedFile] {
				if len(item) == 0 {
					continue
				}
				// NUL separators keep the key unambiguous for arbitrary paths
				// and item names (avoids "ab"+"c" == "a"+"bc" collisions).
				id := file + "\x00" + importedFile + "\x00" + item
				if _, dup := edgeIDs[id]; dup {
					continue
				}
				edgeIDs[id] = struct{}{}

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

// sortedKeys returns the file paths of a dependency graph in lexical order.
func sortedKeys(deps DependencyGraph) []string {
	keys := make([]string, 0, len(deps))
	for k := range deps {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// sortedImportKeys returns the imported paths of a file in lexical order.
func sortedImportKeys(imports map[string][]string) []string {
	keys := make([]string, 0, len(imports))
	for k := range imports {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}
