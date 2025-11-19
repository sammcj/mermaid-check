package parser_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/parser"
)

func TestGitGraphParser_Parse(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantErr bool
		check   func(t *testing.T, diagram ast.Diagram)
	}{
		{
			name: "basic git graph with commits",
			input: `gitGraph
	commit id: "Initial"
	commit id: "Second"`,
			wantErr: false,
			check: func(t *testing.T, diagram ast.Diagram) {
				t.Helper()
				d, ok := diagram.(*ast.GitGraphDiagram)
				if !ok {
					t.Fatal("expected *ast.GitGraphDiagram")
				}
				if len(d.Operations) != 2 {
					t.Errorf("expected 2 operations, got %d", len(d.Operations))
				}
				if d.Operations[0].Type != "commit" || d.Operations[0].ID != "Initial" {
					t.Errorf("unexpected first commit: %+v", d.Operations[0])
				}
			},
		},
		{
			name: "git graph with branches",
			input: `gitGraph
	commit id: "Initial"
	branch develop
	checkout develop
	commit id: "Feature 1"`,
			wantErr: false,
			check: func(t *testing.T, diagram ast.Diagram) {
				t.Helper()
				d, ok := diagram.(*ast.GitGraphDiagram)
				if !ok {
					t.Fatal("expected *ast.GitGraphDiagram")
				}
				if len(d.Operations) != 4 {
					t.Errorf("expected 4 operations, got %d", len(d.Operations))
				}
				if d.Operations[1].Type != "branch" || d.Operations[1].BranchName != "develop" {
					t.Errorf("unexpected branch operation: %+v", d.Operations[1])
				}
				if d.Operations[2].Type != "checkout" || d.Operations[2].BranchName != "develop" {
					t.Errorf("unexpected checkout operation: %+v", d.Operations[2])
				}
			},
		},
		{
			name: "git graph with merge",
			input: `gitGraph
	commit id: "Initial"
	branch develop
	checkout develop
	commit id: "Feature"
	checkout main
	merge develop tag: "v1.0"`,
			wantErr: false,
			check: func(t *testing.T, diagram ast.Diagram) {
				t.Helper()
				d, ok := diagram.(*ast.GitGraphDiagram)
				if !ok {
					t.Fatal("expected *ast.GitGraphDiagram")
				}
				if len(d.Operations) != 6 {
					t.Errorf("expected 6 operations, got %d", len(d.Operations))
				}
				mergeOp := d.Operations[5]
				if mergeOp.Type != "merge" || mergeOp.BranchName != "develop" || mergeOp.Tag != "v1.0" {
					t.Errorf("unexpected merge operation: %+v", mergeOp)
				}
			},
		},
		{
			name: "git graph with cherry-pick",
			input: `gitGraph
	commit id: "Initial"
	commit id: "Feature"
	branch hotfix
	checkout hotfix
	cherry-pick id: "Feature" tag: "fix"`,
			wantErr: false,
			check: func(t *testing.T, diagram ast.Diagram) {
				t.Helper()
				d, ok := diagram.(*ast.GitGraphDiagram)
				if !ok {
					t.Fatal("expected *ast.GitGraphDiagram")
				}
				if len(d.Operations) != 5 {
					t.Errorf("expected 5 operations, got %d", len(d.Operations))
				}
				cherryOp := d.Operations[4]
				if cherryOp.Type != "cherry-pick" || cherryOp.ParentID != "Feature" || cherryOp.Tag != "fix" {
					t.Errorf("unexpected cherry-pick operation: %+v", cherryOp)
				}
			},
		},
		{
			name: "git graph with commit types",
			input: `gitGraph
	commit id: "Normal" type: NORMAL
	commit id: "Reverse" type: REVERSE
	commit id: "Highlight" type: HIGHLIGHT`,
			wantErr: false,
			check: func(t *testing.T, diagram ast.Diagram) {
				t.Helper()
				d, ok := diagram.(*ast.GitGraphDiagram)
				if !ok {
					t.Fatal("expected *ast.GitGraphDiagram")
				}
				if len(d.Operations) != 3 {
					t.Errorf("expected 3 operations, got %d", len(d.Operations))
				}
				if d.Operations[0].CommitType != "NORMAL" {
					t.Errorf("expected NORMAL commit type, got %s", d.Operations[0].CommitType)
				}
				if d.Operations[1].CommitType != "REVERSE" {
					t.Errorf("expected REVERSE commit type, got %s", d.Operations[1].CommitType)
				}
				if d.Operations[2].CommitType != "HIGHLIGHT" {
					t.Errorf("expected HIGHLIGHT commit type, got %s", d.Operations[2].CommitType)
				}
			},
		},
		{
			name: "git graph with branch order",
			input: `gitGraph
	commit id: "Initial"
	branch feature order: 2
	branch bugfix order: 1`,
			wantErr: false,
			check: func(t *testing.T, diagram ast.Diagram) {
				t.Helper()
				d, ok := diagram.(*ast.GitGraphDiagram)
				if !ok {
					t.Fatal("expected *ast.GitGraphDiagram")
				}
				if len(d.Operations) != 3 {
					t.Errorf("expected 3 operations, got %d", len(d.Operations))
				}
				if d.Operations[1].Order != 2 {
					t.Errorf("expected order 2, got %d", d.Operations[1].Order)
				}
				if d.Operations[2].Order != 1 {
					t.Errorf("expected order 1, got %d", d.Operations[2].Order)
				}
			},
		},
		{
			name: "git graph with theme configuration",
			input: `%%{init: { 'theme': 'base' } }%%
gitGraph
	commit id: "Initial"`,
			wantErr: false,
			check: func(t *testing.T, diagram ast.Diagram) {
				t.Helper()
				d, ok := diagram.(*ast.GitGraphDiagram)
				if !ok {
					t.Fatal("expected *ast.GitGraphDiagram")
				}
				if d.Theme != "base" {
					t.Errorf("expected theme 'base', got %s", d.Theme)
				}
			},
		},
		{
			name: "git graph with main branch options",
			input: `gitGraph
	mainBranchName: "trunk"
	mainBranchOrder: 0
	commit id: "Initial"`,
			wantErr: false,
			check: func(t *testing.T, diagram ast.Diagram) {
				t.Helper()
				d, ok := diagram.(*ast.GitGraphDiagram)
				if !ok {
					t.Fatal("expected *ast.GitGraphDiagram")
				}
				if d.MainBranchName != "trunk" {
					t.Errorf("expected mainBranchName 'trunk', got %s", d.MainBranchName)
				}
				if d.MainBranchOrder != 0 {
					t.Errorf("expected mainBranchOrder 0, got %d", d.MainBranchOrder)
				}
			},
		},
		{
			name: "git graph with comments",
			input: `gitGraph
	%% This is a comment
	commit id: "Initial"
	%% Another comment
	commit id: "Second"`,
			wantErr: false,
			check: func(t *testing.T, diagram ast.Diagram) {
				t.Helper()
				d, ok := diagram.(*ast.GitGraphDiagram)
				if !ok {
					t.Fatal("expected *ast.GitGraphDiagram")
				}
				if len(d.Operations) != 2 {
					t.Errorf("expected 2 operations (comments should be ignored), got %d", len(d.Operations))
				}
			},
		},
		{
			name: "commit without ID",
			input: `gitGraph
	commit
	commit tag: "v1.0"`,
			wantErr: false,
			check: func(t *testing.T, diagram ast.Diagram) {
				t.Helper()
				d, ok := diagram.(*ast.GitGraphDiagram)
				if !ok {
					t.Fatal("expected *ast.GitGraphDiagram")
				}
				if len(d.Operations) != 2 {
					t.Errorf("expected 2 operations, got %d", len(d.Operations))
				}
				if d.Operations[0].ID != "" {
					t.Errorf("expected empty ID for first commit, got %s", d.Operations[0].ID)
				}
				if d.Operations[1].Tag != "v1.0" {
					t.Errorf("expected tag 'v1.0', got %s", d.Operations[1].Tag)
				}
			},
		},
		{
			name:    "invalid header",
			input:   `git-graph\ncommit`,
			wantErr: true,
		},
		{
			name:    "empty diagram",
			input:   `gitGraph`,
			wantErr: true,
		},
		{
			name: "invalid syntax",
			input: `gitGraph
	invalid operation`,
			wantErr: true,
		},
		{
			name: "invalid branch order",
			input: `gitGraph
	branch feature order: invalid`,
			wantErr: true,
		},
		{
			name: "invalid main branch order",
			input: `gitGraph
	mainBranchOrder: invalid
	commit`,
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := parser.NewGitGraphParser()
			diagram, err := p.Parse(tt.input)

			if tt.wantErr {
				if err == nil {
					t.Error("expected error, got nil")
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			if diagram.GetType() != "gitGraph" {
				t.Errorf("expected type 'gitGraph', got %s", diagram.GetType())
			}

			if tt.check != nil {
				tt.check(t, diagram)
			}
		})
	}
}

func TestGitGraphParser_SupportedTypes(t *testing.T) {
	p := parser.NewGitGraphParser()
	types := p.SupportedTypes()
	if len(types) != 1 || types[0] != "gitGraph" {
		t.Errorf("expected [gitGraph], got %v", types)
	}
}
