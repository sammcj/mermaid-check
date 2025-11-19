package validator_test

import (	"testing"

	"github.com/sammcj/go-mermaid/ast"

	"github.com/sammcj/go-mermaid/validator"
)

func TestNoDuplicateElementIDsRule(t *testing.T) {
	tests := []struct {
		name      string
		diagram   *ast.C4Diagram
		wantCount int
	}{
		{
			name: "no duplicates",
			diagram: &ast.C4Diagram{
				Elements: []ast.C4Element{
					{ID: "elem1", Pos: ast.Position{Line: 1}},
					{ID: "elem2", Pos: ast.Position{Line: 2}},
				},
			},
			wantCount: 0,
		},
		{
			name: "duplicate element IDs",
			diagram: &ast.C4Diagram{
				Elements: []ast.C4Element{
					{ID: "elem1", Pos: ast.Position{Line: 1}},
					{ID: "elem1", Pos: ast.Position{Line: 2}},
				},
			},
			wantCount: 1,
		},
		{
			name: "duplicate across element and boundary",
			diagram: &ast.C4Diagram{
				Elements: []ast.C4Element{
					{ID: "elem1", Pos: ast.Position{Line: 1}},
				},
				Boundaries: []ast.C4Boundary{
					{ID: "elem1", Pos: ast.Position{Line: 3}},
				},
			},
			wantCount: 1,
		},
		{
			name: "duplicate in nested boundary",
			diagram: &ast.C4Diagram{
				Boundaries: []ast.C4Boundary{
					{
						ID:  "boundary1",
						Pos: ast.Position{Line: 1},
						Elements: []ast.C4Element{
							{ID: "elem1", Pos: ast.Position{Line: 2}},
							{ID: "elem1", Pos: ast.Position{Line: 3}},
						},
					},
				},
			},
			wantCount: 1,
		},
		{
			name: "multiple duplicates",
			diagram: &ast.C4Diagram{
				Elements: []ast.C4Element{
					{ID: "elem1", Pos: ast.Position{Line: 1}},
					{ID: "elem1", Pos: ast.Position{Line: 2}},
					{ID: "elem2", Pos: ast.Position{Line: 3}},
					{ID: "elem2", Pos: ast.Position{Line: 4}},
				},
			},
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &validator.NoDuplicateElementIDsRule{}
			errors := rule.Validate(tt.diagram)

			if len(errors) != tt.wantCount {
				t.Errorf("expected %d errors, got %d", tt.wantCount, len(errors))
			}

			for _, err := range errors {
				if err.Severity != validator.SeverityError {
					t.Errorf("expected severity Error, got %v", err.Severity)
				}
			}
		})
	}
}

func TestC4ValidRelationshipReferencesRule(t *testing.T) {
	tests := []struct {
		name      string
		diagram   *ast.C4Diagram
		wantCount int
	}{
		{
			name: "valid references",
			diagram: &ast.C4Diagram{
				Elements: []ast.C4Element{
					{ID: "elem1"},
					{ID: "elem2"},
				},
				Relationships: []ast.C4Relationship{
					{From: "elem1", To: "elem2", Pos: ast.Position{Line: 3}},
				},
			},
			wantCount: 0,
		},
		{
			name: "undefined from reference",
			diagram: &ast.C4Diagram{
				Elements: []ast.C4Element{
					{ID: "elem1"},
				},
				Relationships: []ast.C4Relationship{
					{From: "undefined", To: "elem1", Pos: ast.Position{Line: 2}},
				},
			},
			wantCount: 1,
		},
		{
			name: "undefined to reference",
			diagram: &ast.C4Diagram{
				Elements: []ast.C4Element{
					{ID: "elem1"},
				},
				Relationships: []ast.C4Relationship{
					{From: "elem1", To: "undefined", Pos: ast.Position{Line: 2}},
				},
			},
			wantCount: 1,
		},
		{
			name: "both references undefined",
			diagram: &ast.C4Diagram{
				Elements: []ast.C4Element{
					{ID: "elem1"},
				},
				Relationships: []ast.C4Relationship{
					{From: "undef1", To: "undef2", Pos: ast.Position{Line: 2}},
				},
			},
			wantCount: 2,
		},
		{
			name: "reference to element in boundary",
			diagram: &ast.C4Diagram{
				Boundaries: []ast.C4Boundary{
					{
						ID: "boundary1",
						Elements: []ast.C4Element{
							{ID: "elem1"},
						},
					},
				},
				Relationships: []ast.C4Relationship{
					{From: "boundary1", To: "elem1", Pos: ast.Position{Line: 3}},
				},
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &validator.C4ValidRelationshipReferencesRule{}
			errors := rule.Validate(tt.diagram)

			if len(errors) != tt.wantCount {
				t.Errorf("expected %d errors, got %d", tt.wantCount, len(errors))
			}
		})
	}
}

