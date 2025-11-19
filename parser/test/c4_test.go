package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestC4ContextParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(*testing.T, *ast.C4Diagram)
	}{
		{
			name: "simple context diagram",
			input: `C4Context
    title System Context for Internet Banking
    Person(customer, "Customer", "A customer of the bank")
    System(banking, "Internet Banking", "Allows customers to view information")
    System_Ext(email, "Email System", "Sends emails")
    Rel(customer, banking, "Uses")
    Rel(banking, email, "Sends emails using")`,
			wantErr: false,
			check: func(t *testing.T, d *ast.C4Diagram) {
				if d.DiagramType != "c4Context" {
					t.Errorf("expected diagram type 'c4Context', got '%s'", d.DiagramType)
				}
				if d.Title != "System Context for Internet Banking" {
					t.Errorf("expected title 'System Context for Internet Banking', got '%s'", d.Title)
				}
				if len(d.Elements) != 3 {
					t.Errorf("expected 3 elements, got %d", len(d.Elements))
				}
				if len(d.Relationships) != 2 {
					t.Errorf("expected 2 relationships, got %d", len(d.Relationships))
				}

				// Check Person element
				person := d.Elements[0]
				if person.ElementType != "Person" || person.ID != "customer" {
					t.Errorf("expected Person with ID 'customer', got %s with ID '%s'", person.ElementType, person.ID)
				}
				if person.External {
					t.Error("expected Person to not be external")
				}

				// Check System_Ext element
				emailSys := d.Elements[2]
				if emailSys.ElementType != "System" || !emailSys.External {
					t.Errorf("expected external System, got %s (external: %v)", emailSys.ElementType, emailSys.External)
				}
			},
		},
		{
			name: "context with boundaries",
			input: `C4Context
    System_Boundary(b1, "Boundary 1") {
        System(sys1, "System 1", "Description")
        System(sys2, "System 2", "Description")
    }
    Rel(sys1, sys2, "Uses")`,
			wantErr: false,
			check: func(t *testing.T, d *ast.C4Diagram) {
				if len(d.Boundaries) != 1 {
					t.Errorf("expected 1 boundary, got %d", len(d.Boundaries))
				}
				boundary := d.Boundaries[0]
				if boundary.BoundaryType != "System_Boundary" {
					t.Errorf("expected System_Boundary, got %s", boundary.BoundaryType)
				}
				if len(boundary.Elements) != 2 {
					t.Errorf("expected 2 elements in boundary, got %d", len(boundary.Elements))
				}
				if len(d.Relationships) != 1 {
					t.Errorf("expected 1 relationship, got %d", len(d.Relationships))
				}
			},
		},
		{
			name: "context with nested boundaries",
			input: `C4Context
    Enterprise_Boundary(e1, "Enterprise") {
        System_Boundary(s1, "System 1") {
            System(sys1, "System 1.1")
        }
        System_Boundary(s2, "System 2") {
            System(sys2, "System 2.1")
        }
    }`,
			wantErr: false,
			check: func(t *testing.T, d *ast.C4Diagram) {
				if len(d.Boundaries) != 1 {
					t.Errorf("expected 1 top-level boundary, got %d", len(d.Boundaries))
				}
				enterprise := d.Boundaries[0]
				if len(enterprise.Boundaries) != 2 {
					t.Errorf("expected 2 nested boundaries, got %d", len(enterprise.Boundaries))
				}
			},
		},
		{
			name: "context with all relationship types",
			input: `C4Context
    System(s1, "System 1")
    System(s2, "System 2")
    System(s3, "System 3")
    Rel(s1, s2, "Uses")
    Rel_Back(s1, s2, "Returns")
    Rel_Neighbor(s1, s3, "Communicates")
    Rel_Down(s1, s2, "Calls")
    BiRel(s2, s3, "Exchanges")`,
			wantErr: false,
			check: func(t *testing.T, d *ast.C4Diagram) {
				if len(d.Relationships) != 5 {
					t.Errorf("expected 5 relationships, got %d", len(d.Relationships))
				}
				relTypes := map[string]bool{
					"Rel": false, "Rel_Back": false, "Rel_Neighbor": false,
					"Rel_Down": false, "BiRel": false,
				}
				for _, rel := range d.Relationships {
					relTypes[rel.RelType] = true
				}
				for relType, found := range relTypes {
					if !found {
						t.Errorf("expected to find relationship type %s", relType)
					}
				}
			},
		},
		{
			name: "context with styles",
			input: `C4Context
    System(sys1, "System 1")
    System(sys2, "System 2")
    Rel(sys1, sys2, "Uses")
    UpdateElementStyle(sys1, "#ff0000", "#ffffff", "#000000", "true", "RoundedBoxShape")
    UpdateRelStyle(sys1, sys2, "#0000ff", "#00ff00", "0", "0")`,
			wantErr: false,
			check: func(t *testing.T, d *ast.C4Diagram) {
				if len(d.Styles) != 2 {
					t.Errorf("expected 2 styles, got %d", len(d.Styles))
				}
				elemStyle := d.Styles[0]
				if elemStyle.StyleType != "UpdateElementStyle" {
					t.Errorf("expected UpdateElementStyle, got %s", elemStyle.StyleType)
				}
				if elemStyle.BgColor != "#ff0000" {
					t.Errorf("expected bg colour #ff0000, got %s", elemStyle.BgColor)
				}
				relStyle := d.Styles[1]
				if relStyle.StyleType != "UpdateRelStyle" {
					t.Errorf("expected UpdateRelStyle, got %s", relStyle.StyleType)
				}
			},
		},
		{
			name: "context with comments",
			input: `C4Context
    %% This is a comment
    System(sys1, "System 1")
    %% Another comment
    System(sys2, "System 2")`,
			wantErr: false,
			check: func(t *testing.T, d *ast.C4Diagram) {
				if len(d.Elements) != 2 {
					t.Errorf("expected 2 elements (comments ignored), got %d", len(d.Elements))
				}
			},
		},
		{
			name: "context with all optional parameters",
			input: `C4Context
    Person(p1, "User", "A user of the system", "person", "tag1,tag2", "https://example.com")
    System(s1, "System", "A system", "sprite", "tag3", "https://system.com")`,
			wantErr: false,
			check: func(t *testing.T, d *ast.C4Diagram) {
				person := d.Elements[0]
				if person.Description != "A user of the system" {
					t.Errorf("expected description 'A user of the system', got '%s'", person.Description)
				}
				if person.Sprite != "person" {
					t.Errorf("expected sprite 'person', got '%s'", person.Sprite)
				}
				if person.Tags != "tag1,tag2" {
					t.Errorf("expected tags 'tag1,tag2', got '%s'", person.Tags)
				}
				if person.Link != "https://example.com" {
					t.Errorf("expected link 'https://example.com', got '%s'", person.Link)
				}
			},
		},
		{
			name:    "invalid header",
			input:   "C4Container\n    System(s1, \"System 1\")",
			wantErr: true,
		},
		{
			name:    "empty diagram",
			input:   "",
			wantErr: true,
		},
		{
			name: "unclosed boundary",
			input: `C4Context
    System_Boundary(b1, "Boundary") {
        System(s1, "System 1")`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewC4ContextParser()
			diagram, err := p.Parse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			c4Diagram, ok := diagram.(*ast.C4Diagram)
			if !ok {
				t.Errorf("expected *ast.C4Diagram, got %T", diagram)
				return
			}

			if tt.check != nil {
				tt.check(t, c4Diagram)
			}
		})
	}
}

