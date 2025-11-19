package validator_test

import (	"testing"

	"github.com/sammcj/go-mermaid/ast"

	"github.com/sammcj/go-mermaid/validator"
)

func TestXYChartXAxisDefinedRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.XYChartDiagram
		wantErr bool
	}{
		{
			name: "categorical x-axis defined",
			diagram: &ast.XYChartDiagram{
				XAxis: ast.XYChartAxis{
					Categories: []string{"a", "b", "c"},
					IsNumeric:  false,
				},
			},
			wantErr: false,
		},
		{
			name: "numeric x-axis defined",
			diagram: &ast.XYChartDiagram{
				XAxis: ast.XYChartAxis{
					Min:       0,
					Max:       100,
					IsNumeric: true,
				},
			},
			wantErr: false,
		},
		{
			name: "x-axis not defined",
			diagram: &ast.XYChartDiagram{
				XAxis: ast.XYChartAxis{
					IsNumeric: false,
				},
			},
			wantErr: true,
		},
	}

	rule := &validator.XYChartXAxisDefinedRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("XYChartXAxisDefinedRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestXYChartYAxisDefinedRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.XYChartDiagram
		wantErr bool
	}{
		{
			name: "categorical y-axis defined",
			diagram: &ast.XYChartDiagram{
				YAxis: ast.XYChartAxis{
					Categories: []string{"a", "b", "c"},
					IsNumeric:  false,
				},
			},
			wantErr: false,
		},
		{
			name: "numeric y-axis defined",
			diagram: &ast.XYChartDiagram{
				YAxis: ast.XYChartAxis{
					Min:       0,
					Max:       100,
					IsNumeric: true,
				},
			},
			wantErr: false,
		},
		{
			name: "y-axis not defined",
			diagram: &ast.XYChartDiagram{
				YAxis: ast.XYChartAxis{
					IsNumeric: false,
				},
			},
			wantErr: true,
		},
	}

	rule := &validator.XYChartYAxisDefinedRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("XYChartYAxisDefinedRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestXYChartMinimumSeriesRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.XYChartDiagram
		wantErr bool
	}{
		{
			name: "one series",
			diagram: &ast.XYChartDiagram{
				Series: []ast.XYChartSeries{
					{Type: "bar", Values: []float64{1, 2, 3}},
				},
			},
			wantErr: false,
		},
		{
			name: "multiple series",
			diagram: &ast.XYChartDiagram{
				Series: []ast.XYChartSeries{
					{Type: "bar", Values: []float64{1, 2, 3}},
					{Type: "line", Values: []float64{4, 5, 6}},
				},
			},
			wantErr: false,
		},
		{
			name: "no series",
			diagram: &ast.XYChartDiagram{
				Series: []ast.XYChartSeries{},
			},
			wantErr: true,
		},
	}

	rule := &validator.XYChartMinimumSeriesRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("XYChartMinimumSeriesRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestXYChartValidSeriesLengthRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.XYChartDiagram
		wantErr bool
	}{
		{
			name: "all series same length",
			diagram: &ast.XYChartDiagram{
				Series: []ast.XYChartSeries{
					{Type: "bar", Values: []float64{1, 2, 3}, Pos: ast.Position{Line: 5, Column: 1}},
					{Type: "line", Values: []float64{4, 5, 6}, Pos: ast.Position{Line: 6, Column: 1}},
				},
			},
			wantErr: false,
		},
		{
			name: "mismatched series lengths",
			diagram: &ast.XYChartDiagram{
				Series: []ast.XYChartSeries{
					{Type: "bar", Values: []float64{1, 2, 3}, Pos: ast.Position{Line: 5, Column: 1}},
					{Type: "line", Values: []float64{4, 5}, Pos: ast.Position{Line: 6, Column: 1}},
				},
			},
			wantErr: true,
		},
		{
			name: "three series with one mismatch",
			diagram: &ast.XYChartDiagram{
				Series: []ast.XYChartSeries{
					{Type: "bar", Values: []float64{1, 2, 3}, Pos: ast.Position{Line: 5, Column: 1}},
					{Type: "line", Values: []float64{4, 5, 6}, Pos: ast.Position{Line: 6, Column: 1}},
					{Type: "bar", Values: []float64{7, 8, 9, 10}, Pos: ast.Position{Line: 7, Column: 1}},
				},
			},
			wantErr: true,
		},
		{
			name: "single series",
			diagram: &ast.XYChartDiagram{
				Series: []ast.XYChartSeries{
					{Type: "bar", Values: []float64{1, 2, 3}, Pos: ast.Position{Line: 5, Column: 1}},
				},
			},
			wantErr: false,
		},
		{
			name:    "no series",
			diagram: &ast.XYChartDiagram{},
			wantErr: false,
		},
	}

	rule := &validator.XYChartValidSeriesLengthRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("XYChartValidSeriesLengthRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestXYChartValidSeriesLengthRule_WithCategoricalAxis(t *testing.T) {
	tests := []struct {
		name         string
		diagram      *ast.XYChartDiagram
		wantWarnings bool
	}{
		{
			name: "series matches x-axis categories",
			diagram: &ast.XYChartDiagram{
				XAxis: ast.XYChartAxis{
					Categories: []string{"a", "b", "c"},
					IsNumeric:  false,
					Pos:        ast.Position{Line: 3, Column: 1},
				},
				Series: []ast.XYChartSeries{
					{Type: "bar", Values: []float64{1, 2, 3}, Pos: ast.Position{Line: 5, Column: 1}},
				},
			},
			wantWarnings: false,
		},
		{
			name: "series length mismatch with x-axis categories",
			diagram: &ast.XYChartDiagram{
				XAxis: ast.XYChartAxis{
					Categories: []string{"a", "b", "c", "d"},
					IsNumeric:  false,
					Pos:        ast.Position{Line: 3, Column: 1},
				},
				Series: []ast.XYChartSeries{
					{Type: "bar", Values: []float64{1, 2, 3}, Pos: ast.Position{Line: 5, Column: 1}},
				},
			},
			wantWarnings: true,
		},
		{
			name: "series matches y-axis categories",
			diagram: &ast.XYChartDiagram{
				YAxis: ast.XYChartAxis{
					Categories: []string{"a", "b"},
					IsNumeric:  false,
					Pos:        ast.Position{Line: 4, Column: 1},
				},
				Series: []ast.XYChartSeries{
					{Type: "line", Values: []float64{1, 2}, Pos: ast.Position{Line: 5, Column: 1}},
				},
			},
			wantWarnings: false,
		},
	}

	rule := &validator.XYChartValidSeriesLengthRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			hasWarnings := false
			for _, err := range errors {
				if err.Severity == validator.SeverityWarning {
					hasWarnings = true
					break
				}
			}
			if hasWarnings != tt.wantWarnings {
				t.Errorf("XYChartValidSeriesLengthRule.Validate() warnings = %v, wantWarnings %v", errors, tt.wantWarnings)
			}
		})
	}
}

func TestXYChartValidOrientationRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.XYChartDiagram
		wantErr bool
	}{
		{
			name: "vertical orientation",
			diagram: &ast.XYChartDiagram{
				Orientation: "vertical",
			},
			wantErr: false,
		},
		{
			name: "horizontal orientation",
			diagram: &ast.XYChartDiagram{
				Orientation: "horizontal",
			},
			wantErr: false,
		},
		{
			name: "invalid orientation",
			diagram: &ast.XYChartDiagram{
				Orientation: "diagonal",
				Pos:         ast.Position{Line: 1, Column: 1},
			},
			wantErr: true,
		},
		{
			name: "empty orientation",
			diagram: &ast.XYChartDiagram{
				Orientation: "",
				Pos:         ast.Position{Line: 1, Column: 1},
			},
			wantErr: true,
		},
	}

	rule := &validator.XYChartValidOrientationRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("XYChartValidOrientationRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestXYChartDefaultRules(t *testing.T) {
	rules := validator.XYChartDefaultRules()
	if len(rules) == 0 {
		t.Error("validator.XYChartDefaultRules() returned empty slice")
	}
	expectedRuleCount := 5 // XAxisDefined, YAxisDefined, MinimumSeries, ValidSeriesLength, ValidOrientation
	if len(rules) != expectedRuleCount {
		t.Errorf("expected %d rules, got %d", expectedRuleCount, len(rules))
	}
}

func TestXYChartStrictRules(t *testing.T) {
	rules := validator.XYChartStrictRules()
	if len(rules) == 0 {
		t.Error("validator.XYChartStrictRules() returned empty slice")
	}
}

func TestValidateXYChart(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.XYChartDiagram
		strict  bool
		wantErr bool
	}{
		{
			name: "valid chart",
			diagram: &ast.XYChartDiagram{
				Type:        "xyChart",
				Orientation: "vertical",
				XAxis: ast.XYChartAxis{
					Categories: []string{"a", "b", "c"},
					IsNumeric:  false,
				},
				YAxis: ast.XYChartAxis{
					Min:       0,
					Max:       100,
					IsNumeric: true,
				},
				Series: []ast.XYChartSeries{
					{Type: "bar", Values: []float64{10, 20, 30}},
				},
			},
			strict:  false,
			wantErr: false,
		},
		{
			name: "chart with missing x-axis",
			diagram: &ast.XYChartDiagram{
				Type:        "xyChart",
				Orientation: "vertical",
				YAxis: ast.XYChartAxis{
					Min:       0,
					Max:       100,
					IsNumeric: true,
				},
				Series: []ast.XYChartSeries{
					{Type: "bar", Values: []float64{10, 20, 30}},
				},
			},
			strict:  false,
			wantErr: true,
		},
		{
			name: "chart with mismatched series lengths",
			diagram: &ast.XYChartDiagram{
				Type:        "xyChart",
				Orientation: "vertical",
				XAxis: ast.XYChartAxis{
					Categories: []string{"a", "b"},
					IsNumeric:  false,
				},
				YAxis: ast.XYChartAxis{
					Min:       0,
					Max:       100,
					IsNumeric: true,
				},
				Series: []ast.XYChartSeries{
					{Type: "bar", Values: []float64{10, 20}, Pos: ast.Position{Line: 5, Column: 1}},
					{Type: "line", Values: []float64{30, 40, 50}, Pos: ast.Position{Line: 6, Column: 1}},
				},
			},
			strict:  false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.ValidateXYChart(tt.diagram, tt.strict)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("validator.ValidateXYChart() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}
