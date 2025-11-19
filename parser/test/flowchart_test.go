package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
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
