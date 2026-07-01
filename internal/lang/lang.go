package lang

import (
	"regexp"
	"strings"

	"github.com/henryhale/depgraph/internal/util"
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

// Analyze runs a language's rules over a single file's source code and returns
// the imports and exports it declares. The code is first preprocessed to strip
// comments so commented-out statements are not mistaken for real ones. This is
// the core per-file extraction step, factored out so it can be tested directly.
func Analyze(l Language, code string, file string, replacers *map[string]string) SourceFile {
	result := SourceFile{
		Imports: make(map[string][]string),
		Exports: []string{},
		Local:   true,
	}

	src := util.Preprocess(code, l.Comments)

	opts := &ExtractorOptions{
		Result:    &result,
		File:      &file,
		Replacers: replacers,
	}

	for i := range l.Rules {
		rule := l.Rules[i]
		re := regexp.MustCompile(rule.RegExp)
		matches := re.FindAllStringSubmatch(src, -1)
		for _, match := range matches {
			m := match
			opts.Rule = &rule
			opts.Match = &m
			l.Extract(opts)
		}
	}

	return result
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
