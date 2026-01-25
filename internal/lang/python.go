package lang

import (
	"strings"

	"github.com/henryhale/depgraph/internal/util"
)

// Python

// extensions
var pyExts = []string{".py"}

var Python = Language{
	Extensions: pyExts,
	Rules: []Rule{
		// import
		{`^import\s+([\w\.]+)`, 1, -1, false},
		// import as
		{`^import\s+([\w\.]+)\s+as\s+(\w+)`, 1, -1, false},
		// from import
		{`^from\s+([\w\.]+)\s+import\s+([\w\*,\s]+)`, 1, 2, false},
		// from import as
		{`^from\s+([\w\.]+)\s+import\s+(\w+)\s+as\s+(\w+)`, 1, 2, false},
		// def export
		{`^def\s+([a-zA-Z_]\w*)\s*\(`, -1, 1, true},
		// class export
		{`^class\s+([A-Z][\w]*)`, -1, 1, true},
	},
	Comments: &[]string{"#"},
	Extract: func(options *ExtractorOptions) {
		rule := options.Rule
		match := *options.Match
		result := options.Result
		file := *options.File
		replacers := options.Replacers

		// exports
		if rule.Export && rule.Items > 0 {
			result.AddExport(match[rule.Items])
			return
		}

		// imports
		if !rule.Export && rule.File > 0 {
			importpath := util.FullPath(match[rule.File], file, replacers)
			var items []string
			if rule.Items > 0 {
				// split by comma, trim spaces
				parts := strings.Split(match[rule.Items], ",")
				for _, part := range parts {
					item := strings.TrimSpace(part)
					if item != "" {
						items = append(items, item)
					}
				}
			} else {
				// for simple import, perhaps add the module as item? But typically, it's the module
				// for now, empty items
			}
			result.AddImport(importpath, items)
		}
	},
}
