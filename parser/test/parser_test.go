package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestParse(t *testing.T) {
	tests := []struct {
		name            string
		source          string
		expectedType    string
		expectError     bool
		expectGeneric   bool
		expectFlowchart bool
		expectSequence  bool
		expectClass     bool
		expectState     bool
	}{
		{
			name: "flowchart diagram",
			source: `flowchart TD
    A --> B`,
			expectedType:    "flowchart",
			expectFlowchart: true,
		},
		{
			name: "graph diagram",
			source: `graph LR
    A --> B`,
			expectedType:    "graph",
			expectFlowchart: true,
		},
		{
			name: "sequence diagram",
			source: `sequenceDiagram
    Alice->>Bob: Hello`,
			expectedType:   "sequence",
			expectSequence: true,
		},
		{
			name: "class diagram",
			source: `classDiagram
    class Animal`,
			expectedType: "class",
			expectClass:  true,
		},
		{
			name: "state diagram",
			source: `stateDiagram
    [*] --> State1`,
			expectedType: "state",
			expectState:  true,
		},
		{
			name: "state diagram v2",
			source: `stateDiagram-v2
    [*] --> State1`,
			expectedType:  "stateDiagram-v2",
			expectState:   true,
		},
		{
			name: "er diagram",
			source: `erDiagram
    CUSTOMER ||--o{ ORDER : places`,
			expectedType: "er",
		},
		{
			name: "gantt chart",
			source: `gantt
    title Project`,
			expectedType: "gantt",
		},
		{
			name: "pie chart",
			source: `pie title Pets
    "Dogs" : 386`,
			expectedType: "pie",
		},
		{
			name: "journey diagram",
			source: `journey
    title My Day`,
			expectedType: "journey",
		},
		{
			name: "gitGraph",
			source: `gitGraph
    commit`,
			expectedType: "gitGraph",
		},
		{
			name: "mindmap",
			source: `mindmap
    Root`,
			expectedType: "mindmap",
		},
		{
			name: "timeline",
			source: `timeline
    title History
    2024 : Event`,
			expectedType: "timeline",
		},
		{
			name: "sankey",
			source: `sankey-beta
    A,B,10`,
			expectedType: "sankey",
		},
		{
			name: "quadrant chart",
			source: `quadrantChart
    title Reach
    x-axis Low --> High
    y-axis Low --> High
    Point A: [0.5, 0.5]`,
			expectedType: "quadrantChart",
		},
		{
			name: "xyChart",
			source: `xychart-beta
    title "Sales"
    x-axis [Q1, Q2]
    y-axis "Revenue" 0 --> 100
    bar [50, 75]`,
			expectedType: "xyChart",
		},
		{
			name: "C4 Context",
			source: `C4Context
    title System`,
			expectedType: "c4Context",
		},
		{
			name: "C4 Container",
			source: `C4Container
    title Container`,
			expectedType: "c4Container",
		},
		{
			name: "C4 Component",
			source: `C4Component
    title Component`,
			expectedType: "c4Component",
		},
		{
			name: "C4 Dynamic",
			source: `C4Dynamic
    title Dynamic`,
			expectedType: "c4Dynamic",
		},
		{
			name: "C4 Deployment",
			source: `C4Deployment
    title Deployment`,
			expectedType: "c4Deployment",
		},
		{
			name:        "empty source",
			source:      "",
			expectError: true,
		},
		{
			name:        "whitespace only",
			source:      "   \n\n  ",
			expectError: true,
		},
		{
			name: "sequence with leading comments",
			source: `%% This is a comment
%% Another comment
sequenceDiagram
    Alice->>Bob: Hello`,
			expectedType:   "sequence",
			expectSequence: true,
		},
		{
			name: "sequence with leading empty lines",
			source: `

sequenceDiagram
    Alice->>Bob: Hello`,
			expectedType:   "sequence",
			expectSequence: true,
		},
		{
			name: "unknown diagram type",
			source: `unknownDiagram
    something`,
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram, err := parser.Parse(tt.source)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error, got none")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diagram == nil {
				t.Fatal("diagram is nil")
			}

			gotType := diagram.GetType()
			if gotType != tt.expectedType {
				t.Errorf("expected type %q, got %q", tt.expectedType, gotType)
			}

			if tt.expectFlowchart {
				if _, ok := diagram.(*ast.Flowchart); !ok {
					t.Errorf("expected *ast.Flowchart, got %T", diagram)
				}
			}

			if tt.expectSequence {
				if _, ok := diagram.(*ast.SequenceDiagram); !ok {
					t.Errorf("expected *ast.SequenceDiagram, got %T", diagram)
				}
			}

			if tt.expectClass {
				if _, ok := diagram.(*ast.ClassDiagram); !ok {
					t.Errorf("expected *ast.ClassDiagram, got %T", diagram)
				}
			}

			if tt.expectState {
				if _, ok := diagram.(*ast.StateDiagram); !ok {
					t.Errorf("expected *ast.StateDiagram, got %T", diagram)
				}
			}

			if tt.expectGeneric {
				if _, ok := diagram.(*ast.GenericDiagram); !ok {
					t.Errorf("expected *ast.GenericDiagram, got %T", diagram)
				}
			}
		})
	}
}

