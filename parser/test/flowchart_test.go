package parser_test

import (
	"maps"
	"os"
	"path/filepath"
	"slices"
	"testing"

	"github.com/sammcj/mermaid-check/ast"
	"github.com/sammcj/mermaid-check/parser"
)

func TestNewFlowchartParser(t *testing.T) {
	p := parser.NewFlowchartParser()
	if p == nil {
		t.Fatal("parser is nil")
	}
}

func TestParseSimple(t *testing.T) {
	p := parser.NewFlowchartParser()

	source := `flowchart TD
    A --> B`

	d, err := p.Parse(source)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	if d == nil {
		t.Fatal("diagram is nil")
	}

	diagram, ok := d.(*ast.Flowchart)
	if !ok {
		t.Fatalf("expected *ast.Flowchart, got %T", d)
	}

	if diagram.Type != "flowchart" {
		t.Errorf("expected type 'flowchart', got %q", diagram.Type)
	}

	if diagram.Direction != "TD" {
		t.Errorf("expected direction 'TD', got %q", diagram.Direction)
	}
}

func TestParseTestDataFiles(t *testing.T) {
	p := parser.NewFlowchartParser()

	testDataDir := "../../testdata/flowchart"
	files := []string{
		"valid-flowchart-1.mmd",
		"valid-flowchart-2.mmd",
		"valid-graph-lr-1.mmd",
		"valid-graph-tb-1.mmd",
		"valid-graph-td-1.mmd",
	}

	for _, filename := range files {
		t.Run(filename, func(t *testing.T) {
			path := filepath.Join(testDataDir, filename)
			data, err := os.ReadFile(path) //nolint:gosec // Test file paths are safe
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			diagram, err := p.Parse(string(data))
			if err != nil {
				t.Errorf("failed to parse %s: %v", filename, err)
				t.Logf("Source:\n%s", string(data))
			} else if diagram == nil {
				t.Errorf("diagram is nil for %s", filename)
			}
		})
	}
}

func TestParseSubgraphTitle(t *testing.T) {
	p := parser.NewFlowchartParser()
	tests := []struct {
		name   string
		src    string
		want   string
		wantID string
	}{
		{"id with bracket label", "flowchart TD\n subgraph one[Group One]\n a --> b\n end", "Group One", "one"},
		{"bare id", "flowchart TD\n subgraph one\n a --> b\n end", "one", "one"},
		{"quoted title", "flowchart TD\n subgraph \"My Group\"\n a --> b\n end", "My Group", ""},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := p.Parse(tt.src)
			if err != nil {
				t.Fatalf("Parse error: %v", err)
			}
			fc, ok := d.(*ast.Flowchart)
			if !ok {
				t.Fatalf("expected *ast.Flowchart, got %T", d)
			}
			var sg *ast.Subgraph
			for _, s := range fc.Statements {
				if g, ok := s.(*ast.Subgraph); ok {
					sg = g
					break
				}
			}
			if sg == nil {
				t.Fatal("no subgraph statement found")
			}
			if sg.Title != tt.want {
				t.Errorf("subgraph title = %q, want %q", sg.Title, tt.want)
			}
			if sg.ID != tt.wantID {
				t.Errorf("subgraph id = %q, want %q", sg.ID, tt.wantID)
			}
		})
	}
}

func TestParseStyleStatement(t *testing.T) {
	source := `flowchart LR
    A --> B
    style A fill:#f9f,stroke:#333,stroke-width:2px`

	d, err := parser.NewFlowchartParser().Parse(source)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	var got *ast.Style
	for _, s := range d.(*ast.Flowchart).Statements {
		if st, ok := s.(*ast.Style); ok {
			got = st
		}
	}
	if got == nil {
		t.Fatal("no *ast.Style statement produced")
	}
	if got.Target != "A" {
		t.Errorf("Target = %q, want %q", got.Target, "A")
	}
	want := map[string]string{"fill": "#f9f", "stroke": "#333", "stroke-width": "2px"}
	if !maps.Equal(got.Styles, want) {
		t.Errorf("Styles = %v, want %v", got.Styles, want)
	}
}

