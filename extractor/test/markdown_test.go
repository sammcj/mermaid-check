package extractor_test

import (
	"fmt"
	"testing"

	"github.com/sammcj/go-mermaid/extractor"
)

func TestExtractFromMarkdown_SingleBlock(t *testing.T) {
	markdown := `# My Document

Here is a diagram:

` + "```mermaid" + `
flowchart TD
    A --> B
` + "```" + `

End of document.
`

	blocks, err := extractor.ExtractFromMarkdown(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	block := blocks[0]
	expected := "flowchart TD\n    A --> B"
	if block.Source != expected {
		t.Errorf("unexpected source:\nwant: %q\ngot:  %q", expected, block.Source)
	}

	if block.LineOffset != 6 {
		t.Errorf("expected line offset 6, got %d", block.LineOffset)
	}

	if block.DiagramType != "flowchart" {
		t.Errorf("expected diagram type 'flowchart', got %q", block.DiagramType)
	}
}

func TestExtractFromMarkdown_MultipleBlocks(t *testing.T) {
	markdown := `# Documentation

First diagram:

` + "```mermaid" + `
flowchart LR
    Start --> End
` + "```" + `

Some text here.

Second diagram:

` + "```mermaid" + `
graph TD
    A --> B
    B --> C
` + "```" + `
`

	blocks, err := extractor.ExtractFromMarkdown(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 2 {
		t.Fatalf("expected 2 blocks, got %d", len(blocks))
	}

	// Check first block
	if blocks[0].DiagramType != "flowchart" {
		t.Errorf("first block: expected type 'flowchart', got %q", blocks[0].DiagramType)
	}
	if blocks[0].LineOffset != 6 {
		t.Errorf("first block: expected line offset 6, got %d", blocks[0].LineOffset)
	}

	// Check second block
	if blocks[1].DiagramType != "graph" {
		t.Errorf("second block: expected type 'graph', got %q", blocks[1].DiagramType)
	}
	if blocks[1].LineOffset != 15 {
		t.Errorf("second block: expected line offset 15, got %d", blocks[1].LineOffset)
	}
}

func TestExtractFromMarkdown_EmptyMarkdown(t *testing.T) {
	blocks, err := extractor.ExtractFromMarkdown("")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 0 {
		t.Errorf("expected 0 blocks, got %d", len(blocks))
	}
}

func TestExtractFromMarkdown_NoMermaidBlocks(t *testing.T) {
	markdown := `# Regular Markdown

This is just regular markdown content.

` + "```python" + `
print("Hello, world!")
` + "```" + `
`

	blocks, err := extractor.ExtractFromMarkdown(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 0 {
		t.Errorf("expected 0 blocks, got %d", len(blocks))
	}
}

func TestExtractFromMarkdown_EmptyMermaidBlock(t *testing.T) {
	markdown := `# Document

` + "```mermaid" + `
` + "```" + `
`

	blocks, err := extractor.ExtractFromMarkdown(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Empty blocks should not be included
	if len(blocks) != 0 {
		t.Errorf("expected 0 blocks (empty blocks should be skipped), got %d", len(blocks))
	}
}

func TestExtractFromMarkdown_WhitespaceOnlyBlock(t *testing.T) {
	markdown := `# Document

` + "```mermaid" + `


` + "```" + `
`

	blocks, err := extractor.ExtractFromMarkdown(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	// Whitespace-only blocks should not be included
	if len(blocks) != 0 {
		t.Errorf("expected 0 blocks (whitespace-only blocks should be skipped), got %d", len(blocks))
	}
}

func TestExtractFromMarkdown_WithComments(t *testing.T) {
	markdown := `# Diagram with Comments

` + "```mermaid" + `
%% This is a comment
flowchart TD
    %% Another comment
    A --> B
` + "```" + `
`

	blocks, err := extractor.ExtractFromMarkdown(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	if blocks[0].DiagramType != "flowchart" {
		t.Errorf("expected diagram type 'flowchart', got %q", blocks[0].DiagramType)
	}
}

func TestExtractFromMarkdown_DifferentDiagramTypes(t *testing.T) {
	tests := []struct {
		name         string
		source       string
		expectedType string
	}{
		{
			name: "sequence diagram",
			source: `sequenceDiagram
    Alice->>Bob: Hello`,
			expectedType: "sequence",
		},
		{
			name: "class diagram",
			source: `classDiagram
    class Animal`,
			expectedType: "class",
		},
		{
			name: "state diagram",
			source: `stateDiagram
    [*] --> State1`,
			expectedType: "state",
		},
		{
			name: "state diagram v2",
			source: `stateDiagram-v2
    [*] --> State1`,
			expectedType: "stateDiagram-v2",
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
    title A Gantt Diagram`,
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
    title My working day`,
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
    title History of Social Media`,
			expectedType: "timeline",
		},
		{
			name: "sankey",
			source: `sankey-beta
    Agricultural 'waste',Bio-conversion,124.729`,
			expectedType: "sankey",
		},
		{
			name: "quadrant chart",
			source: `quadrantChart
    title Reach and engagement of campaigns`,
			expectedType: "quadrantChart",
		},
		{
			name: "xyChart",
			source: `xychart-beta
    title "Sales Revenue"`,
			expectedType: "xyChart",
		},
		{
			name: "C4 Context",
			source: `C4Context
    title System Context diagram`,
			expectedType: "c4Context",
		},
		{
			name: "C4 Container",
			source: `C4Container
    title Container diagram`,
			expectedType: "c4Container",
		},
		{
			name: "C4 Component",
			source: `C4Component
    title Component diagram`,
			expectedType: "c4Component",
		},
		{
			name: "C4 Dynamic",
			source: `C4Dynamic
    title Dynamic diagram`,
			expectedType: "c4Dynamic",
		},
		{
			name: "C4 Deployment",
			source: `C4Deployment
    title Deployment Diagram`,
			expectedType: "c4Deployment",
		},
		{
			name: "unknown type",
			source: `unknown diagram type`,
			expectedType: "unknown",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			markdown := "```mermaid\n" + tt.source + "\n```"
			blocks, err := extractor.ExtractFromMarkdown(markdown)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if len(blocks) != 1 {
				t.Fatalf("expected 1 block, got %d", len(blocks))
			}

			if blocks[0].DiagramType != tt.expectedType {
				t.Errorf("expected type %q, got %q", tt.expectedType, blocks[0].DiagramType)
			}
		})
	}
}

func TestExtractFromMarkdown_NestedCodeBlocks(t *testing.T) {
	// Markdown with code block inside a quoted section
	markdown := `# Document

Normal text.

` + "```mermaid" + `
flowchart TD
    A[Note with code: ` + "`go test`" + `] --> B
` + "```" + `
`

	blocks, err := extractor.ExtractFromMarkdown(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}
}

func TestExtractFromMarkdown_LanguageVariant(t *testing.T) {
	// Some tools use ```mermaid with a space or additional info
	markdown := "```mermaid showLineNumbers\nflowchart TD\n    A --> B\n```"

	blocks, err := extractor.ExtractFromMarkdown(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	if blocks[0].DiagramType != "flowchart" {
		t.Errorf("expected diagram type 'flowchart', got %q", blocks[0].DiagramType)
	}
}

func TestExtractFromMarkdown_OnlyComments(t *testing.T) {
	markdown := "```mermaid\n%% Comment 1\n%% Comment 2\n%% Comment 3\n```"

	blocks, err := extractor.ExtractFromMarkdown(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	if blocks[0].DiagramType != "unknown" {
		t.Errorf("expected 'unknown' for comment-only source, got %q", blocks[0].DiagramType)
	}
}

func TestExtractFromMarkdown_EmptyLinesBeforeDiagram(t *testing.T) {
	markdown := "```mermaid\n\n\nflowchart TD\n    A --> B\n```"

	blocks, err := extractor.ExtractFromMarkdown(markdown)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(blocks) != 1 {
		t.Fatalf("expected 1 block, got %d", len(blocks))
	}

	if blocks[0].DiagramType != "flowchart" {
		t.Errorf("expected 'flowchart', got %q", blocks[0].DiagramType)
	}
}

func TestExtractFromMarkdown_EscapedBackticks(t *testing.T) {
	tests := []struct {
		name     string
		markdown string
		wantLine int
		wantErr  bool
	}{
		{
			name:     "escaped mermaid fence at start of line",
			markdown: "# Document\n\n\\`\\`\\`mermaid\nflowchart TD\n    A --> B\n\\`\\`\\`",
			wantLine: 3,
			wantErr:  true,
		},
		{
			name:     "escaped closing fence at start of line",
			markdown: "```mermaid\nflowchart TD\n    A --> B\n\\`\\`\\`",
			wantLine: 4,
			wantErr:  true,
		},
		{
			name:     "escaped fence in middle of document",
			markdown: "# Document\n\nSome text\n\\`\\`\\`mermaid\nMore text",
			wantLine: 4,
			wantErr:  true,
		},
		{
			name:     "escaped backticks in inline code - should NOT error",
			markdown: "# Document\n\nDetects escaped backticks (e.g., `\\`\\`\\`mermaid` instead of proper fences)\n\n```mermaid\nflowchart TD\n    A --> B\n```",
			wantErr:  false,
		},
		{
			name:     "escaped backticks in documentation example - should NOT error",
			markdown: "Use proper fences like ```mermaid, not `\\`\\`\\`mermaid` with backslashes.\n\n```mermaid\nflowchart TD\n    A --> B\n```",
			wantErr:  false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			blocks, err := extractor.ExtractFromMarkdown(tt.markdown)

			if tt.wantErr {
				if err == nil {
					t.Fatalf("expected error for escaped backticks, got nil (found %d blocks)", len(blocks))
				}

				// Check that error message mentions the correct line number
				expectedMsg := fmt.Sprintf("line %d", tt.wantLine)
				if !contains(err.Error(), expectedMsg) {
					t.Errorf("expected error to mention %q, got: %v", expectedMsg, err)
				}

				// Check that error message mentions escaped backticks
				if !contains(err.Error(), "escaped backticks") {
					t.Errorf("expected error to mention 'escaped backticks', got: %v", err)
				}
			} else {
				if err != nil {
					t.Fatalf("expected no error for inline code examples, got: %v", err)
				}
				// Should have extracted the valid diagram
				if len(blocks) != 1 {
					t.Errorf("expected 1 valid diagram block, got %d", len(blocks))
				}
			}
		})
	}
}

// Helper function for string contains check
func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
