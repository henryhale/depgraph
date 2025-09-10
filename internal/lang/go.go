package lang

import (
	"regexp"
	"strings"

	"github.com/henryhale/depgraph/internal/util"
)

// Golang
var GO = Language{
	Extensions: []string{".go"},
	Rules: []Rule{
		// single import
		{`import\s+"([^"]+)"`, 1, -1, false},
		// multiple imports
		{`import\s+\(\s*([^)]*)\s*\)`, 1, -1, false},
		// exported function
		{`func\s+([A-Z]+\w*)\s*\(.*\)\s*?(\w*)\s*{`, -1, 1, true},
		// exported variable constant
		{`(?:var|const)\s+([A-Z]\w*)\s*?=`, -1, 1, true},
		// exported type
		{`type\s+([A-Z]\w*)\s+`, -1, 1, true},
	},
	Comments: &util.Comments,
	Extract: func(options *ExtractorOptions) {
		rule := options.Rule
		match := *options.Match
		result := options.Result
		file := *options.File
		replacers := options.Replacers

		// exports
		if rule.Export && rule.Items > 0 {
			result.AddExport(*util.Explode(match[rule.Items])...)
			return
		}

		// imports
		if !rule.Export && rule.File > 0 {
			partial := match[rule.File]

			imports := extractImports(partial)

			checkStandardLib, _ := regexp.Compile(`[/\\]`)

			for _, importpath := range imports {
				if checkStandardLib.Match([]byte(importpath)) {
					result.AddImport(importpath, []string{"*"})
					return
				}

				ipath := util.FullPath(importpath, file, replacers)
				result.AddImport(ipath, []string{"*"})
			}

		}

	},
}

func extractImports(s string) []string {
	s = strings.TrimSpace(s)
	badChars := []string{" ", "\t", "\r", "\""}
	for _, char := range badChars {
		s = strings.ReplaceAll(s, char, "")
	}
	s = strings.ReplaceAll(s, "\n\n", "\n")
	list := strings.Split(s, "\n")

	// fmt.Println(len(list), list)
	return list
}