func TestParseInlineClassShorthand(t *testing.T) {
	source := `flowchart LR
    A[Start]:::hot --> B[End]
    C:::cold
    D["a:::b"] --> E`

	d, err := parser.NewFlowchartParser().Parse(source)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}
	diagram := d.(*ast.Flowchart)

	var links int
	var nodes []string
	assigned := map[string]string{}
	for _, s := range diagram.Statements {
		switch v := s.(type) {
		case *ast.Link:
			links++
		case *ast.NodeDef:
			nodes = append(nodes, v.ID+"="+v.Label)
		case *ast.ClassAssignment:
			for _, id := range v.NodeIDs {
				assigned[id] = v.ClassName
			}
		}
	}

	if links != 2 {
		t.Errorf("got %d links, want 2 (a :::class must not drop the statement)", links)
	}
	// Bare and unlabeled node references (C, E) produce no NodeDef; only the
	// labeled ones do.
	wantNodes := []string{"A=Start", "B=End", `D="a:::b"`}
	for _, want := range wantNodes {
		if !slices.Contains(nodes, want) {
			t.Errorf("missing node %q; got %v", want, nodes)
		}
	}
	wantAssigned := map[string]string{"A": "hot", "C": "cold"}
	if !maps.Equal(assigned, wantAssigned) {
		t.Errorf("class assignments = %v, want %v", assigned, wantAssigned)
	}
}

// A `:::` sequence that is part of prose or an edge label is not a class
// shorthand: peeling it off would silently rewrite the diagram's text.
func TestInlineClassShorthandLeavesTextAlone(t *testing.T) {
	source := `flowchart LR
    %% styling note: A:::hot is applied below
    A -->|weight a:::b| B
    A:::hot`

	d, err := parser.NewFlowchartParser().Parse(source)
	if err != nil {
		t.Fatalf("failed to parse: %v", err)
	}

	var comments, labels []string
	assigned := map[string]string{}
	for _, s := range d.(*ast.Flowchart).Statements {
		switch v := s.(type) {
		case *ast.Comment:
			comments = append(comments, v.Text)
		case *ast.Link:
			labels = append(labels, v.Label)
		case *ast.ClassAssignment:
			for _, id := range v.NodeIDs {
				assigned[id] = v.ClassName
			}
		}
	}

	wantComments := []string{"styling note: A:::hot is applied below"}
	if !slices.Equal(comments, wantComments) {
		t.Errorf("comments = %v, want %v", comments, wantComments)
	}
	wantLabels := []string{"weight a:::b"}
	if !slices.Equal(labels, wantLabels) {
		t.Errorf("link labels = %v, want %v", labels, wantLabels)
	}
	wantAssigned := map[string]string{"A": "hot"}
	if !maps.Equal(assigned, wantAssigned) {
		t.Errorf("class assignments = %v, want %v", assigned, wantAssigned)
	}
}

func TestParseInlineClassEdgeCases(t *testing.T) {
	tests := []struct {
		name     string
		source   string
		links    int
		assigned map[string]string
	}{
		{
			// nodeDefPattern accepts the asymmetric `A>Start]` shape, so the
			// peel has to step over it the same way it steps over `A[Start]`.
			name:     "asymmetric node shape",
			source:   "flowchart LR\n    A>Start]:::hot --> B",
			links:    1,
			assigned: map[string]string{"A": "hot"},
		},
		{
			// Mermaid's idString admits MINUS, so the class name runs past it.
			name:     "hyphenated class name",
			source:   "flowchart LR\n    A[x]:::warning-high --> B",
			links:    1,
			assigned: map[string]string{"A": "warning-high"},
		},
		{
			// ...but a trailing `-` belongs to the arrow, not the class name.
			name:     "class name butted against an arrow",
			source:   "flowchart LR\n    A[x]:::hot-->B",
			links:    1,
			assigned: map[string]string{"A": "hot"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d, err := parser.NewFlowchartParser().Parse(tt.source)
			if err != nil {
				t.Fatalf("failed to parse: %v", err)
			}
			links := 0
			assigned := map[string]string{}
			for _, s := range d.(*ast.Flowchart).Statements {
				switch v := s.(type) {
				case *ast.Link:
					links++
				case *ast.ClassAssignment:
					for _, id := range v.NodeIDs {
						assigned[id] = v.ClassName
					}
				}
			}
			if links != tt.links {
				t.Errorf("got %d links, want %d", links, tt.links)
			}
			if !maps.Equal(assigned, tt.assigned) {
				t.Errorf("class assignments = %v, want %v", assigned, tt.assigned)
			}
		})
	}
}