func TestValidBoundaryIDsRule(t *testing.T) {
	tests := []struct {
		name      string
		diagram   *ast.C4Diagram
		wantCount int
	}{
		{
			name: "no duplicate boundaries",
			diagram: &ast.C4Diagram{
				Boundaries: []ast.C4Boundary{
					{ID: "b1", Pos: ast.Position{Line: 1}},
					{ID: "b2", Pos: ast.Position{Line: 2}},
				},
			},
			wantCount: 0,
		},
		{
			name: "duplicate boundary IDs",
			diagram: &ast.C4Diagram{
				Boundaries: []ast.C4Boundary{
					{ID: "b1", Pos: ast.Position{Line: 1}},
					{ID: "b1", Pos: ast.Position{Line: 2}},
				},
			},
			wantCount: 1,
		},
		{
			name: "duplicate in nested boundaries",
			diagram: &ast.C4Diagram{
				Boundaries: []ast.C4Boundary{
					{
						ID:  "b1",
						Pos: ast.Position{Line: 1},
						Boundaries: []ast.C4Boundary{
							{ID: "b2", Pos: ast.Position{Line: 2}},
							{ID: "b2", Pos: ast.Position{Line: 3}},
						},
					},
				},
			},
			wantCount: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &validator.ValidBoundaryIDsRule{}
			errors := rule.Validate(tt.diagram)

			if len(errors) != tt.wantCount {
				t.Errorf("expected %d errors, got %d", tt.wantCount, len(errors))
			}
		})
	}
}

func TestValidStyleReferencesRule(t *testing.T) {
	tests := []struct {
		name      string
		diagram   *ast.C4Diagram
		wantCount int
	}{
		{
			name: "valid element style reference",
			diagram: &ast.C4Diagram{
				Elements: []ast.C4Element{
					{ID: "elem1"},
				},
				Styles: []ast.C4Style{
					{StyleType: "UpdateElementStyle", ElementID: "elem1", Pos: ast.Position{Line: 2}},
				},
			},
			wantCount: 0,
		},
		{
			name: "invalid element style reference",
			diagram: &ast.C4Diagram{
				Elements: []ast.C4Element{
					{ID: "elem1"},
				},
				Styles: []ast.C4Style{
					{StyleType: "UpdateElementStyle", ElementID: "undefined", Pos: ast.Position{Line: 2}},
				},
			},
			wantCount: 1,
		},
		{
			name: "valid relationship style reference",
			diagram: &ast.C4Diagram{
				Elements: []ast.C4Element{
					{ID: "elem1"},
					{ID: "elem2"},
				},
				Styles: []ast.C4Style{
					{StyleType: "UpdateRelStyle", From: "elem1", To: "elem2", Pos: ast.Position{Line: 3}},
				},
			},
			wantCount: 0,
		},
		{
			name: "invalid relationship style from reference",
			diagram: &ast.C4Diagram{
				Elements: []ast.C4Element{
					{ID: "elem1"},
				},
				Styles: []ast.C4Style{
					{StyleType: "UpdateRelStyle", From: "undefined", To: "elem1", Pos: ast.Position{Line: 2}},
				},
			},
			wantCount: 1,
		},
		{
			name: "invalid relationship style to reference",
			diagram: &ast.C4Diagram{
				Elements: []ast.C4Element{
					{ID: "elem1"},
				},
				Styles: []ast.C4Style{
					{StyleType: "UpdateRelStyle", From: "elem1", To: "undefined", Pos: ast.Position{Line: 2}},
				},
			},
			wantCount: 1,
		},
		{
			name: "reference to element in boundary",
			diagram: &ast.C4Diagram{
				Boundaries: []ast.C4Boundary{
					{
						ID: "boundary1",
						Elements: []ast.C4Element{
							{ID: "elem1"},
						},
					},
				},
				Styles: []ast.C4Style{
					{StyleType: "UpdateElementStyle", ElementID: "elem1", Pos: ast.Position{Line: 3}},
				},
			},
			wantCount: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &validator.ValidStyleReferencesRule{}
			errors := rule.Validate(tt.diagram)

			if len(errors) != tt.wantCount {
				t.Errorf("expected %d errors, got %d", tt.wantCount, len(errors))
				for _, err := range errors {
					t.Logf("  error: %s", err.Message)
				}
			}
		})
	}
}

func TestValidateC4(t *testing.T) {
	tests := []struct {
		name      string
		diagram   *ast.C4Diagram
		rules     []validator.C4Rule
		wantCount int
	}{
		{
			name: "valid diagram with default rules",
			diagram: &ast.C4Diagram{
				Elements: []ast.C4Element{
					{ID: "elem1"},
					{ID: "elem2"},
				},
				Relationships: []ast.C4Relationship{
					{From: "elem1", To: "elem2"},
				},
			},
			rules:     validator.DefaultC4Rules(),
			wantCount: 0,
		},
		{
			name: "invalid diagram with multiple issues",
			diagram: &ast.C4Diagram{
				Elements: []ast.C4Element{
					{ID: "elem1", Pos: ast.Position{Line: 1}},
					{ID: "elem1", Pos: ast.Position{Line: 2}},
				},
				Relationships: []ast.C4Relationship{
					{From: "elem1", To: "undefined", Pos: ast.Position{Line: 3}},
				},
			},
			rules:     validator.DefaultC4Rules(),
			wantCount: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.ValidateC4(tt.diagram, tt.rules)

			if len(errors) != tt.wantCount {
				t.Errorf("expected %d errors, got %d", tt.wantCount, len(errors))
				for _, err := range errors {
					t.Logf("  error at line %d: %s", err.Line, err.Message)
				}
			}
		})
	}
}

func TestDefaultC4Rules(t *testing.T) {
	rules := validator.DefaultC4Rules()
	if len(rules) != 4 {
		t.Errorf("expected 4 default rules, got %d", len(rules))
	}
}

func TestStrictC4Rules(t *testing.T) {
	rules := validator.StrictC4Rules()
	if len(rules) != 4 {
		t.Errorf("expected 4 strict rules, got %d", len(rules))
	}
}
