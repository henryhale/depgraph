package output

import (
	"encoding/json"
	"log"

	"github.com/henryhale/depgraph/internal/graph"
)

func JSON(deps *graph.DependencyGraph) string {
	graph := graph.GenerateGraphData(deps)

	output, err := json.MarshalIndent(graph, "", "  ")

	if err != nil {
		log.Fatal("error: failed to marshal output - json")
	}

	return string(output)

}
