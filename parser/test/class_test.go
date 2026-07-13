package parser_test

import (
	"testing"

	"github.com/sammcj/mermaid-check/ast"
	"github.com/sammcj/mermaid-check/parser"
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
			name:    "empty diagram",
			source:  ``,
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

func TestClassParser_Relationships(t *testing.T) {
	source := `classDiagram
    Animal <|-- Dog
    Duck ..|> Flyer
    Car *-- Wheel
    Company o-- Employee
    Order --> Customer
    Client ..> Service`

	p := parser.NewClassParser()
	diagram, err := p.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	cd := diagram.(*ast.ClassDiagram)
	var rels []*ast.Relationship
	for _, stmt := range cd.Statements {
		if r, ok := stmt.(*ast.Relationship); ok {
			rels = append(rels, r)
		}
	}

	want := []struct {
		from, to, typ string
	}{
		{"Animal", "Dog", "inheritance"},
		{"Duck", "Flyer", "realization"},
		{"Car", "Wheel", "composition"},
		{"Company", "Employee", "aggregation"},
		{"Order", "Customer", "association"},
		{"Client", "Service", "dependency"},
	}

	if len(rels) != len(want) {
		t.Fatalf("parsed %d relationships, want %d", len(rels), len(want))
	}

	for i, w := range want {
		if rels[i].From != w.from || rels[i].To != w.to || rels[i].Type != w.typ {
			t.Errorf("relationship %d = {From:%q To:%q Type:%q}, want {From:%q To:%q Type:%q}",
				i, rels[i].From, rels[i].To, rels[i].Type, w.from, w.to, w.typ)
		}
	}
}

func TestClassParser_MemberOrdering(t *testing.T) {
	source := `classDiagram
    class Animal {
        +String name
        +int age
        +makeSound()
        +area() float
        +legs
    }`

	p := parser.NewClassParser()
	diagram, err := p.Parse(source)
	if err != nil {
		t.Fatalf("Parse() error = %v", err)
	}

	cd := diagram.(*ast.ClassDiagram)
	var class *ast.Class
	for _, stmt := range cd.Statements {
		if c, ok := stmt.(*ast.Class); ok {
			class = c
		}
	}
	if class == nil {
		t.Fatal("no class parsed")
	}

	want := []ast.ClassMember{
		{Visibility: "+", Name: "name", Type: "String", IsMethod: false},
		{Visibility: "+", Name: "age", Type: "int", IsMethod: false},
		{Visibility: "+", Name: "makeSound", Type: "", IsMethod: true},
		{Visibility: "+", Name: "area", Type: "float", IsMethod: true},
		{Visibility: "+", Name: "legs", Type: "", IsMethod: false},
	}

	if len(class.Members) != len(want) {
		t.Fatalf("parsed %d members, want %d", len(class.Members), len(want))
	}

	for i, w := range want {
		m := class.Members[i]
		if m.Name != w.Name || m.Type != w.Type || m.IsMethod != w.IsMethod {
			t.Errorf("member %d = {Name:%q Type:%q IsMethod:%v}, want {Name:%q Type:%q IsMethod:%v}",
				i, m.Name, m.Type, m.IsMethod, w.Name, w.Type, w.IsMethod)
		}
	}
}

func TestClassParser_SupportedTypes(t *testing.T) {
	p := parser.NewClassParser()
	types := p.SupportedTypes()
	if len(types) != 1 || types[0] != "class" {
		t.Errorf("SupportedTypes() = %v, want [\"class\"]", types)
	}
}
