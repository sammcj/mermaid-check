package mermaid

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
)

// TestASTInterfaceMethods tests that all AST types properly implement the Diagram interface.
func TestASTInterfaceMethods(t *testing.T) {
	tests := []struct {
		name    string
		diagram ast.Diagram
		expType string
	}{
		{
			name: "Flowchart",
			diagram: &ast.Flowchart{
				Type:   "flowchart",
				Source: "test",
				Pos:    ast.Position{Line: 1, Column: 1},
			},
			expType: "flowchart",
		},
		{
			name: "SequenceDiagram",
			diagram: &ast.SequenceDiagram{
				Type:   "sequence",
				Source: "test",
				Pos:    ast.Position{Line: 1, Column: 1},
			},
			expType: "sequence",
		},
		{
			name: "ClassDiagram",
			diagram: &ast.ClassDiagram{
				Type:   "class",
				Source: "test",
				Pos:    ast.Position{Line: 1, Column: 1},
			},
			expType: "class",
		},
		{
			name: "StateDiagram",
			diagram: &ast.StateDiagram{
				Type:   "stateDiagram-v2",
				Source: "test",
				Pos:    ast.Position{Line: 1, Column: 1},
			},
			expType: "stateDiagram-v2",
		},
		{
			name: "ERDiagram",
			diagram: &ast.ERDiagram{
				Type:   "er",
				Source: "test",
				Pos:    ast.Position{Line: 1, Column: 1},
			},
			expType: "er",
		},
		{
			name: "PieDiagram",
			diagram: &ast.PieDiagram{
				Type:   "pie",
				Source: "test",
				Pos:    ast.Position{Line: 1, Column: 1},
			},
			expType: "pie",
		},
		{
			name: "GanttDiagram",
			diagram: &ast.GanttDiagram{
				Type:   "gantt",
				Source: "test",
				Pos:    ast.Position{Line: 1, Column: 1},
			},
			expType: "gantt",
		},
		{
			name: "JourneyDiagram",
			diagram: &ast.JourneyDiagram{
				Type:   "journey",
				Source: "test",
				Pos:    ast.Position{Line: 1, Column: 1},
			},
			expType: "journey",
		},
		{
			name: "GitGraphDiagram",
			diagram: &ast.GitGraphDiagram{
				Type:   "gitGraph",
				Source: "test",
				Pos:    ast.Position{Line: 1, Column: 1},
			},
			expType: "gitGraph",
		},
		{
			name: "MindmapDiagram",
			diagram: &ast.MindmapDiagram{
				Type:   "mindmap",
				Source: "test",
				Pos:    ast.Position{Line: 1, Column: 1},
			},
			expType: "mindmap",
		},
		{
			name: "TimelineDiagram",
			diagram: &ast.TimelineDiagram{
				Type:   "timeline",
				Source: "test",
				Pos:    ast.Position{Line: 1, Column: 1},
			},
			expType: "timeline",
		},
		{
			name: "SankeyDiagram",
			diagram: &ast.SankeyDiagram{
				Type:   "sankey-beta",
				Source: "test",
				Pos:    ast.Position{Line: 1, Column: 1},
			},
			expType: "sankey-beta",
		},
		{
			name: "QuadrantDiagram",
			diagram: &ast.QuadrantDiagram{
				Type:   "quadrantChart",
				Source: "test",
				Pos:    ast.Position{Line: 1, Column: 1},
			},
			expType: "quadrantChart",
		},
		{
			name: "XYChartDiagram",
			diagram: &ast.XYChartDiagram{
				Type:   "xychart-beta",
				Source: "test",
				Pos:    ast.Position{Line: 1, Column: 1},
			},
			expType: "xychart-beta",
		},
		{
			name: "C4Diagram",
			diagram: &ast.C4Diagram{
				DiagramType: "c4Context",
				Source:      "test",
				Pos:         ast.Position{Line: 1, Column: 1},
			},
			expType: "c4Context",
		},
		{
			name: "GenericDiagram",
			diagram: &ast.GenericDiagram{
				DiagramType: "unknown",
				Source:      "test",
				Pos:         ast.Position{Line: 1, Column: 1},
			},
			expType: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test GetType
			if gotType := tt.diagram.GetType(); gotType != tt.expType {
				t.Errorf("GetType() = %v, want %v", gotType, tt.expType)
			}

			// Test GetPosition
			pos := tt.diagram.GetPosition()
			if pos.Line != 1 || pos.Column != 1 {
				t.Errorf("GetPosition() = %+v, want Line:1, Column:1", pos)
			}
		})
	}
}

// TestGenericDiagramCreation tests the GenericDiagram constructor.
func TestGenericDiagramCreation(t *testing.T) {
	source := "test diagram\nline 2\nline 3"
	diagram := ast.NewGenericDiagram("custom", source, ast.Position{Line: 1, Column: 1})

	if diagram.DiagramType != "custom" {
		t.Errorf("DiagramType = %v, want %v", diagram.DiagramType, "custom")
	}

	if diagram.Source != source {
		t.Errorf("Source = %v, want %v", diagram.Source, source)
	}

	if len(diagram.Lines) != 3 {
		t.Errorf("Lines count = %v, want %v", len(diagram.Lines), 3)
	}
}
