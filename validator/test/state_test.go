package validator_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/validator"
)

func TestNoDuplicateStates(t *testing.T) {
	tests := []struct {
		name       string
		diagram    *ast.StateDiagram
		wantErrors int
	}{
		{
			name: "no duplicates",
			diagram: &ast.StateDiagram{
				Type: "state",
				Statements: []ast.StateStmt{
					&ast.State{ID: "Still", Pos: ast.Position{Line: 2, Column: 1}},
					&ast.State{ID: "Moving", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrors: 0,
		},
		{
			name: "with duplicates",
			diagram: &ast.StateDiagram{
				Type: "state",
				Statements: []ast.StateStmt{
					&ast.State{ID: "Still", Pos: ast.Position{Line: 2, Column: 1}},
					&ast.State{ID: "Still", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrors: 1,
		},
	}

	rule := &validator.NoDuplicateStates{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.ValidateState(tt.diagram)
			if len(errors) != tt.wantErrors {
				t.Errorf("ValidateState() errors = %d, want %d", len(errors), tt.wantErrors)
			}
		})
	}
}

func TestStateDefaultRules(t *testing.T) {
	rules := validator.StateDefaultRules()
	if len(rules) == 0 {
		t.Error("StateDefaultRules() returned empty rules")
	}
}

func TestStateStrictRules(t *testing.T) {
	rules := validator.StateStrictRules()
	if len(rules) == 0 {
		t.Error("StateStrictRules() returned empty rules")
	}
}
