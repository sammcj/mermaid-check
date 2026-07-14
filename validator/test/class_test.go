package validator_test

import (
	"testing"

	"github.com/sammcj/mermaid-check/ast"
	"github.com/sammcj/mermaid-check/validator"
)

func TestNoDuplicateClasses(t *testing.T) {
	tests := []struct {
		name       string
		diagram    *ast.ClassDiagram
		wantErrors int
	}{
		{
			name: "no duplicates",
			diagram: &ast.ClassDiagram{
				Type: "class",
				Statements: []ast.ClassStmt{
					&ast.Class{Name: "Animal", Pos: ast.Position{Line: 2, Column: 1}},
					&ast.Class{Name: "Dog", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrors: 0,
		},
		{
			name: "with duplicates",
			diagram: &ast.ClassDiagram{
				Type: "class",
				Statements: []ast.ClassStmt{
					&ast.Class{Name: "Animal", Pos: ast.Position{Line: 2, Column: 1}},
					&ast.Class{Name: "Animal", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrors: 1,
		},
	}

	rule := &validator.NoDuplicateClasses{}

	if rule.Name() != "no-duplicate-classes" {
		t.Errorf("Name() = %q, want %q", rule.Name(), "no-duplicate-classes")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.ValidateClass(tt.diagram)
			if len(errors) != tt.wantErrors {
				t.Errorf("ValidateClass() errors = %d, want %d", len(errors), tt.wantErrors)
			}
		})
	}
}

func TestValidClassReferences(t *testing.T) {
	tests := []struct {
		name       string
		diagram    *ast.ClassDiagram
		wantErrors int
	}{
		{
			name: "valid note reference",
			diagram: &ast.ClassDiagram{
				Type: "class",
				Statements: []ast.ClassStmt{
					&ast.Class{Name: "Animal", Pos: ast.Position{Line: 2, Column: 1}},
					&ast.ClassNote{ClassName: "Animal", Text: "Note text", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrors: 0,
		},
		{
			name: "note references undefined class",
			diagram: &ast.ClassDiagram{
				Type: "class",
				Statements: []ast.ClassStmt{
					&ast.ClassNote{ClassName: "UndefinedClass", Text: "Note text", Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErrors: 1,
		},
		{
			name: "standalone note needs no class reference",
			diagram: &ast.ClassDiagram{
				Type: "class",
				Statements: []ast.ClassStmt{
					&ast.ClassNote{ClassName: "", Text: "a floating note", Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErrors: 0,
		},
		{
			name: "targeted note still errors on undefined class alongside a standalone note",
			diagram: &ast.ClassDiagram{
				Type: "class",
				Statements: []ast.ClassStmt{
					&ast.ClassNote{ClassName: "", Text: "a floating note", Pos: ast.Position{Line: 2, Column: 1}},
					&ast.ClassNote{ClassName: "UndefinedClass", Text: "Note text", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrors: 1,
		},
		{
			name: "implicit class from relationship",
			diagram: &ast.ClassDiagram{
				Type: "class",
				Statements: []ast.ClassStmt{
					&ast.Relationship{From: "A", To: "B", Type: "inheritance", Pos: ast.Position{Line: 2, Column: 1}},
					&ast.ClassNote{ClassName: "A", Text: "Note", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrors: 0,
		},
	}

	rule := &validator.ValidClassReferences{}

	if rule.Name() != "valid-class-references" {
		t.Errorf("Name() = %q, want %q", rule.Name(), "valid-class-references")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.ValidateClass(tt.diagram)
			if len(errors) != tt.wantErrors {
				t.Errorf("ValidateClass() errors = %d, want %d", len(errors), tt.wantErrors)
			}
		})
	}
}

func TestValidMemberVisibility(t *testing.T) {
	tests := []struct {
		name       string
		diagram    *ast.ClassDiagram
		wantErrors int
	}{
		{
			name: "valid visibility modifiers",
			diagram: &ast.ClassDiagram{
				Type: "class",
				Statements: []ast.ClassStmt{
					&ast.Class{
						Name: "Animal",
						Members: []ast.ClassMember{
							{Visibility: "+", Name: "name", Pos: ast.Position{Line: 3, Column: 5}},
							{Visibility: "-", Name: "age", Pos: ast.Position{Line: 4, Column: 5}},
							{Visibility: "#", Name: "weight", Pos: ast.Position{Line: 5, Column: 5}},
							{Visibility: "~", Name: "height", Pos: ast.Position{Line: 6, Column: 5}},
						},
						Pos: ast.Position{Line: 2, Column: 1},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "invalid visibility modifier",
			diagram: &ast.ClassDiagram{
				Type: "class",
				Statements: []ast.ClassStmt{
					&ast.Class{
						Name: "Animal",
						Members: []ast.ClassMember{
							{Visibility: "*", Name: "name", Pos: ast.Position{Line: 3, Column: 5}},
						},
						Pos: ast.Position{Line: 2, Column: 1},
					},
				},
			},
			wantErrors: 1,
		},
	}

	rule := &validator.ValidMemberVisibility{}

	if rule.Name() != "valid-member-visibility" {
		t.Errorf("Name() = %q, want %q", rule.Name(), "valid-member-visibility")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.ValidateClass(tt.diagram)
			if len(errors) != tt.wantErrors {
				t.Errorf("ValidateClass() errors = %d, want %d", len(errors), tt.wantErrors)
			}
		})
	}
}

func TestValidRelationshipType(t *testing.T) {
	tests := []struct {
		name       string
		diagram    *ast.ClassDiagram
		wantErrors int
	}{
		{
			name: "valid relationship types",
			diagram: &ast.ClassDiagram{
				Type: "class",
				Statements: []ast.ClassStmt{
					&ast.Relationship{From: "A", To: "B", Type: "inheritance", Pos: ast.Position{Line: 2, Column: 1}},
					&ast.Relationship{From: "B", To: "C", Type: "composition", Pos: ast.Position{Line: 3, Column: 1}},
					&ast.Relationship{From: "C", To: "D", Type: "aggregation", Pos: ast.Position{Line: 4, Column: 1}},
					&ast.Relationship{From: "D", To: "E", Type: "association", Pos: ast.Position{Line: 5, Column: 1}},
					&ast.Relationship{From: "E", To: "F", Type: "dependency", Pos: ast.Position{Line: 6, Column: 1}},
					&ast.Relationship{From: "F", To: "G", Type: "realization", Pos: ast.Position{Line: 7, Column: 1}},
				},
			},
			wantErrors: 0,
		},
		{
			name: "invalid relationship type",
			diagram: &ast.ClassDiagram{
				Type: "class",
				Statements: []ast.ClassStmt{
					&ast.Relationship{From: "A", To: "B", Type: "invalid", Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErrors: 1,
		},
	}

	rule := &validator.ValidRelationshipType{}

	if rule.Name() != "valid-relationship-type" {
		t.Errorf("Name() = %q, want %q", rule.Name(), "valid-relationship-type")
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.ValidateClass(tt.diagram)
			if len(errors) != tt.wantErrors {
				t.Errorf("ValidateClass() errors = %d, want %d", len(errors), tt.wantErrors)
			}
		})
	}
}

func TestClassDefaultRules(t *testing.T) {
	rules := validator.ClassDefaultRules()
	if len(rules) != 4 {
		t.Errorf("ClassDefaultRules() returned %d rules, want 4", len(rules))
	}
}

func TestClassStrictRules(t *testing.T) {
	rules := validator.ClassStrictRules()
	if len(rules) != 4 {
		t.Errorf("ClassStrictRules() returned %d rules, want 4", len(rules))
	}
}

func TestNewClass(t *testing.T) {
	rule := &validator.NoDuplicateClasses{}
	v := validator.NewClass(rule)
	if v == nil {
		t.Error("NewClass() returned nil")
	}
}
