package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestNewC4DynamicParser(t *testing.T) {
	p := parser.NewC4DynamicParser()
	if p == nil {
		t.Fatal("parser is nil")
	}
}

func TestParseC4Dynamic(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
	}{
		{
			name: "valid dynamic diagram",
			source: `C4Dynamic
    title Dynamic Diagram
    Container(web, "Web App")
    Container(api, "API")
    ContainerDb(db, "Database")
    Rel(web, api, "Calls", "HTTPS")
    Rel(api, db, "Queries", "SQL")`,
			wantErr: false,
		},
		{
			name: "valid dynamic with bidirectional",
			source: `C4Dynamic
    Container(a, "Service A")
    Container(b, "Service B")
    BiRel(a, b, "Communicates")`,
			wantErr: false,
		},
		{
			name: "invalid header",
			source: `C4Context
    Rel(a, b, "calls")`,
			wantErr: true,
		},
		{
			name: "empty diagram",
			source: `C4Dynamic`,
			wantErr: false,
		},
	}

	p := parser.NewC4DynamicParser()

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

				if c4Diagram.DiagramType != "c4Dynamic" {
					t.Errorf("expected DiagramType 'c4Dynamic', got %q", c4Diagram.DiagramType)
				}
			}
		})
	}
}
