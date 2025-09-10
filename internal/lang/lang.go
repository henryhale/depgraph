package lang

import (
	"strings"
)

type SourceFile struct {
	Imports map[string][]string
	Exports []string
	Local   bool
}

func (r *SourceFile) AddExport(exports ...string) {
	r.Exports = append(r.Exports, exports...)
}

func (r *SourceFile) AddImport(path string, items []string) {
	_, exists := r.Imports[path]
	if exists {
		r.Imports[path] = append(r.Imports[path], items...)
	} else {
		r.Imports[path] = items
	}
}

type Rule struct {
	RegExp string
	File   int
	Items  int
	Export bool
}

type ExtractorOptions struct {
	Rule      *Rule
	Match     *[]string
	Result    *SourceFile
	File      *string
	Replacers *map[string]string
}

type Language struct {
	// valid extensions
	Extensions []string
	// import and export statement rules
	Rules []Rule
	// comments
	Comments *[]string
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
	case "c":
		lang = CC
	case "cpp":
		lang = CC

	// go.go
	case "go":
		lang = GO

	default:
		supported = false
	}

	return
}
