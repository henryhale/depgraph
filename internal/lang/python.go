package lang

import (
	"path/filepath"
	"strings"
)

// Python

// extensions
var pyExts = []string{".py"}

var Python = Language{
	Extensions: pyExts,
	Rules: []Rule{
		// import
		{`(?m)^import\s+([^\n]+)`, 1, -1, false},
		// import as
		{`(?m)^import\s+([^\n]+)`, 1, -1, false}, // duplicate, remove later
		// from import
		{`(?m)^from\s+([\w\.]+)\s+import\s+([^\n]+)`, 1, 2, false},
		// from import as
		{`(?m)^from\s+([\w\.]+)\s+import\s+([^\n]+)`, 1, 2, false}, // duplicate
		// def export
		{`(?m)^def\s+([a-zA-Z_]\w*)\s*\(`, -1, 1, true},
		// class export
		{`(?m)^class\s+([A-Z][\w]*)`, -1, 1, true},
	},
	Comments: &[]string{
		`(?m)^#.*$`,       // single-line comments
		`(?s)"""(.*?)"""`, // triple double-quote docstrings
		`(?s)'''(.*?)'''`, // triple single-quote docstrings
	},
	Extract: func(options *ExtractorOptions) {
		rule := options.Rule
		match := *options.Match
		result := options.Result
		file := *options.File
		// replacers := options.Replacers

		// exports
		if rule.Export && rule.Items > 0 {
			result.AddExport(match[rule.Items])
			return
		}

		// imports
		if !rule.Export && rule.File > 0 {
			rawImport := match[rule.File]
			// handle "as" by taking the module name
			importpath := strings.Split(rawImport, " as ")[0]
			// resolve relative modules
			if !filepath.IsAbs(importpath) && !strings.Contains(importpath, "/") {
				importpath = filepath.Join(filepath.Dir(file), importpath+".py")
			}
			var items []string
			if rule.Items > 0 {
				rawItems := match[rule.Items]
				// split by comma, handle "as"
				parts := strings.Split(rawItems, ",")
				for _, part := range parts {
					item := strings.TrimSpace(strings.Split(part, " as ")[0])
					if item != "" && item != "*" {
						items = append(items, item)
					}
				}
			}
			result.AddImport(importpath, items)
		}
	},
}
