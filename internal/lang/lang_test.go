package lang

import (
	"reflect"
	"sort"
	"testing"
)

func TestGet(t *testing.T) {
	cases := []struct {
		ext       string
		supported bool
	}{
		{"js", true},
		{"ts", true},
		{"JS", true}, // case-insensitive
		{"c", true},
		{"cpp", true},
		{"go", true},
		{"py", false},
		{"", false},
	}

	for _, tc := range cases {
		t.Run(tc.ext, func(t *testing.T) {
			_, supported := Get(tc.ext)
			if supported != tc.supported {
				t.Errorf("Get(%q) supported = %v, want %v", tc.ext, supported, tc.supported)
			}
		})
	}
}

// analyzeExports is a helper that returns a file's exports sorted, for
// order-independent comparison.
func analyzeExports(l Language, code string) []string {
	r := Analyze(l, code, "file", &map[string]string{})
	got := append([]string{}, r.Exports...)
	sort.Strings(got)
	return got
}

func TestAnalyzeJavaScriptExports(t *testing.T) {
	cases := []struct {
		name string
		code string
		want []string
	}{
		{"es named", `export { a, b };`, []string{"a", "b"}},
		{"es const", `export const foo = 1`, []string{"foo"}},
		{"es function", `export function bar() {}`, []string{"bar"}},
		{"es default", `export default thing;`, []string{"thing"}},
		{"cjs property", `exports.baz = 5;`, []string{"baz"}},
		{"commented out", "// export const nope = 1\nexport const yes = 2", []string{"yes"}},
		{"block commented", "/* export const nope = 1 */\nexport const yes = 2", []string{"yes"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := analyzeExports(JavaScript, tc.code)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("exports = %v, want %v", got, tc.want)
			}
		})
	}
}

func TestAnalyzeJavaScriptImports(t *testing.T) {
	// Use non-existent relative targets so path resolution is deterministic
	// (findJsFile falls back to the joined path when no file is found).
	cases := []struct {
		name string
		code string
		path string   // expected resolved import path
		want []string // expected imported items
	}{
		{"named", `import { a, b } from "./missing"`, "dir/missing", []string{"a", "b"}},
		{"default", `import foo from "./missing"`, "dir/missing", []string{"foo"}},
		{"external pkg", `import { x } from "react"`, "react", []string{"x"}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			r := Analyze(JavaScript, tc.code, "dir/file.js", &map[string]string{})
			items, ok := r.Imports[tc.path]
			if !ok {
				t.Fatalf("import %q not found; imports = %v", tc.path, r.Imports)
			}
			sort.Strings(items)
			if !reflect.DeepEqual(items, tc.want) {
				t.Errorf("items = %v, want %v", items, tc.want)
			}
		})
	}
}

func TestAnalyzeGoImports(t *testing.T) {
	code := `package main

import "fmt"

import (
	"strings"
	"github.com/foo/bar"
)
`
	r := Analyze(GO, code, "main.go", &map[string]string{})
	for _, want := range []string{"fmt", "strings", "github.com/foo/bar"} {
		if _, ok := r.Imports[want]; !ok {
			t.Errorf("expected import %q; imports = %v", want, r.Imports)
		}
	}
}

func TestAnalyzeGoExports(t *testing.T) {
	code := `package main

func Exported() {}
func unexported() {}
type Widget struct{}
const Answer = 42
`
	got := analyzeExports(GO, code)
	want := []string{"Answer", "Exported", "Widget"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("exports = %v, want %v", got, want)
	}
}

func TestAnalyzeCIncludes(t *testing.T) {
	code := `#include <stdio.h>
#include "calc.h"
`
	r := Analyze(CC, code, "src/main.c", &map[string]string{})
	if _, ok := r.Imports["stdio.h"]; !ok {
		t.Errorf("expected standard include stdio.h; imports = %v", r.Imports)
	}
	// user include gets a "./" prefix then joined against the file's dir
	if _, ok := r.Imports["src/calc.h"]; !ok {
		t.Errorf("expected user include src/calc.h; imports = %v", r.Imports)
	}
}
