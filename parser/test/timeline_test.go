package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestTimelineParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
		check   func(*testing.T, ast.Diagram)
	}{
		{
			name: "simple timeline with title",
			source: `timeline
    title Project Timeline
    2024 Q1 : Planning phase
    2024 Q2 : Development`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				diagram := d.(*ast.TimelineDiagram)
				if diagram.Title != "Project Timeline" {
					t.Errorf("expected title 'Project Timeline', got %q", diagram.Title)
				}
				if len(diagram.Sections) != 1 {
					t.Fatalf("expected 1 section, got %d", len(diagram.Sections))
				}
				if len(diagram.Sections[0].Periods) != 2 {
					t.Errorf("expected 2 periods, got %d", len(diagram.Sections[0].Periods))
				}
			},
		},
		{
			name: "timeline with multiple events inline",
			source: `timeline
    1940s : ENIAC : Von Neumann architecture
    1950s : FORTRAN language : Transistor computers`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				diagram := d.(*ast.TimelineDiagram)
				if len(diagram.Sections[0].Periods) != 2 {
					t.Fatalf("expected 2 periods, got %d", len(diagram.Sections[0].Periods))
				}
				period := diagram.Sections[0].Periods[0]
				if len(period.Events) != 2 {
					t.Errorf("expected 2 events, got %d", len(period.Events))
				}
				if period.Events[0] != "ENIAC" {
					t.Errorf("expected first event 'ENIAC', got %q", period.Events[0])
				}
				if period.Events[1] != "Von Neumann architecture" {
					t.Errorf("expected second event 'Von Neumann architecture', got %q", period.Events[1])
				}
			},
		},
		{
			name: "timeline with vertical event continuation",
			source: `timeline
    1940s : ENIAC
          : Von Neumann architecture
    1950s : FORTRAN language`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				diagram := d.(*ast.TimelineDiagram)
				period := diagram.Sections[0].Periods[0]
				if len(period.Events) != 2 {
					t.Errorf("expected 2 events, got %d", len(period.Events))
				}
				if period.Events[1] != "Von Neumann architecture" {
					t.Errorf("expected second event 'Von Neumann architecture', got %q", period.Events[1])
				}
			},
		},
		{
			name: "timeline with sections",
			source: `timeline
    title History of Computing
    section Early Years
        1940s : ENIAC
        1950s : FORTRAN language
    section Modern Era
        1970s : Personal computers
        1980s : IBM PC`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				diagram := d.(*ast.TimelineDiagram)
				if len(diagram.Sections) != 2 {
					t.Fatalf("expected 2 sections, got %d", len(diagram.Sections))
				}
				if diagram.Sections[0].Name != "Early Years" {
					t.Errorf("expected first section 'Early Years', got %q", diagram.Sections[0].Name)
				}
				if diagram.Sections[1].Name != "Modern Era" {
					t.Errorf("expected second section 'Modern Era', got %q", diagram.Sections[1].Name)
				}
				if len(diagram.Sections[0].Periods) != 2 {
					t.Errorf("expected 2 periods in first section, got %d", len(diagram.Sections[0].Periods))
				}
			},
		},
		{
			name: "timeline with HTML line breaks",
			source: `timeline
    1940s : ENIAC<br>First electronic computer
    1950s : FORTRAN language`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				diagram := d.(*ast.TimelineDiagram)
				event := diagram.Sections[0].Periods[0].Events[0]
				if event != "ENIAC<br>First electronic computer" {
					t.Errorf("expected event with <br> tag, got %q", event)
				}
			},
		},
		{
			name: "timeline without title",
			source: `timeline
    2024 Q1 : Planning
    2024 Q2 : Development`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				diagram := d.(*ast.TimelineDiagram)
				if diagram.Title != "" {
					t.Errorf("expected empty title, got %q", diagram.Title)
				}
			},
		},
		{
			name: "complex timeline with mixed formats",
			source: `timeline
    title Project Milestones
    section Phase 1
        Week 1 : Research : Analysis : Planning
        Week 2 : Design
               : Prototyping
    section Phase 2
        Week 3 : Development : Testing
        Week 4 : Deployment`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				diagram := d.(*ast.TimelineDiagram)
				if len(diagram.Sections) != 2 {
					t.Fatalf("expected 2 sections, got %d", len(diagram.Sections))
				}
				// Check section 1 has periods
				if len(diagram.Sections[0].Periods) < 1 {
					t.Fatalf("expected at least 1 period in section 1, got %d", len(diagram.Sections[0].Periods))
				}
				// Week 1 should have 3 inline events
				week1 := diagram.Sections[0].Periods[0]
				if len(week1.Events) != 3 {
					t.Errorf("expected 3 events in Week 1, got %d", len(week1.Events))
				}
				// Week 2 should have 2 events (1 inline + 1 continuation)
				if len(diagram.Sections[0].Periods) > 1 {
					week2 := diagram.Sections[0].Periods[1]
					if len(week2.Events) != 2 {
						t.Errorf("expected 2 events in Week 2, got %d", len(week2.Events))
					}
				}
			},
		},
		{
			name:    "invalid header",
			source:  `timelines`,
			wantErr: true,
		},
		{
			name: "empty timeline",
			source: `timeline
    title Empty Timeline`,
			wantErr: true,
		},
		{
			name: "event continuation without period",
			source: `timeline
    : Event without period`,
			wantErr: true,
		},
		{
			name: "empty time period",
			source: `timeline
    : Event`,
			wantErr: true,
		},
		{
			name: "period without events",
			source: `timeline
    2024 :`,
			wantErr: true,
		},
		{
			name: "invalid syntax",
			source: `timeline
    some random text`,
			wantErr: true,
		},
		{
			name:    "empty source",
			source:  ``,
			wantErr: true,
		},
		{
			name: "timeline with comments",
			source: `timeline
    title Test Timeline
    %% This is a comment
    2024 : Event 1
    %% Another comment
    2025 : Event 2`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				diagram := d.(*ast.TimelineDiagram)
				if len(diagram.Sections[0].Periods) != 2 {
					t.Errorf("expected 2 periods, got %d (comments should be ignored)", len(diagram.Sections[0].Periods))
				}
			},
		},
		{
			name: "timeline with empty lines",
			source: `timeline
    title Test Timeline

    2024 : Event 1

    2025 : Event 2`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				diagram := d.(*ast.TimelineDiagram)
				if len(diagram.Sections[0].Periods) != 2 {
					t.Errorf("expected 2 periods, got %d (empty lines should be ignored)", len(diagram.Sections[0].Periods))
				}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewTimelineParser()
			diagram, err := p.Parse(tt.source)

			if tt.wantErr {
				if err == nil {
					t.Errorf("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diagram.GetType() != "timeline" {
				t.Errorf("expected type 'timeline', got %q", diagram.GetType())
			}

			if tt.check != nil {
				tt.check(t, diagram)
			}
		})
	}
}

func TestTimelineParser_SupportedTypes(t *testing.T) {
	p := parser.NewTimelineParser()
	types := p.SupportedTypes()

	if len(types) != 1 {
		t.Errorf("expected 1 supported type, got %d", len(types))
	}

	if types[0] != "timeline" {
		t.Errorf("expected supported type 'timeline', got %q", types[0])
	}
}
