package validator_test

import (	"testing"

	"github.com/sammcj/go-mermaid/ast"

	"github.com/sammcj/go-mermaid/validator"
)

func TestValidateMindmap(t *testing.T) {
	tests := []struct {
		name     string
		diagram  *ast.MindmapDiagram
		strict   bool
		wantErrs int
	}{
		{
			name: "valid mindmap",
			diagram: &ast.MindmapDiagram{
				Type: "mindmap",
				Root: &ast.MindmapNode{
					Text:  "Root",
					Shape: "(())",
					Level: 0,
					Children: []*ast.MindmapNode{
						{Text: "Child 1", Shape: "", Level: 1},
						{Text: "Child 2", Shape: "[]", Level: 1},
					},
				},
			},
			strict:   false,
			wantErrs: 0,
		},
		{
			name: "missing root node",
			diagram: &ast.MindmapDiagram{
				Type: "mindmap",
				Root: nil,
			},
			strict:   false,
			wantErrs: 1,
		},
		{
			name: "empty node text",
			diagram: &ast.MindmapDiagram{
				Type: "mindmap",
				Root: &ast.MindmapNode{
					Text:  "",
					Shape: "(())",
					Level: 0,
					Pos:   ast.Position{Line: 2, Column: 1},
				},
			},
			strict:   false,
			wantErrs: 1,
		},
		{
			name: "invalid shape",
			diagram: &ast.MindmapDiagram{
				Type: "mindmap",
				Root: &ast.MindmapNode{
					Text:  "Root",
					Shape: "<>",
					Level: 0,
					Pos:   ast.Position{Line: 2, Column: 1},
				},
			},
			strict:   false,
			wantErrs: 1,
		},
		{
			name: "multiple validation errors",
			diagram: &ast.MindmapDiagram{
				Type: "mindmap",
				Root: &ast.MindmapNode{
					Text:  "Root",
					Shape: "(())",
					Level: 0,
					Children: []*ast.MindmapNode{
						{Text: "", Shape: "[]", Level: 1, Pos: ast.Position{Line: 3, Column: 1}},
						{Text: "Valid", Shape: "invalid", Level: 1, Pos: ast.Position{Line: 4, Column: 1}},
					},
				},
			},
			strict:   false,
			wantErrs: 2,
		},
		{
			name: "valid shapes",
			diagram: &ast.MindmapDiagram{
				Type: "mindmap",
				Root: &ast.MindmapNode{
					Text:  "Root",
					Shape: "(())",
					Level: 0,
					Children: []*ast.MindmapNode{
						{Text: "Round", Shape: "()", Level: 1},
						{Text: "Square", Shape: "[]", Level: 1},
						{Text: "Cloud", Shape: "{{}}", Level: 1},
						{Text: "Hexagon", Shape: "))((", Level: 1},
						{Text: "Plain", Shape: "", Level: 1},
					},
				},
			},
			strict:   false,
			wantErrs: 0,
		},
		{
			name: "nested empty nodes",
			diagram: &ast.MindmapDiagram{
				Type: "mindmap",
				Root: &ast.MindmapNode{
					Text:  "Root",
					Shape: "",
					Level: 0,
					Children: []*ast.MindmapNode{
						{
							Text:  "Parent",
							Shape: "",
							Level: 1,
							Children: []*ast.MindmapNode{
								{Text: "", Shape: "", Level: 2, Pos: ast.Position{Line: 5, Column: 1}},
							},
						},
					},
				},
			},
			strict:   false,
			wantErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.ValidateMindmap(tt.diagram, tt.strict)
			if len(errors) != tt.wantErrs {
				t.Errorf("validator.ValidateMindmap() returned %d errors, want %d", len(errors), tt.wantErrs)
				for _, err := range errors {
					t.Logf("  Error: %s (line %d)", err.Message, err.Line)
				}
			}
		})
	}
}

func TestRootNodeExistsRule(t *testing.T) {
	rule := &validator.RootNodeExistsRule{}

	tests := []struct {
		name     string
		diagram  *ast.MindmapDiagram
		wantErrs int
	}{
		{
			name: "root exists",
			diagram: &ast.MindmapDiagram{
				Root: &ast.MindmapNode{Text: "Root"},
			},
			wantErrs: 0,
		},
		{
			name: "root missing",
			diagram: &ast.MindmapDiagram{
				Root: nil,
			},
			wantErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if len(errors) != tt.wantErrs {
				t.Errorf("Validate() returned %d errors, want %d", len(errors), tt.wantErrs)
			}
		})
	}
}

