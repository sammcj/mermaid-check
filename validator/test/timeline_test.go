package validator_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/validator"
)

func TestValidateTimeline(t *testing.T) {
	tests := []struct {
		name       string
		diagram    *ast.TimelineDiagram
		strict     bool
		wantErrors int
	}{
		{
			name: "valid timeline",
			diagram: &ast.TimelineDiagram{
				Type:  "timeline",
				Title: "Test Timeline",
				Sections: []ast.TimelineSection{
					{
						Name: "Section 1",
						Periods: []ast.TimelinePeriod{
							{
								TimePeriod: "2024",
								Events:     []string{"Event 1", "Event 2"},
								Pos:        ast.Position{Line: 2, Column: 1},
							},
						},
					},
				},
			},
			strict:     false,
			wantErrors: 0,
		},
		{
			name: "period without events",
			diagram: &ast.TimelineDiagram{
				Type: "timeline",
				Sections: []ast.TimelineSection{
					{
						Name: "Section 1",
						Periods: []ast.TimelinePeriod{
							{
								TimePeriod: "2024",
								Events:     []string{},
								Pos:        ast.Position{Line: 2, Column: 1},
							},
						},
					},
				},
			},
			strict:     false,
			wantErrors: 1,
		},
		{
			name: "empty period name",
			diagram: &ast.TimelineDiagram{
				Type: "timeline",
				Sections: []ast.TimelineSection{
					{
						Name: "Section 1",
						Periods: []ast.TimelinePeriod{
							{
								TimePeriod: "",
								Events:     []string{"Event 1"},
								Pos:        ast.Position{Line: 2, Column: 1},
							},
						},
					},
				},
			},
			strict:     false,
			wantErrors: 1,
		},
		{
			name: "empty event in period",
			diagram: &ast.TimelineDiagram{
				Type: "timeline",
				Sections: []ast.TimelineSection{
					{
						Name: "Section 1",
						Periods: []ast.TimelinePeriod{
							{
								TimePeriod: "2024",
								Events:     []string{"Event 1", "", "Event 2"},
								Pos:        ast.Position{Line: 2, Column: 1},
							},
						},
					},
				},
			},
			strict:     false,
			wantErrors: 1,
		},
		{
			name: "multiple validation errors",
			diagram: &ast.TimelineDiagram{
				Type: "timeline",
				Sections: []ast.TimelineSection{
					{
						Name: "Section 1",
						Periods: []ast.TimelinePeriod{
							{
								TimePeriod: "",
								Events:     []string{},
								Pos:        ast.Position{Line: 2, Column: 1},
							},
							{
								TimePeriod: "2024",
								Events:     []string{"", "Event"},
								Pos:        ast.Position{Line: 3, Column: 1},
							},
						},
					},
				},
			},
			strict:     false,
			wantErrors: 3, // empty period name, no events, empty event
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.ValidateTimeline(tt.diagram, tt.strict)

			if len(errors) != tt.wantErrors {
				t.Errorf("expected %d errors, got %d", tt.wantErrors, len(errors))
				for _, err := range errors {
					t.Logf("  - %s", err.Message)
				}
			}
		})
	}
}

