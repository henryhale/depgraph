package output

import (
	"encoding/json"
	"log"
	"math"
	"strconv"

	"github.com/henryhale/depgraph/internal/graph"
)

type jcNode struct {
	ID     string  `json:"id"`
	Parent string  `json:"parent,omitempty"`
	Type   string  `json:"type"`
	X      float64 `json:"x"`
	Y      float64 `json:"y"`
	Width  float64 `json:"width"`
	Height float64 `json:"height"`
	Text   string  `json:"text,omitempty"`
	Label  string  `json:"label,omitempty"`
}

type jcEdge struct {
	ID       string `json:"id"`
	Label    string `json:"label"`
	FromNode string `json:"fromNode"`
	FromEnd  string `json:"fromEnd"`
	ToNode   string `json:"toNode"`
	ToEnd    string `json:"toEnd"`
}

type jcGraph struct {
	Nodes []jcNode `json:"nodes"`
	Edges []jcEdge `json:"edges,omitempty"`
}

func JSONCanvas(deps *graph.DependencyGraph) string {
	graph := graph.GenerateGraphData(deps)

	output := jcGraph{}

	const (
		nodeWidth  = 100.00
		nodeHeight = 50.0
		padding    = 50.0
	)

	// calculate grid positions
	cols := int(math.Ceil(math.Sqrt(float64(len(graph.Nodes)))))
	rows := (len(graph.Nodes) + cols - 1) / cols

	// transform nodes
	for i, node := range graph.Nodes {
		col := i % cols
		row := i / rows

		posX := float64(col)*(nodeWidth+padding) + padding
		posY := float64(row)*(nodeHeight+padding) + padding

		jnode := jcNode{
			ID:     node.ID,
			Parent: node.Parent,
			Type:   node.Type,
			X:      posX,
			Y:      posY,
			Width:  nodeWidth,
			Height: nodeHeight,
		}

		if node.Type == "group" {
			jnode.Label = node.Label
		} else {
			jnode.Text = node.Label
		}

		output.Nodes = append(output.Nodes, jnode)
	}

	// transform edges
	for i, edge := range graph.Edges {
		jedge := jcEdge{
			ID:       strconv.Itoa(i),
			Label:    edge.Label,
			FromNode: edge.From,
			FromEnd:  "none",
			ToNode:   edge.To,
			ToEnd:    "arrow",
		}

		output.Edges = append(output.Edges, jedge)
	}

	outputJSON, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		log.Fatal("error: failed to marshal output - jsoncanvas")
	}

	return string(outputJSON)

}
