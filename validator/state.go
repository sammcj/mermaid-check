package validator

import (
	"fmt"

	"github.com/sammcj/go-mermaid/ast"
)

// StateRule defines a validation rule for state diagrams.
type StateRule interface {
	Name() string
	ValidateState(diagram *ast.StateDiagram) []ValidationError
}

// NoDuplicateStates checks for duplicate state IDs.
type NoDuplicateStates struct{}

// Name returns the rule name.
func (r *NoDuplicateStates) Name() string {
	return "no-duplicate-states"
}

// ValidateState validates the state diagram.
func (r *NoDuplicateStates) ValidateState(diagram *ast.StateDiagram) []ValidationError {
	var errors []ValidationError
	seen := make(map[string]ast.Position)

	for _, stmt := range diagram.Statements {
		if state, ok := stmt.(*ast.State); ok {
			if pos, exists := seen[state.ID]; exists {
				errors = append(errors, ValidationError{
					Line:     state.Pos.Line,
					Column:   state.Pos.Column,
					Message:  fmt.Sprintf("duplicate state ID %q (first defined at line %d)", state.ID, pos.Line),
					Severity: SeverityError,
				})
			} else {
				seen[state.ID] = state.Pos
			}
		}
	}

	return errors
}

// ValidStateReferences checks that all states referenced in transitions exist.
type ValidStateReferences struct{}

// Name returns the rule name.
func (r *ValidStateReferences) Name() string {
	return "valid-state-references"
}

// ValidateState validates the state diagram.
func (r *ValidStateReferences) ValidateState(diagram *ast.StateDiagram) []ValidationError {
	var errors []ValidationError

	// Collect defined states
	definedStates := make(map[string]bool)
	for _, stmt := range diagram.Statements {
		switch s := stmt.(type) {
		case *ast.State:
			definedStates[s.ID] = true
		case *ast.Fork:
			definedStates[s.ID] = true
		case *ast.Join:
			definedStates[s.ID] = true
		case *ast.Choice:
			definedStates[s.ID] = true
		}
	}

	// Check transitions
	for _, stmt := range diagram.Statements {
		if trans, ok := stmt.(*ast.Transition); ok {
			if !definedStates[trans.From] {
				errors = append(errors, ValidationError{
					Line:     trans.Pos.Line,
					Column:   trans.Pos.Column,
					Message:  fmt.Sprintf("transition references undefined state %q", trans.From),
					Severity: SeverityError,
				})
			}
			if !definedStates[trans.To] {
				errors = append(errors, ValidationError{
					Line:     trans.Pos.Line,
					Column:   trans.Pos.Column,
					Message:  fmt.Sprintf("transition references undefined state %q", trans.To),
					Severity: SeverityError,
				})
			}
		}
	}

	return errors
}

// StateDefaultRules returns the default set of validation rules for state diagrams.
func StateDefaultRules() []StateRule {
	return []StateRule{
		&NoDuplicateStates{},
		&ValidStateReferences{},
	}
}

// StateStrictRules returns a strict set of validation rules for state diagrams.
func StateStrictRules() []StateRule {
	return StateDefaultRules()
}

// NewState creates a new state diagram validator with the given rules.
func NewState(rules ...StateRule) *Validator {
	return &Validator{
		stateRules: rules,
	}
}
