package lang

import (
	"github.com/henryhale/depgraph/util"
)

// JavaScript
var JavaScript = Language{
	Extensions: []string{".js", ".mjs", ".cjs", ".ts", ".mts", ".cts"},
	LocateImports: false,
	Rules: []Rule{
		// importESNamed
		Rule{`import\s*\{([^}]*)\}\s*from\s*['"](.*)['"];?`, 2, 1, false},
		// importESDefault
		Rule{`import\s+(\w+)\s+from\s*['"](.*)['"];?`, 2, 1, false},
		// importESNamespace
		Rule{`import\s+(\*)\s+as\s+\w+\s+from\s*['"](.*)['"];?`, 2, 1, false},
		// importESSideEffect
		Rule{`import\s*['"](.*)['"];?`, 1, -1, false},
		// importCJSNamed
		Rule{`const\s+(\w+)\s*=\s*require\s*\(['"](.*)['"]\);?`, 2, 1, false},
		// importCJSDestructured
		Rule{`const\s*\{([^}]*)\}\s*=\s*require\s*\(['"](.*)['"]\);?`, 2, 1, false},
		// exportESNamed
		Rule{`export\s*\{([^}]*)\};?`, -1, 1, true},
		// exportESConstFunc
		Rule{`export\s+(const|class|function|let|var)\s+(\w+)`, -1, 2, true},
		// exportESDefault
		Rule{`export\s+default\s+(\w+);?`, -1, 1, true},
		// exportCJSObject
		Rule{`module\.exports\s*=\s*\{([^}]*)\};?`, -1, 1, true},
		// exportCJSNamed
		Rule{`module\.exports\s*=\s*(\w+);?`, -1, 1, true},
		// exportCJSProperty
		Rule{`exports\.(\w+)\s*=\s*.+;?`, -1, 1, true},
		// importTSType
		Rule{`import\s+type\s+\{([^}]*)\}\s+from\s*['"](.*)['"];?`, 2, 1, false},
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
			result.AddImport(importpath, *util.Explode(match[rule.Items]))
		}
	},
}
