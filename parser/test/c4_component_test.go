package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestNewC4ComponentParser(t *testing.T) {
	p := parser.NewC4ComponentParser()
	if p == nil {
		t.Fatal("parser is nil")
	}
}

func TestParseC4Component(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
	}{
		{
			name: "valid component diagram",
			source: `C4Component
    title Component Diagram
    Component(api, "API Controller", "Go", "Handles requests")
    ComponentDb(cache, "Cache", "Redis", "Caches data")
    Rel(api, cache, "Uses", "TCP")`,
			wantErr: false,
		},
		{
			name: "valid component with queue",
			source: `C4Component
    Component(worker, "Worker", "Go")
    ComponentQueue(tasks, "Task Queue", "Redis")`,
			wantErr: false,
		},
		{
			name: "invalid header",
			source: `C4Container
    Component(api, "API")`,
			wantErr: true,
		},
		{
			name: "empty diagram",
			source: `C4Component`,
			wantErr: false,
		},
	}

	p := parser.NewC4ComponentParser()

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

				if c4Diagram.DiagramType != "c4Component" {
					t.Errorf("expected DiagramType 'c4Component', got %q", c4Diagram.DiagramType)
				}
			}
		})
	}
}

func TestParseC4ComponentWithBoundaries(t *testing.T) {
	source := `C4Component
    title API Components
    Container_Boundary(api_boundary, "API Layer") {
        Component(controller, "Controller", "Go")
        Component(service, "Service", "Go")
    }
    ComponentDb(db, "Database", "PostgreSQL")`

	p := parser.NewC4ComponentParser()
	diagram, err := p.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	c4Diagram, ok := diagram.(*ast.C4Diagram)
	if !ok {
		t.Fatalf("expected *ast.C4Diagram, got %T", diagram)
	}

	if c4Diagram.Title != "API Components" {
		t.Errorf("expected title 'API Components', got %q", c4Diagram.Title)
	}
}

func TestParseC4ComponentWithRelationships(t *testing.T) {
	source := `C4Component
    Component(frontend, "Frontend", "React")
    Component(backend, "Backend", "Go")
    ComponentDb(db, "DB", "PostgreSQL")
    Rel(frontend, backend, "Calls", "HTTPS")
    Rel(backend, db, "Reads/Writes", "TCP")`

	p := parser.NewC4ComponentParser()
	diagram, err := p.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	c4Diagram, ok := diagram.(*ast.C4Diagram)
	if !ok {
		t.Fatalf("expected *ast.C4Diagram, got %T", diagram)
	}

	if len(c4Diagram.Elements) < 3 {
		t.Errorf("expected at least 3 elements, got %d", len(c4Diagram.Elements))
	}
}
