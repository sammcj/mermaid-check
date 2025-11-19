package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestGanttParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, d ast.Diagram)
	}{
		{
			name: "basic gantt with title and sections",
			input: `gantt
    title Project Timeline
    dateFormat YYYY-MM-DD
    section Design
        Research : des1, 2024-01-01, 2024-01-10
        Wireframes : des2, 2024-01-11, 10d`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				t.Helper()
				gantt, ok := d.(*ast.GanttDiagram)
				if !ok {
					t.Fatal("expected *ast.GanttDiagram")
				}
				if gantt.Title != "Project Timeline" {
					t.Errorf("expected title 'Project Timeline', got %q", gantt.Title)
				}
				if gantt.DateFormat != "YYYY-MM-DD" {
					t.Errorf("expected dateFormat 'YYYY-MM-DD', got %q", gantt.DateFormat)
				}
				if len(gantt.Sections) != 1 {
					t.Fatalf("expected 1 section, got %d", len(gantt.Sections))
				}
				if gantt.Sections[0].Name != "Design" {
					t.Errorf("expected section name 'Design', got %q", gantt.Sections[0].Name)
				}
				if len(gantt.Sections[0].Tasks) != 2 {
					t.Fatalf("expected 2 tasks, got %d", len(gantt.Sections[0].Tasks))
				}
			},
		},
		{
			name: "tasks with status and dependencies",
			input: `gantt
    title Development Plan
    section Backend
        API Design : done, api1, 2024-01-01, 5d
        Implementation : active, api2, after api1, 10d
        Testing : crit, api3, after api2, 5d`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				t.Helper()
				gantt, ok := d.(*ast.GanttDiagram)
				if !ok {
					t.Fatal("expected *ast.GanttDiagram")
				}
				tasks := gantt.Sections[0].Tasks
				if tasks[0].Status != "done" {
					t.Errorf("task 0: expected status 'done', got %q", tasks[0].Status)
				}
				if tasks[0].ID != "api1" {
					t.Errorf("task 0: expected ID 'api1', got %q", tasks[0].ID)
				}
				if tasks[1].Status != "active" {
					t.Errorf("task 1: expected status 'active', got %q", tasks[1].Status)
				}
				if len(tasks[1].Dependencies) != 1 || tasks[1].Dependencies[0] != "api1" {
					t.Errorf("task 1: expected dependency 'api1', got %v", tasks[1].Dependencies)
				}
				if tasks[2].Status != "crit" {
					t.Errorf("task 2: expected status 'crit', got %q", tasks[2].Status)
				}
			},
		},
		{
			name: "multiple sections",
			input: `gantt
    dateFormat YYYY-MM-DD
    section Phase 1
        Task A : a1, 2024-01-01, 10d
    section Phase 2
        Task B : b1, after a1, 5d
    section Phase 3
        Task C : c1, after b1, 3d`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				t.Helper()
				gantt, ok := d.(*ast.GanttDiagram)
				if !ok {
					t.Fatal("expected *ast.GanttDiagram")
				}
				if len(gantt.Sections) != 3 {
					t.Fatalf("expected 3 sections, got %d", len(gantt.Sections))
				}
				if gantt.Sections[0].Name != "Phase 1" {
					t.Errorf("section 0: expected name 'Phase 1', got %q", gantt.Sections[0].Name)
				}
				if gantt.Sections[1].Name != "Phase 2" {
					t.Errorf("section 1: expected name 'Phase 2', got %q", gantt.Sections[1].Name)
				}
				if gantt.Sections[2].Name != "Phase 3" {
					t.Errorf("section 2: expected name 'Phase 3', got %q", gantt.Sections[2].Name)
				}
			},
		},
		{
			name: "excludes and todayMarker",
			input: `gantt
    title Project with Excludes
    dateFormat YYYY-MM-DD
    excludes weekends
    todayMarker off
    section Work
        Task : t1, 2024-01-01, 5d`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				t.Helper()
				gantt, ok := d.(*ast.GanttDiagram)
				if !ok {
					t.Fatal("expected *ast.GanttDiagram")
				}
				if gantt.Excludes != "weekends" {
					t.Errorf("expected excludes 'weekends', got %q", gantt.Excludes)
				}
				if gantt.TodayMarker != "off" {
					t.Errorf("expected todayMarker 'off', got %q", gantt.TodayMarker)
				}
			},
		},
		{
			name: "todayMarker with colour",
			input: `gantt
    todayMarker #FF0000
    section Tasks
        Task : t1, 2024-01-01, 1d`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				t.Helper()
				gantt, ok := d.(*ast.GanttDiagram)
				if !ok {
					t.Fatal("expected *ast.GanttDiagram")
				}
				if gantt.TodayMarker != "#FF0000" {
					t.Errorf("expected todayMarker '#FF0000', got %q", gantt.TodayMarker)
				}
			},
		},
		{
			name: "axisFormat",
			input: `gantt
    axisFormat %d/%m
    section Work
        Task : t1, 2024-01-01, 5d`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				t.Helper()
				gantt, ok := d.(*ast.GanttDiagram)
				if !ok {
					t.Fatal("expected *ast.GanttDiagram")
				}
				if gantt.AxisFormat != "%d/%m" {
					t.Errorf("expected axisFormat '%%d/%%m', got %q", gantt.AxisFormat)
				}
			},
		},
		{
			name: "milestone status",
			input: `gantt
    section Milestones
        Release v1.0 : milestone, m1, 2024-01-31, 0d`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				t.Helper()
				gantt, ok := d.(*ast.GanttDiagram)
				if !ok {
					t.Fatal("expected *ast.GanttDiagram")
				}
				if gantt.Sections[0].Tasks[0].Status != "milestone" {
					t.Errorf("expected status 'milestone', got %q", gantt.Sections[0].Tasks[0].Status)
				}
			},
		},
		{
			name: "task with only start and duration (no ID)",
			input: `gantt
    section Simple
        Simple Task : 2024-01-01, 10d`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				t.Helper()
				gantt, ok := d.(*ast.GanttDiagram)
				if !ok {
					t.Fatal("expected *ast.GanttDiagram")
				}
				task := gantt.Sections[0].Tasks[0]
				if task.Name != "Simple Task" {
					t.Errorf("expected name 'Simple Task', got %q", task.Name)
				}
				if task.ID != "" {
					t.Errorf("expected empty ID, got %q", task.ID)
				}
				if task.StartDate != "2024-01-01" {
					t.Errorf("expected start '2024-01-01', got %q", task.StartDate)
				}
				if task.EndDate != "10d" {
					t.Errorf("expected end '10d', got %q", task.EndDate)
				}
			},
		},
		{
			name: "comments ignored",
			input: `gantt
    %% This is a comment
    title Test
    section Work
        %% Another comment
        Task : t1, 2024-01-01, 1d`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				t.Helper()
				gantt, ok := d.(*ast.GanttDiagram)
				if !ok {
					t.Fatal("expected *ast.GanttDiagram")
				}
				if gantt.Title != "Test" {
					t.Errorf("expected title 'Test', got %q", gantt.Title)
				}
			},
		},
		{
			name:    "empty gantt",
			input:   "gantt",
			wantErr: true,
		},
		{
			name: "task outside section",
			input: `gantt
    Task : t1, 2024-01-01, 1d`,
			wantErr: true,
		},
		{
			name: "invalid header",
			input: `graph TD
    A --> B`,
			wantErr: true,
		},
		{
			name: "task missing parameters",
			input: `gantt
    section Work
        Task : t1`,
			wantErr: true,
		},
		{
			name: "invalid syntax",
			input: `gantt
    section Work
        invalid line without colon`,
			wantErr: true,
		},
	}

	p := parser.NewGanttParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram, err := p.Parse(tt.input)
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

func TestGanttParser_SupportedTypes(t *testing.T) {
	p := parser.NewGanttParser()
	types := p.SupportedTypes()
	if len(types) != 1 || types[0] != "gantt" {
		t.Errorf("expected supported types [gantt], got %v", types)
	}
}