func TestC4ContainerParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(*testing.T, *ast.C4Diagram)
	}{
		{
			name: "simple container diagram",
			input: `C4Container
    title Container Diagram for Internet Banking
    Person(customer, "Customer")
    System_Boundary(c1, "Internet Banking") {
        Container(web, "Web Application", "Java, Spring MVC", "Delivers content")
        ContainerDb(db, "Database", "MySQL", "Stores user data")
        ContainerQueue(queue, "Message Queue", "RabbitMQ", "Handles async")
    }
    Rel(customer, web, "Uses", "HTTPS")
    Rel(web, db, "Reads/Writes", "JDBC")`,
			wantErr: false,
			check: func(t *testing.T, d *ast.C4Diagram) {
				if d.DiagramType != "c4Container" {
					t.Errorf("expected diagram type 'c4Container', got '%s'", d.DiagramType)
				}
				if len(d.Boundaries) != 1 {
					t.Errorf("expected 1 boundary, got %d", len(d.Boundaries))
				}
				boundary := d.Boundaries[0]
				if len(boundary.Elements) != 3 {
					t.Errorf("expected 3 containers in boundary, got %d", len(boundary.Elements))
				}

				// Check Container types
				web := boundary.Elements[0]
				if web.ElementType != "Container" || web.Database || web.Queue {
					t.Error("expected regular Container")
				}
				db := boundary.Elements[1]
				if !db.Database {
					t.Error("expected ContainerDb to have Database flag")
				}
				queue := boundary.Elements[2]
				if !queue.Queue {
					t.Error("expected ContainerQueue to have Queue flag")
				}

				// Check technology field
				if web.Technology != "Java, Spring MVC" {
					t.Errorf("expected technology 'Java, Spring MVC', got '%s'", web.Technology)
				}
			},
		},
		{
			name:    "invalid header",
			input:   "C4Context\n    Container(c1, \"Container 1\")",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewC4ContainerParser()
			diagram, err := p.Parse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			c4Diagram, ok := diagram.(*ast.C4Diagram)
			if !ok {
				t.Errorf("expected *ast.C4Diagram, got %T", diagram)
				return
			}

			if tt.check != nil {
				tt.check(t, c4Diagram)
			}
		})
	}
}