// NOTE: TestDetectDiagramType is commented out because detectDiagramType is an unexported function
// and this file uses black-box testing (package parser_test).
// This test should be moved to a white-box test file if needed.
/*
func TestDetectDiagramType(t *testing.T) {
	tests := []struct {
		name         string
		source       string
		expectedType string
	}{
		{
			name:         "flowchart",
			source:       "flowchart TD\n    A --> B",
			expectedType: "flowchart",
		},
		{
			name:         "graph",
			source:       "graph LR\n    A --> B",
			expectedType: "graph",
		},
		{
			name:         "sequence",
			source:       "sequenceDiagram\n    Alice->>Bob: Hi",
			expectedType: "sequence",
		},
		{
			name:         "class",
			source:       "classDiagram\n    class Animal",
			expectedType: "class",
		},
		{
			name:         "state",
			source:       "stateDiagram\n    [*] --> S1",
			expectedType: "state",
		},
		{
			name:         "stateDiagram-v2 priority",
			source:       "stateDiagram-v2\n    [*] --> S1",
			expectedType: "stateDiagram-v2",
		},
		{
			name:         "er",
			source:       "erDiagram\n    CUSTOMER ||--o{ ORDER : places",
			expectedType: "er",
		},
		{
			name:         "gantt",
			source:       "gantt\n    title Project",
			expectedType: "gantt",
		},
		{
			name:         "pie",
			source:       "pie\n    \"A\" : 10",
			expectedType: "pie",
		},
		{
			name:         "journey",
			source:       "journey\n    title My Day",
			expectedType: "journey",
		},
		{
			name:         "gitGraph",
			source:       "gitGraph\n    commit",
			expectedType: "gitGraph",
		},
		{
			name:         "mindmap",
			source:       "mindmap\n    Root",
			expectedType: "mindmap",
		},
		{
			name:         "timeline",
			source:       "timeline\n    title History",
			expectedType: "timeline",
		},
		{
			name:         "sankey",
			source:       "sankey-beta\n    A,B,10",
			expectedType: "sankey",
		},
		{
			name:         "quadrantChart",
			source:       "quadrantChart\n    title Chart",
			expectedType: "quadrantChart",
		},
		{
			name:         "xyChart",
			source:       "xychart-beta\n    title Sales",
			expectedType: "xyChart",
		},
		{
			name:         "C4Context",
			source:       "C4Context\n    title System",
			expectedType: "c4Context",
		},
		{
			name:         "C4Container",
			source:       "C4Container\n    title Container",
			expectedType: "c4Container",
		},
		{
			name:         "C4Component",
			source:       "C4Component\n    title Component",
			expectedType: "c4Component",
		},
		{
			name:         "C4Dynamic",
			source:       "C4Dynamic\n    title Dynamic",
			expectedType: "c4Dynamic",
		},
		{
			name:         "C4Deployment",
			source:       "C4Deployment\n    title Deployment",
			expectedType: "c4Deployment",
		},
		{
			name:         "with leading comments",
			source:       "%% Comment\n%% Another\nflowchart TD\n    A",
			expectedType: "flowchart",
		},
		{
			name:         "with leading empty lines",
			source:       "\n\n\nflowchart TD\n    A",
			expectedType: "flowchart",
		},
		{
			name:         "only comments",
			source:       "%% Comment 1\n%% Comment 2",
			expectedType: "unknown",
		},
		{
			name:         "only whitespace",
			source:       "   \n\n   ",
			expectedType: "unknown",
		},
		{
			name:         "unknown type",
			source:       "unknownDiagram\n    content",
			expectedType: "unknown",
		},
		{
			name:         "empty string",
			source:       "",
			expectedType: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := detectDiagramType(tt.source)
			if got != tt.expectedType {
				t.Errorf("expected %q, got %q", tt.expectedType, got)
			}
		})
	}
}
*/

func TestParseWithRealFlowchart(t *testing.T) {
	source := `flowchart TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Do Something]
    B -->|No| D[Do Something Else]
    C --> E[End]
    D --> E`

	diagram, err := parser.Parse(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	flowchart, ok := diagram.(*ast.Flowchart)
	if !ok {
		t.Fatalf("expected *ast.Flowchart, got %T", diagram)
	}

	if flowchart.Direction != "TD" {
		t.Errorf("expected direction TD, got %s", flowchart.Direction)
	}

	if len(flowchart.Statements) == 0 {
		t.Error("expected statements, got none")
	}
}

func TestParseWithRealSequenceDiagram(t *testing.T) {
	source := `sequenceDiagram
    participant Alice
    participant Bob
    Alice->>Bob: Hello Bob, how are you?
    Bob-->>Alice: Great!`

	diagram, err := parser.Parse(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	sequence, ok := diagram.(*ast.SequenceDiagram)
	if !ok {
		t.Fatalf("expected *ast.SequenceDiagram, got %T", diagram)
	}

	if sequence.Type != "sequence" {
		t.Errorf("expected type 'sequence', got %s", sequence.Type)
	}

	if len(sequence.Statements) == 0 {
		t.Error("expected statements, got none")
	}
}

func TestParseWithRealERDiagram(t *testing.T) {
	source := `erDiagram
    CUSTOMER ||--o{ ORDER : places
    ORDER ||--o{ LINE_ITEM : contains`

	diagram, err := parser.Parse(source)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	erDiagram, ok := diagram.(*ast.ERDiagram)
	if !ok {
		t.Fatalf("expected *ast.ERDiagram, got %T", diagram)
	}

	if erDiagram.Type != "er" {
		t.Errorf("expected type 'er', got %s", erDiagram.Type)
	}

	if len(erDiagram.Relationships) == 0 {
		t.Error("expected relationships, got none")
	}

	if erDiagram.Source == "" {
		t.Error("expected source to be populated")
	}
}
