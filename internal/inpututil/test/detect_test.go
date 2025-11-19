package inpututil_test

import (
	"testing"

	"github.com/sammcj/go-mermaid/internal/inpututil"
)

func TestDetectFileType(t *testing.T) {
	tests := []struct {
		filename string
		want     inpututil.FileType
	}{
		{"diagram.mmd", inpututil.FileTypeMermaid},
		{"diagram.MMD", inpututil.FileTypeMermaid},
		{"README.md", inpututil.FileTypeMarkdown},
		{"README.MD", inpututil.FileTypeMarkdown},
		{"doc.markdown", inpututil.FileTypeMarkdown},
		{"doc.MARKDOWN", inpututil.FileTypeMarkdown},
		{"component.mdx", inpututil.FileTypeMarkdown},
		{"component.MDX", inpututil.FileTypeMarkdown},
		{"test.txt", inpututil.FileTypeUnknown},
		{"test.go", inpututil.FileTypeUnknown},
		{"noextension", inpututil.FileTypeUnknown},
		{"/path/to/diagram.mmd", inpututil.FileTypeMermaid},
		{"/path/to/README.md", inpututil.FileTypeMarkdown},
	}

	for _, tt := range tests {
		t.Run(tt.filename, func(t *testing.T) {
			got := inpututil.DetectFileType(tt.filename)
			if got != tt.want {
				t.Errorf("DetectFileType(%q) = %v, want %v", tt.filename, got, tt.want)
			}
		})
	}
}

func TestIsMarkdownExtension(t *testing.T) {
	tests := []struct {
		ext  string
		want bool
	}{
		{".md", true},
		{".MD", true},
		{".markdown", true},
		{".MARKDOWN", true},
		{".mdx", true},
		{".MDX", true},
		{".mmd", false},
		{".txt", false},
		{".go", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			got := inpututil.IsMarkdownExtension(tt.ext)
			if got != tt.want {
				t.Errorf("IsMarkdownExtension(%q) = %v, want %v", tt.ext, got, tt.want)
			}
		})
	}
}

func TestIsMermaidExtension(t *testing.T) {
	tests := []struct {
		ext  string
		want bool
	}{
		{".mmd", true},
		{".MMD", true},
		{".md", false},
		{".markdown", false},
		{".txt", false},
		{"", false},
	}

	for _, tt := range tests {
		t.Run(tt.ext, func(t *testing.T) {
			got := inpututil.IsMermaidExtension(tt.ext)
			if got != tt.want {
				t.Errorf("IsMermaidExtension(%q) = %v, want %v", tt.ext, got, tt.want)
			}
		})
	}
}
