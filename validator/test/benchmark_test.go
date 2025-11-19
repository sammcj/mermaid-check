package validator_test

import (	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"

	"github.com/sammcj/go-mermaid/validator"
)

// Benchmark validation for various diagram types using DefaultRules
func BenchmarkValidateFlowchart(b *testing.B) {
	source := `flowchart TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Process 1]
    B -->|No| D[Process 2]
    C --> E[End]
    D --> E`

	diagram, err := parser.Parse(source)
	if err != nil {
		b.Fatal(err)
	}

	rules := validator.DefaultRules()
	validator := validator.New(rules...)

	b.ReportAllocs()
	for b.Loop() {
		_ = validator.ValidateDiagram(diagram)
	}
}

func BenchmarkValidateSequence(b *testing.B) {
	source := `sequenceDiagram
    participant Alice
    participant Bob
    Alice->>Bob: Hello
    Bob-->>Alice: Hi
    loop Every minute
        Alice->>Bob: Ping
        Bob-->>Alice: Pong
    end`

	diagram, err := parser.Parse(source)
	if err != nil {
		b.Fatal(err)
	}

	rules := validator.SequenceDefaultRules()
	validator := validator.NewSequence(rules...)

	b.ReportAllocs()
	for b.Loop() {
		_ = validator.ValidateDiagram(diagram)
	}
}

func BenchmarkValidateClass(b *testing.B) {
	source := `classDiagram
    class Animal {
        +String name
        +void eat()
    }
    class Dog {
        +void bark()
    }
    Animal <|-- Dog`

	diagram, err := parser.Parse(source)
	if err != nil {
		b.Fatal(err)
	}

	rules := validator.ClassDefaultRules()
	validator := validator.NewClass(rules...)

	b.ReportAllocs()
	for b.Loop() {
		_ = validator.ValidateDiagram(diagram)
	}
}

func BenchmarkValidateState(b *testing.B) {
	source := `stateDiagram-v2
    [*] --> Still
    Still --> Moving
    Moving --> [*]`

	diagram, err := parser.Parse(source)
	if err != nil {
		b.Fatal(err)
	}

	rules := validator.StateDefaultRules()
	validator := validator.NewState(rules...)

	b.ReportAllocs()
	for b.Loop() {
		_ = validator.ValidateDiagram(diagram)
	}
}

func BenchmarkDuplicateChecker(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		checker := validator.NewDuplicateChecker("test item")
		for j := range 100 {
			itemName := "item"
			if j == 50 {
				itemName = "duplicate"
			}
			pos := ast.Position{Line: 1, Column: 1}
			_ = checker.Check(itemName, pos)
		}
	}
}

func BenchmarkReferenceChecker(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		checker := validator.NewReferenceChecker("test item")
		for j := range 100 {
			itemID := "item" + string(rune('A'+j%26))
			checker.Add(itemID)
		}
		for j := range 100 {
			itemID := "item" + string(rune('A'+j%26))
			pos := ast.Position{Line: 1, Column: 1}
			_ = checker.Check(itemID, pos, "test context")
		}
	}
}

func BenchmarkEnumValidator(b *testing.B) {
	allowed := []string{"option1", "option2", "option3", "option4", "option5"}
	b.ReportAllocs()
	for b.Loop() {
		validator := validator.NewEnumValidator("test field", allowed...)
		for j := range 100 {
			pos := ast.Position{Line: 1, Column: 1}
			_ = validator.Check(allowed[j%len(allowed)], pos)
		}
	}
}
