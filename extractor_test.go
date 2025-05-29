package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseDetailLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected DetailLevel
	}{
		{"minimal", Minimal},
		{"standard", Standard},
		{"full", Full},
		{"", Standard},
		{"invalid", Standard},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseDetailLevel(tt.input)
			if result != tt.expected {
				t.Errorf("parseDetailLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestFindFiles(t *testing.T) {
	// Create test directory structure
	testDir := t.TempDir()
	
	// Create test files
	testFiles := []string{
		"main.go",
		"utils.go",
		"src/server.go",
		"src/client.go",
		"test/main_test.go",
		"docs/readme.md",
	}
	
	for _, file := range testFiles {
		path := filepath.Join(testDir, file)
		dir := filepath.Dir(path)
		if err := os.MkdirAll(dir, 0755); err != nil {
			t.Fatal(err)
		}
		if err := os.WriteFile(path, []byte("test content"), 0644); err != nil {
			t.Fatal(err)
		}
	}
	
	tests := []struct {
		pattern  string
		expected int
	}{
		{filepath.Join(testDir, "*.go"), 2},
		{filepath.Join(testDir, "**/*.go"), 5},
		{filepath.Join(testDir, "src/*.go"), 2},
		{filepath.Join(testDir, "**/*_test.go"), 1},
	}
	
	for _, tt := range tests {
		t.Run(tt.pattern, func(t *testing.T) {
			files, err := findFiles(tt.pattern)
			if err != nil {
				t.Fatalf("findFiles(%q) error = %v", tt.pattern, err)
			}
			if len(files) != tt.expected {
				t.Errorf("findFiles(%q) returned %d files, want %d", tt.pattern, len(files), tt.expected)
			}
		})
	}
}

func TestExtractSymbols(t *testing.T) {
	// Create a test Go file
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "test.go")
	
	testCode := `package main

import "fmt"

type Server struct {
	host string
	port int
}

func (s *Server) Start() error {
	fmt.Printf("Starting server on %s:%d\n", s.host, s.port)
	return nil
}

func main() {
	server := &Server{
		host: "localhost",
		port: 8080,
	}
	server.Start()
}

const Version = "1.0.0"

var Config = map[string]string{
	"env": "production",
}
`

	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		t.Fatal(err)
	}
	
	// Test different detail levels
	tests := []struct {
		detail   string
		contains []string
	}{
		{
			detail:   "minimal",
			contains: []string{"type_spec: Server", "method_declaration: Start", "function_declaration: main"},
		},
		{
			detail:   "standard",
			contains: []string{"type_spec: Server struct", "func (s *Server) Start() error", "func main()"},
		},
		{
			detail:   "full",
			contains: []string{"Server struct", "func (s *Server) Start() error", "```"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.detail, func(t *testing.T) {
			result, err := extractSymbols(testFile, tt.detail)
			if err != nil {
				t.Fatalf("extractSymbols error = %v", err)
			}
			
			// Check if we have expected content
			if result == "No symbols found" {
				t.Errorf("No symbols were extracted")
			}
			
			// For debugging
			t.Logf("Result for %s:\n%s", tt.detail, result)
			
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Result does not contain expected string %q", expected)
				}
			}
		})
	}
}

func TestGetParserForFile(t *testing.T) {
	tests := []struct {
		filePath string
		wantNil  bool
	}{
		{"main.go", false},
		{"app.js", false},
		{"index.ts", false},
		{"script.py", false},
		{"style.css", true},
		{"readme.md", true},
	}
	
	for _, tt := range tests {
		t.Run(tt.filePath, func(t *testing.T) {
			lang, parser := getParserForFile(tt.filePath)
			if tt.wantNil {
				if lang != nil || parser != nil {
					t.Errorf("Expected nil parser for %s", tt.filePath)
				}
			} else {
				if lang == nil || parser == nil {
					t.Errorf("Expected non-nil parser for %s", tt.filePath)
				}
			}
		})
	}
}