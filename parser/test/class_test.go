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

func TestClassParser_Relationships(t *testing.T) {
	src := "classDiagram\n" +
		"    Animal <|-- Dog\n" +
		"    Duck ..|> Flyer\n" +
		"    Car *-- Wheel\n" +
		"    House o-- Room\n" +
		"    A --> B\n" +
		"    X ..> Y"
	d, err := parser.NewClassParser().Parse(src)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	cd, ok := d.(*ast.ClassDiagram)
	if !ok {
		t.Fatalf("Parse() = %T, want *ast.ClassDiagram", d)
	}
	var rels []*ast.Relationship
	for _, s := range cd.Statements {
		if r, ok := s.(*ast.Relationship); ok {
			rels = append(rels, r)
		}
	}
	want := []ast.Relationship{
		{From: "Animal", To: "Dog", Type: "inheritance"},
		{From: "Duck", To: "Flyer", Type: "realization"},
		{From: "Car", To: "Wheel", Type: "composition"},
		{From: "House", To: "Room", Type: "aggregation"},
		{From: "A", To: "B", Type: "association"},
		{From: "X", To: "Y", Type: "dependency"},
	}
	if len(rels) != len(want) {
		t.Fatalf("got %d relationships, want %d: %+v", len(rels), len(want), rels)
	}
	for i, w := range want {
		g := rels[i]
		if g.From != w.From || g.To != w.To || g.Type != w.Type {
			t.Errorf("rel %d = {From:%q To:%q Type:%q}, want {From:%q To:%q Type:%q}", i, g.From, g.To, g.Type, w.From, w.To, w.Type)
		}
	}
}

func TestClassParser_Notes(t *testing.T) {
	src := "classDiagram\n" +
		"    note for Animal \"An animal note\"\n" +
		"    note \"a floating note\""
	d, err := parser.NewClassParser().Parse(src)
	if err != nil {
		t.Fatalf("Parse() error: %v", err)
	}
	cd, ok := d.(*ast.ClassDiagram)
	if !ok {
		t.Fatalf("Parse() = %T, want *ast.ClassDiagram", d)
	}
	var notes []*ast.ClassNote
	for _, s := range cd.Statements {
		if n, ok := s.(*ast.ClassNote); ok {
			notes = append(notes, n)
		}
	}
	want := []ast.ClassNote{
		{ClassName: "Animal", Text: "An animal note"},
		{ClassName: "", Text: "a floating note"},
	}
	if len(notes) != len(want) {
		t.Fatalf("got %d notes, want %d: %+v", len(notes), len(want), notes)
	}
	for i, w := range want {
		g := notes[i]
		if g.ClassName != w.ClassName || g.Text != w.Text {
			t.Errorf("note %d = {ClassName:%q Text:%q}, want {ClassName:%q Text:%q}", i, g.ClassName, g.Text, w.ClassName, w.Text)
		}
	}
}
