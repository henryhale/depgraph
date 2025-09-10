package lang

import (
	"path/filepath"
	"slices"

	"github.com/henryhale/depgraph/internal/util"
)

// JavaScript

// extensions
var jsExts = []string{".js", ".mjs", ".cjs", ".ts", ".mts", ".cts", ".jsx", ".tsx"}

var JavaScript = Language{
	Extensions: jsExts,
	Rules: []Rule{
		// importESNamed
		{`import\s*\{([^}]*)\}\s*from\s*['"](.*)['"];?`, 2, 1, false},
		// importESDefault
		{`import\s+(\w+)\s+from\s*['"](.*)['"];?`, 2, 1, false},
		// importESNamespace
		{`import\s+(\*)\s+as\s+\w+\s+from\s*['"](.*)['"];?`, 2, 1, false},
		// importESSideEffect
		{`import\s*['"](.*)['"];?`, 1, -1, false},
		// importCJSNamed
		{`const\s+(\w+)\s*=\s*require\s*\(['"](.*)['"]\);?`, 2, 1, false},
		// importCJSDestructured
		{`const\s*\{([^}]*)\}\s*=\s*require\s*\(['"](.*)['"]\);?`, 2, 1, false},
		// exportESNamed
		{`export\s*\{([^}]*)\};?`, -1, 1, true},
		// exportESConstFunc
		{`export\s+(const|class|function|let|var)\s+(\w+)`, -1, 2, true},
		// exportESDefault
		{`export\s+default\s+(\w+);?`, -1, 1, true},
		// exportCJSObject
		{`module\.exports\s*=\s*\{([^}]*)\};?`, -1, 1, true},
		// exportCJSNamed
		{`module\.exports\s*=\s*(\w+);?`, -1, 1, true},
		// exportCJSProperty
		{`exports\.(\w+)\s*=\s*.+;?`, -1, 1, true},
		// importTSType
		{`import\s+type\s+\{([^}]*)\}\s+from\s*['"](.*)['"];?`, 2, 1, false},
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
		if !rule.Export && rule.Items > 0 && rule.File > 0 {
			importpath := util.FullPath(match[rule.File], file, replacers)
			importpath = findJsFile(importpath)
			result.AddImport(importpath, *util.Explode(match[rule.Items]))
		}
	},
}

/*
in case of 'demo/app' check for:
- 'demo/app.{js,ts,mjs,cjs,jsx,tsx}'
- 'demo/app/index.{js,ts,mjs,cjs,jsx,tsx}'
*/
func findJsFile(path string) string {
	ext := filepath.Ext(path)
	if len(ext) > 0 && slices.Contains(jsExts, ext) {
		return path
	}

	for _, ext := range jsExts {
		if util.FileExists(path + ext) {
			return path + ext
		}
		index := filepath.Join(path, "index"+ext)
		if util.FileExists(index) {
			return index
		}
	}

	return path
}
