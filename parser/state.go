package parser

import (
	"fmt"
	"regexp"
	"strings"

	"github.com/sammcj/mermaid-check/ast"
)

var (
	// State diagram patterns
	stateHeaderPattern  = regexp.MustCompile(`^(stateDiagram|stateDiagram-v2)\s*$`)
	stateCommentPattern = regexp.MustCompile(`^%%(.*)$`)

	// State declaration patterns
	stateDefPattern = regexp.MustCompile(`^state\s+"([^"]+)"\s+as\s+(\w+)\s*$`)

	// Composite state patterns (a state containing nested statements).
	// Both `state Name {` and `state "Description" as Name {` are supported.
	compositeStartPattern     = regexp.MustCompile(`^state\s+(\w+)\s*\{\s*$`)
	compositeStartDescPattern = regexp.MustCompile(`^state\s+"([^"]+)"\s+as\s+(\w+)\s*\{\s*$`)
	stateBodyEndPattern       = regexp.MustCompile(`^\}\s*$`)

	// Transition patterns
	transitionPattern = regexp.MustCompile(`^(\w+|\[\*\])\s+-->\s+(\w+|\[\*\])(?:\s*:\s*(.+))?\s*$`)

	// Special state patterns
	forkPattern   = regexp.MustCompile(`^state\s+(\w+)\s+<<fork>>\s*$`)
	joinPattern   = regexp.MustCompile(`^state\s+(\w+)\s+<<join>>\s*$`)
	choicePattern = regexp.MustCompile(`^state\s+(\w+)\s+<<choice>>\s*$`)

	// Note patterns
	stateNotePattern = regexp.MustCompile(`^note\s+(left|right)\s+of\s+(\w+)\s*:\s*(.+)\s*$`)
)

// StateParser parses Mermaid state diagrams.
type StateParser struct{}

// NewStateParser creates a new state diagram parser.
func NewStateParser() *StateParser {
	return &StateParser{}
}

// SupportedTypes returns the diagram types this parser handles.
func (p *StateParser) SupportedTypes() []string {
	return []string{"state", "stateDiagram-v2"}
}

// Parse parses a Mermaid state diagram from a string.
func (p *StateParser) Parse(source string) (ast.Diagram, error) {
	lines := strings.Split(source, "\n")
	if len(lines) == 0 {
		return nil, fmt.Errorf("empty diagram")
	}

	// Parse header
	header := strings.TrimSpace(lines[0])
	matches := stateHeaderPattern.FindStringSubmatch(header)
	if matches == nil {
		return nil, fmt.Errorf("invalid state diagram header: expected 'stateDiagram' or 'stateDiagram-v2'")
	}

	diagramType := "state"
	if matches[1] == "stateDiagram-v2" {
		diagramType = "stateDiagram-v2"
	}

	diagram := &ast.StateDiagram{
		Type:   diagramType,
		Source: source,
		Pos:    ast.Position{Line: 1, Column: 1},
	}

	// Parse statements
	statements := p.parseStatements(lines[1:], 1)
	diagram.Statements = statements

	return diagram, nil
}

