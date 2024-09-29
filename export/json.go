package export

import (
	"encoding/json"
	"log"

	"github.com/henryhale/depgraph/lang"
)

func JSON(deps *lang.DependencyGraph) string {
	graph := GenerateGraphData(deps)

	output, err := json.MarshalIndent(graph, "", "  ")

	if err != nil {
		log.Fatal("error: failed to marshal output - json")
	}

	return string(output)

}
