package validator_test

import (	"testing"

	"github.com/sammcj/go-mermaid/ast"

	"github.com/sammcj/go-mermaid/validator"
)

func TestValidateGantt(t *testing.T) {
	tests := []struct {
		name       string
		diagram    *ast.GanttDiagram
		strict     bool
		wantErrors int
	}{
		{
			name: "valid gantt diagram",
			diagram: &ast.GanttDiagram{
				Type:       "gantt",
				Title:      "Project",
				DateFormat: "YYYY-MM-DD",
				Sections: []ast.GanttSection{
					{
						Name: "Phase 1",
						Tasks: []ast.GanttTask{
							{
								Name:      "Task A",
								ID:        "a1",
								Status:    "done",
								StartDate: "2024-01-01",
								EndDate:   "10d",
								Pos:       ast.Position{Line: 3, Column: 1},
							},
							{
								Name:         "Task B",
								ID:           "b1",
								Status:       "active",
								Dependencies: []string{"a1"},
								StartDate:    "after a1",
								EndDate:      "5d",
								Pos:          ast.Position{Line: 4, Column: 1},
							},
						},
					},
				},
			},
			strict:     false,
			wantErrors: 0,
		},
		{
			name: "duplicate task IDs",
			diagram: &ast.GanttDiagram{
				Type: "gantt",
				Sections: []ast.GanttSection{
					{
						Name: "Section 1",
						Tasks: []ast.GanttTask{
							{
								Name: "Task A",
								ID:   "t1",
								Pos:  ast.Position{Line: 3, Column: 1},
							},
							{
								Name: "Task B",
								ID:   "t1",
								Pos:  ast.Position{Line: 4, Column: 1},
							},
						},
					},
				},
			},
			strict:     false,
			wantErrors: 1,
		},
		{
			name: "undefined task reference",
			diagram: &ast.GanttDiagram{
				Type: "gantt",
				Sections: []ast.GanttSection{
					{
						Name: "Section 1",
						Tasks: []ast.GanttTask{
							{
								Name:         "Task A",
								ID:           "a1",
								Dependencies: []string{"undefined_task"},
								Pos:          ast.Position{Line: 3, Column: 1},
							},
						},
					},
				},
			},
			strict:     false,
			wantErrors: 1,
		},
		{
			name: "invalid task status",
			diagram: &ast.GanttDiagram{
				Type: "gantt",
				Sections: []ast.GanttSection{
					{
						Name: "Section 1",
						Tasks: []ast.GanttTask{
							{
								Name:   "Task A",
								Status: "invalid_status",
								Pos:    ast.Position{Line: 3, Column: 1},
							},
						},
					},
				},
			},
			strict:     false,
			wantErrors: 1,
		},
		{
			name: "valid task statuses",
			diagram: &ast.GanttDiagram{
				Type: "gantt",
				Sections: []ast.GanttSection{
					{
						Name: "Section 1",
						Tasks: []ast.GanttTask{
							{Name: "Task 1", Status: "done", Pos: ast.Position{Line: 3, Column: 1}},
							{Name: "Task 2", Status: "active", Pos: ast.Position{Line: 4, Column: 1}},
							{Name: "Task 3", Status: "crit", Pos: ast.Position{Line: 5, Column: 1}},
							{Name: "Task 4", Status: "milestone", Pos: ast.Position{Line: 6, Column: 1}},
						},
					},
				},
			},
			strict:     false,
			wantErrors: 0,
		},
		{
			name: "valid date formats",
			diagram: &ast.GanttDiagram{
				Type:       "gantt",
				DateFormat: "YYYY-MM-DD",
				Sections: []ast.GanttSection{
					{
						Name: "Section 1",
						Tasks: []ast.GanttTask{
							{Name: "Task", Pos: ast.Position{Line: 3, Column: 1}},
						},
					},
				},
			},
			strict:     false,
			wantErrors: 0,
		},
		{
			name: "multiple dependencies",
			diagram: &ast.GanttDiagram{
				Type: "gantt",
				Sections: []ast.GanttSection{
					{
						Name: "Section 1",
						Tasks: []ast.GanttTask{
							{
								Name: "Task A",
								ID:   "a1",
								Pos:  ast.Position{Line: 3, Column: 1},
							},
							{
								Name: "Task B",
								ID:   "b1",
								Pos:  ast.Position{Line: 4, Column: 1},
							},
							{
								Name:         "Task C",
								ID:           "c1",
								Dependencies: []string{"a1", "b1"},
								Pos:          ast.Position{Line: 5, Column: 1},
							},
						},
					},
				},
			},
			strict:     false,
			wantErrors: 0,
		},
		{
			name: "tasks without IDs are valid",
			diagram: &ast.GanttDiagram{
				Type: "gantt",
				Sections: []ast.GanttSection{
					{
						Name: "Section 1",
						Tasks: []ast.GanttTask{
							{
								Name: "Task A",
								Pos:  ast.Position{Line: 3, Column: 1},
							},
							{
								Name: "Task B",
								Pos:  ast.Position{Line: 4, Column: 1},
							},
						},
					},
				},
			},
			strict:     false,
			wantErrors: 0,
		},
		{
			name: "duplicate IDs across sections",
			diagram: &ast.GanttDiagram{
				Type: "gantt",
				Sections: []ast.GanttSection{
					{
						Name: "Section 1",
						Tasks: []ast.GanttTask{
							{
								Name: "Task A",
								ID:   "t1",
								Pos:  ast.Position{Line: 3, Column: 1},
							},
						},
					},
					{
						Name: "Section 2",
						Tasks: []ast.GanttTask{
							{
								Name: "Task B",
								ID:   "t1",
								Pos:  ast.Position{Line: 6, Column: 1},
							},
						},
					},
				},
			},
			strict:     false,
			wantErrors: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.ValidateGantt(tt.diagram, tt.strict)
			if len(errors) != tt.wantErrors {
				t.Errorf("validator.ValidateGantt() returned %d errors, want %d", len(errors), tt.wantErrors)
				for _, err := range errors {
					t.Logf("  - %s", err.Message)
				}
			}
		})
	}
}