func (p *StateParser) parseStatements(lines []string, startLine int) []ast.StateStmt {
	var statements []ast.StateStmt
	lineNum := startLine

	for i := 0; i < len(lines); i++ {
		lineNum++
		line := lines[i]
		trimmed := strings.TrimSpace(line)

		// Skip empty lines
		if trimmed == "" {
			continue
		}

		// Handle comments
		if matches := stateCommentPattern.FindStringSubmatch(trimmed); matches != nil {
			statements = append(statements, &ast.StateComment{
				Text: strings.TrimSpace(matches[1]),
				Pos:  ast.Position{Line: lineNum, Column: 1},
			})
			continue
		}

		// Handle composite states (a state containing nested statements).
		if state, consumed := p.parseComposite(trimmed, lines[i+1:], lineNum); state != nil {
			statements = append(statements, state)
			i += consumed
			lineNum += consumed
			continue
		}

		// Handle fork
		if matches := forkPattern.FindStringSubmatch(trimmed); matches != nil {
			statements = append(statements, &ast.Fork{
				ID:  matches[1],
				Pos: ast.Position{Line: lineNum, Column: 1},
			})
			continue
		}

		// Handle join
		if matches := joinPattern.FindStringSubmatch(trimmed); matches != nil {
			statements = append(statements, &ast.Join{
				ID:  matches[1],
				Pos: ast.Position{Line: lineNum, Column: 1},
			})
			continue
		}

		// Handle choice
		if matches := choicePattern.FindStringSubmatch(trimmed); matches != nil {
			statements = append(statements, &ast.Choice{
				ID:  matches[1],
				Pos: ast.Position{Line: lineNum, Column: 1},
			})
			continue
		}

		// Handle state with description
		if matches := stateDefPattern.FindStringSubmatch(trimmed); matches != nil {
			statements = append(statements, &ast.State{
				ID:          matches[2],
				Description: matches[1],
				Pos:         ast.Position{Line: lineNum, Column: 1},
			})
			continue
		}

		// Handle transitions
		if matches := transitionPattern.FindStringSubmatch(trimmed); matches != nil {
			from := matches[1]
			to := matches[2]
			label := ""
			if len(matches) > 3 {
				label = matches[3]
			}

			// Handle start state
			if from == "[*]" {
				statements = append(statements, &ast.StartState{
					To:  to,
					Pos: ast.Position{Line: lineNum, Column: 1},
				})
				continue
			}

			// Handle end state
			if to == "[*]" {
				statements = append(statements, &ast.EndState{
					From: from,
					Pos:  ast.Position{Line: lineNum, Column: 1},
				})
				continue
			}

			// Regular transition
			statements = append(statements, &ast.Transition{
				From:  from,
				To:    to,
				Label: label,
				Pos:   ast.Position{Line: lineNum, Column: 1},
			})
			continue
		}

		// Handle notes
		if matches := stateNotePattern.FindStringSubmatch(trimmed); matches != nil {
			statements = append(statements, &ast.StateNote{
				Position: matches[1] + " of",
				StateID:  matches[2],
				Text:     matches[3],
				Pos:      ast.Position{Line: lineNum, Column: 1},
			})
			continue
		}

		// Skip lines we can't parse
		continue
	}

	return statements
}

// parseComposite parses a composite state (`state Name {` ... `}`), returning the
// resulting State and the number of lines consumed after the opening line
// (including the closing brace). It returns nil, 0 when the header line does not
// open a composite state.
func (p *StateParser) parseComposite(header string, rest []string, headerLine int) (*ast.State, int) {
	var id, description string
	switch {
	case compositeStartDescPattern.MatchString(header):
		m := compositeStartDescPattern.FindStringSubmatch(header)
		description, id = m[1], m[2]
	case compositeStartPattern.MatchString(header):
		id = compositeStartPattern.FindStringSubmatch(header)[1]
	default:
		return nil, 0
	}

	// Find the matching closing brace, accounting for nested composite states.
	depth := 1
	end := -1
	for j, line := range rest {
		trimmed := strings.TrimSpace(line)
		switch {
		case compositeStartPattern.MatchString(trimmed), compositeStartDescPattern.MatchString(trimmed):
			depth++
		case stateBodyEndPattern.MatchString(trimmed):
			depth--
			if depth == 0 {
				end = j
			}
		}
		if end >= 0 {
			break
		}
	}

	if end < 0 {
		// Unclosed composite state: treat the remainder as its body.
		end = len(rest)
	}

	nested := p.parseStatements(rest[:end], headerLine)

	return &ast.State{
		ID:          id,
		Description: description,
		IsComposite: true,
		Nested:      nested,
		Pos:         ast.Position{Line: headerLine, Column: 1},
	}, end + 1
}
