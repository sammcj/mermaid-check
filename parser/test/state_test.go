package parser_test

import (
	"testing"

	"github.com/sammcj/mermaid-check/ast"
	"github.com/sammcj/mermaid-check/parser"
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
			name:    "empty diagram",
			source:  ``,
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

func TestStateParser_CompositeState(t *testing.T) {
	source := `stateDiagram-v2
    [*] --> First
    state First {
        [*] --> s2
        s2 --> [*]
    }
    First --> [*]`

	p := parser.NewStateParser()
	diagram, err := p.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	sd := diagram.(*ast.StateDiagram)
	var composite *ast.State
	for _, stmt := range sd.Statements {
		if s, ok := stmt.(*ast.State); ok && s.ID == "First" {
			composite = s
		}
	}

	if composite == nil {
		t.Fatal("composite state \"First\" not found in top-level statements")
	}
	if !composite.IsComposite {
		t.Error("expected First.IsComposite = true")
	}
	if len(composite.Nested) != 2 {
		t.Fatalf("expected 2 nested statements, got %d", len(composite.Nested))
	}
	if _, ok := composite.Nested[0].(*ast.StartState); !ok {
		t.Errorf("nested[0] = %T, want *ast.StartState", composite.Nested[0])
	}
	if _, ok := composite.Nested[1].(*ast.EndState); !ok {
		t.Errorf("nested[1] = %T, want *ast.EndState", composite.Nested[1])
	}
}

func TestStateParser_NestedCompositeState(t *testing.T) {
	source := `stateDiagram-v2
    state Outer {
        [*] --> Inner
        state Inner {
            [*] --> done
        }
    }`

	p := parser.NewStateParser()
	diagram, err := p.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	sd := diagram.(*ast.StateDiagram)
	if len(sd.Statements) != 1 {
		t.Fatalf("expected 1 top-level statement, got %d", len(sd.Statements))
	}
	outer, ok := sd.Statements[0].(*ast.State)
	if !ok || outer.ID != "Outer" {
		t.Fatalf("top-level = %T, want composite state Outer", sd.Statements[0])
	}

	var inner *ast.State
	for _, stmt := range outer.Nested {
		if s, ok := stmt.(*ast.State); ok && s.ID == "Inner" {
			inner = s
		}
	}
	if inner == nil {
		t.Fatal("nested composite state \"Inner\" not found")
	}
	if !inner.IsComposite || len(inner.Nested) != 1 {
		t.Errorf("Inner = {IsComposite:%v Nested:%d}, want {true 1}", inner.IsComposite, len(inner.Nested))
	}
}

func TestStateParser_SupportedTypes(t *testing.T) {
	p := parser.NewStateParser()
	types := p.SupportedTypes()
	if len(types) != 2 {
		t.Errorf("SupportedTypes() = %v, want 2 types", types)
	}
}
