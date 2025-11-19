package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestXYChartParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
		check   func(*testing.T, ast.Diagram)
	}{
		{
			name: "complete chart with categorical x-axis",
			source: `xychart-beta
    title "Sales Revenue"
    x-axis [jan, feb, mar, apr, may, jun, jul, aug, sep, oct, nov, dec]
    y-axis "Revenue (in $)" 4000 --> 11000
    bar [5000, 6000, 7500, 8200, 9500, 10500, 11000, 10200, 9200, 8500, 7000, 6000]
    line [5000, 6000, 7500, 8200, 9500, 10500, 11000, 10200, 9200, 8500, 7000, 6000]`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				chart, ok := d.(*ast.XYChartDiagram)
				if !ok {
					t.Fatalf("expected *ast.XYChartDiagram, got %T", d)
				}
				if chart.Title != "Sales Revenue" {
					t.Errorf("expected title 'Sales Revenue', got %q", chart.Title)
				}
				if chart.Orientation != "vertical" {
					t.Errorf("expected orientation 'vertical', got %q", chart.Orientation)
				}
				if len(chart.XAxis.Categories) != 12 {
					t.Errorf("expected 12 x-axis categories, got %d", len(chart.XAxis.Categories))
				}
				if chart.YAxis.Label != "Revenue (in $)" {
					t.Errorf("expected y-axis label 'Revenue (in $)', got %q", chart.YAxis.Label)
				}
				if chart.YAxis.Min != 4000 {
					t.Errorf("expected y-axis min 4000, got %f", chart.YAxis.Min)
				}
				if chart.YAxis.Max != 11000 {
					t.Errorf("expected y-axis max 11000, got %f", chart.YAxis.Max)
				}
				if len(chart.Series) != 2 {
					t.Fatalf("expected 2 series, got %d", len(chart.Series))
				}
				if chart.Series[0].Type != "bar" {
					t.Errorf("expected first series type 'bar', got %q", chart.Series[0].Type)
				}
				if chart.Series[1].Type != "line" {
					t.Errorf("expected second series type 'line', got %q", chart.Series[1].Type)
				}
				if len(chart.Series[0].Values) != 12 {
					t.Errorf("expected 12 values in first series, got %d", len(chart.Series[0].Values))
				}
			},
		},
		{
			name: "horizontal orientation",
			source: `xychart-beta horizontal
    title "Test Chart"
    x-axis "Time" 0 --> 100
    y-axis [category1, category2, category3]
    bar [10, 20, 30]`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				chart, ok := d.(*ast.XYChartDiagram)
				if !ok {
					t.Fatalf("expected *ast.XYChartDiagram, got %T", d)
				}
				if chart.Orientation != "horizontal" {
					t.Errorf("expected orientation 'horizontal', got %q", chart.Orientation)
				}
				if !chart.XAxis.IsNumeric {
					t.Error("expected x-axis to be numeric")
				}
				if chart.YAxis.IsNumeric {
					t.Error("expected y-axis to be categorical")
				}
			},
		},
		{
			name: "numeric axes",
			source: `xychart-beta
    x-axis "Temperature" -10 --> 40
    y-axis "Pressure" 0 --> 100
    line [10, 20, 30, 40, 50]`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				chart, ok := d.(*ast.XYChartDiagram)
				if !ok {
					t.Fatalf("expected *ast.XYChartDiagram, got %T", d)
				}
				if !chart.XAxis.IsNumeric {
					t.Error("expected x-axis to be numeric")
				}
				if !chart.YAxis.IsNumeric {
					t.Error("expected y-axis to be numeric")
				}
				if chart.XAxis.Min != -10 {
					t.Errorf("expected x-axis min -10, got %f", chart.XAxis.Min)
				}
				if chart.XAxis.Max != 40 {
					t.Errorf("expected x-axis max 40, got %f", chart.XAxis.Max)
				}
			},
		},
		{
			name: "chart with comments",
			source: `xychart-beta
    %% This is a comment
    title "Test"
    x-axis [a, b, c]
    %% Another comment
    y-axis "Values" 0 --> 100
    bar [10, 20, 30]`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				chart, ok := d.(*ast.XYChartDiagram)
				if !ok {
					t.Fatalf("expected *ast.XYChartDiagram, got %T", d)
				}
				if len(chart.Series) != 1 {
					t.Errorf("expected 1 series, got %d", len(chart.Series))
				}
			},
		},
		{
			name: "multiple series of same type",
			source: `xychart-beta
    x-axis [q1, q2, q3, q4]
    y-axis "Revenue" 0 --> 100
    bar [25, 30, 35, 40]
    bar [20, 25, 30, 35]
    bar [15, 20, 25, 30]`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				chart, ok := d.(*ast.XYChartDiagram)
				if !ok {
					t.Fatalf("expected *ast.XYChartDiagram, got %T", d)
				}
				if len(chart.Series) != 3 {
					t.Errorf("expected 3 series, got %d", len(chart.Series))
				}
				for i, series := range chart.Series {
					if series.Type != "bar" {
						t.Errorf("expected series %d to be 'bar', got %q", i, series.Type)
					}
				}
			},
		},
		{
			name: "chart without title",
			source: `xychart-beta
    x-axis [a, b]
    y-axis "Y" 0 --> 10
    line [5, 8]`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				chart, ok := d.(*ast.XYChartDiagram)
				if !ok {
					t.Fatalf("expected *ast.XYChartDiagram, got %T", d)
				}
				if chart.Title != "" {
					t.Errorf("expected no title, got %q", chart.Title)
				}
			},
		},
		{
			name:    "invalid header",
			source:  "notxychart\n",
			wantErr: true,
		},
		{
			name: "missing x-axis",
			source: `xychart-beta
    title "Test"
    y-axis "Y" 0 --> 10
    bar [5, 8]`,
			wantErr: true,
		},
		{
			name: "missing y-axis",
			source: `xychart-beta
    title "Test"
    x-axis [a, b]
    bar [5, 8]`,
			wantErr: true,
		},
		{
			name: "no data series",
			source: `xychart-beta
    x-axis [a, b]
    y-axis "Y" 0 --> 10`,
			wantErr: true,
		},
		{
			name: "duplicate x-axis definition",
			source: `xychart-beta
    x-axis [a, b]
    x-axis [c, d]
    y-axis "Y" 0 --> 10
    bar [5, 8]`,
			wantErr: true,
		},
		{
			name: "duplicate y-axis definition",
			source: `xychart-beta
    x-axis [a, b]
    y-axis "Y" 0 --> 10
    y-axis "Z" 0 --> 20
    bar [5, 8]`,
			wantErr: true,
		},
		{
			name: "invalid numeric value in series",
			source: `xychart-beta
    x-axis [a, b]
    y-axis "Y" 0 --> 10
    bar [5, invalid]`,
			wantErr: true,
		},
		{
			name: "invalid axis range format",
			source: `xychart-beta
    x-axis "X" 0 -> 10
    y-axis "Y" 0 --> 10
    bar [5, 8]`,
			wantErr: true,
		},
		{
			name: "unrecognised line",
			source: `xychart-beta
    x-axis [a, b]
    y-axis "Y" 0 --> 10
    invalid line here
    bar [5, 8]`,
			wantErr: true,
		},
	}

	p := parser.NewXYChartParser()
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

