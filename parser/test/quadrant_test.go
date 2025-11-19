package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestQuadrantParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
		check   func(*testing.T, ast.Diagram)
	}{
		{
			name: "complete quadrant chart with all elements",
			source: `quadrantChart
    title Reach and engagement of campaigns
    x-axis Low Reach --> High Reach
    y-axis Low Engagement --> High Engagement
    quadrant-1 We should expand
    quadrant-2 Need to promote
    quadrant-3 Re-evaluate
    quadrant-4 May be improved
    Campaign A: [0.3, 0.6]
    Campaign B: [0.45, 0.23]
    Campaign C: [0.57, 0.69]`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				quad, ok := d.(*ast.QuadrantDiagram)
				if !ok {
					t.Fatalf("expected *ast.QuadrantDiagram, got %T", d)
				}
				if quad.Title != "Reach and engagement of campaigns" {
					t.Errorf("expected title 'Reach and engagement of campaigns', got %q", quad.Title)
				}
				if quad.XAxis.Min != "Low Reach" || quad.XAxis.Max != "High Reach" {
					t.Errorf("expected x-axis 'Low Reach' --> 'High Reach', got %q --> %q", quad.XAxis.Min, quad.XAxis.Max)
				}
				if quad.YAxis.Min != "Low Engagement" || quad.YAxis.Max != "High Engagement" {
					t.Errorf("expected y-axis 'Low Engagement' --> 'High Engagement', got %q --> %q", quad.YAxis.Min, quad.YAxis.Max)
				}
				if quad.QuadrantLabels[0] != "We should expand" {
					t.Errorf("expected quadrant-1 'We should expand', got %q", quad.QuadrantLabels[0])
				}
				if quad.QuadrantLabels[1] != "Need to promote" {
					t.Errorf("expected quadrant-2 'Need to promote', got %q", quad.QuadrantLabels[1])
				}
				if quad.QuadrantLabels[2] != "Re-evaluate" {
					t.Errorf("expected quadrant-3 'Re-evaluate', got %q", quad.QuadrantLabels[2])
				}
				if quad.QuadrantLabels[3] != "May be improved" {
					t.Errorf("expected quadrant-4 'May be improved', got %q", quad.QuadrantLabels[3])
				}
				if len(quad.Points) != 3 {
					t.Errorf("expected 3 points, got %d", len(quad.Points))
				}
				if quad.Points[0].Name != "Campaign A" || quad.Points[0].X != 0.3 || quad.Points[0].Y != 0.6 {
					t.Errorf("expected Campaign A [0.3, 0.6], got %s [%f, %f]", quad.Points[0].Name, quad.Points[0].X, quad.Points[0].Y)
				}
			},
		},
		{
			name: "quadrant chart without title",
			source: `quadrantChart
    x-axis Low --> High
    y-axis Bottom --> Top
    Point A: [0.5, 0.5]`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				quad, ok := d.(*ast.QuadrantDiagram)
				if !ok {
					t.Fatalf("expected *ast.QuadrantDiagram, got %T", d)
				}
				if quad.Title != "" {
					t.Errorf("expected no title, got %q", quad.Title)
				}
				if len(quad.Points) != 1 {
					t.Errorf("expected 1 point, got %d", len(quad.Points))
				}
			},
		},
		{
			name: "quadrant chart without quadrant labels",
			source: `quadrantChart
    title Test Chart
    x-axis Left --> Right
    y-axis Bottom --> Top
    Point A: [0.2, 0.8]
    Point B: [0.7, 0.3]`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				quad, ok := d.(*ast.QuadrantDiagram)
				if !ok {
					t.Fatalf("expected *ast.QuadrantDiagram, got %T", d)
				}
				if quad.QuadrantLabels[0] != "" || quad.QuadrantLabels[1] != "" || quad.QuadrantLabels[2] != "" || quad.QuadrantLabels[3] != "" {
					t.Error("expected empty quadrant labels")
				}
				if len(quad.Points) != 2 {
					t.Errorf("expected 2 points, got %d", len(quad.Points))
				}
			},
		},
		{
			name: "quadrant chart with partial quadrant labels",
			source: `quadrantChart
    x-axis Left --> Right
    y-axis Bottom --> Top
    quadrant-1 First
    quadrant-3 Third
    Point A: [0.5, 0.5]`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				quad, ok := d.(*ast.QuadrantDiagram)
				if !ok {
					t.Fatalf("expected *ast.QuadrantDiagram, got %T", d)
				}
				if quad.QuadrantLabels[0] != "First" {
					t.Errorf("expected quadrant-1 'First', got %q", quad.QuadrantLabels[0])
				}
				if quad.QuadrantLabels[1] != "" {
					t.Errorf("expected quadrant-2 empty, got %q", quad.QuadrantLabels[1])
				}
				if quad.QuadrantLabels[2] != "Third" {
					t.Errorf("expected quadrant-3 'Third', got %q", quad.QuadrantLabels[2])
				}
			},
		},
		{
			name: "quadrant chart with comments",
			source: `quadrantChart
    %% This is a comment
    title Test Chart
    x-axis Left --> Right
    %% Another comment
    y-axis Bottom --> Top
    %% Data points follow
    Point A: [0.5, 0.5]`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				quad, ok := d.(*ast.QuadrantDiagram)
				if !ok {
					t.Fatalf("expected *ast.QuadrantDiagram, got %T", d)
				}
				if len(quad.Points) != 1 {
					t.Errorf("expected 1 point, got %d", len(quad.Points))
				}
			},
		},
		{
			name: "quadrant chart with decimal coordinates",
			source: `quadrantChart
    x-axis Left --> Right
    y-axis Bottom --> Top
    Point A: [0.123, 0.456]
    Point B: [0.789, 0.012]
    Point C: [1.0, 0.0]`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				quad, ok := d.(*ast.QuadrantDiagram)
				if !ok {
					t.Fatalf("expected *ast.QuadrantDiagram, got %T", d)
				}
				if len(quad.Points) != 3 {
					t.Errorf("expected 3 points, got %d", len(quad.Points))
				}
				if quad.Points[0].X != 0.123 || quad.Points[0].Y != 0.456 {
					t.Errorf("expected Point A [0.123, 0.456], got [%f, %f]", quad.Points[0].X, quad.Points[0].Y)
				}
			},
		},
		{
			name: "quadrant chart with spaces in point names",
			source: `quadrantChart
    x-axis Left --> Right
    y-axis Bottom --> Top
    Campaign A - Q1: [0.3, 0.6]
    Campaign B - Q2: [0.7, 0.4]`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				quad, ok := d.(*ast.QuadrantDiagram)
				if !ok {
					t.Fatalf("expected *ast.QuadrantDiagram, got %T", d)
				}
				if quad.Points[0].Name != "Campaign A - Q1" {
					t.Errorf("expected point name 'Campaign A - Q1', got %q", quad.Points[0].Name)
				}
			},
		},
		{
			name:    "invalid header",
			source:  "notquadrant\n",
			wantErr: true,
		},
		{
			name: "missing x-axis",
			source: `quadrantChart
    y-axis Bottom --> Top
    Point A: [0.5, 0.5]`,
			wantErr: true,
		},
		{
			name: "missing y-axis",
			source: `quadrantChart
    x-axis Left --> Right
    Point A: [0.5, 0.5]`,
			wantErr: true,
		},
		{
			name: "no data points",
			source: `quadrantChart
    x-axis Left --> Right
    y-axis Bottom --> Top`,
			wantErr: true,
		},
		{
			name: "invalid point format - missing colon",
			source: `quadrantChart
    x-axis Left --> Right
    y-axis Bottom --> Top
    Point A [0.5, 0.5]`,
			wantErr: true,
		},
		{
			name: "invalid point format - missing brackets",
			source: `quadrantChart
    x-axis Left --> Right
    y-axis Bottom --> Top
    Point A: 0.5, 0.5`,
			wantErr: true,
		},
		{
			name: "invalid coordinate - not a number",
			source: `quadrantChart
    x-axis Left --> Right
    y-axis Bottom --> Top
    Point A: [abc, 0.5]`,
			wantErr: true,
		},
		{
			name: "invalid x-axis format",
			source: `quadrantChart
    x-axis Left Right
    y-axis Bottom --> Top
    Point A: [0.5, 0.5]`,
			wantErr: true,
		},
		{
			name: "invalid y-axis format",
			source: `quadrantChart
    x-axis Left --> Right
    y-axis Bottom Top
    Point A: [0.5, 0.5]`,
			wantErr: true,
		},
		{
			name: "invalid quadrant number",
			source: `quadrantChart
    x-axis Left --> Right
    y-axis Bottom --> Top
    quadrant-5 Invalid
    Point A: [0.5, 0.5]`,
			wantErr: true,
		},
		{
			name:    "empty source",
			source:  "",
			wantErr: true,
		},
	}

	p := parser.NewQuadrantParser()
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

func TestQuadrantParser_SupportedTypes(t *testing.T) {
	p := parser.NewQuadrantParser()
	types := p.SupportedTypes()
	if len(types) != 1 || types[0] != "quadrantChart" {
		t.Errorf("expected [quadrantChart], got %v", types)
	}
}
