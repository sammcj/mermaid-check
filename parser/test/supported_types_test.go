package parser_test

import (
	"slices"
	"testing"

	"github.com/sammcj/go-mermaid/parser"
)

// TestParserSupportedTypes ensures parsers report their supported diagram types correctly.
func TestParserSupportedTypes(t *testing.T) {
	tests := []struct {
		name         string
		parser       interface{ SupportedTypes() []string }
		wantContains string
	}{
		{
			name:         "FlowchartParser",
			parser:       parser.NewFlowchartParser(),
			wantContains: "flowchart",
		},
		{
			name:         "C4ComponentParser",
			parser:       parser.NewC4ComponentParser(),
			wantContains: "c4Component",
		},
		{
			name:         "C4ContainerParser",
			parser:       parser.NewC4ContainerParser(),
			wantContains: "c4Container",
		},
		{
			name:         "C4DeploymentParser",
			parser:       parser.NewC4DeploymentParser(),
			wantContains: "c4Deployment",
		},
		{
			name:         "C4DynamicParser",
			parser:       parser.NewC4DynamicParser(),
			wantContains: "c4Dynamic",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			types := tt.parser.SupportedTypes()
			if len(types) == 0 {
				t.Error("SupportedTypes() returned empty slice")
				return
			}

			found := slices.Contains(types, tt.wantContains)

			if !found {
				t.Errorf("SupportedTypes() = %v, want to contain %v", types, tt.wantContains)
			}
		})
	}
}
