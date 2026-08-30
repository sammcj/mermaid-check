package parser_test

import (
	"strings"
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

func TestStateParser_SupportedTypes(t *testing.T) {
	p := parser.NewStateParser()
	types := p.SupportedTypes()
	if len(types) != 2 {
		t.Errorf("SupportedTypes() = %v, want 2 types", types)
	}
}

func TestStateParser_CompositeState(t *testing.T) {
	src := "stateDiagram-v2\n" +
		"    [*] --> Active\n" +
		"    state Active {\n" +
		"        [*] --> Idle\n" +
		"        Idle --> Running: go\n" +
		"        Running --> [*]\n" +
		"    }\n" +
		"    Active --> [*]"
	d, err := parser.NewStateParser().Parse(src)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	sd, ok := d.(*ast.StateDiagram)
	if !ok {
		t.Fatalf("Parse() = %T, want *ast.StateDiagram", d)
	}
	if len(sd.Statements) != 3 {
		t.Fatalf("top-level statements = %d, want 3: %+v", len(sd.Statements), sd.Statements)
	}
	comp, ok := sd.Statements[1].(*ast.State)
	if !ok {
		t.Fatalf("statement[1] = %T, want *ast.State", sd.Statements[1])
	}
	if comp.ID != "Active" || !comp.IsComposite {
		t.Errorf("composite = {ID:%q IsComposite:%v}, want {Active true}", comp.ID, comp.IsComposite)
	}
	if len(comp.Nested) != 3 {
		t.Fatalf("nested statements = %d, want 3: %+v", len(comp.Nested), comp.Nested)
	}
	if _, ok := comp.Nested[0].(*ast.StartState); !ok {
		t.Errorf("nested[0] = %T, want *ast.StartState", comp.Nested[0])
	}
	if tr, ok := comp.Nested[1].(*ast.Transition); !ok {
		t.Errorf("nested[1] = %T, want *ast.Transition", comp.Nested[1])
	} else if tr.From != "Idle" || tr.To != "Running" || tr.Label != "go" {
		t.Errorf("nested[1] = %+v, want Idle->Running 'go'", tr)
	}
	if _, ok := comp.Nested[2].(*ast.EndState); !ok {
		t.Errorf("nested[2] = %T, want *ast.EndState", comp.Nested[2])
	}
}

func TestStateParser_UnclosedCompositeState(t *testing.T) {
	src := "stateDiagram-v2\n" +
		"    state Broken {\n" +
		"        a --> b\n" +
		"    Foo --> Bar"
	if _, err := parser.NewStateParser().Parse(src); err == nil {
		t.Fatal("Parse() error = nil, want unclosed composite state error")
	} else if !strings.Contains(err.Error(), "unclosed composite state") {
		t.Errorf("Parse() error = %q, want it to mention an unclosed composite state", err)
	}
}

func TestStateParser_UnclosedNestedCompositeState(t *testing.T) {
	src := "stateDiagram-v2\n" +
		"    state Outer {\n" +
		"        state Inner {\n" +
		"            a --> b\n" +
		"    }"
	if _, err := parser.NewStateParser().Parse(src); err == nil {
		t.Fatal("Parse() error = nil, want unclosed composite state error")
	}
}
