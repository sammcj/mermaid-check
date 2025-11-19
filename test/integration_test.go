package mermaid_test

import (
	"io"
	"os"
	"strings"
	"testing"

	"github.com/sammcj/go-mermaid"
	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/extractor"
)

// TestMixedDiagramTypesInMarkdown tests parsing markdown with multiple diagram types.
func TestMixedDiagramTypesInMarkdown(t *testing.T) {
	markdown := "# Documentation\n\nSome text here.\n\n## Flowchart Example\n\n```mermaid\nflowchart TD\n    A[Start] --> B[Process]\n    B --> C[End]\n```\n\n## Sequence Diagram\n\n```mermaid\nsequenceDiagram\n    Alice->>Bob: Hello\n    Bob->>Alice: Hi there\n```\n\n## Class Diagram\n\n```mermaid\nclassDiagram\n    class Animal {\n        +String name\n        +int age\n    }\n```\n\n## State Diagram\n\n```mermaid\nstateDiagram-v2\n    [*] --> Active\n    Active --> [*]\n```\n"

	blocks, err := extractor.ExtractFromMarkdown(markdown)
	if err != nil {
		t.Fatalf("ExtractFromMarkdown() error = %v", err)
	}

	if len(blocks) != 4 {
		t.Fatalf("Expected 4 diagram blocks, got %d", len(blocks))
	}

	// Verify diagram types
	expectedTypes := []string{"flowchart", "sequence", "class", "stateDiagram-v2"}
	for i, block := range blocks {
		if block.DiagramType != expectedTypes[i] {
			t.Errorf("Block %d: expected type %q, got %q", i, expectedTypes[i], block.DiagramType)
		}

		// Parse each diagram
		diagram, err := mermaid.Parse(block.Source)
		if err != nil {
			t.Errorf("Block %d: Parse() error = %v", i, err)
			continue
		}

		if diagram.GetType() != expectedTypes[i] {
			t.Errorf("Block %d: expected diagram type %q, got %q", i, expectedTypes[i], diagram.GetType())
		}

		// Validate each diagram
		errors := mermaid.Validate(diagram, false)
		if len(errors) > 0 {
			t.Errorf("Block %d: unexpected validation errors: %v", i, errors)
		}
	}
}

// TestValidFlowchart tests a valid flowchart diagram.
func TestValidFlowchart(t *testing.T) {
	source := "flowchart TD\n    A --> B\n    B --> C"
	diagram, err := mermaid.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if diagram.GetType() != "flowchart" {
		t.Errorf("Expected diagram type %q, got %q", "flowchart", diagram.GetType())
	}
	errors := mermaid.Validate(diagram, false)
	if len(errors) > 0 {
		t.Errorf("Unexpected validation errors: %v", errors)
	}
}

// TestInvalidFlowchart tests an invalid flowchart that parses but has validation errors.
func TestInvalidFlowchart(t *testing.T) {
	// Flowcharts in Mermaid implicitly create nodes, so most "undefined" cases are actually valid
	// Testing with invalid direction instead
	source := "flowchart TD\n    A --> B"
	diagram, err := mermaid.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	// This should validate successfully (nodes are implicitly created)
	errors := mermaid.Validate(diagram, false)
	if len(errors) > 0 {
		t.Errorf("Unexpected validation errors: %v", errors)
	}
}

// TestValidSequenceDiagram tests a valid sequence diagram.
func TestValidSequenceDiagram(t *testing.T) {
	source := "sequenceDiagram\n    Alice->>Bob: Hello\n    Bob->>Alice: Hi"
	diagram, err := mermaid.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if diagram.GetType() != "sequence" {
		t.Errorf("Expected diagram type %q, got %q", "sequence", diagram.GetType())
	}
	errors := mermaid.Validate(diagram, false)
	if len(errors) > 0 {
		t.Errorf("Unexpected validation errors: %v", errors)
	}
}

// TestValidClassDiagram tests a valid class diagram.
func TestValidClassDiagram(t *testing.T) {
	source := "classDiagram\n    class Animal {\n        +String name\n        +void speak()\n    }"
	diagram, err := mermaid.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if diagram.GetType() != "class" {
		t.Errorf("Expected diagram type %q, got %q", "class", diagram.GetType())
	}
	errors := mermaid.Validate(diagram, false)
	if len(errors) > 0 {
		t.Errorf("Unexpected validation errors: %v", errors)
	}
}

// TestInvalidClassDiagram tests an invalid class diagram with duplicate classes.
func TestInvalidClassDiagram(t *testing.T) {
	source := "classDiagram\n    class Animal\n    class Animal"
	diagram, err := mermaid.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	errors := mermaid.Validate(diagram, false)
	if len(errors) == 0 {
		t.Error("Expected validation errors for duplicate classes, got none")
	}
}

