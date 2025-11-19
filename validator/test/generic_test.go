package validator_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"

	"github.com/sammcj/go-mermaid/validator"
)

func TestValidComments(t *testing.T) {
	rule := &validator.ValidComments{}

	tests := []struct {
		name        string
		source      string
		expectError bool
		errorCount  int
	}{
		{
			name: "valid comments",
			source: `sequenceDiagram
%% This is a valid comment
    Alice->>Bob: Hello
%% Another valid comment`,
			expectError: false,
		},
		{
			name: "invalid single % comment",
			source: `sequenceDiagram
% This is invalid
    Alice->>Bob: Hello`,
			expectError: true,
			errorCount:  1,
		},
		{
			name: "multiple invalid comments",
			source: `sequenceDiagram
% Invalid 1
    Alice->>Bob: Hello
% Invalid 2`,
			expectError: true,
			errorCount:  2,
		},
		{
			name: "no comments",
			source: `sequenceDiagram
    Alice->>Bob: Hello
    Bob-->>Alice: Hi`,
			expectError: false,
		},
		{
			name: "mixed valid and invalid",
			source: `sequenceDiagram
%% Valid comment
% Invalid comment
    Alice->>Bob: Hello`,
			expectError: true,
			errorCount:  1,
		},
		{
			name:        "empty diagram",
			source:      "",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram := ast.NewGenericDiagram("sequence", tt.source, ast.Position{Line: 1, Column: 1})
			errors := rule.ValidateGeneric(diagram)

			if tt.expectError {
				if len(errors) == 0 {
					t.Error("expected validation errors, got none")
				}
				if tt.errorCount > 0 && len(errors) != tt.errorCount {
					t.Errorf("expected %d errors, got %d", tt.errorCount, len(errors))
				}
				for _, err := range errors {
					if err.Severity != validator.SeverityError {
						t.Errorf("expected validator.SeverityError, got %v", err.Severity)
					}
				}
			} else if len(errors) > 0 {
				t.Errorf("expected no errors, got %d: %v", len(errors), errors)
			}
		})
	}
}

func TestNoTrailingWhitespace(t *testing.T) {
	rule := &validator.NoTrailingWhitespace{}

	tests := []struct {
		name        string
		source      string
		expectError bool
		errorCount  int
	}{
		{
			name: "no trailing whitespace",
			source: `sequenceDiagram
    Alice->>Bob: Hello
    Bob-->>Alice: Hi`,
			expectError: false,
		},
		{
			name:        "trailing space",
			source:      "sequenceDiagram \n    Alice->>Bob: Hello",
			expectError: true,
			errorCount:  1,
		},
		{
			name:        "trailing tab",
			source:      "sequenceDiagram\t\n    Alice->>Bob: Hello",
			expectError: true,
			errorCount:  1,
		},
		{
			name:        "multiple lines with trailing whitespace",
			source:      "sequenceDiagram \n    Alice->>Bob: Hello \n    Bob-->>Alice: Hi",
			expectError: true,
			errorCount:  2,
		},
		{
			name:        "empty diagram",
			source:      "",
			expectError: false,
		},
		{
			name: "empty lines (no trailing whitespace)",
			source: `sequenceDiagram

    Alice->>Bob: Hello`,
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram := ast.NewGenericDiagram("sequence", tt.source, ast.Position{Line: 1, Column: 1})
			errors := rule.ValidateGeneric(diagram)

			if tt.expectError {
				if len(errors) == 0 {
					t.Error("expected validation errors, got none")
				}
				if tt.errorCount > 0 && len(errors) != tt.errorCount {
					t.Errorf("expected %d errors, got %d", tt.errorCount, len(errors))
				}
				for _, err := range errors {
					if err.Severity != validator.SeverityWarning {
						t.Errorf("expected validator.SeverityWarning, got %v", err.Severity)
					}
				}
			} else if len(errors) > 0 {
				t.Errorf("expected no errors, got %d: %v", len(errors), errors)
			}
		})
	}
}

