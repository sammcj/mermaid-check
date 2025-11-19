package validator_test

import (	"testing"

	"github.com/sammcj/go-mermaid/ast"

	"github.com/sammcj/go-mermaid/validator"
)

func TestValidateGitGraph(t *testing.T) {
	tests := []struct {
		name      string
		diagram   *ast.GitGraphDiagram
		strict    bool
		wantErrs  int
		checkErrs func(t *testing.T, errs []validator.ValidationError)
	}{
		{
			name: "valid git graph",
			diagram: &ast.GitGraphDiagram{
				Type: "gitGraph",
				Operations: []ast.GitOperation{
					{Type: "commit", ID: "Initial", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "branch", BranchName: "develop", Pos: ast.Position{Line: 3, Column: 1}},
					{Type: "checkout", BranchName: "develop", Pos: ast.Position{Line: 4, Column: 1}},
					{Type: "commit", ID: "Feature", Pos: ast.Position{Line: 5, Column: 1}},
					{Type: "checkout", BranchName: "main", Pos: ast.Position{Line: 6, Column: 1}},
					{Type: "merge", BranchName: "develop", Pos: ast.Position{Line: 7, Column: 1}},
				},
			},
			wantErrs: 0,
		},
		{
			name: "valid git graph with cherry-pick",
			diagram: &ast.GitGraphDiagram{
				Type: "gitGraph",
				Operations: []ast.GitOperation{
					{Type: "commit", ID: "Initial", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "commit", ID: "Feature", Pos: ast.Position{Line: 3, Column: 1}},
					{Type: "branch", BranchName: "hotfix", Pos: ast.Position{Line: 4, Column: 1}},
					{Type: "checkout", BranchName: "hotfix", Pos: ast.Position{Line: 5, Column: 1}},
					{Type: "cherry-pick", ParentID: "Feature", Pos: ast.Position{Line: 6, Column: 1}},
				},
			},
			wantErrs: 0,
		},
		{
			name: "valid commit types",
			diagram: &ast.GitGraphDiagram{
				Type: "gitGraph",
				Operations: []ast.GitOperation{
					{Type: "commit", ID: "C1", CommitType: "NORMAL", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "commit", ID: "C2", CommitType: "REVERSE", Pos: ast.Position{Line: 3, Column: 1}},
					{Type: "commit", ID: "C3", CommitType: "HIGHLIGHT", Pos: ast.Position{Line: 4, Column: 1}},
				},
			},
			wantErrs: 0,
		},
		{
			name: "duplicate branch names",
			diagram: &ast.GitGraphDiagram{
				Type: "gitGraph",
				Operations: []ast.GitOperation{
					{Type: "branch", BranchName: "feature", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "branch", BranchName: "feature", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrs: 1,
			checkErrs: func(t *testing.T, errs []validator.ValidationError) {
				t.Helper()
				if errs[0].Line != 3 {
					t.Errorf("expected error on line 3, got line %d", errs[0].Line)
				}
				if errs[0].Severity != validator.SeverityError {
					t.Errorf("expected error severity, got %s", errs[0].Severity)
				}
			},
		},
		{
			name: "checkout non-existent branch",
			diagram: &ast.GitGraphDiagram{
				Type: "gitGraph",
				Operations: []ast.GitOperation{
					{Type: "commit", ID: "Initial", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "checkout", BranchName: "nonexistent", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrs: 1,
			checkErrs: func(t *testing.T, errs []validator.ValidationError) {
				t.Helper()
				if errs[0].Line != 3 {
					t.Errorf("expected error on line 3, got line %d", errs[0].Line)
				}
			},
		},
		{
			name: "merge non-existent branch",
			diagram: &ast.GitGraphDiagram{
				Type: "gitGraph",
				Operations: []ast.GitOperation{
					{Type: "commit", ID: "Initial", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "merge", BranchName: "nonexistent", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrs: 1,
			checkErrs: func(t *testing.T, errs []validator.ValidationError) {
				t.Helper()
				if errs[0].Line != 3 {
					t.Errorf("expected error on line 3, got line %d", errs[0].Line)
				}
			},
		},
		{
			name: "cherry-pick non-existent commit",
			diagram: &ast.GitGraphDiagram{
				Type: "gitGraph",
				Operations: []ast.GitOperation{
					{Type: "commit", ID: "Initial", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "cherry-pick", ParentID: "nonexistent", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrs: 1,
			checkErrs: func(t *testing.T, errs []validator.ValidationError) {
				t.Helper()
				if errs[0].Line != 3 {
					t.Errorf("expected error on line 3, got line %d", errs[0].Line)
				}
			},
		},
		{
			name: "invalid commit type",
			diagram: &ast.GitGraphDiagram{
				Type: "gitGraph",
				Operations: []ast.GitOperation{
					{Type: "commit", ID: "C1", CommitType: "INVALID", Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErrs: 1,
			checkErrs: func(t *testing.T, errs []validator.ValidationError) {
				t.Helper()
				if errs[0].Line != 2 {
					t.Errorf("expected error on line 2, got line %d", errs[0].Line)
				}
			},
		},
		{
			name: "checkout main branch (should be valid)",
			diagram: &ast.GitGraphDiagram{
				Type: "gitGraph",
				Operations: []ast.GitOperation{
					{Type: "commit", ID: "Initial", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "checkout", BranchName: "main", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrs: 0,
		},
		{
			name: "checkout custom main branch",
			diagram: &ast.GitGraphDiagram{
				Type:           "gitGraph",
				MainBranchName: "trunk",
				Operations: []ast.GitOperation{
					{Type: "commit", ID: "Initial", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "checkout", BranchName: "trunk", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrs: 0,
		},
		{
			name: "multiple validation errors",
			diagram: &ast.GitGraphDiagram{
				Type: "gitGraph",
				Operations: []ast.GitOperation{
					{Type: "branch", BranchName: "feature", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "branch", BranchName: "feature", Pos: ast.Position{Line: 3, Column: 1}},
					{Type: "checkout", BranchName: "nonexistent", Pos: ast.Position{Line: 4, Column: 1}},
					{Type: "commit", ID: "C1", CommitType: "INVALID", Pos: ast.Position{Line: 5, Column: 1}},
				},
			},
			wantErrs: 3,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errsPtr := validator.ValidateGitGraph(tt.diagram, tt.strict)
			errs := make([]validator.ValidationError, len(errsPtr))
			for i, e := range errsPtr {
				errs[i] = *e
			}

			if len(errs) != tt.wantErrs {
				t.Errorf("expected %d errors, got %d", tt.wantErrs, len(errs))
				for i, err := range errs {
					t.Logf("  error %d: line %d: %s", i+1, err.Line, err.Message)
				}
			}

			if tt.checkErrs != nil && len(errs) > 0 {
				tt.checkErrs(t, errs)
			}
		})
	}
}

func TestNoDuplicateBranchNamesRule(t *testing.T) {
	rule := &validator.NoDuplicateBranchNamesRule{}

	tests := []struct {
		name     string
		diagram  *ast.GitGraphDiagram
		wantErrs int
	}{
		{
			name: "no duplicates",
			diagram: &ast.GitGraphDiagram{
				Operations: []ast.GitOperation{
					{Type: "branch", BranchName: "feature1", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "branch", BranchName: "feature2", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrs: 0,
		},
		{
			name: "duplicate branch names",
			diagram: &ast.GitGraphDiagram{
				Operations: []ast.GitOperation{
					{Type: "branch", BranchName: "feature", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "branch", BranchName: "feature", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrs: 1,
		},
		{
			name: "multiple duplicates",
			diagram: &ast.GitGraphDiagram{
				Operations: []ast.GitOperation{
					{Type: "branch", BranchName: "feature", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "branch", BranchName: "feature", Pos: ast.Position{Line: 3, Column: 1}},
					{Type: "branch", BranchName: "bugfix", Pos: ast.Position{Line: 4, Column: 1}},
					{Type: "branch", BranchName: "bugfix", Pos: ast.Position{Line: 5, Column: 1}},
				},
			},
			wantErrs: 2,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := rule.Validate(tt.diagram)
			if len(errs) != tt.wantErrs {
				t.Errorf("expected %d errors, got %d", tt.wantErrs, len(errs))
			}
		})
	}
}

func TestValidBranchReferencesRule(t *testing.T) {
	rule := &validator.ValidBranchReferencesRule{}

	tests := []struct {
		name     string
		diagram  *ast.GitGraphDiagram
		wantErrs int
	}{
		{
			name: "valid references",
			diagram: &ast.GitGraphDiagram{
				Operations: []ast.GitOperation{
					{Type: "branch", BranchName: "develop", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "checkout", BranchName: "develop", Pos: ast.Position{Line: 3, Column: 1}},
					{Type: "checkout", BranchName: "main", Pos: ast.Position{Line: 4, Column: 1}},
					{Type: "merge", BranchName: "develop", Pos: ast.Position{Line: 5, Column: 1}},
				},
			},
			wantErrs: 0,
		},
		{
			name: "invalid checkout reference",
			diagram: &ast.GitGraphDiagram{
				Operations: []ast.GitOperation{
					{Type: "checkout", BranchName: "nonexistent", Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErrs: 1,
		},
		{
			name: "invalid merge reference",
			diagram: &ast.GitGraphDiagram{
				Operations: []ast.GitOperation{
					{Type: "merge", BranchName: "nonexistent", Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErrs: 1,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := rule.Validate(tt.diagram)
			if len(errs) != tt.wantErrs {
				t.Errorf("expected %d errors, got %d", tt.wantErrs, len(errs))
			}
		})
	}
}

func TestValidCommitReferencesRule(t *testing.T) {
	rule := &validator.ValidCommitReferencesRule{}

	tests := []struct {
		name     string
		diagram  *ast.GitGraphDiagram
		wantErrs int
	}{
		{
			name: "valid cherry-pick reference",
			diagram: &ast.GitGraphDiagram{
				Operations: []ast.GitOperation{
					{Type: "commit", ID: "Feature", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "cherry-pick", ParentID: "Feature", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrs: 0,
		},
		{
			name: "invalid cherry-pick reference",
			diagram: &ast.GitGraphDiagram{
				Operations: []ast.GitOperation{
					{Type: "cherry-pick", ParentID: "nonexistent", Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErrs: 1,
		},
		{
			name: "cherry-pick referencing merge commit",
			diagram: &ast.GitGraphDiagram{
				Operations: []ast.GitOperation{
					{Type: "merge", BranchName: "develop", ID: "MergeCommit", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "cherry-pick", ParentID: "MergeCommit", Pos: ast.Position{Line: 3, Column: 1}},
				},
			},
			wantErrs: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := rule.Validate(tt.diagram)
			if len(errs) != tt.wantErrs {
				t.Errorf("expected %d errors, got %d", tt.wantErrs, len(errs))
			}
		})
	}
}

func TestValidCommitTypeRule(t *testing.T) {
	rule := &validator.ValidCommitTypeRule{}

	tests := []struct {
		name     string
		diagram  *ast.GitGraphDiagram
		wantErrs int
	}{
		{
			name: "valid commit types",
			diagram: &ast.GitGraphDiagram{
				Operations: []ast.GitOperation{
					{Type: "commit", CommitType: "NORMAL", Pos: ast.Position{Line: 2, Column: 1}},
					{Type: "commit", CommitType: "REVERSE", Pos: ast.Position{Line: 3, Column: 1}},
					{Type: "commit", CommitType: "HIGHLIGHT", Pos: ast.Position{Line: 4, Column: 1}},
				},
			},
			wantErrs: 0,
		},
		{
			name: "invalid commit type",
			diagram: &ast.GitGraphDiagram{
				Operations: []ast.GitOperation{
					{Type: "commit", CommitType: "INVALID", Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErrs: 1,
		},
		{
			name: "empty commit type (should be valid)",
			diagram: &ast.GitGraphDiagram{
				Operations: []ast.GitOperation{
					{Type: "commit", CommitType: "", Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErrs: 0,
		},
		{
			name: "valid merge commit type",
			diagram: &ast.GitGraphDiagram{
				Operations: []ast.GitOperation{
					{Type: "merge", BranchName: "develop", CommitType: "HIGHLIGHT", Pos: ast.Position{Line: 2, Column: 1}},
				},
			},
			wantErrs: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			errs := rule.Validate(tt.diagram)
			if len(errs) != tt.wantErrs {
				t.Errorf("expected %d errors, got %d", tt.wantErrs, len(errs))
			}
		})
	}
}
