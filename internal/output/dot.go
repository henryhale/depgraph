package output

import (
	"fmt"

	"github.com/henryhale/depgraph/internal/graph"
)

func DOT(deps *graph.DependencyGraph) string {
	nl := "\n"
	tab := "  "
	indent1 := nl + tab
	indent2 := indent1 + tab

	nodes := ""
	edges := ""
	subgraphs := ""

	// to avoid duplicates, everything has a unique id
	nxtID := 0
	ids := make(map[string]int)
	idExists := func(id string) bool {
		_, found := ids[id]
		return found
	}

	// node ids
	nid := 0
	addNode := func(file string, label string) string {
		nid = nid + 1
		ids[file+label] = nid
		snid := fmt.Sprintf("%v", nid)
		style := ""
		if label == file {
			style += " shape=ellipse style=filled fillcolor=lightgray penwidth=0.2"
		}
		nodes += indent1 + snid + " [label=\"" + label + "\"" + style + "];"
		return indent2 + snid + ";"
	}

	for file, analysis := range *deps {
		if idExists(file) {
			continue
		}
		nxtID = nxtID + 1
		ids[file] = nxtID

		subgraph := fmt.Sprintf(indent1+"subgraph cluster_%v {"+indent2+"style=dashed;", nxtID)

		// add a default node to represent internals/file itself
		internalUse := addNode(file, file)
		subgraph += internalUse

		for _, export := range analysis.Exports {
			subgraph += addNode(file, export)
		}

		subgraph += indent1 + "}"

		subgraphs += subgraph
	}

	for file, analysis := range *deps {
		for importFile, items := range analysis.Imports {
			for _, item := range items {
				fromNode, found := ids[file+file]
				if !found {
					continue
				}
				toNode, found := ids[importFile+item]
				if !found {
					continue
				}
				edges += fmt.Sprintf(indent1+"%v->%v [label=\"imports\"];", fromNode, toNode)
			}
		}
	}

	config := indent1 + `graph [bgcolor="#ffffff", rankdir=LR];`
	config += indent1 + `node [shape=box, style=filled, fontcolor="#333333", fillcolor="#ffe082"];`

	return "digraph depgraph {" + config + nodes + edges + subgraphs + nl + "}" + nl
}
