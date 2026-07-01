package output

import (
	"flag"
	"os"
	"path/filepath"
	"testing"

	"github.com/henryhale/depgraph/internal/graph"
	"github.com/henryhale/depgraph/internal/lang"
)

var update = flag.Bool("update", false, "update golden files")

func fixture() graph.DependencyGraph {
	return graph.DependencyGraph{
		"a.js": lang.SourceFile{
			Imports: map[string][]string{"b.js": {"foo", "bar"}},
			Exports: []string{"main"},
			Local:   true,
		},
		"b.js": lang.SourceFile{
			Imports: map[string][]string{},
			Exports: []string{"foo", "bar"},
			Local:   true,
		},
	}
}

func TestFormattersGolden(t *testing.T) {
	deps := fixture()

	cases := []struct {
		name string
		fn   func(*graph.DependencyGraph) string
	}{
		{"mermaid", Mermaid},
		{"dot", DOT},
		{"json", JSON},
		{"jsoncanvas", JSONCanvas},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			got := tc.fn(&deps)
			golden := filepath.Join("testdata", tc.name+".golden")

			if *update {
				if err := os.WriteFile(golden, []byte(got), 0644); err != nil {
					t.Fatalf("writing golden: %v", err)
				}
			}

			want, err := os.ReadFile(golden)
			if err != nil {
				t.Fatalf("reading golden (run with -update to create): %v", err)
			}
			if got != string(want) {
				t.Errorf("%s output mismatch\n--- got ---\n%s\n--- want ---\n%s", tc.name, got, string(want))
			}
		})
	}
}

// Formatters must be deterministic across repeated calls.
func TestFormattersDeterministic(t *testing.T) {
	deps := fixture()
	fns := map[string]func(*graph.DependencyGraph) string{
		"mermaid": Mermaid, "dot": DOT, "json": JSON, "jsoncanvas": JSONCanvas,
	}
	for name, fn := range fns {
		first := fn(&deps)
		for i := 0; i < 20; i++ {
			if fn(&deps) != first {
				t.Fatalf("%s not deterministic on iteration %d", name, i)
			}
		}
	}
}
