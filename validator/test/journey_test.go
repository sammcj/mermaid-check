package validator_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/validator"
)

func TestValidateJourney(t *testing.T) {
	tests := []struct {
		name       string
		diagram    *ast.JourneyDiagram
		strict     bool
		wantErrors int
	}{
		{
			name: "valid journey",
			diagram: &ast.JourneyDiagram{
				Type:  "journey",
				Title: "Test Journey",
				Sections: []ast.Section{
					{
						Name: "Section 1",
						Tasks: []ast.Task{
							{
								Name:   "Task 1",
								Score:  3,
								Actors: []string{"Actor1"},
								Pos:    ast.Position{Line: 3, Column: 1},
							},
							{
								Name:   "Task 2",
								Score:  5,
								Actors: []string{"Actor1", "Actor2"},
								Pos:    ast.Position{Line: 4, Column: 1},
							},
						},
						Pos: ast.Position{Line: 2, Column: 1},
					},
				},
				Pos: ast.Position{Line: 1, Column: 1},
			},
			strict:     false,
			wantErrors: 0,
		},
		{
			name: "all valid scores 1-5",
			diagram: &ast.JourneyDiagram{
				Type: "journey",
				Sections: []ast.Section{
					{
						Name: "Scores",
						Tasks: []ast.Task{
							{Name: "Task 1", Score: 1, Actors: []string{"A"}, Pos: ast.Position{Line: 2, Column: 1}},
							{Name: "Task 2", Score: 2, Actors: []string{"A"}, Pos: ast.Position{Line: 3, Column: 1}},
							{Name: "Task 3", Score: 3, Actors: []string{"A"}, Pos: ast.Position{Line: 4, Column: 1}},
							{Name: "Task 4", Score: 4, Actors: []string{"A"}, Pos: ast.Position{Line: 5, Column: 1}},
							{Name: "Task 5", Score: 5, Actors: []string{"A"}, Pos: ast.Position{Line: 6, Column: 1}},
						},
						Pos: ast.Position{Line: 1, Column: 1},
					},
				},
				Pos: ast.Position{Line: 1, Column: 1},
			},
			strict:     false,
			wantErrors: 0,
		},
		{
			name: "invalid score - too low",
			diagram: &ast.JourneyDiagram{
				Type: "journey",
				Sections: []ast.Section{
					{
						Name: "Section 1",
						Tasks: []ast.Task{
							{
								Name:   "Task 1",
								Score:  0,
								Actors: []string{"Actor1"},
								Pos:    ast.Position{Line: 3, Column: 1},
							},
						},
						Pos: ast.Position{Line: 2, Column: 1},
					},
				},
				Pos: ast.Position{Line: 1, Column: 1},
			},
			strict:     false,
			wantErrors: 1,
		},
		{
			name: "invalid score - too high",
			diagram: &ast.JourneyDiagram{
				Type: "journey",
				Sections: []ast.Section{
					{
						Name: "Section 1",
						Tasks: []ast.Task{
							{
								Name:   "Task 1",
								Score:  6,
								Actors: []string{"Actor1"},
								Pos:    ast.Position{Line: 3, Column: 1},
							},
						},
						Pos: ast.Position{Line: 2, Column: 1},
					},
				},
				Pos: ast.Position{Line: 1, Column: 1},
			},
			strict:     false,
			wantErrors: 1,
		},
		{
			name: "multiple invalid scores",
			diagram: &ast.JourneyDiagram{
				Type: "journey",
				Sections: []ast.Section{
					{
						Name: "Section 1",
						Tasks: []ast.Task{
							{
								Name:   "Task 1",
								Score:  0,
								Actors: []string{"Actor1"},
								Pos:    ast.Position{Line: 3, Column: 1},
							},
							{
								Name:   "Task 2",
								Score:  10,
								Actors: []string{"Actor1"},
								Pos:    ast.Position{Line: 4, Column: 1},
							},
						},
						Pos: ast.Position{Line: 2, Column: 1},
					},
				},
				Pos: ast.Position{Line: 1, Column: 1},
			},
			strict:     false,
			wantErrors: 2,
		},
		{
			name: "task without actors",
			diagram: &ast.JourneyDiagram{
				Type: "journey",
				Sections: []ast.Section{
					{
						Name: "Section 1",
						Tasks: []ast.Task{
							{
								Name:   "Task 1",
								Score:  3,
								Actors: []string{},
								Pos:    ast.Position{Line: 3, Column: 1},
							},
						},
						Pos: ast.Position{Line: 2, Column: 1},
					},
				},
				Pos: ast.Position{Line: 1, Column: 1},
			},
			strict:     false,
			wantErrors: 1,
		},
		{
			name: "multiple validation errors",
			diagram: &ast.JourneyDiagram{
				Type: "journey",
				Sections: []ast.Section{
					{
						Name: "Section 1",
						Tasks: []ast.Task{
							{
								Name:   "Task 1",
								Score:  0,
								Actors: []string{},
								Pos:    ast.Position{Line: 3, Column: 1},
							},
							{
								Name:   "Task 2",
								Score:  10,
								Actors: []string{"Actor"},
								Pos:    ast.Position{Line: 4, Column: 1},
							},
						},
						Pos: ast.Position{Line: 2, Column: 1},
					},
				},
				Pos: ast.Position{Line: 1, Column: 1},
			},
			strict:     false,
			wantErrors: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.ValidateJourney(tt.diagram, tt.strict)
			if len(errors) != tt.wantErrors {
				t.Errorf("ValidateJourney() returned %d errors, want %d", len(errors), tt.wantErrors)
				for _, err := range errors {
					t.Logf("  - %s", err.Message)
				}
			}
		})
	}
}

