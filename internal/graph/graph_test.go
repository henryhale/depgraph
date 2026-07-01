package graph

import (
	"reflect"
	"testing"

	"github.com/henryhale/depgraph/internal/lang"
)

func fixture() DependencyGraph {
	return DependencyGraph{
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

// Regression test: distinct imported items between the same pair of files must
// each produce their own edge, and must not be collapsed by the dedup key.
func TestGenerateGraphDataDistinctEdges(t *testing.T) {
	deps := fixture()
	g := GenerateGraphData(&deps)

	var count int
	for _, e := range g.Edges {
		if e.From == "a.js" && e.To == "b.js" {
			count++
		}
	}
	if count != 2 {
		t.Errorf("expected 2 distinct a.js->b.js edges, got %d (%+v)", count, g.Edges)
	}
}

// Regression test: repeated imports of the same item collapse to one edge.
func TestGenerateGraphDataDeduplicates(t *testing.T) {
	deps := DependencyGraph{
		"a.js": lang.SourceFile{
			Imports: map[string][]string{"b.js": {"foo", "foo"}},
		},
		"b.js": lang.SourceFile{Imports: map[string][]string{}},
	}
	g := GenerateGraphData(&deps)

	var count int
	for _, e := range g.Edges {
		if e.From == "a.js" && e.To == "b.js" {
			count++
		}
	}
	if count != 1 {
		t.Errorf("expected duplicate items to collapse to 1 edge, got %d", count)
	}
}

// The generated graph must be identical across runs despite Go's randomized
// map iteration order.
func TestGenerateGraphDataDeterministic(t *testing.T) {
	deps := fixture()
	first := GenerateGraphData(&deps)
	for i := 0; i < 20; i++ {
		next := GenerateGraphData(&deps)
		if !reflect.DeepEqual(first, next) {
			t.Fatalf("graph output not deterministic on iteration %d", i)
		}
	}
}
