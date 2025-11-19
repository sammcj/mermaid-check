package parser_test

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestSequenceParser_Parse(t *testing.T) {
	p := parser.NewSequenceParser()

	tests := []struct {
		name    string
		source  string
		wantErr bool
	}{
		{
			name: "simple message",
			source: `sequenceDiagram
    Alice->>John: Hello John`,
			wantErr: false,
		},
		{
			name: "with participant",
			source: `sequenceDiagram
    participant Alice
    Alice->>Bob: Hi`,
			wantErr: false,
		},
		{
			name: "with loop",
			source: `sequenceDiagram
    Alice->>Bob: Start
    loop Every minute
        Bob->>Alice: Ping
    end`,
			wantErr: false,
		},
		{
			name: "with alt",
			source: `sequenceDiagram
    Alice->>Bob: Request
    alt Success
        Bob->>Alice: OK
    else Failure
        Bob->>Alice: Error
    end`,
			wantErr: false,
		},
		{
			name: "invalid header",
			source: `sequence
    Alice->>Bob: Hi`,
			wantErr: true,
		},
		{
			name: "empty diagram",
			source: ``,
			wantErr: true,
		},
	}

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
				seqDiagram, ok := diagram.(*ast.SequenceDiagram)
				if !ok {
					t.Errorf("Parse() returned wrong type: %T", diagram)
				}
				if seqDiagram.Type != "sequence" {
					t.Errorf("Parse() diagram type = %s, want 'sequence'", seqDiagram.Type)
				}
			}
		})
	}
}

func TestSequenceParser_ParseTestDataFiles(t *testing.T) {
	testDataDir := filepath.Join("../../testdata", "sequence")

	files, err := os.ReadDir(testDataDir)
	if err != nil {
		if os.IsNotExist(err) {
			t.Skip("testdata/sequence directory does not exist")
		}
		t.Fatalf("failed to read testdata directory: %v", err)
	}

	p := parser.NewSequenceParser()

	for _, file := range files {
		if file.IsDir() || filepath.Ext(file.Name()) != ".mmd" {
			continue
		}

		t.Run(file.Name(), func(t *testing.T) {
			path := filepath.Join(testDataDir, file.Name())
			data, err := os.ReadFile(path) //nolint:gosec // Test file paths are safe
			if err != nil {
				t.Fatalf("failed to read file: %v", err)
			}

			diagram, err := p.Parse(string(data))
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			if diagram == nil {
				t.Fatal("Parse() returned nil diagram")
			}

			seqDiagram, ok := diagram.(*ast.SequenceDiagram)
			if !ok {
				t.Fatalf("Parse() returned wrong type: %T", diagram)
			}

			if len(seqDiagram.Statements) == 0 {
				t.Error("Parse() returned diagram with no statements")
			}
		})
	}
}

func TestSequenceParser_Messages(t *testing.T) {
	p := parser.NewSequenceParser()

	tests := []struct {
		name      string
		source    string
		wantArrow string
	}{
		{"solid arrow", "sequenceDiagram\n    A->>B: Hi", "->>"},
		{"dotted arrow", "sequenceDiagram\n    A-->>B: Hi", "-->>"},
		{"solid no arrow", "sequenceDiagram\n    A->B: Hi", "->"},
		{"dotted no arrow", "sequenceDiagram\n    A-->B: Hi", "-->"},
		{"cross arrow", "sequenceDiagram\n    A-xB: Hi", "-x"},
		{"async arrow", "sequenceDiagram\n    A-)B: Hi", "-)"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram, err := p.Parse(tt.source)
			if err != nil {
				t.Fatalf("Parse() error = %v", err)
			}

			seqDiagram := diagram.(*ast.SequenceDiagram)
			if len(seqDiagram.Statements) == 0 {
				t.Fatal("no statements parsed")
			}

			msg, ok := seqDiagram.Statements[0].(*ast.Message)
			if !ok {
				t.Fatalf("first statement is not a message: %T", seqDiagram.Statements[0])
			}

			if msg.Arrow != tt.wantArrow {
				t.Errorf("arrow = %s, want %s", msg.Arrow, tt.wantArrow)
			}
		})
	}
}