func TestC4ComponentParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(*testing.T, *ast.C4Diagram)
	}{
		{
			name: "simple component diagram",
			input: `C4Component
    title Component Diagram for Web Application
    Container_Boundary(c1, "Web Application") {
        Component(controller, "MVC Controller", "Spring MVC", "Handles requests")
        Component(service, "Business Logic", "Java", "Implements logic")
        ComponentDb(cache, "Cache", "Redis", "Caches data")
    }
    Rel(controller, service, "Uses")
    Rel(service, cache, "Reads/Writes")`,
			wantErr: false,
			check: func(t *testing.T, d *ast.C4Diagram) {
				if d.DiagramType != "c4Component" {
					t.Errorf("expected diagram type 'c4Component', got '%s'", d.DiagramType)
				}
				boundary := d.Boundaries[0]
				if boundary.BoundaryType != "Container_Boundary" {
					t.Errorf("expected Container_Boundary, got %s", boundary.BoundaryType)
				}

				// Check Component types
				controller := boundary.Elements[0]
				if controller.ElementType != "Component" {
					t.Errorf("expected Component, got %s", controller.ElementType)
				}
				cache := boundary.Elements[2]
				if !cache.Database {
					t.Error("expected ComponentDb to have Database flag")
				}
			},
		},
		{
			name:    "invalid header",
			input:   "C4Container\n    Component(c1, \"Component 1\")",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewC4ComponentParser()
			diagram, err := p.Parse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			c4Diagram, ok := diagram.(*ast.C4Diagram)
			if !ok {
				t.Errorf("expected *ast.C4Diagram, got %T", diagram)
				return
			}

			if tt.check != nil {
				tt.check(t, c4Diagram)
			}
		})
	}
}