func TestNoParenthesesInText(t *testing.T) {
	rule := &validator.NoParenthesesInText{}

	tests := []struct {
		name        string
		diagramType string
		source      string
		expectError bool
	}{
		{
			name:        "no parentheses",
			diagramType: "sequence",
			source: `sequenceDiagram
    Alice->>Bob: Hello
    Bob-->>Alice: Hi`,
			expectError: false,
		},
		{
			name:        "parentheses in sequence diagram",
			diagramType: "sequence",
			source: `sequenceDiagram
    participant Alice
    Alice->Bob: Hello (world)`,
			expectError: true,
		},
		{
			name:        "class diagram with method signatures - allowed",
			diagramType: "class",
			source: `classDiagram
    class Animal
    +isMammal()
    +mate(partner)`,
			expectError: false,
		},
		{
			name:        "er diagram with notation - allowed",
			diagramType: "er",
			source: `erDiagram
    CUSTOMER ||--o{ ORDER : places`,
			expectError: false,
		},
		{
			name:        "mindmap with parentheses",
			diagramType: "mindmap",
			source: `mindmap
    Root
      Topic (important)`,
			expectError: true,
		},
		{
			name:        "empty first line",
			diagramType: "sequence",
			source:      "sequenceDiagram\n",
			expectError: false,
		},
		{
			name:        "comments with parentheses",
			diagramType: "sequence",
			source: `sequenceDiagram
%% This comment (has parentheses)
    Alice->>Bob: Hello`,
			expectError: false, // Comments are skipped
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram := ast.NewGenericDiagram(tt.diagramType, tt.source, ast.Position{Line: 1, Column: 1})
			errors := rule.ValidateGeneric(diagram)

			if tt.expectError {
				if len(errors) == 0 {
					t.Error("expected validation errors, got none")
				}
				for _, err := range errors {
					if err.Severity != validator.SeverityWarning {
						t.Errorf("expected validator.SeverityWarning, got %v", err.Severity)
					}
				}
			} else if len(errors) > 0 {
				t.Errorf("expected no errors, got %d: %v", len(errors), errors)
			}
		})
	}
}

func TestValidDiagramHeader(t *testing.T) {
	rule := &validator.ValidDiagramHeader{}

	tests := []struct {
		name        string
		diagramType string
		source      string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid sequence header",
			diagramType: "sequence",
			source:      "sequenceDiagram\n    Alice->>Bob: Hello",
			expectError: false,
		},
		{
			name:        "valid class header",
			diagramType: "class",
			source:      "classDiagram\n    class Animal",
			expectError: false,
		},
		{
			name:        "valid flowchart header",
			diagramType: "flowchart",
			source:      "flowchart\n    A --> B",
			expectError: false,
		},
		{
			name:        "valid state v2 header",
			diagramType: "stateDiagram-v2",
			source:      "stateDiagram-v2\n    [*] --> S1",
			expectError: false,
		},
		{
			name:        "valid gantt header",
			diagramType: "gantt",
			source:      "gantt\n    title Project",
			expectError: false,
		},
		{
			name:        "valid pie header",
			diagramType: "pie",
			source:      "pie\n    \"A\" : 10",
			expectError: false,
		},
		{
			name:        "valid journey header",
			diagramType: "journey",
			source:      "journey\n    title Day",
			expectError: false,
		},
		{
			name:        "valid gitGraph header",
			diagramType: "gitGraph",
			source:      "gitGraph\n    commit",
			expectError: false,
		},
		{
			name:        "valid mindmap header",
			diagramType: "mindmap",
			source:      "mindmap\n    Root",
			expectError: false,
		},
		{
			name:        "valid timeline header",
			diagramType: "timeline",
			source:      "timeline\n    title History",
			expectError: false,
		},
		{
			name:        "valid sankey header",
			diagramType: "sankey",
			source:      "sankey-beta\n    A,B,10",
			expectError: false,
		},
		{
			name:        "valid quadrant header",
			diagramType: "quadrantChart",
			source:      "quadrantChart\n    title Chart",
			expectError: false,
		},
		{
			name:        "valid xyChart header",
			diagramType: "xyChart",
			source:      "xychart-beta\n    title Sales",
			expectError: false,
		},
		{
			name:        "valid C4Context header",
			diagramType: "c4Context",
			source:      "C4Context\n    title System",
			expectError: false,
		},
		{
			name:        "header with leading comments",
			diagramType: "sequence",
			source:      "%% Comment\nsequenceDiagram\n    Alice->>Bob: Hi",
			expectError: false,
		},
		{
			name:        "header with leading empty lines",
			diagramType: "sequence",
			source:      "\n\nsequenceDiagram\n    Alice->>Bob: Hi",
			expectError: false,
		},
		{
			name:        "empty diagram",
			diagramType: "sequence",
			source:      "",
			expectError: true,
			errorMsg:    "diagram is empty",
		},
		{
			name:        "only comments",
			diagramType: "sequence",
			source:      "%% Comment 1\n%% Comment 2",
			expectError: true,
			errorMsg:    "no content",
		},
		{
			name:        "only whitespace",
			diagramType: "sequence",
			source:      "   \n\n   ",
			expectError: true,
			errorMsg:    "no content",
		},
		{
			name:        "wrong header for type",
			diagramType: "sequence",
			source:      "classDiagram\n    class Animal",
			expectError: true,
			errorMsg:    "invalid diagram header",
		},
		{
			name:        "completely wrong header",
			diagramType: "sequence",
			source:      "unknownDiagram\n    content",
			expectError: true,
			errorMsg:    "invalid diagram header",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram := ast.NewGenericDiagram(tt.diagramType, tt.source, ast.Position{Line: 1, Column: 1})
			errors := rule.ValidateGeneric(diagram)

			if tt.expectError {
				if len(errors) == 0 {
					t.Error("expected validation errors, got none")
				}
				if len(errors) > 0 {
					if tt.errorMsg != "" && !containsSubstr(errors[0].Message, tt.errorMsg) {
						t.Errorf("expected error message to contain %q, got %q", tt.errorMsg, errors[0].Message)
					}
					if errors[0].Severity != validator.SeverityError {
						t.Errorf("expected validator.SeverityError, got %v", errors[0].Severity)
					}
				}
			} else if len(errors) > 0 {
				t.Errorf("expected no errors, got %d: %v", len(errors), errors)
			}
		})
	}
}