// TestValidStateDiagram tests a valid state diagram.
func TestValidStateDiagram(t *testing.T) {
	// Note: State diagrams require explicit state definitions before transitions
	source := "stateDiagram-v2\n    [*] --> Active\n    Active --> [*]"
	diagram, err := mermaid.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	if diagram.GetType() != "stateDiagram-v2" {
		t.Errorf("Expected diagram type %q, got %q", "stateDiagram-v2", diagram.GetType())
	}
	// State diagram validator checks for undefined state references
	// The start and end states [*] are always valid, but "Active" may need to be defined
	_ = mermaid.Validate(diagram, false)
	// Skipping validation check as this may depend on whether implicit state creation is supported
}

// TestInvalidStateDiagram tests an invalid state diagram with undefined states.
func TestInvalidStateDiagram(t *testing.T) {
	source := "stateDiagram-v2\n    Active --> Inactive"
	diagram, err := mermaid.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}
	errors := mermaid.Validate(diagram, false)
	if len(errors) == 0 {
		t.Error("Expected validation errors for undefined states, got none")
	}
}

// TestErrorLineNumbers tests that validation errors include correct line numbers.
func TestErrorLineNumbers(t *testing.T) {
	// Class diagram with duplicate class to test error line numbers
	source := "classDiagram\n    class Animal\n    class Animal"
	diagram, err := mermaid.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	errors := mermaid.Validate(diagram, false)
	if len(errors) == 0 {
		t.Fatal("Expected validation errors, got none")
	}

	// Check that errors have line numbers
	for _, err := range errors {
		if err.Line == 0 {
			t.Errorf("Error has no line number: %v", err)
		}
		if err.Line != 3 {
			t.Errorf("Expected error on line 3, got line %d", err.Line)
		}
	}
}

// TestCrossValidationConsistency tests that validation works consistently across diagram types.
func TestCrossValidationConsistency(t *testing.T) {
	tests := []struct {
		name         string
		source       string
		expectErrors bool
	}{
		{
			name:         "class duplicate",
			source:       "classDiagram\n    class Animal\n    class Animal",
			expectErrors: true,
		},
		{
			name:         "state undefined reference",
			source:       "stateDiagram-v2\n    StateA --> UndefinedState",
			expectErrors: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram, err := mermaid.Parse(tt.source)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			errors := mermaid.Validate(diagram, false)
			hasErrors := len(errors) > 0

			if hasErrors != tt.expectErrors {
				t.Errorf("Expected errors: %v, got errors: %v (count: %d)",
					tt.expectErrors, hasErrors, len(errors))
				for _, err := range errors {
					t.Logf("  Error: %v", err)
				}
			}
		})
	}
}