func TestValidTaskScoresRule(t *testing.T) {
	rule := &validator.ValidTaskScoresRule{}

	tests := []struct {
		name       string
		diagram    *ast.JourneyDiagram
		wantErrors int
	}{
		{
			name: "all valid scores",
			diagram: &ast.JourneyDiagram{
				Sections: []ast.Section{
					{
						Tasks: []ast.Task{
							{Name: "T1", Score: 1, Actors: []string{"A"}, Pos: ast.Position{Line: 1, Column: 1}},
							{Name: "T2", Score: 3, Actors: []string{"A"}, Pos: ast.Position{Line: 2, Column: 1}},
							{Name: "T3", Score: 5, Actors: []string{"A"}, Pos: ast.Position{Line: 3, Column: 1}},
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "score too low",
			diagram: &ast.JourneyDiagram{
				Sections: []ast.Section{
					{
						Tasks: []ast.Task{
							{Name: "T1", Score: 0, Actors: []string{"A"}, Pos: ast.Position{Line: 1, Column: 1}},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "score too high",
			diagram: &ast.JourneyDiagram{
				Sections: []ast.Section{
					{
						Tasks: []ast.Task{
							{Name: "T1", Score: 6, Actors: []string{"A"}, Pos: ast.Position{Line: 1, Column: 1}},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "multiple invalid scores",
			diagram: &ast.JourneyDiagram{
				Sections: []ast.Section{
					{
						Tasks: []ast.Task{
							{Name: "T1", Score: -1, Actors: []string{"A"}, Pos: ast.Position{Line: 1, Column: 1}},
							{Name: "T2", Score: 100, Actors: []string{"A"}, Pos: ast.Position{Line: 2, Column: 1}},
						},
					},
				},
			},
			wantErrors: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if len(errors) != tt.wantErrors {
				t.Errorf("Validate() returned %d errors, want %d", len(errors), tt.wantErrors)
			}
		})
	}
}

func TestTasksHaveActorsRule(t *testing.T) {
	rule := &validator.TasksHaveActorsRule{}

	tests := []struct {
		name       string
		diagram    *ast.JourneyDiagram
		wantErrors int
	}{
		{
			name: "all tasks have actors",
			diagram: &ast.JourneyDiagram{
				Sections: []ast.Section{
					{
						Tasks: []ast.Task{
							{Name: "T1", Score: 3, Actors: []string{"A"}, Pos: ast.Position{Line: 1, Column: 1}},
							{Name: "T2", Score: 3, Actors: []string{"A", "B"}, Pos: ast.Position{Line: 2, Column: 1}},
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "task without actors",
			diagram: &ast.JourneyDiagram{
				Sections: []ast.Section{
					{
						Tasks: []ast.Task{
							{Name: "T1", Score: 3, Actors: []string{}, Pos: ast.Position{Line: 1, Column: 1}},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "multiple tasks without actors",
			diagram: &ast.JourneyDiagram{
				Sections: []ast.Section{
					{
						Tasks: []ast.Task{
							{Name: "T1", Score: 3, Actors: []string{}, Pos: ast.Position{Line: 1, Column: 1}},
							{Name: "T2", Score: 3, Actors: []string{"A"}, Pos: ast.Position{Line: 2, Column: 1}},
							{Name: "T3", Score: 3, Actors: []string{}, Pos: ast.Position{Line: 3, Column: 1}},
						},
					},
				},
			},
			wantErrors: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if len(errors) != tt.wantErrors {
				t.Errorf("Validate() returned %d errors, want %d", len(errors), tt.wantErrors)
			}
		})
	}
}