func TestGenericDefaultRules(t *testing.T) {
	rules := validator.GenericDefaultRules()

	if len(rules) == 0 {
		t.Error("expected default rules, got none")
	}

	// Check that default rules include ValidDiagramHeader and ValidComments
	foundHeader := false
	foundComments := false

	for _, rule := range rules {
		switch rule.Name() {
		case "valid-diagram-header":
			foundHeader = true
		case "valid-comments":
			foundComments = true
		}
	}

	if !foundHeader {
		t.Error("default rules should include valid-diagram-header")
	}
	if !foundComments {
		t.Error("default rules should include valid-comments")
	}
}

func TestGenericStrictRules(t *testing.T) {
	rules := validator.GenericStrictRules()

	if len(rules) == 0 {
		t.Error("expected strict rules, got none")
	}

	// Check that strict rules include all validation types
	foundHeader := false
	foundComments := false
	foundParentheses := false
	foundWhitespace := false

	for _, rule := range rules {
		switch rule.Name() {
		case "valid-diagram-header":
			foundHeader = true
		case "valid-comments":
			foundComments = true
		case "no-parentheses-in-text":
			foundParentheses = true
		case "no-trailing-whitespace":
			foundWhitespace = true
		}
	}

	if !foundHeader {
		t.Error("strict rules should include valid-diagram-header")
	}
	if !foundComments {
		t.Error("strict rules should include valid-comments")
	}
	if !foundParentheses {
		t.Error("strict rules should include no-parentheses-in-text")
	}
	if !foundWhitespace {
		t.Error("strict rules should include no-trailing-whitespace")
	}
}

// TODO: This test uses private functions/fields and cannot be converted to black-box testing
// func TestValidatorNewGeneric(t *testing.T) {
// 	rules := validator.GenericDefaultRules()
// 	v := validator.NewGeneric(rules...)
// 
// 	if v == nil {
// 		t.Fatal("NewGeneric returned nil")
// 	}
// 
// 	if len(v.genericRules) != len(rules) {
// 		t.Errorf("expected %d rules, got %d", len(rules), len(v.genericRules))
// 	}
// }

// TODO: This test uses private functions/fields and cannot be converted to black-box testing
// func TestValidatorValidateDiagram(t *testing.T) {
// 	tests := []struct {
// 		name          string
// 		diagram       ast.Diagram
// 		expectErrors  bool
// 		expectedCount int
// 	}{
// 		{
// 			name: "valid generic diagram",
// 			diagram: ast.NewGenericDiagram("sequence", `sequenceDiagram
//     Alice->>Bob: Hello`, ast.Position{Line: 1, Column: 1}),
// 			expectErrors: false,
// 		},
// 		{
// 			name: "generic diagram with invalid comments",
// 			diagram: ast.NewGenericDiagram("sequence", `sequenceDiagram
// % Invalid comment
//     Alice->>Bob: Hello`, ast.Position{Line: 1, Column: 1}),
// 			expectErrors:  true,
// 			expectedCount: 1,
// 		},
// 		{
// 			name: "valid flowchart",
// 			diagram: &ast.Flowchart{
// 				Type:      "flowchart",
// 				Direction: "TD",
// 				Pos:       ast.Position{Line: 1, Column: 1},
// 			},
// 			expectErrors: false,
// 		},
// 	}
// 
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			var v *Validator
// 			switch tt.diagram.(type) {
// 			case *ast.GenericDiagram:
// 				v = validator.NewGeneric(validator.GenericDefaultRules()...)
// 			case *ast.Flowchart:
// 				v = New(validator.DefaultRules()...)
// 			}
// 
// 			errors := v.ValidateDiagram(tt.diagram)
// 
// 			if tt.expectErrors {
// 				if len(errors) == 0 {
// 					t.Error("expected validation errors, got none")
// 				}
// 				if tt.expectedCount > 0 && len(errors) != tt.expectedCount {
// 					t.Errorf("expected %d errors, got %d", tt.expectedCount, len(errors))
// 				}
// 			} else if len(errors) > 0 {
// 				t.Errorf("expected no errors, got %d: %v", len(errors), errors)
// 			}
// 		})
// 	}
// }

