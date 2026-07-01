package util

import (
	"strings"
	"testing"
)

func TestPreprocess(t *testing.T) {
	cases := []struct {
		name     string
		code     string
		contains []string // substrings that must survive
		absent   []string // substrings that must be stripped
	}{
		{
			name:     "line comment stripped",
			code:     "const x = 1 // drop me\nconst y = 2",
			contains: []string{"const x = 1", "const y = 2"},
			absent:   []string{"drop me"},
		},
		{
			name:     "line comment at eof without newline",
			code:     "const x = 1 // drop me",
			contains: []string{"const x = 1"},
			absent:   []string{"drop me"},
		},
		{
			name:     "block comment non-greedy keeps code between",
			code:     "/* a */ keep /* b */",
			contains: []string{"keep"},
			absent:   []string{"a", "b"},
		},
		{
			name:     "url in double-quoted string preserved",
			code:     `const u = "https://example.com/path"`,
			contains: []string{"https://example.com/path"},
		},
		{
			name:     "block-comment marker inside string preserved",
			code:     `const s = "a /* not a comment */ b"`,
			contains: []string{"a /* not a comment */ b"},
		},
		{
			name:     "escaped quote inside string does not end it early",
			code:     `const s = "he said \"// hi\"" ; const y = 2`,
			contains: []string{`const y = 2`, `// hi`},
		},
		{
			name:     "backtick string preserved",
			code:     "const s = `http://x // y`",
			contains: []string{"http://x // y"},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := Preprocess(tc.code, &Comments)
			for _, s := range tc.contains {
				if !strings.Contains(got, s) {
					t.Errorf("expected result to contain %q\n got: %q", s, got)
				}
			}
			for _, s := range tc.absent {
				if strings.Contains(got, s) {
					t.Errorf("expected result to NOT contain %q\n got: %q", s, got)
				}
			}
		})
	}
}

func TestExplode(t *testing.T) {
	got := *Explode(` a, b ,c , "d" `)
	want := []string{"a", "b", "c", "d"}
	if len(got) != len(want) {
		t.Fatalf("Explode len = %d (%v), want %d (%v)", len(got), got, len(want), want)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("Explode[%d] = %q, want %q", i, got[i], want[i])
		}
	}
}
