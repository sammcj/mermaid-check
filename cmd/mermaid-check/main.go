// Command mermaid-check provides a CLI tool for parsing, validating, and linting Mermaid diagrams.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"
	"strings"

	"github.com/fatih/color"
	mermaid "github.com/sammcj/mermaid-check"
	"github.com/sammcj/mermaid-check/ast"
	"github.com/sammcj/mermaid-check/extractor"
	"github.com/sammcj/mermaid-check/internal/inpututil"
)

const version = "0.1.0"

var (
	// Colour definitions for clean, modern output
	green  = color.New(color.FgGreen).SprintFunc()
	red    = color.New(color.FgRed).SprintFunc()
	yellow = color.New(color.FgYellow).SprintFunc()
	orange = color.New(color.FgYellow, color.Bold).SprintFunc()
	cyan   = color.New(color.FgCyan).SprintFunc()
	bold   = color.New(color.Bold).SprintFunc()
	dim    = color.New(color.Faint).SprintFunc()
)

func main() {
	// Define flags
	var (
		strict       = flag.Bool("strict", false, "use strict validation rules")
		formatFlag   = flag.String("format", "", "force input format (mermaid or markdown)")
		errorOnEmpty = flag.Bool("error-on-empty", false, "treat files with no Mermaid diagrams as errors")
		showHelp     = flag.Bool("help", false, "show help message")
		showVersion  = flag.Bool("version", false, "show version")
	)

	flag.Parse()

	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("mermaid-check version %s\n", version)
		os.Exit(0)
	}

	// Determine input source
	args := flag.Args()
	var exitCode int

	if len(args) == 0 {
		// Read from stdin
		exitCode = processStdin(*formatFlag, *strict, *errorOnEmpty)
	} else {
		// Process files
		exitCode = processFiles(args, *strict, *errorOnEmpty)
	}

	os.Exit(exitCode)
}

func processStdin(format string, strict bool, errorOnEmpty bool) int {
	data, err := io.ReadAll(os.Stdin)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading stdin: %v\n", err)
		return 1
	}

	content := string(data)

	// Determine format
	isMarkdown := format == "markdown"
	if format == "" {
		// Try to auto-detect - if it looks like markdown (has code blocks), treat as markdown
		isMarkdown = containsCodeBlocks(content)
	}

	var hasErrors bool

	if isMarkdown {
		// Extract and validate Mermaid blocks from markdown
		blocks, err := extractor.ExtractFromMarkdown(content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Error extracting Mermaid blocks: %v\n", err)
			return 1
		}

		if len(blocks) == 0 {
			fmt.Fprintf(os.Stderr, "No Mermaid diagrams found in markdown\n")
			fmt.Fprintf(os.Stderr, "Hint: Ensure code blocks use proper markdown fences: ```mermaid\n")
			if errorOnEmpty {
				return 1
			}
			return 0
		}

		fmt.Printf("Found %d Mermaid diagram(s)\n", len(blocks))

		// Collect statistics
		stats := make(map[string]int)
		for i, block := range blocks {
			displayName := diagramTypeDisplayName(block.DiagramType)
			fmt.Printf("\n--- Diagram %d - %s (%s, line %d) ---\n", i+1, displayName, block.DiagramType, block.LineOffset)
			stats[block.DiagramType]++
			if processBlock(&block, strict) {
				hasErrors = true
			}
		}

		// Print summary statistics
		if len(blocks) > 1 {
			fmt.Printf("\nDiagram type distribution:\n")
			for diagramType, count := range stats {
				fmt.Printf("  %s: %d\n", diagramType, count)
			}
		}
	} else {
		// Parse as raw Mermaid
		diagram, err := mermaid.Parse(content)
		if err != nil {
			fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
			return 1
		}

		diagramType := diagram.GetType()
		displayName := diagramTypeDisplayName(diagramType)
		fmt.Printf("Diagram type: %s (%s)\n", displayName, diagramType)
		if validateDiagram(diagram, strict, "") {
			hasErrors = true
		}
	}

	if hasErrors {
		return 1
	}
	return 0
}

// fileResult represents the validation result for a single file
type fileResult struct {
	path         string
	resultType   resultType
	diagramCount int
	blocks       []blockResult
	stats        map[string]int
	errorMsg     string
}

