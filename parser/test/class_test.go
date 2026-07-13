package parser_test

import (
	"strings"
	"testing"

	"github.com/noamsto/mermaid-check/ast"
	"github.com/noamsto/mermaid-check/parser"
)

func TestClassParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
	}{
		{
			name: "simple class",
			source: `classDiagram
    class Animal`,
			wantErr: false,
		},
		{
			name: "class with body",
			source: `classDiagram
    class Animal {
        +name
        +age
        +makeSound()
    }`,
			wantErr: false,
		},
		{
			name: "class with relationship",
			source: `classDiagram
    class Animal
    class Dog
    Animal <|-- Dog`,
			wantErr: false,
		},
		{
			name: "invalid header",
			source: `class
    class Animal`,
			wantErr: true,
		},
		{
			name: "empty diagram",
			source: ``,
			wantErr: true,
		},
	}

	p := parser.NewClassParser()

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
				classDiagram, ok := diagram.(*ast.ClassDiagram)
				if !ok {
					t.Errorf("Parse() returned wrong type: %T", diagram)
				}
				if classDiagram.Type != "class" {
					t.Errorf("Parse() diagram type = %s, want 'class'", classDiagram.Type)
				}
			}
		})
	}
}

func TestClassParser_SupportedTypes(t *testing.T) {
	p := parser.NewClassParser()
	types := p.SupportedTypes()
	if len(types) != 1 || types[0] != "class" {
		t.Errorf("SupportedTypes() = %v, want [\"class\"]", types)
	}
}

func TestClassParser_Members(t *testing.T) {
	src := "classDiagram\n" +
		"    class Animal {\n" +
		"        +String name\n" +
		"        -int age\n" +
		"        +makeSound(String kind) bool\n" +
		"        +run()\n" +
		"    }"
	d, err := parser.NewClassParser().Parse(src)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	cd, ok := d.(*ast.ClassDiagram)
	if !ok {
		t.Fatalf("Parse() = %T, want *ast.ClassDiagram", d)
	}
	var class *ast.Class
	for _, s := range cd.Statements {
		if c, ok := s.(*ast.Class); ok {
			class = c
			break
		}
	}
	if class == nil {
		t.Fatal("no class statement found")
	}

	want := []ast.ClassMember{
		{Visibility: "+", Name: "name", Type: "String", IsMethod: false},
		{Visibility: "-", Name: "age", Type: "int", IsMethod: false},
		{Visibility: "+", Name: "makeSound", Type: "bool", IsMethod: true, Parameters: []string{"String kind"}},
		{Visibility: "+", Name: "run", Type: "", IsMethod: true},
	}
	if len(class.Members) != len(want) {
		t.Fatalf("got %d members, want %d: %+v", len(class.Members), len(want), class.Members)
	}
	for i, w := range want {
		g := class.Members[i]
		if g.Visibility != w.Visibility || g.Name != w.Name || g.Type != w.Type || g.IsMethod != w.IsMethod {
			t.Errorf("member %d = %+v, want vis=%q name=%q type=%q method=%v", i, g, w.Visibility, w.Name, w.Type, w.IsMethod)
		}
		if strings.Join(g.Parameters, ",") != strings.Join(w.Parameters, ",") {
			t.Errorf("member %d params = %v, want %v", i, g.Parameters, w.Parameters)
		}
	}
}