func TestC4DynamicParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(*testing.T, *ast.C4Diagram)
	}{
		{
			name: "simple dynamic diagram",
			input: `C4Dynamic
    title Dynamic Diagram for User Login
    Person(user, "User")
    System(auth, "Auth Service")
    System(db, "Database")
    Rel(user, auth, "1. Login request")
    Rel(auth, db, "2. Verify credentials")
    Rel(db, auth, "3. Return result")
    Rel(auth, user, "4. Login response")`,
			wantErr: false,
			check: func(t *testing.T, d *ast.C4Diagram) {
				if d.DiagramType != "c4Dynamic" {
					t.Errorf("expected diagram type 'c4Dynamic', got '%s'", d.DiagramType)
				}
				if len(d.Relationships) != 4 {
					t.Errorf("expected 4 relationships, got %d", len(d.Relationships))
				}
			},
		},
		{
			name:    "invalid header",
			input:   "C4Context\n    System(s1, \"System 1\")",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewC4DynamicParser()
			diagram, err := p.Parse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			c4Diagram, ok := diagram.(*ast.C4Diagram)
			if !ok {
				t.Errorf("expected *ast.C4Diagram, got %T", diagram)
				return
			}

			if tt.check != nil {
				tt.check(t, c4Diagram)
			}
		})
	}
}

func TestC4DeploymentParser(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(*testing.T, *ast.C4Diagram)
	}{
		{
			name: "simple deployment diagram",
			input: `C4Deployment
    title Deployment Diagram for Production
    Deployment_Node(aws, "AWS Cloud", "Cloud Provider") {
        Deployment_Node(region, "EU-West-1", "AWS Region") {
            Node(server, "Web Server", "EC2", "Runs application")
            Node(db, "Database", "RDS", "Stores data")
        }
    }
    Rel(server, db, "Connects to")`,
			wantErr: false,
			check: func(t *testing.T, d *ast.C4Diagram) {
				if d.DiagramType != "c4Deployment" {
					t.Errorf("expected diagram type 'c4Deployment', got '%s'", d.DiagramType)
				}
				if len(d.Boundaries) != 1 {
					t.Errorf("expected 1 top-level boundary, got %d", len(d.Boundaries))
				}
				topBoundary := d.Boundaries[0]
				if len(topBoundary.Boundaries) != 1 {
					t.Errorf("expected 1 nested boundary, got %d", len(topBoundary.Boundaries))
				}
				nestedBoundary := topBoundary.Boundaries[0]
				if len(nestedBoundary.Elements) != 2 {
					t.Errorf("expected 2 nodes in nested boundary, got %d", len(nestedBoundary.Elements))
				}
			},
		},
		{
			name:    "invalid header",
			input:   "C4Context\n    Node(n1, \"Node 1\")",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewC4DeploymentParser()
			diagram, err := p.Parse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			c4Diagram, ok := diagram.(*ast.C4Diagram)
			if !ok {
				t.Errorf("expected *ast.C4Diagram, got %T", diagram)
				return
			}

			if tt.check != nil {
				tt.check(t, c4Diagram)
			}
		})
	}
}

// NOTE: TestParseC4Parameters is commented out because parseC4Parameters is an unexported function
// and this file uses black-box testing (package parser_test).
// This test should be moved to a white-box test file if needed.
/*
func TestParseC4Parameters(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected []string
	}{
		{
			name:     "simple parameters",
			input:    `id, "Label", "Description"`,
			expected: []string{"id", "Label", "Description"},
		},
		{
			name:     "parameters with commas in quotes",
			input:    `id, "Label, with comma", "Description"`,
			expected: []string{"id", "Label, with comma", "Description"},
		},
		{
			name:     "parameters with escaped quotes",
			input:    `id, "Label with \"quotes\"", "Description"`,
			expected: []string{"id", `Label with "quotes"`, "Description"},
		},
		{
			name:     "empty optional parameters",
			input:    `id, "Label", "", "", "tag"`,
			expected: []string{"id", "Label", "", "", "tag"},
		},
		{
			name:     "single parameter",
			input:    `id`,
			expected: []string{"id"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseC4Parameters(tt.input)
			if len(result) != len(tt.expected) {
				t.Errorf("expected %d parameters, got %d", len(tt.expected), len(result))
				return
			}
			for i, expected := range tt.expected {
				if result[i] != expected {
					t.Errorf("parameter %d: expected '%s', got '%s'", i, expected, result[i])
				}
			}
		})
	}
}
*/