// blockResult represents the validation result for a single diagram block
type blockResult struct {
	diagramType string
	lineRange   string
	isValid     bool
	errors      []string
	blockNum    int
}

type resultType int

const (
	resultNoDiagrams resultType = iota
	resultSuccess
	resultValidationError
	resultParseError
	resultFileError
	resultUnsupportedType
)

func processFiles(paths []string, strict bool, errorOnEmpty bool) int {
	var hasErrors bool
	results := make([]fileResult, 0, len(paths))

	// Collect all results first
	for _, path := range paths {
		result := fileResult{
			path:   path,
			blocks: make([]blockResult, 0),
		}

		// Read file content
		data, err := os.ReadFile(path)
		if err != nil {
			result.resultType = resultFileError
			result.errorMsg = err.Error()
			results = append(results, result)
			hasErrors = true
			continue
		}

		content := string(data)
		fileType := inpututil.DetectFileType(path)

		// Check if .mmd file contains markdown code fences
		if fileType == inpututil.FileTypeMermaid && containsMarkdownFences(content) {
			fileType = inpututil.FileTypeMarkdown
		}

		switch fileType {
		case inpututil.FileTypeMarkdown:
			// Extract blocks from markdown to preserve line information
			blocks, err := extractor.ExtractFromMarkdown(content)
			if err != nil {
				result.resultType = resultParseError
				result.errorMsg = err.Error()
				results = append(results, result)
				hasErrors = true
				continue
			}

			if len(blocks) == 0 {
				result.resultType = resultNoDiagrams
				results = append(results, result)
				// Only treat as error if errorOnEmpty flag is set (markdown files are optional)
				if errorOnEmpty {
					hasErrors = true
				}
				continue
			}

			result.diagramCount = len(blocks)
			result.stats = make(map[string]int)
			hasValidationErrors := false

			for i, block := range blocks {
				result.stats[block.DiagramType]++
				lineRange := fmt.Sprintf("(L%d-L%d)", block.LineOffset, block.EndLine)

				blockRes := blockResult{
					diagramType: block.DiagramType,
					lineRange:   lineRange,
					blockNum:    i + 1,
				}

				diagram, err := mermaid.Parse(block.Source)
				if err != nil {
					blockRes.isValid = false
					blockRes.errors = []string{fmt.Sprintf("parse error: %v", err)}
					result.blocks = append(result.blocks, blockRes)
					hasValidationErrors = true
					continue
				}

				validationErrors := mermaid.Validate(diagram, strict)
				if len(validationErrors) == 0 {
					blockRes.isValid = true
				} else {
					blockRes.isValid = false
					for _, ve := range validationErrors {
						blockRes.errors = append(blockRes.errors, ve.Error())
					}
					hasValidationErrors = true
				}

				result.blocks = append(result.blocks, blockRes)
			}

			if hasValidationErrors {
				result.resultType = resultValidationError
				hasErrors = true
			} else {
				result.resultType = resultSuccess
			}

		case inpututil.FileTypeMermaid:
			// For .mmd files, check if content is empty or whitespace-only
			var trimmedContent strings.Builder
			for _, ch := range content {
				if ch != ' ' && ch != '\t' && ch != '\n' && ch != '\r' {
					trimmedContent.WriteString(string(ch))
				}
			}
			if trimmedContent.String() == "" {
				result.resultType = resultNoDiagrams
				result.errorMsg = "empty .mmd file"
				results = append(results, result)
				hasErrors = true // .mmd files should always contain Mermaid
				continue
			}

			// Parse as raw Mermaid
			diagram, err := mermaid.Parse(content)
			if err != nil {
				result.resultType = resultParseError
				result.errorMsg = err.Error()
				results = append(results, result)
				hasErrors = true
				continue
			}

			diagramType := diagram.GetType()
			result.stats = map[string]int{diagramType: 1}
			result.diagramCount = 1

			blockRes := blockResult{
				diagramType: diagramType,
				blockNum:    1,
			}

			validationErrors := mermaid.Validate(diagram, strict)
			if len(validationErrors) == 0 {
				blockRes.isValid = true
				result.resultType = resultSuccess
			} else {
				blockRes.isValid = false
				for _, ve := range validationErrors {
					blockRes.errors = append(blockRes.errors, ve.Error())
				}
				result.resultType = resultValidationError
				hasErrors = true
			}

			result.blocks = append(result.blocks, blockRes)

		default:
			result.resultType = resultUnsupportedType
			results = append(results, result)
			hasErrors = true
			continue
		}

		results = append(results, result)
	}

	// Output results grouped by type
	printGroupedResults(results, errorOnEmpty)

	if hasErrors {
		return 1
	}
	return 0
}

