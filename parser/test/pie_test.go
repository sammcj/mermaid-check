package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestPieParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
		check   func(*testing.T, ast.Diagram)
	}{
		{
			name: "simple pie chart",
			source: `pie
    "Apples" : 42.5
    "Oranges" : 30.0
    "Bananas" : 27.5`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				pie, ok := d.(*ast.PieDiagram)
				if !ok {
					t.Fatalf("expected *ast.PieDiagram, got %T", d)
				}
				if len(pie.DataEntries) != 3 {
					t.Errorf("expected 3 entries, got %d", len(pie.DataEntries))
				}
				if pie.Title != "" {
					t.Errorf("expected no title, got %q", pie.Title)
				}
				if pie.ShowData {
					t.Error("expected ShowData to be false")
				}
			},
		},
		{
			name: "pie chart with title",
			source: `pie title Sales Distribution
    "Product A" : 42.5
    "Product B" : 57.5`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				pie, ok := d.(*ast.PieDiagram)
				if !ok {
					t.Fatalf("expected *ast.PieDiagram, got %T", d)
				}
				if pie.Title != "Sales Distribution" {
					t.Errorf("expected title 'Sales Distribution', got %q", pie.Title)
				}
			},
		},
		{
			name: "pie chart with showData",
			source: `pie showData
    "Category A" : 60
    "Category B" : 40`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				pie, ok := d.(*ast.PieDiagram)
				if !ok {
					t.Fatalf("expected *ast.PieDiagram, got %T", d)
				}
				if !pie.ShowData {
					t.Error("expected ShowData to be true")
				}
			},
		},
		{
			name: "pie chart with title and showData",
			source: `pie showData title Revenue Breakdown
    "Q1" : 25
    "Q2" : 25
    "Q3" : 25
    "Q4" : 25`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				pie, ok := d.(*ast.PieDiagram)
				if !ok {
					t.Fatalf("expected *ast.PieDiagram, got %T", d)
				}
				if pie.Title != "Revenue Breakdown" {
					t.Errorf("expected title 'Revenue Breakdown', got %q", pie.Title)
				}
				if !pie.ShowData {
					t.Error("expected ShowData to be true")
				}
				if len(pie.DataEntries) != 4 {
					t.Errorf("expected 4 entries, got %d", len(pie.DataEntries))
				}
			},
		},
		{
			name: "pie chart with comments",
			source: `pie title Test
    %% This is a comment
    "Item A" : 50
    %% Another comment
    "Item B" : 50`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				pie, ok := d.(*ast.PieDiagram)
				if !ok {
					t.Fatalf("expected *ast.PieDiagram, got %T", d)
				}
				if len(pie.DataEntries) != 2 {
					t.Errorf("expected 2 entries, got %d", len(pie.DataEntries))
				}
			},
		},
		{
			name:    "invalid header",
			source:  "notpie\n",
			wantErr: true,
		},
		{
			name:    "empty diagram",
			source:  "pie\n",
			wantErr: true,
		},
		{
			name: "negative value",
			source: `pie
    "Item" : -10`,
			wantErr: true,
		},
		{
			name: "zero value",
			source: `pie
    "Item" : 0`,
			wantErr: true,
		},
		{
			name: "invalid entry format",
			source: `pie
    Invalid entry without quotes`,
			wantErr: true,
		},
	}

	p := parser.NewPieParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram, err := p.Parse(tt.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, diagram)
			}
		})
	}
}

func TestPieParser_SupportedTypes(t *testing.T) {
	p := parser.NewPieParser()
	types := p.SupportedTypes()
	if len(types) != 1 || types[0] != "pie" {
		t.Errorf("expected [pie], got %v", types)
	}
}
