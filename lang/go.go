package lang

// Go
var GO = Language{
	Extensions: []string{".go"},
	LocateImports: true,
	Rules: []Rule{
		// single import
		Rule{`import\s+"([^"]+)"`, 1, -1, false},
		// need to split multiple imports durring analysis
		// multiple imports
		Rule{`import\s+\(\s*([^)]*)\s*\)`, 1, -1, false},
		// exported function
		Rule{`func\s+([A-Z]\w*)\s*\(.*\)\s*{`, -1, 1, true},
		// exported variable/constant
		Rule{`(?:var|const)\s+([A-Z]\w*)\s*=`, -1, 1, true},
		// exported type
		Rule{`type\s+([A-Z]\w*)\s+`, -1, 1, true},
	},
}
