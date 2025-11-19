package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestNewC4DeploymentParser(t *testing.T) {
	p := parser.NewC4DeploymentParser()
	if p == nil {
		t.Fatal("parser is nil")
	}
}

func TestParseC4Deployment(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
	}{
		{
			name: "valid deployment diagram",
			source: `C4Deployment
    title Deployment Diagram
    Deployment_Node(aws, "AWS Cloud", "Cloud Platform") {
        Container(app, "Web App", "React")
        ContainerDb(db, "Database", "PostgreSQL")
    }`,
			wantErr: false,
		},
		{
			name: "valid deployment with node",
			source: `C4Deployment
    Deployment_Node(cloud, "Cloud Provider") {
        Node(server, "Server", "Ubuntu")
    }`,
			wantErr: false,
		},
		{
			name: "invalid header",
			source: `C4Context
    Deployment_Node(aws, "AWS")`,
			wantErr: true,
		},
		{
			name: "empty diagram",
			source: `C4Deployment`,
			wantErr: false,
		},
	}

	p := parser.NewC4DeploymentParser()

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

				if c4Diagram.DiagramType != "c4Deployment" {
					t.Errorf("expected DiagramType 'c4Deployment', got %q", c4Diagram.DiagramType)
				}
			}
		})
	}
}
