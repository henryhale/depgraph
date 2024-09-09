package lang

// PHP
var PHP = Language{
	Extensions: []string{".php"},
	LocateImports: false,
	Rules: []Rule{
		// namespace
		Rule{`namespace\s+([a-zA-Z0-9_\\]+)\s*;`, 1, -1, false},
		// use statement
		Rule{`use\s+(?:(function|const)\s+)?([a-zA-Z0-9_\\]+)\s*;`, 2, -1, false},
		// file import
		Rule{`(require|include|require_once|include_once)\s*\(\s*['"]([^'"]+)['"]\s*\)\s*;`, 2, -1, false},
		// file import no parentheses
		Rule{`(require|include|require_once|include_once)\s+['"]([^'"]+)['"]\s*;`, 2, -1, false},
		// class definition
		Rule{`class\s+([A-Z][a-zA-Z0-9_]*)`, -1, 1, true},
		// function definition
		Rule{`function\s+([a-zA-Z0-9_]+)\s*\([^)]*\)\s*{`, -1, 1, true},
		// constant definition
		Rule{`const\s+([A-Z][A-Z0-9_]*)\s*=`, -1, 1, true},
	},
}
