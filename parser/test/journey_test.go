package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestJourneyParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
		check   func(*testing.T, ast.Diagram)
	}{
		{
			name: "simple journey with one section",
			source: `journey
    title My Shopping Journey
    section Browse Products
        View catalogue: 5: Customer
        Search for item: 4: Customer`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				journey, ok := d.(*ast.JourneyDiagram)
				if !ok {
					t.Fatalf("expected *ast.JourneyDiagram, got %T", d)
				}
				if journey.Title != "My Shopping Journey" {
					t.Errorf("expected title 'My Shopping Journey', got %q", journey.Title)
				}
				if len(journey.Sections) != 1 {
					t.Fatalf("expected 1 section, got %d", len(journey.Sections))
				}
				section := journey.Sections[0]
				if section.Name != "Browse Products" {
					t.Errorf("expected section name 'Browse Products', got %q", section.Name)
				}
				if len(section.Tasks) != 2 {
					t.Fatalf("expected 2 tasks, got %d", len(section.Tasks))
				}
				task1 := section.Tasks[0]
				if task1.Name != "View catalogue" {
					t.Errorf("expected task name 'View catalogue', got %q", task1.Name)
				}
				if task1.Score != 5 {
					t.Errorf("expected score 5, got %d", task1.Score)
				}
				if len(task1.Actors) != 1 || task1.Actors[0] != "Customer" {
					t.Errorf("expected actors [Customer], got %v", task1.Actors)
				}
			},
		},
		{
			name: "journey with multiple sections and actors",
			source: `journey
    title Shopping Experience
    section Browse
        View items: 5: Customer
    section Purchase
        Add to cart: 4: Customer
        Checkout: 3: Customer, System
    section Delivery
        Ship order: 4: Warehouse, Courier`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				journey, ok := d.(*ast.JourneyDiagram)
				if !ok {
					t.Fatalf("expected *ast.JourneyDiagram, got %T", d)
				}
				if len(journey.Sections) != 3 {
					t.Fatalf("expected 3 sections, got %d", len(journey.Sections))
				}
				// Check multi-actor task
				purchaseSection := journey.Sections[1]
				checkoutTask := purchaseSection.Tasks[1]
				if len(checkoutTask.Actors) != 2 {
					t.Fatalf("expected 2 actors, got %d", len(checkoutTask.Actors))
				}
				if checkoutTask.Actors[0] != "Customer" || checkoutTask.Actors[1] != "System" {
					t.Errorf("expected actors [Customer, System], got %v", checkoutTask.Actors)
				}
			},
		},
		{
			name: "journey without title",
			source: `journey
    section Main
        Task one: 3: Actor`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				journey, ok := d.(*ast.JourneyDiagram)
				if !ok {
					t.Fatalf("expected *ast.JourneyDiagram, got %T", d)
				}
				if journey.Title != "" {
					t.Errorf("expected no title, got %q", journey.Title)
				}
				if len(journey.Sections) != 1 {
					t.Fatalf("expected 1 section, got %d", len(journey.Sections))
				}
			},
		},
		{
			name: "journey with comments",
			source: `journey
    title Test Journey
    %% This is a comment
    section Section One
        %% Another comment
        Task A: 4: Actor1
        Task B: 3: Actor2`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				journey, ok := d.(*ast.JourneyDiagram)
				if !ok {
					t.Fatalf("expected *ast.JourneyDiagram, got %T", d)
				}
				if len(journey.Sections) != 1 {
					t.Fatalf("expected 1 section, got %d", len(journey.Sections))
				}
				if len(journey.Sections[0].Tasks) != 2 {
					t.Errorf("expected 2 tasks, got %d", len(journey.Sections[0].Tasks))
				}
			},
		},
		{
			name: "journey with empty lines",
			source: `journey
    title Test

    section Section One

        Task A: 4: Actor1

        Task B: 3: Actor2`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				journey, ok := d.(*ast.JourneyDiagram)
				if !ok {
					t.Fatalf("expected *ast.JourneyDiagram, got %T", d)
				}
				if len(journey.Sections[0].Tasks) != 2 {
					t.Errorf("expected 2 tasks, got %d", len(journey.Sections[0].Tasks))
				}
			},
		},
		{
			name: "all valid scores 1-5",
			source: `journey
    section Scores
        Task 1: 1: Actor
        Task 2: 2: Actor
        Task 3: 3: Actor
        Task 4: 4: Actor
        Task 5: 5: Actor`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				journey, ok := d.(*ast.JourneyDiagram)
				if !ok {
					t.Fatalf("expected *ast.JourneyDiagram, got %T", d)
				}
				tasks := journey.Sections[0].Tasks
				if len(tasks) != 5 {
					t.Fatalf("expected 5 tasks, got %d", len(tasks))
				}
				for i, task := range tasks {
					expectedScore := i + 1
					if task.Score != expectedScore {
						t.Errorf("task %d: expected score %d, got %d", i, expectedScore, task.Score)
					}
				}
			},
		},
		{
			name:    "invalid header",
			source:  "notjourney\n",
			wantErr: true,
		},
		{
			name:    "empty diagram",
			source:  "journey\n",
			wantErr: true,
		},
		{
			name: "task outside section",
			source: `journey
    title Test
    Task without section: 3: Actor`,
			wantErr: true,
		},
		{
			name: "score too low",
			source: `journey
    section Test
        Task: 0: Actor`,
			wantErr: true,
		},
		{
			name: "score too high",
			source: `journey
    section Test
        Task: 6: Actor`,
			wantErr: true,
		},
		{
			name: "invalid score format",
			source: `journey
    section Test
        Task: abc: Actor`,
			wantErr: true,
		},
		{
			name: "task without actors",
			source: `journey
    section Test
        Task: 3:`,
			wantErr: true,
		},
		{
			name: "invalid task format",
			source: `journey
    section Test
        Invalid task without colons`,
			wantErr: true,
		},
		{
			name: "task with whitespace in actors",
			source: `journey
    section Test
        Task: 3:  Actor One , Actor Two  `,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				journey, ok := d.(*ast.JourneyDiagram)
				if !ok {
					t.Fatalf("expected *ast.JourneyDiagram, got %T", d)
				}
				task := journey.Sections[0].Tasks[0]
				if len(task.Actors) != 2 {
					t.Fatalf("expected 2 actors, got %d", len(task.Actors))
				}
				if task.Actors[0] != "Actor One" || task.Actors[1] != "Actor Two" {
					t.Errorf("expected trimmed actors, got %v", task.Actors)
				}
			},
		},
	}

	p := parser.NewJourneyParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram, err := p.Parse(tt.source)
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

func TestJourneyParser_SupportedTypes(t *testing.T) {
	p := parser.NewJourneyParser()
	types := p.SupportedTypes()
	if len(types) != 1 || types[0] != "journey" {
		t.Errorf("expected [journey], got %v", types)
	}
}
