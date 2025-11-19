package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestMindmapParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
		check   func(*testing.T, ast.Diagram)
	}{
		{
			name: "simple mindmap with root and children",
			source: `mindmap
  root((Central Idea))
    Topic 1
    Topic 2`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				mm, ok := d.(*ast.MindmapDiagram)
				if !ok {
					t.Fatalf("expected *ast.MindmapDiagram, got %T", d)
				}
				if mm.Root == nil {
					t.Fatal("expected root node")
				}
				if mm.Root.Text != "Central Idea" {
					t.Errorf("expected root text 'Central Idea', got %q", mm.Root.Text)
				}
				if mm.Root.Shape != "(())" {
					t.Errorf("expected root shape '(())', got %q", mm.Root.Shape)
				}
				if len(mm.Root.Children) != 2 {
					t.Errorf("expected 2 children, got %d", len(mm.Root.Children))
				}
				if mm.Root.Children[0].Text != "Topic 1" {
					t.Errorf("expected first child 'Topic 1', got %q", mm.Root.Children[0].Text)
				}
			},
		},
		{
			name: "multi-level hierarchy",
			source: `mindmap
  root((Main))
    Branch 1
      Leaf 1.1
      Leaf 1.2
        Detail 1.2.1
    Branch 2
      Leaf 2.1`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				mm, ok := d.(*ast.MindmapDiagram)
				if !ok {
					t.Fatalf("expected *ast.MindmapDiagram, got %T", d)
				}
				if mm.Root == nil {
					t.Fatal("expected root node")
				}
				if len(mm.Root.Children) != 2 {
					t.Fatalf("expected 2 children, got %d", len(mm.Root.Children))
				}
				branch1 := mm.Root.Children[0]
				if len(branch1.Children) != 2 {
					t.Errorf("expected Branch 1 to have 2 children, got %d", len(branch1.Children))
				}
				if len(branch1.Children[1].Children) != 1 {
					t.Errorf("expected Leaf 1.2 to have 1 child, got %d", len(branch1.Children[1].Children))
				}
			},
		},
		{
			name: "different node shapes",
			source: `mindmap
  root((Circle))
    Square[Square node]
    Round(Round node)
    Cloud{{Cloud node}}
    Hexagon))Hexagon node((`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				mm, ok := d.(*ast.MindmapDiagram)
				if !ok {
					t.Fatalf("expected *ast.MindmapDiagram, got %T", d)
				}
				if mm.Root.Shape != "(())" {
					t.Errorf("expected root shape '(())', got %q", mm.Root.Shape)
				}
				if len(mm.Root.Children) != 4 {
					t.Fatalf("expected 4 children, got %d", len(mm.Root.Children))
				}
				if mm.Root.Children[0].Shape != "[]" {
					t.Errorf("expected Square shape '[]', got %q", mm.Root.Children[0].Shape)
				}
				if mm.Root.Children[1].Shape != "()" {
					t.Errorf("expected Round shape '()', got %q", mm.Root.Children[1].Shape)
				}
				if mm.Root.Children[2].Shape != "{{}}" {
					t.Errorf("expected Cloud shape '{{}}', got %q", mm.Root.Children[2].Shape)
				}
				if mm.Root.Children[3].Shape != "))((" {
					t.Errorf("expected Hexagon shape '))((', got %q", mm.Root.Children[3].Shape)
				}
			},
		},
		{
			name: "node with icon",
			source: `mindmap
  root((Books))
    Fiction
      ::icon(fa fa-book)
    Non-Fiction`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				mm, ok := d.(*ast.MindmapDiagram)
				if !ok {
					t.Fatalf("expected *ast.MindmapDiagram, got %T", d)
				}
				if len(mm.Root.Children) != 2 {
					t.Fatalf("expected 2 children, got %d", len(mm.Root.Children))
				}
				fiction := mm.Root.Children[0]
				if fiction.Icon != "fa fa-book" {
					t.Errorf("expected icon 'fa fa-book', got %q", fiction.Icon)
				}
			},
		},
		{
			name: "plain text nodes",
			source: `mindmap
  Root Node
    Child 1
    Child 2`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				mm, ok := d.(*ast.MindmapDiagram)
				if !ok {
					t.Fatalf("expected *ast.MindmapDiagram, got %T", d)
				}
				if mm.Root.Shape != "" {
					t.Errorf("expected no shape, got %q", mm.Root.Shape)
				}
				if mm.Root.Text != "Root Node" {
					t.Errorf("expected root text 'Root Node', got %q", mm.Root.Text)
				}
			},
		},
		{
			name: "mindmap with comments",
			source: `mindmap
  %% This is the root
  root((Main))
    %% This is a child
    Child 1
    Child 2`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				mm, ok := d.(*ast.MindmapDiagram)
				if !ok {
					t.Fatalf("expected *ast.MindmapDiagram, got %T", d)
				}
				if len(mm.Root.Children) != 2 {
					t.Errorf("expected 2 children, got %d", len(mm.Root.Children))
				}
			},
		},
		{
			name: "4-space indentation",
			source: `mindmap
    root((Main))
        Child 1
        Child 2
            Grandchild`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				mm, ok := d.(*ast.MindmapDiagram)
				if !ok {
					t.Fatalf("expected *ast.MindmapDiagram, got %T", d)
				}
				if len(mm.Root.Children) != 2 {
					t.Errorf("expected 2 children, got %d", len(mm.Root.Children))
				}
				if len(mm.Root.Children[1].Children) != 1 {
					t.Errorf("expected Child 2 to have 1 child, got %d", len(mm.Root.Children[1].Children))
				}
			},
		},
		{
			name:    "invalid header",
			source:  "notmindmap\n",
			wantErr: true,
		},
		{
			name:    "missing root node",
			source:  "mindmap\n",
			wantErr: true,
		},
		{
			name: "multiple root nodes",
			source: `mindmap
  root1
  root2`,
			wantErr: true,
		},
		{
			name: "inconsistent indentation",
			source: `mindmap
  root
    Child 1
   Child 2`,
			wantErr: true,
		},
		{
			name: "child before root",
			source: `mindmap
    Child
  Root`,
			wantErr: true,
		},
		{
			name: "empty node text",
			source: `mindmap
  root(())`,
			wantErr: true,
		},
		{
			name: "icon without node",
			source: `mindmap
  ::icon(fa fa-book)`,
			wantErr: true,
		},
	}

	p := parser.NewMindmapParser()
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

