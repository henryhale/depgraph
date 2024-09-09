package lang

// C/C++
var CC = Language{
	Extensions: []string{".c", ".h", ".cpp", ".hpp", ".cc", ".cxx"},
	LocateImports: false,
	Rules: []Rule{
		// standard include
		Rule{`^#include\s*<([^>]+)>$`, 1, -1, false},
		// user-defined include
		Rule{`^#include\s*"([^"]+)"$`, 1, -1, false},
		// function declaration
		Rule{`\n\s*[\w\s\*]+\s+(\w+)\s*\([^)]*\)\s*;`, -1, 1, true},
		// variable declaration
		Rule{`\n\s*extern\s+[\w\s\*]+\s+(\w+)\s*;`, -1, 1, true},
		// function definition
		Rule{`\n\s*[\w\s\*]+\s+(\w+)\s*\([^)]*\)\s*\{[^}]*\}`, -1, 1, true},
	},
}