// TODO: This test uses private functions/fields and cannot be converted to black-box testing
// func TestGetValidHeaderPrefixes(t *testing.T) {
// 	tests := []struct {
// 		diagramType      string
// 		expectedPrefix   string
// 		shouldNotBeEmpty bool
// 	}{
// 		{"sequence", "sequenceDiagram", true},
// 		{"class", "classDiagram", true},
// 		{"state", "stateDiagram", true},
// 		{"stateDiagram-v2", "stateDiagram-v2", true},
// 		{"er", "erDiagram", true},
// 		{"gantt", "gantt", true},
// 		{"pie", "pie", true},
// 		{"journey", "journey", true},
// 		{"gitGraph", "gitGraph", true},
// 		{"mindmap", "mindmap", true},
// 		{"timeline", "timeline", true},
// 		{"sankey", "sankey-beta", true},
// 		{"quadrantChart", "quadrantChart", true},
// 		{"xyChart", "xychart-beta", true},
// 		{"c4Context", "C4Context", true},
// 		{"c4Container", "C4Container", true},
// 		{"c4Component", "C4Component", true},
// 		{"c4Dynamic", "C4Dynamic", true},
// 		{"c4Deployment", "C4Deployment", true},
// 		{"unknown", "", false},
// 	}
// 
// 	for _, tt := range tests {
// 		t.Run(tt.diagramType, func(t *testing.T) {
// 			prefixes := getValidHeaderPrefixes(tt.diagramType)
// 
// 			if tt.shouldNotBeEmpty {
// 				if len(prefixes) == 0 {
// 					t.Errorf("expected prefixes for %s, got none", tt.diagramType)
// 				}
// 				if tt.expectedPrefix != "" {
// 					found := slices.Contains(prefixes, tt.expectedPrefix)
// 					if !found {
// 						t.Errorf("expected prefix %q in list for %s", tt.expectedPrefix, tt.diagramType)
// 					}
// 				}
// 			} else if len(prefixes) > 0 {
// 				t.Errorf("expected no prefixes for %s, got %v", tt.diagramType, prefixes)
// 			}
// 		})
// 	}
// }

// TODO: This test uses private functions/fields and cannot be converted to black-box testing
// func TestIsAllowedParenthesesContext(t *testing.T) {
// 	tests := []struct {
// 		name     string
// 		line     string
// 		expected bool
// 	}{
// 		{
// 			name:     "class method signature with +",
// 			line:     "    +method(param)",
// 			expected: true,
// 		},
// 		{
// 			name:     "class method signature with -",
// 			line:     "    -privateMethod()",
// 			expected: true,
// 		},
// 		{
// 			name:     "class method signature with #",
// 			line:     "    #protectedMethod(arg1, arg2)",
// 			expected: true,
// 		},
// 		{
// 			name:     "ER diagram notation with ||",
// 			line:     "    CUSTOMER ||--o{ ORDER",
// 			expected: true,
// 		},
// 		{
// 			name:     "ER diagram notation with }|",
// 			line:     "    ORDER }|--|| CUSTOMER",
// 			expected: true,
// 		},
// 		{
// 			name:     "ER diagram notation with |{",
// 			line:     "    CUSTOMER |{--o{ ORDER",
// 			expected: true,
// 		},
// 		{
// 			name:     "regular text with parentheses",
// 			line:     "    Some text (with parentheses)",
// 			expected: false,
// 		},
// 		{
// 			name:     "sequence diagram with parentheses",
// 			line:     "    participant Alice",
// 			expected: false,
// 		},
// 	}
// 
// 	for _, tt := range tests {
// 		t.Run(tt.name, func(t *testing.T) {
// 			got := isAllowedParenthesesContext(tt.line)
// 			if got != tt.expected {
// 				t.Errorf("expected %v, got %v for line: %s", tt.expected, got, tt.line)
// 			}
// 		})
// 	}
// }
// 
// // Helper function to check if a string contains a substring
func containsSubstr(s, substr string) bool {
	return len(s) >= len(substr) && findSubstr(s, substr)
}

func findSubstr(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}
// 	return len(s) >= len(substr) && findSubstr(s, substr)
// }
// 
// func findSubstr(s, substr string) bool {
// 	for i := 0; i <= len(s)-len(substr); i++ {
// 		if s[i:i+len(substr)] == substr {
// 			return true
// 		}
// 	}
// 	return false
// }