func TestMindmapParser_SupportedTypes(t *testing.T) {
	p := parser.NewMindmapParser()
	types := p.SupportedTypes()
	if len(types) != 1 || types[0] != "mindmap" {
		t.Errorf("expected [mindmap], got %v", types)
	}
}

func TestMindmapParser_ComplexHierarchy(t *testing.T) {
	source := `mindmap
  root((Software Development))
    Planning[Planning Phase]
      Requirements
      Design
    Development[Development Phase]
      Frontend
        ::icon(fa fa-laptop)
        React
        Vue
      Backend
        ::icon(fa fa-server)
        Node.js
        Python
    Testing[Testing Phase]
      Unit Tests
      Integration Tests
    Deployment))Deployment((
      Staging
      Production`

	diagram, err := parser.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	mm, ok := diagram.(*ast.MindmapDiagram)
	if !ok {
		t.Fatalf("expected *ast.MindmapDiagram, got %T", diagram)
	}

	// Check root
	if mm.Root.Text != "Software Development" {
		t.Errorf("expected root text 'Software Development', got %q", mm.Root.Text)
	}

	// Check top-level children
	if len(mm.Root.Children) != 4 {
		t.Fatalf("expected 4 top-level children, got %d", len(mm.Root.Children))
	}

	// Check Planning phase
	planning := mm.Root.Children[0]
	if planning.Text != "Planning Phase" || planning.Shape != "[]" {
		t.Errorf("unexpected Planning node: text=%q, shape=%q", planning.Text, planning.Shape)
	}
	if len(planning.Children) != 2 {
		t.Errorf("expected Planning to have 2 children, got %d", len(planning.Children))
	}

	// Check Development phase with icons
	development := mm.Root.Children[1]
	if len(development.Children) != 2 {
		t.Fatalf("expected Development to have 2 children, got %d", len(development.Children))
	}

	frontend := development.Children[0]
	if frontend.Icon != "fa fa-laptop" {
		t.Errorf("expected Frontend icon 'fa fa-laptop', got %q", frontend.Icon)
	}
	if len(frontend.Children) != 2 {
		t.Errorf("expected Frontend to have 2 children, got %d", len(frontend.Children))
	}

	backend := development.Children[1]
	if backend.Icon != "fa fa-server" {
		t.Errorf("expected Backend icon 'fa fa-server', got %q", backend.Icon)
	}

	// Check Deployment phase
	deployment := mm.Root.Children[3]
	if deployment.Shape != "))((" {
		t.Errorf("expected Deployment shape '))((', got %q", deployment.Shape)
	}
}
