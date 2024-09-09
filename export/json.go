package export

import (
	"encoding/json"
	"fmt"
)

func JSON(deps *AnalysisResultMap) string {
	graph := GenerateGraphData(deps)

	output, err := json.MarshalIndent(graph, "", "  ")

	if err != nil {
		fmt.Println("error: failed to marshal output - json")
	}

	return string(output)

}