func TestNoDuplicateTaskIDsRule(t *testing.T) {
	rule := &validator.NoDuplicateTaskIDsRule{}

	tests := []struct {
		name       string
		diagram    *ast.GanttDiagram
		wantErrors int
	}{
		{
			name: "no duplicates",
			diagram: &ast.GanttDiagram{
				Sections: []ast.GanttSection{
					{
						Tasks: []ast.GanttTask{
							{ID: "a1", Pos: ast.Position{Line: 1, Column: 1}},
							{ID: "a2", Pos: ast.Position{Line: 2, Column: 1}},
							{ID: "a3", Pos: ast.Position{Line: 3, Column: 1}},
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "with duplicates",
			diagram: &ast.GanttDiagram{
				Sections: []ast.GanttSection{
					{
						Tasks: []ast.GanttTask{
							{ID: "a1", Pos: ast.Position{Line: 1, Column: 1}},
							{ID: "a1", Pos: ast.Position{Line: 2, Column: 1}},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "empty IDs ignored",
			diagram: &ast.GanttDiagram{
				Sections: []ast.GanttSection{
					{
						Tasks: []ast.GanttTask{
							{ID: "", Pos: ast.Position{Line: 1, Column: 1}},
							{ID: "", Pos: ast.Position{Line: 2, Column: 1}},
							{ID: "a1", Pos: ast.Position{Line: 3, Column: 1}},
						},
					},
				},
			},
			wantErrors: 0,
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

func TestValidTaskReferencesRule(t *testing.T) {
	rule := &validator.ValidTaskReferencesRule{}

	tests := []struct {
		name       string
		diagram    *ast.GanttDiagram
		wantErrors int
	}{
		{
			name: "all references valid",
			diagram: &ast.GanttDiagram{
				Sections: []ast.GanttSection{
					{
						Tasks: []ast.GanttTask{
							{ID: "a1", Pos: ast.Position{Line: 1, Column: 1}},
							{ID: "a2", Dependencies: []string{"a1"}, Pos: ast.Position{Line: 2, Column: 1}},
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "undefined reference",
			diagram: &ast.GanttDiagram{
				Sections: []ast.GanttSection{
					{
						Tasks: []ast.GanttTask{
							{ID: "a1", Dependencies: []string{"undefined"}, Pos: ast.Position{Line: 1, Column: 1}},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "multiple undefined references",
			diagram: &ast.GanttDiagram{
				Sections: []ast.GanttSection{
					{
						Tasks: []ast.GanttTask{
							{ID: "a1", Pos: ast.Position{Line: 1, Column: 1}},
							{ID: "a2", Dependencies: []string{"undefined1", "undefined2"}, Pos: ast.Position{Line: 2, Column: 1}},
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

func TestValidTaskStatusRule(t *testing.T) {
	rule := &validator.ValidTaskStatusRule{}

	tests := []struct {
		name       string
		diagram    *ast.GanttDiagram
		wantErrors int
	}{
		{
			name: "all valid statuses",
			diagram: &ast.GanttDiagram{
				Sections: []ast.GanttSection{
					{
						Tasks: []ast.GanttTask{
							{Status: "done", Pos: ast.Position{Line: 1, Column: 1}},
							{Status: "active", Pos: ast.Position{Line: 2, Column: 1}},
							{Status: "crit", Pos: ast.Position{Line: 3, Column: 1}},
							{Status: "milestone", Pos: ast.Position{Line: 4, Column: 1}},
						},
					},
				},
			},
			wantErrors: 0,
		},
		{
			name: "invalid status",
			diagram: &ast.GanttDiagram{
				Sections: []ast.GanttSection{
					{
						Tasks: []ast.GanttTask{
							{Status: "invalid", Pos: ast.Position{Line: 1, Column: 1}},
						},
					},
				},
			},
			wantErrors: 1,
		},
		{
			name: "empty status is valid",
			diagram: &ast.GanttDiagram{
				Sections: []ast.GanttSection{
					{
						Tasks: []ast.GanttTask{
							{Status: "", Pos: ast.Position{Line: 1, Column: 1}},
						},
					},
				},
			},
			wantErrors: 0,
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

func TestValidDateFormatRule(t *testing.T) {
	rule := &validator.ValidDateFormatRule{}

	tests := []struct {
		name       string
		diagram    *ast.GanttDiagram
		wantErrors int
	}{
		{
			name: "valid common formats",
			diagram: &ast.GanttDiagram{
				DateFormat: "YYYY-MM-DD",
				Pos:        ast.Position{Line: 1, Column: 1},
			},
			wantErrors: 0,
		},
		{
			name: "valid DD-MM-YYYY",
			diagram: &ast.GanttDiagram{
				DateFormat: "DD-MM-YYYY",
				Pos:        ast.Position{Line: 1, Column: 1},
			},
			wantErrors: 0,
		},
		{
			name: "valid with time",
			diagram: &ast.GanttDiagram{
				DateFormat: "YYYY-MM-DD HH:mm:ss",
				Pos:        ast.Position{Line: 1, Column: 1},
			},
			wantErrors: 0,
		},
		{
			name: "default format is valid",
			diagram: &ast.GanttDiagram{
				DateFormat: "YYYY-MM-DD",
				Pos:        ast.Position{Line: 1, Column: 1},
			},
			wantErrors: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if len(errors) != tt.wantErrors {
				t.Errorf("Validate() returned %d errors, want %d", len(errors), tt.wantErrors)
				for _, err := range errors {
					t.Logf("  - %s", err.Message)
				}
			}
		})
	}
}
