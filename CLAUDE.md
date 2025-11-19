# CLAUDE.md

Project-specific guidance for AI coding agents working with go-mermaid.

## Project Overview

Pure Go parser, validator, and linter for Mermaid diagrams supporting all 21+ diagram types with full AST parsing and semantic validation. Under active development, not production-ready.

## Essential Commands

```bash
make build     # Build CLI → ./go-mermaid
make test      # Run tests with race detection
make lint      # golangci-lint (must pass with zero warnings)
make coverage  # Generate HTML coverage report
```

## Architecture Overview

**Three-layer design:**
- `ast/` - AST definitions for each diagram type implementing `Diagram` interface
- `parser/` - Auto-detecting parser registry, each type has dedicated parser
- `validator/` - Composable validation rules (DefaultRules/StrictRules per type)

**Key files:**
- `mermaid.go` - Public API (`Parse()`, `Validate()`, `ParseFile()`)
- `parser/parser.go` - Parser registry with `diagramTypeMapping` table
- `validator/utilities.go` - Shared utilities: `DuplicateChecker`, `ReferenceChecker`, `EnumValidator`

**Input handling:**
- `.mmd` files → raw Mermaid (auto-detects markdown fences)
- `.md` files → extracts `​```mermaid` blocks
- `extractor/markdown.go` handles unclosed blocks at EOF

## Critical Patterns

**Parser behaviour:**
- Auto-detects type via first non-comment line against `diagramTypeMapping`
- Returns `ast.Diagram` interface (use type assertion for specific types)
- All parsers implement `DiagramParser` interface
- Position tracking is 1-indexed for line/column

**Validation:**
- Parse errors → `error` (unparseable)
- Validation errors → `[]ValidationError` (parsed but invalid)
- Each diagram type has own rule interface (`Rule`, `SequenceRule`, etc.)

**Testing:**
- Black-box tests in `{module}/test/` directories (e.g., `parser/test/`)
- Table-driven tests preferred
- Test data in `testdata/{type}/` directories
- Maintain 75%+ coverage

## Code Standards

- **British English**: All code, comments, docs (colour, behaviour, optimise)
- **Max complexity**: 10 (enforced by golangci-lint)
- **Go version**: 1.23.4+
- All exported functions require comments
- Zero linter warnings before merge

## Critical Gotchas

**Mermaid syntax:**
- Comments use `%%`, not `//` or `#`
- Position tracking is 1-indexed (Line and Column)
- Subgraphs support both `subgraph "name"` and `subgraph ID[Display Name]` syntax
- BiDir links use different regex than unidirectional
- Cannot use `()` in labels; use `<br>` for line breaks

**Validation:**
- **IMPORTANT**: Validate Mermaid syntax against official docs before "fixing" parsers
- Reference: `$HOME/git/mermaid/packages/mermaid/src/docs/` or https://mermaid.js.org/
- When debugging issues don't assume files containing mermaid diagrams always contain valid Mermaid syntax, validate first before making code changes

**Code structure:**
- GenericDiagram stores both `Source` (raw) and `Lines` (split)
- Handle subgraphs recursively in validators
- Regex patterns cached as package-level variables in parsers

## Adding New Diagram Types

Study existing implementations (sequence.go, gantt.go, c4.go) and follow the pattern:

1. Define AST in `ast/{type}.go` with `Diagram` interface
2. Create parser in `parser/{type}.go` with Position tracking
3. Register in `parser/parser.go` mapping table
4. Add validators in `validator/{type}.go` with rule interface
5. Update `mermaid.go` Validate() switch
6. Add black-box tests in `parser/test/` and `validator/test/`
7. Ensure `make test && make lint` passes

Target: <100ms parsing for diagrams up to 1000 lines.
