package lang

import (
	"github.com/henryhale/depgraph/util"
	"regexp"
)

// C/C++
var CC = Language{
	Extensions:    []string{".c", ".h", ".cpp", ".hpp", ".cc", ".cxx"},
	LocateImports: false,
	Rules: []Rule{
		// standard include
		{`#include\s*<([^>]+)>`, 1, -1, false},
		// user-defined include
		{`#include\s*"([^"]+)"`, 1, -1, false},
		// function declaration
		{`\n\s*[\w\s\*]+\s+(\w+)\s*\([^)]*\)\s*;`, -1, 1, true},
		// variable declaration
		{`\n\s*extern\s+[\w\s\*]+\s+(\w+)\s*;`, -1, 1, true},
		// function definition
		{`\n\s*[\w\s\*]+\s+(\w+)\s*\([^)]*\)\s*\{[^}]*\}`, -1, 1, true},
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
			isStandard, _ := regexp.Match(`<.*>`, []byte(match[0]))
			if isStandard {
				result.AddImport(partial, []string{"*"})
				return
			}
			if len(partial) > 0 && partial[0] != '.' {
				partial = "./" + partial
			}
			importpath := util.FullPath(partial, file, replacers)
			result.AddImport(importpath, []string{"*"})
		}

	},
}
