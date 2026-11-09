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
		{"trailing semicolon", "flowchart LR\n  a -->|call (x)| b;\n"},
		{"endpoint label with parens", "flowchart LR\n  a[\"f(x)\"] -->|call (x)| b\n"},
		{"chained second label", "flowchart LR\n  a -->|ok| b -->|bad (x)| c\n"},
		{"bracket after quoted string", "flowchart LR\n  a -->|\"p\" (x) \"s\"| b\n"},
		{"second quoted string", "flowchart LR\n  a -->|\"p\" now \"s\"| b\n"},
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

// TestQuotedLinkLabelWithBracketsIsAccepted checks that labels mermaid renders
// still parse: brackets inside a leading quoted string, plain text after it,
// and arrow-like text inside quotes.
func TestQuotedLinkLabelWithBracketsIsAccepted(t *testing.T) {
	cases := []struct{ name, source string }{
		{"quoted", "flowchart LR\n  a[x] -->|\"argv: unset $(git rev-parse) so a hook\"| b[y]\n"},
		{"quoted then text", "flowchart LR\n  a -->|\"(x)\" now| b\n"},
		{"pipe inside quotes", "flowchart LR\n  a -->|\"a|b (c)\"| b\n"},
		{"markdown string", "flowchart LR\n  a -->|\"`(x)`\"| b\n"},
		{"arrow inside quoted node label", "flowchart LR\n  a[\"x -->|(y)|\"] --> b\n"},
		{"plain with semicolon", "flowchart LR\n  a -->|plain| b;\n"},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			p := NewFlowchartParser()
			d, err := p.Parse(c.source)
			if err != nil {
				t.Fatalf("refused a label mermaid renders: %v", err)
			}
			if d == nil {
				t.Fatal("no diagram")
			}
		})
	}
}
