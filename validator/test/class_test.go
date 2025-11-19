package validator_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/validator"
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
	if len(rules) == 0 {
		t.Error("ClassDefaultRules() returned empty rules")
	}
}

func TestClassStrictRules(t *testing.T) {
	rules := validator.ClassStrictRules()
	if len(rules) == 0 {
		t.Error("ClassStrictRules() returned empty rules")
	}
}
