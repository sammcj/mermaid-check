package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestNewC4ContextParser(t *testing.T) {
	p := parser.NewC4ContextParser()
	if p == nil {
		t.Fatal("parser is nil")
	}
}

func TestParseC4Context(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
	}{
		{
			name: "valid context with person and system",
			source: `C4Context
    title System Context Diagram
    Person(user, "User", "A user of the system")
    System(app, "Application", "The main app")
    Rel(user, app, "Uses")`,
			wantErr: false,
		},
		{
			name: "valid context with external system",
			source: `C4Context
    Person(admin, "Admin")
    System_Ext(external, "External System")`,
			wantErr: false,
		},
		{
			name: "invalid header",
			source: `flowchart TD
    A --> B`,
			wantErr: true,
		},
		{
			name: "empty diagram",
			source: `C4Context`,
			wantErr: false,
		},
	}

	p := parser.NewC4ContextParser()

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

				if c4Diagram.DiagramType != "c4Context" {
					t.Errorf("expected DiagramType 'c4Context', got %q", c4Diagram.DiagramType)
				}
			}
		})
	}
}
