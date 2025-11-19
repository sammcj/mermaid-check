package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestSankeyParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
		check   func(*testing.T, ast.Diagram)
	}{
		{
			name: "simple sankey diagram",
			source: `sankey-beta

Node1,Node2,10
Node2,Node3,20
Node1,Node3,5`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				sankey, ok := d.(*ast.SankeyDiagram)
				if !ok {
					t.Fatalf("expected *ast.SankeyDiagram, got %T", d)
				}
				if len(sankey.Links) != 3 {
					t.Errorf("expected 3 links, got %d", len(sankey.Links))
				}
				if sankey.Links[0].Source != "Node1" {
					t.Errorf("expected source 'Node1', got %q", sankey.Links[0].Source)
				}
				if sankey.Links[0].Target != "Node2" {
					t.Errorf("expected target 'Node2', got %q", sankey.Links[0].Target)
				}
				if sankey.Links[0].Value != 10 {
					t.Errorf("expected value 10, got %f", sankey.Links[0].Value)
				}
			},
		},
		{
			name: "complex sankey with multiple flows",
			source: `sankey-beta

A,B,100
A,C,50
B,D,75
B,E,25
C,D,30
C,E,20`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				sankey, ok := d.(*ast.SankeyDiagram)
				if !ok {
					t.Fatalf("expected *ast.SankeyDiagram, got %T", d)
				}
				if len(sankey.Links) != 6 {
					t.Errorf("expected 6 links, got %d", len(sankey.Links))
				}
			},
		},
		{
			name: "sankey with decimal values",
			source: `sankey-beta

Source,Target,42.5
Target,Destination,30.75`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				sankey, ok := d.(*ast.SankeyDiagram)
				if !ok {
					t.Fatalf("expected *ast.SankeyDiagram, got %T", d)
				}
				if sankey.Links[0].Value != 42.5 {
					t.Errorf("expected value 42.5, got %f", sankey.Links[0].Value)
				}
				if sankey.Links[1].Value != 30.75 {
					t.Errorf("expected value 30.75, got %f", sankey.Links[1].Value)
				}
			},
		},
		{
			name: "sankey with whitespace",
			source: `sankey-beta

  Node A  ,  Node B  ,  25.5
Node C, Node D, 10`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				sankey, ok := d.(*ast.SankeyDiagram)
				if !ok {
					t.Fatalf("expected *ast.SankeyDiagram, got %T", d)
				}
				if sankey.Links[0].Source != "Node A" {
					t.Errorf("expected source 'Node A', got %q", sankey.Links[0].Source)
				}
				if sankey.Links[0].Target != "Node B" {
					t.Errorf("expected target 'Node B', got %q", sankey.Links[0].Target)
				}
				if sankey.Links[0].Value != 25.5 {
					t.Errorf("expected value 25.5, got %f", sankey.Links[0].Value)
				}
			},
		},
		{
			name: "sankey with comments",
			source: `sankey-beta

%% This is a comment
A,B,100
%% Another comment
B,C,50`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				sankey, ok := d.(*ast.SankeyDiagram)
				if !ok {
					t.Fatalf("expected *ast.SankeyDiagram, got %T", d)
				}
				if len(sankey.Links) != 2 {
					t.Errorf("expected 2 links, got %d", len(sankey.Links))
				}
			},
		},
		{
			name: "sankey with empty lines",
			source: `sankey-beta

A,B,10

B,C,20

`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				sankey, ok := d.(*ast.SankeyDiagram)
				if !ok {
					t.Fatalf("expected *ast.SankeyDiagram, got %T", d)
				}
				if len(sankey.Links) != 2 {
					t.Errorf("expected 2 links, got %d", len(sankey.Links))
				}
			},
		},
		{
			name:    "invalid header",
			source:  "sankey\nA,B,10",
			wantErr: true,
		},
		{
			name:    "wrong header",
			source:  "pie\nA,B,10",
			wantErr: true,
		},
		{
			name:    "empty diagram",
			source:  "sankey-beta\n",
			wantErr: true,
		},
		{
			name: "only comments",
			source: `sankey-beta
%% Only comments here`,
			wantErr: true,
		},
		{
			name: "negative value",
			source: `sankey-beta

A,B,-10`,
			wantErr: true,
		},
		{
			name: "zero value",
			source: `sankey-beta

A,B,0`,
			wantErr: true,
		},
		{
			name: "invalid format - missing value",
			source: `sankey-beta

A,B`,
			wantErr: true,
		},
		{
			name: "invalid format - too many fields",
			source: `sankey-beta

A,B,10,Extra`,
			wantErr: true,
		},
		{
			name: "invalid format - non-numeric value",
			source: `sankey-beta

A,B,invalid`,
			wantErr: true,
		},
		{
			name: "empty source node",
			source: `sankey-beta

,B,10`,
			wantErr: true,
		},
		{
			name: "empty target node",
			source: `sankey-beta

A,,10`,
			wantErr: true,
		},
		{
			name: "self-loop",
			source: `sankey-beta

A,A,10`,
			wantErr: true,
		},
		{
			name: "whitespace-only node name",
			source: `sankey-beta

  , B, 10`,
			wantErr: true,
		},
	}

	p := parser.NewSankeyParser()
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

func TestSankeyParser_SupportedTypes(t *testing.T) {
	p := parser.NewSankeyParser()
	types := p.SupportedTypes()
	if len(types) != 1 || types[0] != "sankey" {
		t.Errorf("expected [sankey], got %v", types)
	}
}
