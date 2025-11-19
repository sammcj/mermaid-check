package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestStateParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
	}{
		{
			name: "simple state diagram",
			source: `stateDiagram
    [*] --> Still
    Still --> Moving`,
			wantErr: false,
		},
		{
			name: "state diagram v2",
			source: `stateDiagram-v2
    [*] --> Still
    Still --> [*]`,
			wantErr: false,
		},
		{
			name: "invalid header",
			source: `state
    [*] --> Still`,
			wantErr: true,
		},
		{
			name: "empty diagram",
			source: ``,
			wantErr: true,
		},
	}

	p := parser.NewStateParser()

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram, err := p.Parse(tt.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && diagram == nil {
				t.Errorf("Parse() returned nil diagram")
			}
			if !tt.wantErr {
				stateDiagram, ok := diagram.(*ast.StateDiagram)
				if !ok {
					t.Errorf("Parse() returned wrong type: %T", diagram)
				}
				if stateDiagram.Type != "state" && stateDiagram.Type != "stateDiagram-v2" {
					t.Errorf("Parse() diagram type = %s, want 'state' or 'stateDiagram-v2'", stateDiagram.Type)
				}
			}
		})
	}
}

func TestStateParser_SupportedTypes(t *testing.T) {
	p := parser.NewStateParser()
	types := p.SupportedTypes()
	if len(types) != 2 {
		t.Errorf("SupportedTypes() = %v, want 2 types", types)
	}
}
