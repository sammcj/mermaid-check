package validator_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/validator"
)

// TestRuleNames ensures all validation rules have proper names.
func TestRuleNames(t *testing.T) {
	tests := []struct {
		name     string
		rule     interface{ Name() string }
		expected string
	}{
		// Sequence rules
		{"NoDuplicateParticipants", &validator.NoDuplicateParticipants{}, "no-duplicate-participants"},
		{"ValidMessageArrows", &validator.ValidMessageArrows{}, "valid-message-arrows"},
		{"ValidNotePositions", &validator.ValidNotePositions{}, "valid-note-positions"},

		// Class rules
		{"NoDuplicateClasses", &validator.NoDuplicateClasses{}, "no-duplicate-classes"},
		{"ValidClassReferences", &validator.ValidClassReferences{}, "valid-class-references"},

		// State rules
		{"NoDuplicateStates", &validator.NoDuplicateStates{}, "no-duplicate-states"},
		{"ValidStateReferences", &validator.ValidStateReferences{}, "valid-state-references"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := tt.rule.Name(); got != tt.expected {
				t.Errorf("Name() = %v, want %v", got, tt.expected)
			}
		})
	}
}
