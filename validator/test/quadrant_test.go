package validator_test

import (	"testing"

	"github.com/sammcj/go-mermaid/ast"

	"github.com/sammcj/go-mermaid/validator"
)

func TestValidCoordinatesRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.QuadrantDiagram
		wantErr bool
	}{
		{
			name: "all valid coordinates",
			diagram: &ast.QuadrantDiagram{
				Points: []ast.QuadrantPoint{
					{Name: "A", X: 0.0, Y: 0.0, Pos: ast.Position{Line: 2, Column: 1}},
					{Name: "B", X: 0.5, Y: 0.5, Pos: ast.Position{Line: 3, Column: 1}},
					{Name: "C", X: 1.0, Y: 1.0, Pos: ast.Position{Line: 4, Column: 1}},
				},
			},
			wantErr: false,
		},
		{
			name: "X coordinate too low",
			diagram: &ast.QuadrantDiagram{
				Points: []ast.QuadrantPoint{
					{Name: "A", X: -0.1, Y: 0.5, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErr: true,
		},
		{
			name: "X coordinate too high",
			diagram: &ast.QuadrantDiagram{
				Points: []ast.QuadrantPoint{
					{Name: "A", X: 1.1, Y: 0.5, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErr: true,
		},
		{
			name: "Y coordinate too low",
			diagram: &ast.QuadrantDiagram{
				Points: []ast.QuadrantPoint{
					{Name: "A", X: 0.5, Y: -0.1, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErr: true,
		},
		{
			name: "Y coordinate too high",
			diagram: &ast.QuadrantDiagram{
				Points: []ast.QuadrantPoint{
					{Name: "A", X: 0.5, Y: 1.1, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErr: true,
		},
		{
			name: "both coordinates invalid",
			diagram: &ast.QuadrantDiagram{
				Points: []ast.QuadrantPoint{
					{Name: "A", X: -0.5, Y: 1.5, Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErr: true,
		},
	}

	rule := &validator.ValidCoordinatesRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("ValidCoordinatesRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestNoDuplicatePointNamesRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.QuadrantDiagram
		wantErr bool
	}{
		{
			name: "no duplicates",
			diagram: &ast.QuadrantDiagram{
				Points: []ast.QuadrantPoint{
					{Name: "Campaign A", X: 0.3, Y: 0.6, Pos: ast.Position{Line: 2, Column: 1}},
					{Name: "Campaign B", X: 0.5, Y: 0.4, Pos: ast.Position{Line: 3, Column: 1}},
					{Name: "Campaign C", X: 0.7, Y: 0.8, Pos: ast.Position{Line: 4, Column: 1}},
				},
			},
			wantErr: false,
		},
		{
			name: "with duplicates",
			diagram: &ast.QuadrantDiagram{
				Points: []ast.QuadrantPoint{
					{Name: "Campaign A", X: 0.3, Y: 0.6, Pos: ast.Position{Line: 2, Column: 1}},
					{Name: "Campaign B", X: 0.5, Y: 0.4, Pos: ast.Position{Line: 3, Column: 1}},
					{Name: "Campaign A", X: 0.7, Y: 0.8, Pos: ast.Position{Line: 4, Column: 1}},
				},
			},
			wantErr: true,
		},
	}

	rule := &validator.NoDuplicatePointNamesRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("NoDuplicatePointNamesRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestQuadrantXAxisDefinedRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.QuadrantDiagram
		wantErr bool
	}{
		{
			name: "x-axis defined",
			diagram: &ast.QuadrantDiagram{
				XAxis: ast.QuadrantAxis{Min: "Low", Max: "High"},
				Pos:   ast.Position{Line: 1, Column: 1},
			},
			wantErr: false,
		},
		{
			name: "x-axis not defined",
			diagram: &ast.QuadrantDiagram{
				XAxis: ast.QuadrantAxis{},
				Pos:   ast.Position{Line: 1, Column: 1},
			},
			wantErr: true,
		},
	}

	rule := &validator.QuadrantXAxisDefinedRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("QuadrantXAxisDefinedRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestQuadrantYAxisDefinedRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.QuadrantDiagram
		wantErr bool
	}{
		{
			name: "y-axis defined",
			diagram: &ast.QuadrantDiagram{
				YAxis: ast.QuadrantAxis{Min: "Bottom", Max: "Top"},
				Pos:   ast.Position{Line: 1, Column: 1},
			},
			wantErr: false,
		},
		{
			name: "y-axis not defined",
			diagram: &ast.QuadrantDiagram{
				YAxis: ast.QuadrantAxis{},
				Pos:   ast.Position{Line: 1, Column: 1},
			},
			wantErr: true,
		},
	}

	rule := &validator.QuadrantYAxisDefinedRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("QuadrantYAxisDefinedRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestMinimumPointsRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.QuadrantDiagram
		wantErr bool
	}{
		{
			name: "has points",
			diagram: &ast.QuadrantDiagram{
				Points: []ast.QuadrantPoint{
					{Name: "A", X: 0.5, Y: 0.5, Pos: ast.Position{Line: 2, Column: 1}},
				},
				Pos: ast.Position{Line: 1, Column: 1},
			},
			wantErr: false,
		},
		{
			name: "no points",
			diagram: &ast.QuadrantDiagram{
				Points: []ast.QuadrantPoint{},
				Pos:    ast.Position{Line: 1, Column: 1},
			},
			wantErr: true,
		},
	}

	rule := &validator.MinimumPointsRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("MinimumPointsRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestValidateQuadrant(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.QuadrantDiagram
		strict  bool
		wantErr bool
	}{
		{
			name: "valid quadrant diagram",
			diagram: &ast.QuadrantDiagram{
				XAxis: ast.QuadrantAxis{Min: "Low", Max: "High"},
				YAxis: ast.QuadrantAxis{Min: "Bottom", Max: "Top"},
				Points: []ast.QuadrantPoint{
					{Name: "A", X: 0.3, Y: 0.6, Pos: ast.Position{Line: 2, Column: 1}},
					{Name: "B", X: 0.7, Y: 0.4, Pos: ast.Position{Line: 3, Column: 1}},
				},
				Pos: ast.Position{Line: 1, Column: 1},
			},
			strict:  false,
			wantErr: false,
		},
		{
			name: "invalid coordinates",
			diagram: &ast.QuadrantDiagram{
				XAxis: ast.QuadrantAxis{Min: "Low", Max: "High"},
				YAxis: ast.QuadrantAxis{Min: "Bottom", Max: "Top"},
				Points: []ast.QuadrantPoint{
					{Name: "A", X: 1.5, Y: 0.5, Pos: ast.Position{Line: 2, Column: 1}},
				},
				Pos: ast.Position{Line: 1, Column: 1},
			},
			strict:  false,
			wantErr: true,
		},
		{
			name: "missing x-axis",
			diagram: &ast.QuadrantDiagram{
				YAxis: ast.QuadrantAxis{Min: "Bottom", Max: "Top"},
				Points: []ast.QuadrantPoint{
					{Name: "A", X: 0.5, Y: 0.5, Pos: ast.Position{Line: 2, Column: 1}},
				},
				Pos: ast.Position{Line: 1, Column: 1},
			},
			strict:  false,
			wantErr: true,
		},
		{
			name: "duplicate point names",
			diagram: &ast.QuadrantDiagram{
				XAxis: ast.QuadrantAxis{Min: "Low", Max: "High"},
				YAxis: ast.QuadrantAxis{Min: "Bottom", Max: "Top"},
				Points: []ast.QuadrantPoint{
					{Name: "A", X: 0.3, Y: 0.6, Pos: ast.Position{Line: 2, Column: 1}},
					{Name: "A", X: 0.7, Y: 0.4, Pos: ast.Position{Line: 3, Column: 1}},
				},
				Pos: ast.Position{Line: 1, Column: 1},
			},
			strict:  false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.ValidateQuadrant(tt.diagram, tt.strict)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("validator.ValidateQuadrant() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestQuadrantDefaultRules(t *testing.T) {
	rules := validator.QuadrantDefaultRules()
	if len(rules) == 0 {
		t.Error("validator.QuadrantDefaultRules() returned empty slice")
	}
}

func TestQuadrantStrictRules(t *testing.T) {
	rules := validator.QuadrantStrictRules()
	if len(rules) == 0 {
		t.Error("validator.QuadrantStrictRules() returned empty slice")
	}
}