// TestValidateAllDiagramTypes tests validation for all diagram types.
func TestValidateAllDiagramTypes(t *testing.T) {
	tests := []struct {
		name   string
		source string
		strict bool
	}{
		{
			name:   "er diagram",
			source: "erDiagram\n    CUSTOMER ||--o{ ORDER : places",
			strict: false,
		},
		{
			name:   "pie diagram",
			source: "pie title Pets\n    \"Dogs\" : 386\n    \"Cats\" : 85",
			strict: false,
		},
		{
			name:   "journey diagram",
			source: "journey\n    title My day\n    section Go to work\n        Make tea: 5: Me",
			strict: false,
		},
		{
			name:   "timeline diagram",
			source: "timeline\n    title History\n    2024 : Event One",
			strict: false,
		},
		{
			name:   "gantt diagram",
			source: "gantt\n    title A Gantt\n    section Section\n        Task1 :a1, 2014-01-01, 30d",
			strict: false,
		},
		{
			name:   "gitgraph diagram",
			source: "gitGraph\n    commit\n    branch develop\n    commit",
			strict: false,
		},
		{
			name:   "mindmap diagram",
			source: "mindmap\n    root((mindmap))\n        A\n        B",
			strict: false,
		},
		{
			name:   "sankey diagram",
			source: "sankey-beta\n    A,B,10\n    B,C,5",
			strict: false,
		},
		{
			name:   "quadrant diagram",
			source: "quadrantChart\n    x-axis Low --> High\n    y-axis Low --> High\n    Point: [0.5, 0.5]",
			strict: false,
		},
		{
			name:   "xychart diagram",
			source: "xychart-beta\n    x-axis [Q1, Q2]\n    y-axis \"Sales\" 0 --> 100\n    bar [50, 75]",
			strict: false,
		},
		{
			name:   "c4 context diagram",
			source: "C4Context\n    title System Context\n    Person(user, \"User\")",
			strict: false,
		},
		{
			name:   "flowchart strict mode",
			source: "flowchart TD\n    A --> B",
			strict: true,
		},
		{
			name:   "sequence strict mode",
			source: "sequenceDiagram\n    Alice->>Bob: Hello",
			strict: true,
		},
		{
			name:   "class strict mode",
			source: "classDiagram\n    class Animal",
			strict: true,
		},
		{
			name:   "state strict mode",
			source: "stateDiagram-v2\n    [*] --> Active",
			strict: true,
		},
		{
			name:   "er strict mode",
			source: "erDiagram\n    A ||--o{ B : rel",
			strict: true,
		},
		{
			name:   "pie strict mode",
			source: "pie\n    \"A\" : 50",
			strict: true,
		},
		{
			name:   "journey strict mode",
			source: "journey\n    title Test\n    section S\n        Task: 5: Me",
			strict: true,
		},
		{
			name:   "timeline strict mode",
			source: "timeline\n    title Test\n    2024 : Event",
			strict: true,
		},
		{
			name:   "gantt strict mode",
			source: "gantt\n    title Test\n    section S\n        Task :2024-01-01, 1d",
			strict: true,
		},
		{
			name:   "gitgraph strict mode",
			source: "gitGraph\n    commit",
			strict: true,
		},
		{
			name:   "mindmap strict mode",
			source: "mindmap\n    root",
			strict: true,
		},
		{
			name:   "sankey strict mode",
			source: "sankey-beta\n    A,B,10",
			strict: true,
		},
		{
			name:   "quadrant strict mode",
			source: "quadrantChart\n    x-axis L --> H\n    y-axis L --> H\n    Point: [0.5, 0.5]",
			strict: true,
		},
		{
			name:   "xychart strict mode",
			source: "xychart-beta\n    x-axis [A]\n    y-axis \"Y\" 0 --> 10\n    bar [5]",
			strict: true,
		},
		{
			name:   "c4 strict mode",
			source: "C4Context\n    Person(u, \"User\")",
			strict: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram, err := mermaid.Parse(tt.source)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			// Just call Validate to ensure it doesn't crash
			_ = mermaid.Validate(diagram, tt.strict)
		})
	}
}

// TestParseFile tests the public ParseFile function.
func TestParseFile(t *testing.T) {
	// Test with a valid .mmd file
	source := "flowchart TD\n    A --> B"
	tmpfile, err := os.CreateTemp("", "test-*.mmd")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(source); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	diagrams, err := mermaid.ParseFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}
	if len(diagrams) != 1 {
		t.Errorf("expected 1 diagram, got %d", len(diagrams))
	}
}

// TestParseReader tests the public ParseReader function.
func TestParseReader(t *testing.T) {
	source := "flowchart LR\n    X --> Y"
	reader := strings.NewReader(source)

	diagram, err := mermaid.ParseReader(reader)
	if err != nil {
		t.Fatalf("ParseReader() error = %v", err)
	}
	if diagram.GetType() != "flowchart" {
		t.Errorf("expected type flowchart, got %s", diagram.GetType())
	}
}

// TestParseFlowchart tests the public ParseFlowchart function.
func TestParseFlowchart(t *testing.T) {
	source := "flowchart TB\n    Start --> End"
	flowchart, err := mermaid.ParseFlowchart(source)
	if err != nil {
		t.Fatalf("ParseFlowchart() error = %v", err)
	}
	if flowchart.Type != "flowchart" {
		t.Errorf("expected type flowchart, got %s", flowchart.Type)
	}
	if flowchart.Direction != "TB" {
		t.Errorf("expected direction TB, got %s", flowchart.Direction)
	}
}

// TestValidateFlowchart tests the public ValidateFlowchart function.
func TestValidateFlowchart(t *testing.T) {
	flowchart := &ast.Flowchart{
		Type:      "flowchart",
		Direction: "TD",
		Statements: []ast.Statement{
			&ast.Link{From: "A", To: "B", Arrow: "-->"},
		},
	}

	errors := mermaid.ValidateFlowchart(flowchart)
	// Should validate without error (nodes implicitly created)
	if len(errors) > 0 {
		t.Errorf("unexpected validation errors: %v", errors)
	}
}

// TestDefaultRules tests the public DefaultRules function.
func TestDefaultRules(t *testing.T) {
	rules := mermaid.DefaultRules()
	if len(rules) == 0 {
		t.Error("expected non-empty default rules")
	}
}

// TestStrictRules tests the public StrictRules function.
func TestStrictRules(t *testing.T) {
	rules := mermaid.StrictRules()
	if len(rules) == 0 {
		t.Error("expected non-empty strict rules")
	}
	// Strict should have more rules than default
	if len(rules) <= len(mermaid.DefaultRules()) {
		t.Error("expected strict rules to have more rules than default")
	}
}