func printGroupedResults(results []fileResult, errorOnEmpty bool) {
	// Group results by type
	noDiagramsInfo := make([]fileResult, 0)  // informational (markdown with no diagrams)
	noDiagramsError := make([]fileResult, 0) // errors (empty .mmd files)
	parseErrors := make([]fileResult, 0)
	fileErrors := make([]fileResult, 0)
	unsupported := make([]fileResult, 0)
	others := make([]fileResult, 0)

	for _, r := range results {
		switch r.resultType {
		case resultNoDiagrams:
			// Empty .mmd files have errorMsg set and should be errors
			if r.errorMsg != "" {
				noDiagramsError = append(noDiagramsError, r)
			} else {
				noDiagramsInfo = append(noDiagramsInfo, r)
			}
		case resultParseError:
			parseErrors = append(parseErrors, r)
		case resultFileError:
			fileErrors = append(fileErrors, r)
		case resultUnsupportedType:
			unsupported = append(unsupported, r)
		default:
			others = append(others, r)
		}
	}

	// Print "no diagrams found" group first (informational for markdown files)
	if len(noDiagramsInfo) > 0 {
		if errorOnEmpty {
			fmt.Printf("\n%s\n", bold(red("Files with no Mermaid diagrams found:")))
		} else {
			fmt.Printf("\n%s\n", orange("Files with no Mermaid diagrams:"))
		}
		for _, r := range noDiagramsInfo {
			if errorOnEmpty {
				fmt.Printf("  %s %s\n", red("✗"), red(r.path))
			} else {
				fmt.Printf("  %s %s\n", orange("⚠"), orange(r.path))
			}
		}
		fmt.Printf("  %s Ensure code blocks use proper markdown fences: ```mermaid\n\n", dim("Hint:"))
	}

	// Print errors group
	hasErrors := len(noDiagramsError) > 0 || len(parseErrors) > 0 || len(fileErrors) > 0 || len(unsupported) > 0
	if hasErrors {
		fmt.Printf("%s\n", bold(red("Errors:")))

		// Print empty .mmd file errors
		if len(noDiagramsError) > 0 {
			for _, r := range noDiagramsError {
				fmt.Printf("  %s %s: %s\n", red("✗"), red(r.path), r.errorMsg)
			}
		}

		// Print parse errors
		if len(parseErrors) > 0 {
			for _, r := range parseErrors {
				fmt.Printf("  %s %s: %s\n", red("✗"), red(r.path), r.errorMsg)
			}
		}

		// Print file errors
		if len(fileErrors) > 0 {
			for _, r := range fileErrors {
				fmt.Printf("  %s %s: %s\n", red("✗"), red(r.path), r.errorMsg)
			}
		}

		// Print unsupported file types
		if len(unsupported) > 0 {
			for _, r := range unsupported {
				fmt.Printf("  %s %s: unsupported file type\n", red("✗"), red(r.path))
			}
		}
		fmt.Println()
	}

	// Print all other results (successful validations and validation errors)
	if len(others) > 0 {
		for _, r := range others {
			fmt.Printf("\n%s %s\n", bold("Validating:"), cyan(r.path))

			if r.diagramCount > 1 {
				fmt.Printf("  %s %d diagrams\n", dim("Found"), r.diagramCount)
			}

			for _, block := range r.blocks {
				var prefix string
				displayName := diagramTypeDisplayName(block.diagramType)
				if r.diagramCount > 1 {
					prefix = fmt.Sprintf("  %s - %s %s: ",
						bold(fmt.Sprintf("Diagram %d", block.blockNum)),
						cyan(displayName),
						dim(block.lineRange))
				} else {
					prefix = fmt.Sprintf("  %s ", cyan(displayName))
					if block.lineRange != "" {
						prefix += fmt.Sprintf("%s - ", dim(block.lineRange))
					}
				}

				if block.isValid {
					fmt.Printf("%s%s %s\n", prefix, green("✓"), dim("Valid"))
				} else {
					fmt.Printf("%s%s %s:\n", prefix, red("✗"), red(fmt.Sprintf("%d validation error(s)", len(block.errors))))
					for _, errMsg := range block.errors {
						fmt.Printf("%s  %s\n", prefix, yellow(errMsg))
					}
				}
			}

			// Print summary statistics for files with multiple diagrams
			if r.diagramCount > 1 {
				fmt.Printf("\n  %s\n", bold("Diagram type distribution:"))
				for diagramType, count := range r.stats {
					fmt.Printf("    %s %d\n", cyan(diagramType+":"), count)
				}
			}
		}
	}
}

