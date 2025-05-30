package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetParserForFile(t *testing.T) {
	tests := []struct {
		filePath string
		wantNil  bool
	}{
		{"main.go", false},
		{"app.js", false},
		{"index.ts", false},
		{"script.py", false},
		{"Main.java", false},
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
