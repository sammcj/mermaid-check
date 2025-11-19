package validator_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/validator"
)

func TestNoDuplicateLabelsRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.PieDiagram
		wantErr bool
	}{
		{
			name: "no duplicates",
			diagram: &ast.PieDiagram{
				DataEntries: []ast.PieEntry{
					{Label: "A", Value: 10, Pos: ast.Position{Line: 2, Column: 1}},
					{Label: "B", Value: 20, Pos: ast.Position{Line: 3, Column: 1}},
					{Label: "C", Value: 30, Pos: ast.Position{Line: 4, Column: 1}},
				},
			},
			wantErr: false,
		},
		{
			name: "with duplicates",
			diagram: &ast.PieDiagram{
				DataEntries: []ast.PieEntry{
					{Label: "A", Value: 10, Pos: ast.Position{Line: 2, Column: 1}},
					{Label: "B", Value: 20, Pos: ast.Position{Line: 3, Column: 1}},
					{Label: "A", Value: 30, Pos: ast.Position{Line: 4, Column: 1}},
				},
			},
			wantErr: true,
		},
	}

	rule := &validator.NoDuplicateLabelsRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("NoDuplicateLabelsRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestPositiveValuesRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.PieDiagram
		wantErr bool
	}{
		{
			name: "all positive values",
			diagram: &ast.PieDiagram{
				DataEntries: []ast.PieEntry{
					{Label: "A", Value: 10.5, Pos: ast.Position{Line: 2, Column: 1}},
					{Label: "B", Value: 20.25, Pos: ast.Position{Line: 3, Column: 1}},
					{Label: "C", Value: 30, Pos: ast.Position{Line: 4, Column: 1}},
				},
			},
			wantErr: false,
		},
		{
			name: "with zero value",
			diagram: &ast.PieDiagram{
				DataEntries: []ast.PieEntry{
					{Label: "A", Value: 0, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErr: true,
		},
		{
			name: "with negative value",
			diagram: &ast.PieDiagram{
				DataEntries: []ast.PieEntry{
					{Label: "A", Value: -10, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErr: true,
		},
	}

	rule := &validator.PositiveValuesRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("PositiveValuesRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestPieDefaultRules(t *testing.T) {
	rules := validator.PieDefaultRules()
	if len(rules) == 0 {
		t.Error("PieDefaultRules() returned empty slice")
	}
}

func TestPieStrictRules(t *testing.T) {
	rules := validator.PieStrictRules()
	if len(rules) == 0 {
		t.Error("PieStrictRules() returned empty slice")
	}
}
