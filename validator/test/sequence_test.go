package validator_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/validator"
)

func TestNoDuplicateParticipants(t *testing.T) {
	tests := []struct {
		name       string
		diagram    *ast.SequenceDiagram
		wantErrors int
	}{
		{
			name: "no duplicates",
			diagram: &ast.SequenceDiagram{
				Type: "sequence",
				Statements: []ast.SeqStmt{
					&ast.Participant{ID: "Alice", Pos: ast.Position{Line: 2, Column: 1}},
					&ast.Participant{ID: "Bob", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrors: 0,
		},
		{
			name: "with duplicates",
			diagram: &ast.SequenceDiagram{
				Type: "sequence",
				Statements: []ast.SeqStmt{
					&ast.Participant{ID: "Alice", Pos: ast.Position{Line: 2, Column: 1}},
					&ast.Participant{ID: "Alice", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrors: 1,
		},
	}

	rule := &validator.NoDuplicateParticipants{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.ValidateSequence(tt.diagram)
			if len(errors) != tt.wantErrors {
				t.Errorf("ValidateSequence() errors = %d, want %d", len(errors), tt.wantErrors)
			}
		})
	}
}

func TestValidMessageArrows(t *testing.T) {
	tests := []struct {
		name       string
		diagram    *ast.SequenceDiagram
		wantErrors int
	}{
		{
			name: "valid arrow",
			diagram: &ast.SequenceDiagram{
				Type: "sequence",
				Statements: []ast.SeqStmt{
					&ast.Message{From: "A", To: "B", Arrow: "->>", Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErrors: 0,
		},
		{
			name: "invalid arrow",
			diagram: &ast.SequenceDiagram{
				Type: "sequence",
				Statements: []ast.SeqStmt{
					&ast.Message{From: "A", To: "B", Arrow: "===", Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErrors: 1,
		},
		{
			name: "bidirectional arrow",
			diagram: &ast.SequenceDiagram{
				Type: "sequence",
				Statements: []ast.SeqStmt{
					&ast.Message{From: "A", To: "B", Arrow: "<<->>", Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErrors: 0,
		},
	}

	rule := &validator.ValidMessageArrows{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.ValidateSequence(tt.diagram)
			if len(errors) != tt.wantErrors {
				t.Errorf("ValidateSequence() errors = %d, want %d", len(errors), tt.wantErrors)
				for _, err := range errors {
					t.Logf("  error: %s", err.Message)
				}
			}
		})
	}
}

func TestValidNotePositions(t *testing.T) {
	tests := []struct {
		name       string
		diagram    *ast.SequenceDiagram
		wantErrors int
	}{
		{
			name: "valid note",
			diagram: &ast.SequenceDiagram{
				Type: "sequence",
				Statements: []ast.SeqStmt{
					&ast.Participant{ID: "Alice", Pos: ast.Position{Line: 2, Column: 1}},
					&ast.Note{Position: "right of", Participants: []string{"Alice"}, Text: "Note", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrors: 0,
		},
		{
			name: "note referencing undefined participant",
			diagram: &ast.SequenceDiagram{
				Type: "sequence",
				Statements: []ast.SeqStmt{
					&ast.Participant{ID: "Alice", Pos: ast.Position{Line: 2, Column: 1}},
					&ast.Note{Position: "right of", Participants: []string{"Bob"}, Text: "Note", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrors: 1,
		},
		{
			name: "note with implicit participant",
			diagram: &ast.SequenceDiagram{
				Type: "sequence",
				Statements: []ast.SeqStmt{
					&ast.Message{From: "Alice", To: "Bob", Arrow: "->>", Pos: ast.Position{Line: 2, Column: 1}},
					&ast.Note{Position: "over", Participants: []string{"Alice", "Bob"}, Text: "Note", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrors: 0,
		},
	}

	rule := &validator.ValidNotePositions{}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.ValidateSequence(tt.diagram)
			if len(errors) != tt.wantErrors {
				t.Errorf("ValidateSequence() errors = %d, want %d", len(errors), tt.wantErrors)
				for _, err := range errors {
					t.Logf("  error: %s", err.Message)
				}
			}
		})
	}
}

func TestSequenceDefaultRules(t *testing.T) {
	rules := validator.SequenceDefaultRules()
	if len(rules) == 0 {
		t.Error("SequenceDefaultRules() returned empty rules")
	}
}

func TestSequenceStrictRules(t *testing.T) {
	rules := validator.SequenceStrictRules()
	if len(rules) == 0 {
		t.Error("SequenceStrictRules() returned empty rules")
	}
}

func TestSequenceValidatorWithComplexDiagram(t *testing.T) {
	diagram := &ast.SequenceDiagram{
		Type: "sequence",
		Statements: []ast.SeqStmt{
			&ast.Participant{ID: "Alice", Pos: ast.Position{Line: 2, Column: 1}},
			&ast.Participant{ID: "Bob", Pos: ast.Position{Line: 3, Column: 1}},
			&ast.Message{From: "Alice", To: "Bob", Arrow: "->>", Text: "Request", Pos: ast.Position{Line: 4, Column: 1}},
			&ast.Activation{Participant: "Bob", Active: true, Pos: ast.Position{Line: 5, Column: 1}},
			&ast.Message{From: "Bob", To: "Alice", Arrow: "-->>", Text: "Response", Pos: ast.Position{Line: 6, Column: 1}},
			&ast.Activation{Participant: "Bob", Active: false, Pos: ast.Position{Line: 7, Column: 1}},
			&ast.Note{Position: "right of", Participants: []string{"Alice"}, Text: "Alice receives", Pos: ast.Position{Line: 8, Column: 1}},
		},
	}

	v := validator.NewSequence(validator.SequenceDefaultRules()...)
	errors := v.ValidateDiagram(diagram)
	if len(errors) > 0 {
		t.Errorf("unexpected validation errors: %v", errors)
	}
}

func TestSequenceValidatorWithLoopAndAlt(t *testing.T) {
	diagram := &ast.SequenceDiagram{
		Type: "sequence",
		Statements: []ast.SeqStmt{
			&ast.Participant{ID: "User", Pos: ast.Position{Line: 2, Column: 1}},
			&ast.Participant{ID: "System", Pos: ast.Position{Line: 3, Column: 1}},
			&ast.Loop{
				Label: "while active",
				Statements: []ast.SeqStmt{
					&ast.Message{From: "User", To: "System", Arrow: "->>", Text: "ping", Pos: ast.Position{Line: 5, Column: 1}},
					&ast.Message{From: "System", To: "User", Arrow: "-->>", Text: "pong", Pos: ast.Position{Line: 6, Column: 1}},
				},
				Pos: ast.Position{Line: 4, Column: 1},
			},
			&ast.Alt{
				Conditions: []ast.AltCondition{
					{
						Label: "successful",
						Statements: []ast.SeqStmt{
							&ast.Message{From: "System", To: "User", Arrow: "->>", Text: "OK", Pos: ast.Position{Line: 9, Column: 1}},
						},
					},
					{
						Label: "else",
						IsElse: true,
						Statements: []ast.SeqStmt{
							&ast.Message{From: "System", To: "User", Arrow: "->>", Text: "Error", Pos: ast.Position{Line: 11, Column: 1}},
						},
					},
				},
				Pos: ast.Position{Line: 8, Column: 1},
			},
		},
	}

	v := validator.NewSequence(validator.SequenceDefaultRules()...)
	errors := v.ValidateDiagram(diagram)
	if len(errors) > 0 {
		t.Errorf("unexpected validation errors: %v", errors)
	}
}

func TestSequenceValidatorWithParBlock(t *testing.T) {
	diagram := &ast.SequenceDiagram{
		Type: "sequence",
		Statements: []ast.SeqStmt{
			&ast.Participant{ID: "A", Pos: ast.Position{Line: 2, Column: 1}},
			&ast.Participant{ID: "B", Pos: ast.Position{Line: 3, Column: 1}},
			&ast.Participant{ID: "C", Pos: ast.Position{Line: 4, Column: 1}},
			&ast.Par{
				Branches: []ast.ParBranch{
					{
						Label: "parallel",
						Statements: []ast.SeqStmt{
							&ast.Message{From: "A", To: "B", Arrow: "->>", Text: "msg1", Pos: ast.Position{Line: 6, Column: 1}},
						},
					},
					{
						Label: "and",
						Statements: []ast.SeqStmt{
							&ast.Message{From: "A", To: "C", Arrow: "->>", Text: "msg2", Pos: ast.Position{Line: 8, Column: 1}},
						},
					},
				},
				Pos: ast.Position{Line: 5, Column: 1},
			},
		},
	}

	v := validator.NewSequence(validator.SequenceDefaultRules()...)
	errors := v.ValidateDiagram(diagram)
	if len(errors) > 0 {
		t.Errorf("unexpected validation errors: %v", errors)
	}
}

func TestSequenceValidatorWithBreakAndOpt(t *testing.T) {
	diagram := &ast.SequenceDiagram{
		Type: "sequence",
		Statements: []ast.SeqStmt{
			&ast.Participant{ID: "Client", Pos: ast.Position{Line: 2, Column: 1}},
			&ast.Participant{ID: "Server", Pos: ast.Position{Line: 3, Column: 1}},
			&ast.Break{
				Label: "when error",
				Statements: []ast.SeqStmt{
					&ast.Message{From: "Server", To: "Client", Arrow: "-->>", Text: "error", Pos: ast.Position{Line: 5, Column: 1}},
				},
				Pos: ast.Position{Line: 4, Column: 1},
			},
			&ast.Opt{
				Label: "if success",
				Statements: []ast.SeqStmt{
					&ast.Message{From: "Server", To: "Client", Arrow: "-->>", Text: "success", Pos: ast.Position{Line: 8, Column: 1}},
				},
				Pos: ast.Position{Line: 7, Column: 1},
			},
		},
	}

	v := validator.NewSequence(validator.SequenceDefaultRules()...)
	errors := v.ValidateDiagram(diagram)
	if len(errors) > 0 {
		t.Errorf("unexpected validation errors: %v", errors)
	}
}

func TestSequenceValidatorWithCritical(t *testing.T) {
	diagram := &ast.SequenceDiagram{
		Type: "sequence",
		Statements: []ast.SeqStmt{
			&ast.Participant{ID: "A", Pos: ast.Position{Line: 2, Column: 1}},
			&ast.Participant{ID: "B", Pos: ast.Position{Line: 3, Column: 1}},
			&ast.Critical{
				Label: "atomic operation",
				Statements: []ast.SeqStmt{
					&ast.Message{From: "A", To: "B", Arrow: "->>", Text: "do", Pos: ast.Position{Line: 5, Column: 1}},
				},
				Pos: ast.Position{Line: 4, Column: 1},
			},
		},
	}

	v := validator.NewSequence(validator.SequenceDefaultRules()...)
	errors := v.ValidateDiagram(diagram)
	if len(errors) > 0 {
		t.Errorf("unexpected validation errors: %v", errors)
	}
}

func TestValidMessageArrowsExtended(t *testing.T) {
	tests := []struct {
		name       string
		arrow      string
		wantErrors int
	}{
		{"solid line", "->", 0},
		{"solid arrow", "->>", 0},
		{"dotted arrow", "-->>", 0},
		{"solid cross", "-x", 0},
		{"dotted cross", "--x", 0},
		{"solid open arrow", "-)", 0},
		{"dotted open arrow", "--)", 0},
		{"bidirectional", "<<->>", 0},
		{"invalid", ">>>", 1},
		{"empty", "", 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram := &ast.SequenceDiagram{
				Type: "sequence",
				Statements: []ast.SeqStmt{
					&ast.Message{From: "A", To: "B", Arrow: tt.arrow, Pos: ast.Position{Line: 2, Column: 1}},
				},
			}
			rule := &validator.ValidMessageArrows{}
			errors := rule.ValidateSequence(diagram)
			if len(errors) != tt.wantErrors {
				t.Errorf("ValidateSequence() errors = %d, want %d for arrow %q", len(errors), tt.wantErrors, tt.arrow)
			}
		})
	}
}

func TestValidNotePositionsExtended(t *testing.T) {
	tests := []struct {
		name         string
		position     string
		participants []string
		wantErrors   int
	}{
		{"right of", "right of", []string{"Alice"}, 0},
		{"left of", "left of", []string{"Alice"}, 0},
		{"over single", "over", []string{"Alice"}, 0},
		{"over multiple", "over", []string{"Alice", "Bob"}, 0},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram := &ast.SequenceDiagram{
				Type: "sequence",
				Statements: []ast.SeqStmt{
					&ast.Participant{ID: "Alice", Pos: ast.Position{Line: 2, Column: 1}},
					&ast.Participant{ID: "Bob", Pos: ast.Position{Line: 3, Column: 1}},
					&ast.Note{Position: tt.position, Participants: tt.participants, Text: "Note", Pos: ast.Position{Line: 4, Column: 1}},
				},
			}
			rule := &validator.ValidNotePositions{}
			errors := rule.ValidateSequence(diagram)
			if len(errors) != tt.wantErrors {
				t.Errorf("ValidateSequence() errors = %d, want %d for position %q", len(errors), tt.wantErrors, tt.position)
			}
		})
	}
}