func processBlock(block *extractor.DiagramBlock, strict bool) bool {
	diagram, err := mermaid.Parse(block.Source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Parse error: %v\n", err)
		return true
	}

	return validateDiagram(diagram, strict, "")
}

func validateDiagram(diagram ast.Diagram, strict bool, prefix string) bool {
	errors := mermaid.Validate(diagram, strict)

	if len(errors) == 0 {
		fmt.Printf("%s%s %s\n", prefix, green("✓"), dim("Valid"))
		return false
	}

	fmt.Printf("%s%s %s:\n", prefix, red("✗"), red(fmt.Sprintf("%d validation error(s)", len(errors))))
	for _, err := range errors {
		fmt.Printf("%s  %s\n", prefix, yellow(fmt.Sprintf("%v", err)))
	}

	return true
}

func containsCodeBlocks(content string) bool {
	return len(content) > 10 && (contains(content, "```mermaid") ||
		contains(content, "```\nmermaid") ||
		contains(content, "# ") || // Markdown heading
		contains(content, "## "))
}

func containsMarkdownFences(content string) bool {
	return len(content) > 10 && (contains(content, "```mermaid") ||
		contains(content, "~~~mermaid") ||
		contains(content, "``` mermaid"))
}

// diagramTypeDisplayName returns a user-friendly display name for a diagram type.
func diagramTypeDisplayName(diagType string) string {
	displayNames := map[string]string{
		"flowchart":       "Flowchart",
		"graph":           "Flow Chart",
		"sequence":        "Sequence Diagram",
		"class":           "Class Diagram",
		"state":           "State Diagram",
		"stateDiagram-v2": "State Diagram",
		"er":              "ER Diagram",
		"gantt":           "Gantt Chart",
		"pie":             "Pie Chart",
		"journey":         "User Journey",
		"timeline":        "Timeline",
		"gitGraph":        "Git Graph",
		"mindmap":         "Mindmap",
		"sankey":          "Sankey Diagram",
		"quadrantChart":   "Quadrant Chart",
		"xyChart":         "XY Chart",
		"c4Context":       "C4 Context Diagram",
		"c4Container":     "C4 Container Diagram",
		"c4Component":     "C4 Component Diagram",
		"c4Dynamic":       "C4 Dynamic Diagram",
		"c4Deployment":    "C4 Deployment Diagram",
	}

	if displayName, ok := displayNames[diagType]; ok {
		return displayName
	}
	// Fallback to the original type if not in the map
	return diagType
}

func contains(s, substr string) bool {
	return len(s) >= len(substr) && findSubstring(s, substr)
}

func findSubstring(s, substr string) bool {
	for i := 0; i <= len(s)-len(substr); i++ {
		if s[i:i+len(substr)] == substr {
			return true
		}
	}
	return false
}

func printHelp() {
	fmt.Print(`mermaid-check - Mermaid diagram validator and linter

Usage:
  mermaid-check [flags] [file...]

Flags:
  --help             Show this help message
  --version          Show version information
  --strict           Use strict validation rules (includes style checks)
  --error-on-empty   Treat files with no Mermaid diagrams as errors
  --format FORMAT    Force input format: 'mermaid' or 'markdown'

Examples:
  # Validate a Mermaid file
  mermaid-check diagram.mmd

  # Validate markdown with Mermaid blocks
  mermaid-check README.md

  # Validate from stdin
  cat diagram.mmd | mermaid-check

  # Force markdown mode for stdin
  cat content.txt | mermaid-check --format markdown

  # Use strict rules
  mermaid-check --strict diagram.mmd

  # Treat empty files as errors
  mermaid-check --error-on-empty *.md

Exit codes:
  0 - All diagrams are valid (or no diagrams found unless --error-on-empty is set)
  1 - Validation errors found or processing failed
`)
}
