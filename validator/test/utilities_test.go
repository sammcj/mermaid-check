package validator_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/validator"
)

func TestDuplicateChecker(t *testing.T) {
	tests := []struct {
		name     string
		itemType string
		checks   []struct {
			id      string
			pos     ast.Position
			wantErr bool
		}
	}{
		{
			name:     "no duplicates",
			itemType: "class",
			checks: []struct {
				id      string
				pos     ast.Position
				wantErr bool
			}{
				{"ClassA", ast.Position{Line: 1, Column: 1}, false},
				{"ClassB", ast.Position{Line: 2, Column: 1}, false},
				{"ClassC", ast.Position{Line: 3, Column: 1}, false},
			},
		},
		{
			name:     "with duplicates",
			itemType: "state",
			checks: []struct {
				id      string
				pos     ast.Position
				wantErr bool
			}{
				{"StateA", ast.Position{Line: 1, Column: 1}, false},
				{"StateB", ast.Position{Line: 2, Column: 1}, false},
				{"StateA", ast.Position{Line: 3, Column: 1}, true}, // duplicate
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dc := validator.NewDuplicateChecker(tt.itemType)
			for _, check := range tt.checks {
				err := dc.Check(check.id, check.pos)
				if (err != nil) != check.wantErr {
					t.Errorf("Check(%q) error = %v, wantErr %v", check.id, err, check.wantErr)
				}
			}
		})
	}
}

func TestReferenceChecker(t *testing.T) {
	tests := []struct {
		name     string
		itemType string
		defined  []string
		checks   []struct {
			id      string
			pos     ast.Position
			context string
			wantErr bool
		}
	}{
		{
			name:     "all references defined",
			itemType: "class",
			defined:  []string{"ClassA", "ClassB", "ClassC"},
			checks: []struct {
				id      string
				pos     ast.Position
				context string
				wantErr bool
			}{
				{"ClassA", ast.Position{Line: 5, Column: 1}, "relationship", false},
				{"ClassB", ast.Position{Line: 6, Column: 1}, "relationship", false},
			},
		},
		{
			name:     "undefined reference",
			itemType: "state",
			defined:  []string{"StateA", "StateB"},
			checks: []struct {
				id      string
				pos     ast.Position
				context string
				wantErr bool
			}{
				{"StateA", ast.Position{Line: 5, Column: 1}, "transition", false},
				{"StateC", ast.Position{Line: 6, Column: 1}, "transition", true}, // undefined
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			rc := validator.NewReferenceChecker(tt.itemType)
			for _, id := range tt.defined {
				rc.Add(id)
			}
			for _, check := range tt.checks {
				err := rc.Check(check.id, check.pos, check.context)
				if (err != nil) != check.wantErr {
					t.Errorf("Check(%q) error = %v, wantErr %v", check.id, err, check.wantErr)
				}
			}
		})
	}
}

func TestEnumValidator(t *testing.T) {
	tests := []struct {
		name       string
		valueType  string
		allowed    []string
		checks     []struct {
			value   string
			pos     ast.Position
			wantErr bool
		}
	}{
		{
			name:      "valid visibility modifiers",
			valueType: "visibility modifier",
			allowed:   []string{"+", "-", "#", "~"},
			checks: []struct {
				value   string
				pos     ast.Position
				wantErr bool
			}{
				{"+", ast.Position{Line: 1, Column: 1}, false},
				{"-", ast.Position{Line: 2, Column: 1}, false},
				{"#", ast.Position{Line: 3, Column: 1}, false},
				{"~", ast.Position{Line: 4, Column: 1}, false},
				{"*", ast.Position{Line: 5, Column: 1}, true}, // invalid
			},
		},
		{
			name:      "valid directions",
			valueType: "direction",
			allowed:   []string{"TB", "TD", "BT", "LR", "RL"},
			checks: []struct {
				value   string
				pos     ast.Position
				wantErr bool
			}{
				{"TB", ast.Position{Line: 1, Column: 1}, false},
				{"LR", ast.Position{Line: 2, Column: 1}, false},
				{"INVALID", ast.Position{Line: 3, Column: 1}, true}, // invalid
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ev := validator.NewEnumValidator(tt.valueType, tt.allowed...)
			for _, check := range tt.checks {
				err := ev.Check(check.value, check.pos)
				if (err != nil) != check.wantErr {
					t.Errorf("Check(%q) error = %v, wantErr %v", check.value, err, check.wantErr)
				}
			}
		})
	}
}
