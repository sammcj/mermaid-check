package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/parser"
)

// Benchmark data for various diagram types
var (
	benchmarkFlowchart = `flowchart TD
    A[Start] --> B{Decision}
    B -->|Yes| C[Process 1]
    B -->|No| D[Process 2]
    C --> E[End]
    D --> E
    subgraph Sub
        C
        D
    end`

	benchmarkSequence = `sequenceDiagram
    participant Alice
    participant Bob
    participant Charlie
    Alice->>Bob: Hello Bob
    Bob->>Charlie: Hello Charlie
    Charlie-->>Bob: Hi Bob
    Bob-->>Alice: Hi Alice
    loop Every minute
        Alice->>Bob: Ping
        Bob-->>Alice: Pong
    end`

	benchmarkClass = `classDiagram
    class Animal {
        +String name
        +int age
        +void eat()
        +void sleep()
    }
    class Dog {
        +String breed
        +void bark()
    }
    class Cat {
        +void meow()
    }
    Animal <|-- Dog
    Animal <|-- Cat
    Dog --> Owner
    class Owner {
        +String name
    }`

	benchmarkState = `stateDiagram-v2
    [*] --> Still
    Still --> [*]
    Still --> Moving
    Moving --> Still
    Moving --> Crash
    Crash --> [*]

    state Moving {
        [*] --> Slow
        Slow --> Fast
        Fast --> Slow
        Slow --> [*]
    }`

	benchmarkER = `erDiagram
    CUSTOMER ||--o{ ORDER : places
    CUSTOMER {
        string name
        string email
        int customerId
    }
    ORDER ||--|{ LINE-ITEM : contains
    ORDER {
        int orderNumber
        date orderDate
    }
    LINE-ITEM {
        string productCode
        int quantity
        float price
    }`

	benchmarkGantt = `gantt
    title Project Timeline
    dateFormat YYYY-MM-DD
    section Planning
    Research           :2024-01-01, 30d
    Design             :2024-01-15, 20d
    section Development
    Backend            :2024-02-01, 45d
    Frontend           :2024-02-15, 40d
    section Testing
    Unit Tests         :2024-03-20, 15d
    Integration Tests  :2024-04-01, 10d`

	benchmarkLargeFlowchart = generateLargeFlowchart(100)
)

func generateLargeFlowchart(nodes int) string {
	diagram := "flowchart TD\n"
	for i := range nodes {
		diagram += "    Node" + string(rune('A'+i%26)) + string(rune('0'+i/26)) + " --> Node" + string(rune('A'+(i+1)%26)) + string(rune('0'+(i+1)/26)) + "\n"
	}
	return diagram
}

func BenchmarkParseFlowchart(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, err := parser.Parse(benchmarkFlowchart)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseSequence(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, err := parser.Parse(benchmarkSequence)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseClass(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, err := parser.Parse(benchmarkClass)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseState(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, err := parser.Parse(benchmarkState)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseER(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, err := parser.Parse(benchmarkER)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseGantt(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, err := parser.Parse(benchmarkGantt)
		if err != nil {
			b.Fatal(err)
		}
	}
}

func BenchmarkParseLargeFlowchart(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_, err := parser.Parse(benchmarkLargeFlowchart)
		if err != nil {
			b.Fatal(err)
		}
	}
}

// NOTE: BenchmarkDetectDiagramType is commented out because detectDiagramType is an unexported function
// and this file uses black-box testing (package parser_test).
// This benchmark should be moved to a white-box test file if needed.
/*
func BenchmarkDetectDiagramType(b *testing.B) {
	b.ReportAllocs()
	for b.Loop() {
		_ = detectDiagramType(benchmarkFlowchart)
	}
}
*/
