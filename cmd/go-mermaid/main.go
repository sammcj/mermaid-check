// Command go-mermaid provides a CLI tool for parsing, validating, and linting Mermaid diagrams.
package main

import (
	"flag"
	"fmt"
	"io"
	"os"

	"github.com/sammcj/go-mermaid"
	"github.com/sammcj/go-mermaid/ast"
	"github.com/sammcj/go-mermaid/extractor"
	"github.com/sammcj/go-mermaid/internal/inpututil"
)

const version = "0.1.0"

func main() {
	// Define flags
	var (
		strict     = flag.Bool("strict", false, "use strict validation rules")
		formatFlag = flag.String("format", "", "force input format (mermaid or markdown)")
		showHelp   = flag.Bool("help", false, "show help message")
		showVersion = flag.Bool("version", false, "show version")
	)

	flag.Parse()

	if *showHelp {
		printHelp()
		os.Exit(0)
	}

	if *showVersion {
		fmt.Printf("go-mermaid version %s\n", version)
		os.Exit(0)
	}

	// Determine input source
	args := flag.Args()
	var exitCode int

	if len(args) == 0 {
		// Read from stdin
		exitCode = processStdin(*formatFlag, *strict)
	} else {
		// Process files
		exitCode = processFiles(args, *strict)
	}

	os.Exit(exitCode)
}

func processStdin(format string, strict bool) int {
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
			fmt.Fprintf(os.Stderr, "Error: No Mermaid diagrams found in markdown\n")
			fmt.Fprintf(os.Stderr, "Hint: Ensure code blocks use proper markdown fences: ```mermaid\n")
			return 1
		}

		fmt.Printf("Found %d Mermaid diagram(s)\n", len(blocks))

		// Collect statistics
		stats := make(map[string]int)
		for i, block := range blocks {
			fmt.Printf("\n--- Diagram %d (%s, line %d) ---\n", i+1, block.DiagramType, block.LineOffset)
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
		fmt.Printf("Diagram type: %s\n", diagramType)
		if validateDiagram(diagram, strict, "") {
			hasErrors = true
		}
	}

	if hasErrors {
		return 1
	}
	return 0
}

func processFiles(paths []string, strict bool) int {
	var hasErrors bool

	for _, path := range paths {
		fmt.Printf("\nValidating: %s\n", path)

		diagrams, err := mermaid.ParseFile(path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "  Error: %v\n", err)
			hasErrors = true
			continue
		}

		// Check if no diagrams were found
		if len(diagrams) == 0 {
			fileType := inpututil.DetectFileType(path)
			if fileType == inpututil.FileTypeMarkdown {
				fmt.Fprintf(os.Stderr, "  Error: No Mermaid diagrams found in markdown file\n")
				fmt.Fprintf(os.Stderr, "  Hint: Ensure code blocks use proper markdown fences: ```mermaid\n")
			} else {
				fmt.Fprintf(os.Stderr, "  Error: No valid Mermaid diagram found in file\n")
			}
			hasErrors = true
			continue
		}

		fileType := inpututil.DetectFileType(path)
		if fileType == inpututil.FileTypeMarkdown && len(diagrams) > 1 {
			fmt.Printf("  Found %d diagrams\n", len(diagrams))
		}

		// Collect statistics
		stats := make(map[string]int)
		for i, diagram := range diagrams {
			prefix := ""
			diagramType := diagram.GetType()
			stats[diagramType]++

			if len(diagrams) > 1 {
				prefix = fmt.Sprintf("  Diagram %d (%s): ", i+1, diagramType)
			} else {
				prefix = fmt.Sprintf("  Type: %s - ", diagramType)
			}

			if validateDiagram(diagram, strict, prefix) {
				hasErrors = true
			}
		}

		// Print summary statistics for files with multiple diagrams
		if len(diagrams) > 1 {
			fmt.Printf("\n  Diagram type distribution:\n")
			for diagramType, count := range stats {
				fmt.Printf("    %s: %d\n", diagramType, count)
			}
		}
	}

	if hasErrors {
		return 1
	}
	return 0
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
		fmt.Printf("%s✓ Valid\n", prefix)
		return false
	}

	fmt.Printf("%s✗ %d validation error(s):\n", prefix, len(errors))
	for _, err := range errors {
		fmt.Printf("%s  %v\n", prefix, err)
	}

	return true
}

func containsCodeBlocks(content string) bool {
	return len(content) > 10 && (
		contains(content, "```mermaid") ||
		contains(content, "```\nmermaid") ||
		contains(content, "# ") || // Markdown heading
		contains(content, "## "))
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
	fmt.Print(`go-mermaid - Mermaid diagram validator and linter

Usage:
  go-mermaid [flags] [file...]

Flags:
  --help             Show this help message
  --version          Show version information
  --strict           Use strict validation rules (includes style checks)
  --format FORMAT    Force input format: 'mermaid' or 'markdown'

Examples:
  # Validate a Mermaid file
  go-mermaid diagram.mmd

  # Validate markdown with Mermaid blocks
  go-mermaid README.md

  # Validate from stdin
  cat diagram.mmd | go-mermaid

  # Force markdown mode for stdin
  cat content.txt | go-mermaid --format markdown

  # Use strict rules
  go-mermaid --strict diagram.mmd

Exit codes:
  0 - All diagrams are valid
  1 - Validation errors found or processing failed
`)
}