func TestPeriodsHaveEventsRule(t *testing.T) {
	tests := []struct {
		name       string
		diagram    *ast.TimelineDiagram
		wantErrors int
	}{
		{
			name: "all periods have events",
			diagram: &ast.TimelineDiagram{
				Sections: []ast.TimelineSection{
					{
						Periods: []ast.TimelinePeriod{
							{
								TimePeriod: "2024",
								Events:     []string{"Event 1"},
								Pos:        ast.Position{Line: 2, Column: 1},
							},
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "period without events",
			diagram: &ast.TimelineDiagram{
				Sections: []ast.TimelineSection{
					{
						Periods: []ast.TimelinePeriod{
							{
								TimePeriod: "2024",
								Events:     []string{},
								Pos:        ast.Position{Line: 2, Column: 1},
							},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "multiple periods, some without events",
			diagram: &ast.TimelineDiagram{
				Sections: []ast.TimelineSection{
					{
						Periods: []ast.TimelinePeriod{
							{
								TimePeriod: "2024",
								Events:     []string{"Event 1"},
								Pos:        ast.Position{Line: 2, Column: 1},
							},
							{
								TimePeriod: "2025",
								Events:     []string{},
								Pos:        ast.Position{Line: 3, Column: 1},
							},
							{
								TimePeriod: "2026",
								Events:     []string{},
								Pos:        ast.Position{Line: 4, Column: 1},
							},
						},
					},
				},
			},
			wantErrors: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &validator.PeriodsHaveEventsRule{}
			errors := rule.Validate(tt.diagram)

			if len(errors) != tt.wantErrors {
				t.Errorf("expected %d errors, got %d", tt.wantErrors, len(errors))
			}

			for _, err := range errors {
				if err.Severity != validator.SeverityError {
					t.Errorf("expected severity Error, got %s", err.Severity)
				}
			}
		})
	}
}

func TestNoEmptyPeriodsRule(t *testing.T) {
	tests := []struct {
		name       string
		diagram    *ast.TimelineDiagram
		wantErrors int
	}{
		{
			name: "valid periods and events",
			diagram: &ast.TimelineDiagram{
				Sections: []ast.TimelineSection{
					{
						Periods: []ast.TimelinePeriod{
							{
								TimePeriod: "2024",
								Events:     []string{"Event 1", "Event 2"},
								Pos:        ast.Position{Line: 2, Column: 1},
							},
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "empty period name",
			diagram: &ast.TimelineDiagram{
				Sections: []ast.TimelineSection{
					{
						Periods: []ast.TimelinePeriod{
							{
								TimePeriod: "",
								Events:     []string{"Event 1"},
								Pos:        ast.Position{Line: 2, Column: 1},
							},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "empty event",
			diagram: &ast.TimelineDiagram{
				Sections: []ast.TimelineSection{
					{
						Periods: []ast.TimelinePeriod{
							{
								TimePeriod: "2024",
								Events:     []string{"Event 1", "", "Event 2"},
								Pos:        ast.Position{Line: 2, Column: 1},
							},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "multiple empty events",
			diagram: &ast.TimelineDiagram{
				Sections: []ast.TimelineSection{
					{
						Periods: []ast.TimelinePeriod{
							{
								TimePeriod: "2024",
								Events:     []string{"", "", "Event"},
								Pos:        ast.Position{Line: 2, Column: 1},
							},
						},
					},
				},
			},
			wantErrors: 2,
		},
		{
			name: "empty period and empty event",
			diagram: &ast.TimelineDiagram{
				Sections: []ast.TimelineSection{
					{
						Periods: []ast.TimelinePeriod{
							{
								TimePeriod: "",
								Events:     []string{"", "Event"},
								Pos:        ast.Position{Line: 2, Column: 1},
							},
						},
					},
				},
			},
			wantErrors: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rule := &validator.NoEmptyPeriodsRule{}
			errors := rule.Validate(tt.diagram)

			if len(errors) != tt.wantErrors {
				t.Errorf("expected %d errors, got %d", tt.wantErrors, len(errors))
				for _, err := range errors {
					t.Logf("  - %s", err.Message)
				}
			}

			for _, err := range errors {
				if err.Severity != validator.SeverityError {
					t.Errorf("expected severity Error, got %s", err.Severity)
				}
			}
		})
	}
}

func TestTimelineDefaultRules(t *testing.T) {
	rules := validator.TimelineDefaultRules()

	expectedRules := 2
	if len(rules) != expectedRules {
		t.Errorf("expected %d default rules, got %d", expectedRules, len(rules))
	}
}

func TestTimelineStrictRules(t *testing.T) {
	rules := validator.TimelineStrictRules()

	minExpectedRules := 2
	if len(rules) < minExpectedRules {
		t.Errorf("expected at least %d strict rules, got %d", minExpectedRules, len(rules))
	}
}