func TestNoEmptyNodesRule(t *testing.T) {
	rule := &validator.NoEmptyNodesRule{}

	tests := []struct {
		name     string
		diagram  *ast.MindmapDiagram
		wantErrs int
	}{
		{
			name: "all nodes have text",
			diagram: &ast.MindmapDiagram{
				Root: &ast.MindmapNode{
					Text: "Root",
					Children: []*ast.MindmapNode{
						{Text: "Child 1"},
						{Text: "Child 2"},
					},
				},
			},
			wantErrs: 0,
		},
		{
			name: "root has empty text",
			diagram: &ast.MindmapDiagram{
				Root: &ast.MindmapNode{
					Text: "",
					Pos:  ast.Position{Line: 2, Column: 1},
				},
			},
			wantErrs: 1,
		},
		{
			name: "child has empty text",
			diagram: &ast.MindmapDiagram{
				Root: &ast.MindmapNode{
					Text: "Root",
					Children: []*ast.MindmapNode{
						{Text: "Child 1"},
						{Text: "", Pos: ast.Position{Line: 4, Column: 1}},
					},
				},
			},
			wantErrs: 1,
		},
		{
			name: "multiple empty nodes",
			diagram: &ast.MindmapDiagram{
				Root: &ast.MindmapNode{
					Text: "Root",
					Children: []*ast.MindmapNode{
						{Text: "", Pos: ast.Position{Line: 3, Column: 1}},
						{
							Text: "Parent",
							Children: []*ast.MindmapNode{
								{Text: "", Pos: ast.Position{Line: 5, Column: 1}},
							},
						},
					},
				},
			},
			wantErrs: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if len(errors) != tt.wantErrs {
				t.Errorf("Validate() returned %d errors, want %d", len(errors), tt.wantErrs)
			}
		})
	}
}

func TestValidShapeRule(t *testing.T) {
	rule := &validator.ValidShapeRule{}

	tests := []struct {
		name     string
		diagram  *ast.MindmapDiagram
		wantErrs int
	}{
		{
			name: "all valid shapes",
			diagram: &ast.MindmapDiagram{
				Root: &ast.MindmapNode{
					Text:  "Root",
					Shape: "(())",
					Children: []*ast.MindmapNode{
						{Text: "Round", Shape: "()"},
						{Text: "Square", Shape: "[]"},
						{Text: "Cloud", Shape: "{{}}"},
						{Text: "Hexagon", Shape: "))(("},
						{Text: "Plain", Shape: ""},
					},
				},
			},
			wantErrs: 0,
		},
		{
			name: "invalid shape",
			diagram: &ast.MindmapDiagram{
				Root: &ast.MindmapNode{
					Text:  "Root",
					Shape: "<>",
					Pos:   ast.Position{Line: 2, Column: 1},
				},
			},
			wantErrs: 1,
		},
		{
			name: "multiple invalid shapes",
			diagram: &ast.MindmapDiagram{
				Root: &ast.MindmapNode{
					Text:  "Root",
					Shape: "(())",
					Children: []*ast.MindmapNode{
						{Text: "Invalid1", Shape: "<>", Pos: ast.Position{Line: 3, Column: 1}},
						{Text: "Valid", Shape: "[]"},
						{Text: "Invalid2", Shape: "//", Pos: ast.Position{Line: 5, Column: 1}},
					},
				},
			},
			wantErrs: 2,
		},
		{
			name: "nested invalid shape",
			diagram: &ast.MindmapDiagram{
				Root: &ast.MindmapNode{
					Text:  "Root",
					Shape: "(())",
					Children: []*ast.MindmapNode{
						{
							Text:  "Parent",
							Shape: "[]",
							Children: []*ast.MindmapNode{
								{Text: "Child", Shape: "***", Pos: ast.Position{Line: 5, Column: 1}},
							},
						},
					},
				},
			},
			wantErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if len(errors) != tt.wantErrs {
				t.Errorf("Validate() returned %d errors, want %d", len(errors), tt.wantErrs)
				for _, err := range errors {
					t.Logf("  Error: %s", err.Message)
				}
			}
		})
	}
}

func TestMindmapDefaultRules(t *testing.T) {
	rules := validator.MindmapDefaultRules()
	if len(rules) != 3 {
		t.Errorf("expected 3 default rules, got %d", len(rules))
	}
}

func TestMindmapStrictRules(t *testing.T) {
	rules := validator.MindmapStrictRules()
	if len(rules) < 3 {
		t.Errorf("expected at least 3 strict rules, got %d", len(rules))
	}
}
