package lang

import (
	"strings"
)

type SourceFile struct {
	Imports map[string][]string
	Exports []string
	Local bool
}

func (r *SourceFile) AddExport(exports ...string) {
	for _, export := range exports {
		r.Exports = append(r.Exports, export)
	}
}

func (r *SourceFile) AddImport(path string, items []string) {
	_, exists := r.Imports[path]
	if exists {
		r.Imports[path] = append(r.Imports[path], items...)
	} else {
		r.Imports[path] = items
	}
}

type DependencyGraph map[string]SourceFile

type Rule struct {
	RegExp string
	File   int
	Items  int
	Export bool
}

type ExtractorOptions struct {
	Rule *Rule
	Match *[]string
	Result *SourceFile
	File *string
	Replacers *map[string]string
}

type Language struct {
	// valid extensions
	Extensions []string
	// import and export statement rules
	Rules []Rule
	// comments
	Comments *[]string
	// for languages using the import statements but
	// don't specify actual imports like go, dart, py
	LocateImports bool
	// imports and exports extractor
	Extract func(*ExtractorOptions)
}

func Get(ext string) (lang Language, supported bool) {
	supported = true
	switch strings.ToLower(ext) {

	// js/ts -> js.go
	case "js":
		lang = JavaScript
	case "ts":
		lang = JavaScript

	// c/c++ -> c.go
	// case "c": lang = CC
	// case "cpp": lang = CC

	// go -> go.go
	// case "go": lang = GO

	// php -> php.go
	// case "php": lang = PHP

	// ...

	default:
		supported = false
	}
	return
}
