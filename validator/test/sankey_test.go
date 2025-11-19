package validator_test

import (	"testing"

	"github.com/sammcj/go-mermaid/ast"

	"github.com/sammcj/go-mermaid/validator"
)

func TestSankeyPositiveValuesRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.SankeyDiagram
		wantErr bool
	}{
		{
			name: "all positive values",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "A", Target: "B", Value: 10.5, Pos: ast.Position{Line: 2, Column: 1}},
					{Source: "B", Target: "C", Value: 20.25, Pos: ast.Position{Line: 3, Column: 1}},
					{Source: "A", Target: "C", Value: 5.0, Pos: ast.Position{Line: 4, Column: 1}},
				},
			},
			wantErr: false,
		},
		{
			name: "with zero value",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "A", Target: "B", Value: 0, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErr: true,
		},
		{
			name: "with negative value",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "A", Target: "B", Value: -10, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErr: true,
		},
		{
			name: "mixed positive and negative",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "A", Target: "B", Value: 10, Pos: ast.Position{Line: 2, Column: 1}},
					{Source: "B", Target: "C", Value: -5, Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErr: true,
		},
	}

	rule := &validator.SankeyPositiveValuesRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("SankeyPositiveValuesRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestSankeyNoSelfLoopsRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.SankeyDiagram
		wantErr bool
	}{
		{
			name: "no self-loops",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "A", Target: "B", Value: 10, Pos: ast.Position{Line: 2, Column: 1}},
					{Source: "B", Target: "C", Value: 20, Pos: ast.Position{Line: 3, Column: 1}},
					{Source: "A", Target: "C", Value: 5, Pos: ast.Position{Line: 4, Column: 1}},
				},
			},
			wantErr: false,
		},
		{
			name: "with self-loop",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "A", Target: "A", Value: 10, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErr: true,
		},
		{
			name: "multiple self-loops",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "A", Target: "B", Value: 10, Pos: ast.Position{Line: 2, Column: 1}},
					{Source: "B", Target: "B", Value: 5, Pos: ast.Position{Line: 3, Column: 1}},
					{Source: "C", Target: "C", Value: 3, Pos: ast.Position{Line: 4, Column: 1}},
				},
			},
			wantErr: true,
		},
	}

	rule := &validator.SankeyNoSelfLoopsRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("SankeyNoSelfLoopsRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestSankeyValidNodeReferencesRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.SankeyDiagram
		wantErr bool
	}{
		{
			name: "all valid node references",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "A", Target: "B", Value: 10, Pos: ast.Position{Line: 2, Column: 1}},
					{Source: "B", Target: "C", Value: 20, Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErr: false,
		},
		{
			name: "empty source node",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "", Target: "B", Value: 10, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErr: true,
		},
		{
			name: "empty target node",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "A", Target: "", Value: 10, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErr: true,
		},
		{
			name: "both source and target empty",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "", Target: "", Value: 10, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErr: true,
		},
	}

	rule := &validator.SankeyValidNodeReferencesRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("SankeyValidNodeReferencesRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestSankeyMinimumLinksRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.SankeyDiagram
		wantErr bool
	}{
		{
			name: "has links",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "A", Target: "B", Value: 10, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple links",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "A", Target: "B", Value: 10, Pos: ast.Position{Line: 2, Column: 1}},
					{Source: "B", Target: "C", Value: 20, Pos: ast.Position{Line: 3, Column: 1}},
					{Source: "A", Target: "C", Value: 5, Pos: ast.Position{Line: 4, Column: 1}},
				},
			},
			wantErr: false,
		},
		{
			name: "no links",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{},
			},
			wantErr: true,
		},
	}

	rule := &validator.SankeyMinimumLinksRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("SankeyMinimumLinksRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestSankeyDefaultRules(t *testing.T) {
	rules := validator.SankeyDefaultRules()
	if len(rules) == 0 {
		t.Error("validator.SankeyDefaultRules() returned empty slice")
	}
	expectedRuleCount := 4
	if len(rules) != expectedRuleCount {
		t.Errorf("expected %d rules, got %d", expectedRuleCount, len(rules))
	}
}

func TestSankeyStrictRules(t *testing.T) {
	rules := validator.SankeyStrictRules()
	if len(rules) == 0 {
		t.Error("validator.SankeyStrictRules() returned empty slice")
	}
}

func TestValidateSankey(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.SankeyDiagram
		strict  bool
		wantErr bool
	}{
		{
			name: "valid sankey diagram",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "A", Target: "B", Value: 10, Pos: ast.Position{Line: 2, Column: 1}},
					{Source: "B", Target: "C", Value: 20, Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			strict:  false,
			wantErr: false,
		},
		{
			name: "invalid sankey - self-loop",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "A", Target: "A", Value: 10, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			strict:  false,
			wantErr: true,
		},
		{
			name: "invalid sankey - negative value",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "A", Target: "B", Value: -10, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			strict:  false,
			wantErr: true,
		},
		{
			name: "invalid sankey - empty links",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{},
			},
			strict:  false,
			wantErr: true,
		},
		{
			name: "valid sankey diagram strict mode",
			diagram: &ast.SankeyDiagram{
				Links: []ast.SankeyLink{
					{Source: "A", Target: "B", Value: 10, Pos: ast.Position{Line: 2, Column: 1}},
					{Source: "B", Target: "C", Value: 20, Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			strict:  true,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.ValidateSankey(tt.diagram, tt.strict)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("validator.ValidateSankey() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}
