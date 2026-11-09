package parser

import (
	"strings"
	"testing"
)

// TestUnquotedLinkLabelWithBracketsIsRejected checks that an unquoted |...|
// link label holding ( ) [ ] or { } is a parse error, as it is in mermaid.
func TestUnquotedLinkLabelWithBracketsIsRejected(t *testing.T) {
	cases := []struct{ name, source string }{
		{"paren", "flowchart LR\n  a[x] -->|argv: unset $(git rev-parse) so a hook| b[y]\n"},
		{"bracket", "flowchart LR\n  a[x] -->|run [--candidate] args| b[y]\n"},
		{"brace", "flowchart LR\n  a[x] -->|a {b} c| b[y]\n"},
		{"bidir paren", "flowchart LR\n  a[x] <-->|call (x)| b[y]\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := NewFlowchartParser()
			_, err := p.Parse(c.source)
			if err == nil {
				t.Fatalf("accepted a link label mermaid refuses:\n%s", strings.TrimSpace(c.source))
			}
		})
	}
}

// TestQuotedLinkLabelWithBracketsIsAccepted checks that the quoted form,
// |"..."|, still parses.
func TestQuotedLinkLabelWithBracketsIsAccepted(t *testing.T) {
	p := NewFlowchartParser()
	d, err := p.Parse("flowchart LR\n  a[x] -->|\"argv: unset $(git rev-parse) so a hook\"| b[y]\n")
	if err != nil {
		t.Fatalf("refused the quoted form mermaid renders: %v", err)
	}
	if d == nil {
		t.Fatal("no diagram")
	}
}