func TestXYChartParser_SupportedTypes(t *testing.T) {
	p := parser.NewXYChartParser()
	types := p.SupportedTypes()
	if len(types) != 1 || types[0] != "xyChart" {
		t.Errorf("expected [xyChart], got %v", types)
	}
}

// NOTE: TestParseCategories is commented out because parseCategories is an unexported function
// and this file uses black-box testing (package parser_test).
// This test should be moved to a white-box test file if needed.
/*
func TestParseCategories(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple categories",
			input:    "a, b, c",
			expected: []string{"a", "b", "c"},
		},
		{
			name:     "categories without spaces",
			input:    "jan,feb,mar",
			expected: []string{"jan", "feb", "mar"},
		},
		{
			name:     "categories with extra spaces",
			input:    "  x  ,  y  ,  z  ",
			expected: []string{"x", "y", "z"},
		},
		{
			name:     "single category",
			input:    "single",
			expected: []string{"single"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseCategories(tt.input)
			if len(result) != len(tt.expected) {
				t.Fatalf("expected %d categories, got %d", len(tt.expected), len(result))
			}
			for i, cat := range result {
				if cat != tt.expected[i] {
					t.Errorf("category %d: expected %q, got %q", i, tt.expected[i], cat)
				}
			}
		})
	}
}
*/

// NOTE: TestParseValues is commented out because parseValues is an unexported function
// and this file uses black-box testing (package parser_test).
// This test should be moved to a white-box test file if needed.
/*
func TestParseValues(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		want    []float64
		wantErr bool
	}{
		{
			name:    "integer values",
			input:   "10, 20, 30",
			want:    []float64{10, 20, 30},
			wantErr: false,
		},
		{
			name:    "decimal values",
			input:   "1.5, 2.75, 3.25",
			want:    []float64{1.5, 2.75, 3.25},
			wantErr: false,
		},
		{
			name:    "negative values",
			input:   "-10, 0, 10",
			want:    []float64{-10, 0, 10},
			wantErr: false,
		},
		{
			name:    "values without spaces",
			input:   "1,2,3",
			want:    []float64{1, 2, 3},
			wantErr: false,
		},
		{
			name:    "invalid value",
			input:   "10, invalid, 30",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := parseValues(tt.input)
			if (err != nil) != tt.wantErr {
				t.Errorf("parseValues() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr {
				if len(result) != len(tt.want) {
					t.Fatalf("expected %d values, got %d", len(tt.want), len(result))
				}
				for i, val := range result {
					if val != tt.want[i] {
						t.Errorf("value %d: expected %f, got %f", i, tt.want[i], val)
					}
				}
			}
		})
	}
}
*/
