package main

import (
	"os"
	"path/filepath"
	"testing"
)

func TestGetLanguageForFile(t *testing.T) {
	tests := []struct {
		filePath string
		wantErr  bool
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
			lang, err := GetLanguageForFile(tt.filePath)
			if tt.wantErr {
				if err == nil {
					t.Errorf("Expected error for %s", tt.filePath)
				}
			} else {
				if err != nil {
					t.Errorf("Unexpected error for %s: %v", tt.filePath, err)
				}
				if lang == nil {
					t.Errorf("Expected non-nil language for %s", tt.filePath)
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
			files, err := FindFiles(tt.pattern)
			if err != nil {
				t.Fatalf("FindFiles(%q) error = %v", tt.pattern, err)
			}
			if len(files) != tt.expected {
				t.Errorf("FindFiles(%q) returned %d files, want %d", tt.pattern, len(files), tt.expected)
			}
		})
	}
}
