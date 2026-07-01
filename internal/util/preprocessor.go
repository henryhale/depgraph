package util

import "strings"

// Comments documents the comment styles handled by Preprocess. It is retained
// for reference (and as the per-language hook in lang.Language); the current
// scanner handles C-style `//` line and `/* */` block comments, which covers
// every supported language (js/ts, c/c++, go).
var Comments = []string{
	// single line comment
	`//`,
	// multi-line comment
	`/* */`,
}

// preprocessor scan states.
const (
	psNormal = iota
	psLineComment
	psBlockComment
	psString
)

// Preprocess strips C-style line (`//`) and block (`/* */`) comments from
// source code before the regex extractors run.
//
// Unlike a naive regex pass it is string-literal aware: `//` and `/*` inside a
// string (e.g. a URL like "https://example.com") are preserved, block comments
// are matched non-greedily so code between two comments survives, and a line
// comment at end-of-file (no trailing newline) is still removed. Double-quoted,
// single-quoted and backtick strings are recognised, with backslash escapes in
// the two quoted forms (backticks are treated as raw, per Go).
//
// The second parameter is unused; it is kept so callers passing a language's
// comment rules continue to compile.
func Preprocess(code string, _ *[]string) string {
	var b strings.Builder
	b.Grow(len(code))

	state := psNormal
	var quote byte

	for i := 0; i < len(code); i++ {
		c := code[i]
		switch state {
		case psNormal:
			switch {
			case c == '/' && i+1 < len(code) && code[i+1] == '/':
				state = psLineComment
				i++
			case c == '/' && i+1 < len(code) && code[i+1] == '*':
				state = psBlockComment
				i++
			case c == '"' || c == '\'' || c == '`':
				state = psString
				quote = c
				b.WriteByte(c)
			default:
				b.WriteByte(c)
			}
		case psLineComment:
			if c == '\n' || c == '\r' {
				state = psNormal
				b.WriteByte(c)
			}
		case psBlockComment:
			if c == '*' && i+1 < len(code) && code[i+1] == '/' {
				state = psNormal
				i++
			}
		case psString:
			b.WriteByte(c)
			if c == '\\' && quote != '`' && i+1 < len(code) {
				// preserve the escaped character verbatim
				b.WriteByte(code[i+1])
				i++
			} else if c == quote {
				state = psNormal
			}
		}
	}

	return b.String()
}
