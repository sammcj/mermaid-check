package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestNewC4ContainerParser(t *testing.T) {
	p := parser.NewC4ContainerParser()
	if p == nil {
		t.Fatal("parser is nil")
	}
}

func TestParseC4Container(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
	}{
		{
			name: "valid container diagram",
			source: `C4Container
    title Container Diagram
    Container(app, "Web App", "React", "Frontend application")
    ContainerDb(db, "Database", "PostgreSQL", "Stores data")
    Rel(app, db, "Reads/Writes", "SQL")`,
			wantErr: false,
		},
		{
			name: "valid container with queue",
			source: `C4Container
    Container(api, "API", "Go")
    ContainerQueue(queue, "Queue", "RabbitMQ")`,
			wantErr: false,
		},
		{
			name: "invalid header",
			source: `C4Context
    Container(app, "App")`,
			wantErr: true,
		},
		{
			name: "empty diagram",
			source: `C4Container`,
			wantErr: false,
		},
	}

	p := parser.NewC4ContainerParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram, err := p.Parse(tt.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			if !tt.wantErr {
				if diagram == nil {
					t.Fatal("diagram is nil")
				}

				c4Diagram, ok := diagram.(*ast.C4Diagram)
				if !ok {
					t.Fatalf("expected *ast.C4Diagram, got %T", diagram)
				}

				if c4Diagram.DiagramType != "c4Container" {
					t.Errorf("expected DiagramType 'c4Container', got %q", c4Diagram.DiagramType)
				}
			}
		})
	}
}

func TestParseC4ContainerWithBoundaries(t *testing.T) {
	source := `C4Container
    title System Containers
    System_Boundary(boundary, "System") {
        Container(web, "Web App", "React")
        Container(api, "API", "Go")
    }
    ContainerDb(db, "DB", "PostgreSQL")`

	p := parser.NewC4ContainerParser()
	diagram, err := p.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	c4Diagram, ok := diagram.(*ast.C4Diagram)
	if !ok {
		t.Fatalf("expected *ast.C4Diagram, got %T", diagram)
	}

	if c4Diagram.Title != "System Containers" {
		t.Errorf("expected title 'System Containers', got %q", c4Diagram.Title)
	}
}

func TestParseC4ContainerWithMultipleElements(t *testing.T) {
	source := `C4Container
    Container(web, "Web", "React")
    Container(api, "API", "Go")
    Container(worker, "Worker", "Python")
    ContainerDb(db, "Database", "PostgreSQL")
    ContainerQueue(queue, "Queue", "RabbitMQ")
    Rel(web, api, "Calls")
    Rel(api, db, "Reads")
    Rel(worker, queue, "Consumes")`

	p := parser.NewC4ContainerParser()
	diagram, err := p.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	c4Diagram, ok := diagram.(*ast.C4Diagram)
	if !ok {
		t.Fatalf("expected *ast.C4Diagram, got %T", diagram)
	}

	if len(c4Diagram.Elements) < 5 {
		t.Errorf("expected at least 5 elements, got %d", len(c4Diagram.Elements))
	}
}
