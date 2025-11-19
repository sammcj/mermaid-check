package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestERParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		source  string
		wantErr bool
		check   func(*testing.T, ast.Diagram)
	}{
		{
			name: "simple ER diagram with relationship",
			source: `erDiagram
    CUSTOMER ||--o{ ORDER : places`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				er, ok := d.(*ast.ERDiagram)
				if !ok {
					t.Fatalf("expected *ast.ERDiagram, got %T", d)
				}
				if len(er.Relationships) != 1 {
					t.Errorf("expected 1 relationship, got %d", len(er.Relationships))
				}
				if er.Relationships[0].From != "CUSTOMER" {
					t.Errorf("expected From 'CUSTOMER', got %q", er.Relationships[0].From)
				}
				if er.Relationships[0].To != "ORDER" {
					t.Errorf("expected To 'ORDER', got %q", er.Relationships[0].To)
				}
				if er.Relationships[0].FromCard != "||" {
					t.Errorf("expected FromCard '||', got %q", er.Relationships[0].FromCard)
				}
				if er.Relationships[0].ToCard != "o{" {
					t.Errorf("expected ToCard 'o{', got %q", er.Relationships[0].ToCard)
				}
				if er.Relationships[0].Label != "places" {
					t.Errorf("expected Label 'places', got %q", er.Relationships[0].Label)
				}
			},
		},
		{
			name: "entity with attributes",
			source: `erDiagram
    CUSTOMER {
        string name PK
        string email UK
        int customerID
    }`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				er, ok := d.(*ast.ERDiagram)
				if !ok {
					t.Fatalf("expected *ast.ERDiagram, got %T", d)
				}
				if len(er.Entities) != 1 {
					t.Errorf("expected 1 entity, got %d", len(er.Entities))
				}
				entity := er.Entities[0]
				if entity.Name != "CUSTOMER" {
					t.Errorf("expected entity name 'CUSTOMER', got %q", entity.Name)
				}
				if len(entity.Attributes) != 3 {
					t.Errorf("expected 3 attributes, got %d", len(entity.Attributes))
				}
				// Check first attribute
				if entity.Attributes[0].Name != "name" {
					t.Errorf("expected attribute name 'name', got %q", entity.Attributes[0].Name)
				}
				if len(entity.Attributes[0].Keys) != 1 || entity.Attributes[0].Keys[0] != "PK" {
					t.Errorf("expected Keys [PK], got %v", entity.Attributes[0].Keys)
				}
			},
		},
		{
			name: "entity with asterisk primary key notation",
			source: `erDiagram
    USER {
        int *userID
        string username
    }`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				er, ok := d.(*ast.ERDiagram)
				if !ok {
					t.Fatalf("expected *ast.ERDiagram, got %T", d)
				}
				if len(er.Entities) != 1 {
					t.Errorf("expected 1 entity, got %d", len(er.Entities))
				}
				entity := er.Entities[0]
				if len(entity.Attributes) != 2 {
					t.Errorf("expected 2 attributes, got %d", len(entity.Attributes))
				}
				// Check asterisk notation was converted
				if entity.Attributes[0].Name != "userID" {
					t.Errorf("expected attribute name 'userID', got %q", entity.Attributes[0].Name)
				}
				if len(entity.Attributes[0].Keys) != 1 || entity.Attributes[0].Keys[0] != "PK" {
					t.Errorf("expected Keys [PK], got %v", entity.Attributes[0].Keys)
				}
			},
		},
		{
			name: "entity with composite keys",
			source: `erDiagram
    ORDER {
        int orderID PK,FK
        string orderDate
    }`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				er, ok := d.(*ast.ERDiagram)
				if !ok {
					t.Fatalf("expected *ast.ERDiagram, got %T", d)
				}
				entity := er.Entities[0]
				if len(entity.Attributes[0].Keys) != 2 {
					t.Errorf("expected 2 keys, got %d", len(entity.Attributes[0].Keys))
				}
			},
		},
		{
			name: "entity with attribute comments",
			source: `erDiagram
    PRODUCT {
        int productID PK "Primary identifier"
        string name "Product name"
    }`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				er, ok := d.(*ast.ERDiagram)
				if !ok {
					t.Fatalf("expected *ast.ERDiagram, got %T", d)
				}
				entity := er.Entities[0]
				if entity.Attributes[0].Comment != "Primary identifier" {
					t.Errorf("expected comment 'Primary identifier', got %q", entity.Attributes[0].Comment)
				}
			},
		},
		{
			name: "complete ER diagram",
			source: `erDiagram
    CUSTOMER ||--o{ ORDER : places
    ORDER ||--|{ LINE_ITEM : contains
    CUSTOMER {
        string name PK
        string email UK
        int customerID
    }
    ORDER {
        int orderID PK
        string orderDate
        int customerID FK
    }`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				er, ok := d.(*ast.ERDiagram)
				if !ok {
					t.Fatalf("expected *ast.ERDiagram, got %T", d)
				}
				if len(er.Entities) != 2 {
					t.Errorf("expected 2 entities, got %d", len(er.Entities))
				}
				if len(er.Relationships) != 2 {
					t.Errorf("expected 2 relationships, got %d", len(er.Relationships))
				}
			},
		},
		{
			name: "ER diagram with direction",
			source: `erDiagram LR
    CUSTOMER ||--o{ ORDER : places`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				er, ok := d.(*ast.ERDiagram)
				if !ok {
					t.Fatalf("expected *ast.ERDiagram, got %T", d)
				}
				if er.Direction != "LR" {
					t.Errorf("expected direction 'LR', got %q", er.Direction)
				}
			},
		},
		{
			name: "entity with alias",
			source: `erDiagram
    CUSTOMER [Customer Entity]`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				er, ok := d.(*ast.ERDiagram)
				if !ok {
					t.Fatalf("expected *ast.ERDiagram, got %T", d)
				}
				if len(er.Entities) != 1 {
					t.Errorf("expected 1 entity, got %d", len(er.Entities))
				}
				if er.Entities[0].Alias != "Customer Entity" {
					t.Errorf("expected alias 'Customer Entity', got %q", er.Entities[0].Alias)
				}
			},
		},
		{
			name: "non-identifying relationship",
			source: `erDiagram
    CUSTOMER }o..o{ ORDER : may_place`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				er, ok := d.(*ast.ERDiagram)
				if !ok {
					t.Fatalf("expected *ast.ERDiagram, got %T", d)
				}
				if er.Relationships[0].Type != ".." {
					t.Errorf("expected Type '..', got %q", er.Relationships[0].Type)
				}
			},
		},
		{
			name: "diagram with comments",
			source: `erDiagram
    %% This is a comment
    CUSTOMER ||--o{ ORDER : places
    %% Another comment`,
			wantErr: false,
			check: func(t *testing.T, d ast.Diagram) {
				er, ok := d.(*ast.ERDiagram)
				if !ok {
					t.Fatalf("expected *ast.ERDiagram, got %T", d)
				}
				if len(er.Relationships) != 1 {
					t.Errorf("expected 1 relationship, got %d", len(er.Relationships))
				}
			},
		},
		{
			name:    "invalid header",
			source:  "notErDiagram\n",
			wantErr: true,
		},
		{
			name:    "empty diagram",
			source:  "erDiagram\n",
			wantErr: false,
		},
		{
			name: "invalid syntax",
			source: `erDiagram
    this is invalid`,
			wantErr: true,
		},
	}

	p := parser.NewERParser()
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			diagram, err := p.Parse(tt.source)
			if (err != nil) != tt.wantErr {
				t.Errorf("Parse() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !tt.wantErr && tt.check != nil {
				tt.check(t, diagram)
			}
		})
	}
}

func TestERParser_SupportedTypes(t *testing.T) {
	p := parser.NewERParser()
	types := p.SupportedTypes()
	if len(types) != 1 || types[0] != "erDiagram" {
		t.Errorf("expected [erDiagram], got %v", types)
	}
}