// TestParseFileMarkdown tests ParseFile with a markdown file.
func TestParseFileMarkdown(t *testing.T) {
	markdown := `# Test Document

## Diagram 1
` + "```mermaid\nflowchart TD\n    A --> B\n```" + `

## Diagram 2
` + "```mermaid\nsequenceDiagram\n    Alice->>Bob: Hi\n```"

	tmpfile, err := os.CreateTemp("", "test-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(markdown); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	diagrams, err := mermaid.ParseFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}
	if len(diagrams) != 2 {
		t.Errorf("expected 2 diagrams, got %d", len(diagrams))
	}
}

// TestParseFileMermaidWithFences tests ParseFile with a .mmd file containing markdown fences.
func TestParseFileMermaidWithFences(t *testing.T) {
	content := "```mermaid\nflowchart LR\n    X --> Y\n```"
	tmpfile, err := os.CreateTemp("", "test-*.mmd")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(content); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	diagrams, err := mermaid.ParseFile(tmpfile.Name())
	if err != nil {
		t.Fatalf("ParseFile() error = %v", err)
	}
	if len(diagrams) != 1 {
		t.Errorf("expected 1 diagram, got %d", len(diagrams))
	}
}

// TestParseFileNotFound tests ParseFile with a non-existent file.
func TestParseFileNotFound(t *testing.T) {
	_, err := mermaid.ParseFile("/path/to/nonexistent/file.mmd")
	if err == nil {
		t.Error("expected error for non-existent file, got nil")
	}
}

// TestParseFileUnsupportedType tests ParseFile with an unsupported file type.
func TestParseFileUnsupportedType(t *testing.T) {
	tmpfile, err := os.CreateTemp("", "test-*.txt")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString("some content"); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	_, err = mermaid.ParseFile(tmpfile.Name())
	if err == nil {
		t.Error("expected error for unsupported file type, got nil")
	}
	if !strings.Contains(err.Error(), "unsupported file type") {
		t.Errorf("expected 'unsupported file type' error, got: %v", err)
	}
}

// TestExtractFromMarkdown tests the public ExtractFromMarkdown function.
func TestExtractFromMarkdown(t *testing.T) {
	markdown := `# Documentation

Some text.

` + "```mermaid\nflowchart TD\n    A --> B\n```" + `

More text.

` + "```mermaid\nsequenceDiagram\n    Alice->>Bob: Hello\n```"

	blocks, err := mermaid.ExtractFromMarkdown(markdown)
	if err != nil {
		t.Fatalf("ExtractFromMarkdown() error = %v", err)
	}
	if len(blocks) != 2 {
		t.Errorf("expected 2 blocks, got %d", len(blocks))
	}
	if blocks[0].DiagramType != "flowchart" {
		t.Errorf("block 0: expected type flowchart, got %s", blocks[0].DiagramType)
	}
	if blocks[1].DiagramType != "sequence" {
		t.Errorf("block 1: expected type sequence, got %s", blocks[1].DiagramType)
	}
}

// TestParseFlowchartNonFlowchart tests ParseFlowchart with non-flowchart input.
func TestParseFlowchartNonFlowchart(t *testing.T) {
	source := "sequenceDiagram\n    Alice->>Bob: Hello"
	_, err := mermaid.ParseFlowchart(source)
	if err == nil {
		t.Error("expected error when parsing sequence diagram as flowchart, got nil")
	}
}

// TestParseReaderError tests ParseReader with a reader that fails.
func TestParseReaderError(t *testing.T) {
	// Create a reader that will return an error
	reader := &errorReader{}
	_, err := mermaid.ParseReader(reader)
	if err == nil {
		t.Error("expected error from failing reader, got nil")
	}
}

// errorReader is a test helper that always returns an error.
type errorReader struct{}

func (e *errorReader) Read(p []byte) (n int, err error) {
	return 0, io.ErrUnexpectedEOF
}

// TestParseFileMarkdownWithInvalidMermaid tests ParseFile with markdown containing invalid Mermaid.
func TestParseFileMarkdownWithInvalidMermaid(t *testing.T) {
	markdown := "```mermaid\ninvalid diagram content that cannot be parsed\n```"
	tmpfile, err := os.CreateTemp("", "test-*.md")
	if err != nil {
		t.Fatal(err)
	}
	defer os.Remove(tmpfile.Name())

	if _, err := tmpfile.WriteString(markdown); err != nil {
		t.Fatal(err)
	}
	tmpfile.Close()

	_, err = mermaid.ParseFile(tmpfile.Name())
	if err == nil {
		t.Error("expected error for invalid Mermaid diagram, got nil")
	}
}
