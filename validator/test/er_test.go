package validator_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/validator"
)

func TestNoDuplicateEntitiesRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.ERDiagram
		wantErr bool
	}{
		{
			name: "no duplicates",
			diagram: &ast.ERDiagram{
				Entities: []ast.EREntity{
					{Name: "CUSTOMER", Pos: ast.Position{Line: 2, Column: 1}},
					{Name: "ORDER", Pos: ast.Position{Line: 3, Column: 1}},
					{Name: "PRODUCT", Pos: ast.Position{Line: 4, Column: 1}},
				},
			},
			wantErr: false,
		},
		{
			name: "with duplicates",
			diagram: &ast.ERDiagram{
				Entities: []ast.EREntity{
					{Name: "CUSTOMER", Pos: ast.Position{Line: 2, Column: 1}},
					{Name: "ORDER", Pos: ast.Position{Line: 3, Column: 1}},
					{Name: "CUSTOMER", Pos: ast.Position{Line: 4, Column: 1}},
				},
			},
			wantErr: true,
		},
	}

	rule := &validator.NoDuplicateEntitiesRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("NoDuplicateEntitiesRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestValidRelationshipReferencesRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.ERDiagram
		wantErr bool
	}{
		{
			name: "valid references",
			diagram: &ast.ERDiagram{
				Entities: []ast.EREntity{
					{Name: "CUSTOMER", Pos: ast.Position{Line: 2, Column: 1}},
					{Name: "ORDER", Pos: ast.Position{Line: 3, Column: 1}},
				},
				Relationships: []ast.ERRelationship{
					{From: "CUSTOMER", To: "ORDER", Pos: ast.Position{Line: 4, Column: 1}},
				},
			},
			wantErr: false,
		},
		{
			name: "undefined source entity",
			diagram: &ast.ERDiagram{
				Entities: []ast.EREntity{
					{Name: "ORDER", Pos: ast.Position{Line: 2, Column: 1}},
				},
				Relationships: []ast.ERRelationship{
					{From: "CUSTOMER", To: "ORDER", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErr: true,
		},
		{
			name: "undefined target entity",
			diagram: &ast.ERDiagram{
				Entities: []ast.EREntity{
					{Name: "CUSTOMER", Pos: ast.Position{Line: 2, Column: 1}},
				},
				Relationships: []ast.ERRelationship{
					{From: "CUSTOMER", To: "ORDER", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErr: true,
		},
		{
			name: "both entities undefined",
			diagram: &ast.ERDiagram{
				Entities: []ast.EREntity{},
				Relationships: []ast.ERRelationship{
					{From: "CUSTOMER", To: "ORDER", Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErr: true,
		},
	}

	rule := &validator.ValidRelationshipReferencesRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("ValidRelationshipReferencesRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestValidAttributeKeysRule(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.ERDiagram
		wantErr bool
	}{
		{
			name: "valid keys",
			diagram: &ast.ERDiagram{
				Entities: []ast.EREntity{
					{
						Name: "CUSTOMER",
						Attributes: []ast.ERAttribute{
							{Name: "id", Keys: []string{"PK"}, Pos: ast.Position{Line: 2, Column: 1}},
							{Name: "email", Keys: []string{"UK"}, Pos: ast.Position{Line: 3, Column: 1}},
							{Name: "orderId", Keys: []string{"FK"}, Pos: ast.Position{Line: 4, Column: 1}},
						},
					},
				},
			},
			wantErr: false,
		},
		{
			name: "invalid key",
			diagram: &ast.ERDiagram{
				Entities: []ast.EREntity{
					{
						Name: "CUSTOMER",
						Attributes: []ast.ERAttribute{
							{Name: "id", Keys: []string{"INVALID"}, Pos: ast.Position{Line: 2, Column: 1}},
						},
					},
				},
			},
			wantErr: true,
		},
		{
			name: "mixed valid and invalid keys",
			diagram: &ast.ERDiagram{
				Entities: []ast.EREntity{
					{
						Name: "CUSTOMER",
						Attributes: []ast.ERAttribute{
							{Name: "id", Keys: []string{"PK", "INVALID"}, Pos: ast.Position{Line: 2, Column: 1}},
						},
					},
				},
			},
			wantErr: true,
		},
	}

	rule := &validator.ValidAttributeKeysRule{}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := rule.Validate(tt.diagram)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("ValidAttributeKeysRule.Validate() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}

func TestERDefaultRules(t *testing.T) {
	rules := validator.ERDefaultRules()
	if len(rules) == 0 {
		t.Error("ERDefaultRules() returned empty slice")
	}
	if len(rules) != 3 {
		t.Errorf("expected 3 default rules, got %d", len(rules))
	}
}

func TestERStrictRules(t *testing.T) {
	rules := validator.ERStrictRules()
	if len(rules) == 0 {
		t.Error("ERStrictRules() returned empty slice")
	}
}

func TestValidateER(t *testing.T) {
	tests := []struct {
		name    string
		diagram *ast.ERDiagram
		strict  bool
		wantErr bool
	}{
		{
			name: "valid diagram",
			diagram: &ast.ERDiagram{
				Entities: []ast.EREntity{
					{
						Name: "CUSTOMER",
						Attributes: []ast.ERAttribute{
							{Name: "id", Keys: []string{"PK"}, Pos: ast.Position{Line: 2, Column: 1}},
						},
						Pos: ast.Position{Line: 2, Column: 1},
					},
					{
						Name: "ORDER",
						Attributes: []ast.ERAttribute{
							{Name: "orderId", Keys: []string{"PK"}, Pos: ast.Position{Line: 5, Column: 1}},
							{Name: "customerId", Keys: []string{"FK"}, Pos: ast.Position{Line: 6, Column: 1}},
						},
						Pos: ast.Position{Line: 4, Column: 1},
					},
				},
				Relationships: []ast.ERRelationship{
					{From: "CUSTOMER", To: "ORDER", Pos: ast.Position{Line: 8, Column: 1}},
				},
			},
			strict:  false,
			wantErr: false,
		},
		{
			name: "duplicate entities",
			diagram: &ast.ERDiagram{
				Entities: []ast.EREntity{
					{Name: "CUSTOMER", Pos: ast.Position{Line: 2, Column: 1}},
					{Name: "CUSTOMER", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			strict:  false,
			wantErr: true,
		},
		{
			name: "undefined entity in relationship",
			diagram: &ast.ERDiagram{
				Entities: []ast.EREntity{
					{Name: "CUSTOMER", Pos: ast.Position{Line: 2, Column: 1}},
				},
				Relationships: []ast.ERRelationship{
					{From: "CUSTOMER", To: "ORDER", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			strict:  false,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errors := validator.ValidateER(tt.diagram, tt.strict)
			if (len(errors) > 0) != tt.wantErr {
				t.Errorf("ValidateER() errors = %v, wantErr %v", errors, tt.wantErr)
			}
		})
	}
}
