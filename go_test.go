package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestGoSymbolExtraction(t *testing.T) {
	extractor := NewSymbolExtractor()

	tests := []struct {
		name     string
		file     string
		expected map[string][]string // symbol type -> list of expected names
	}{
		{
			name: "BasicGo",
			file: "testdata/go_basic.go.txt",
			expected: map[string][]string{
				"const":     {"Version", "MaxSize", "DefaultPort", "StatusPending", "StatusRunning", "StatusComplete"},
				"var":       {"GlobalCounter", "ServerName", "isDebug"},
				"type":      {"UserID", "Config", "Status", "Handler", "Logger", "Server", "Response"},
				"struct":    {"Config", "Server", "Response"},
				"interface": {"Handler", "Logger"},
				"func":      {"main", "NewServer", "processRequest"},
				"method":    {"Start", "Stop", "GetConfig", "SetLogger"},
			},
		},
		{
			name: "Generics",
			file: "testdata/go_generics.go.txt",
			expected: map[string][]string{
				"type":      {"Stack", "Pair", "Result", "Comparable", "Container", "Ordered", "Numeric", "Cache"},
				"struct":    {"Stack", "Pair", "Result", "Cache"},
				"interface": {"Comparable", "Container", "Ordered", "Numeric"},
				"func":      {"Map", "Filter", "Reduce", "Max", "Sum", "NewCache", "ProcessWithContext"},
				"method":    {"Push", "Pop", "Peek", "Size", "String", "Set", "Get", "Delete", "Keys"},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Test that the file exists
			if _, err := os.Stat(tt.file); os.IsNotExist(err) {
				t.Fatalf("Test file does not exist: %s", tt.file)
			}

			// Extract symbols
			symbols, err := extractor.ExtractFromFile(tt.file, Standard)
			if err != nil {
				t.Fatalf("Failed to extract symbols from %s: %v", tt.file, err)
			}

			if len(symbols) == 0 {
				t.Fatalf("No symbols extracted from %s", tt.file)
			}

			// Group symbols by kind
			symbolsByKind := make(map[string][]string)
			for _, symbol := range symbols {
				symbolsByKind[symbol.Kind] = append(symbolsByKind[symbol.Kind], symbol.Name)
			}

			// Check expected symbols
			for expectedKind, expectedNames := range tt.expected {
				actualNames, found := symbolsByKind[expectedKind]
				if !found {
					t.Errorf("Expected symbol kind %s not found in %s", expectedKind, tt.file)
					continue
				}

				for _, expectedName := range expectedNames {
					if !contains(actualNames, expectedName) {
						t.Errorf("Expected %s symbol '%s' not found in %s. Found: %v",
							expectedKind, expectedName, tt.file, actualNames)
					}
				}
			}

			// Log the results for debugging
			result := FormatSymbols(symbols, Standard)
			t.Logf("Symbols extracted from %s:\n%s", tt.file, result)
		})
	}
}

func TestGoDetailLevels(t *testing.T) {
	extractor := NewSymbolExtractor()
	testFile := "testdata/go_basic.go.txt"

	// Test different detail levels
	detailTests := []struct {
		level    DetailLevel
		contains []string
	}{
		{
			level:    Minimal,
			contains: []string{"func: main", "type: Config", "method: Start"},
		},
		{
			level:    Standard,
			contains: []string{"func: func main()", "type: Config struct", "method: func (s *Server) Start() error"},
		},
		{
			level:    Full,
			contains: []string{"func (lines", "```", "type (lines"},
		},
	}

	for _, tt := range detailTests {
		t.Run(tt.level.String(), func(t *testing.T) {
			symbols, err := extractor.ExtractFromFile(testFile, tt.level)
			if err != nil {
				t.Fatalf("Failed to extract symbols: %v", err)
			}

			result := FormatSymbols(symbols, tt.level)

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Expected %s detail level to contain %q, but it didn't.\nResult:\n%s",
						tt.level.String(), expected, result)
				}
			}
		})
	}
}

func TestGoFilePatterns(t *testing.T) {
	// Test that our Go files can be found with glob patterns
	pattern := filepath.Join("testdata", "go_*.go.txt")
	files, err := FindFiles(pattern)
	if err != nil {
		t.Fatalf("Failed to find Go test files: %v", err)
	}

	expectedFiles := []string{
		"go_basic.go.txt",
		"go_generics.go.txt",
	}

	if len(files) < len(expectedFiles) {
		t.Errorf("Expected at least %d Go test files, found %d", len(expectedFiles), len(files))
	}

	for _, expectedFile := range expectedFiles {
		found := false
		for _, file := range files {
			if strings.HasSuffix(file, expectedFile) {
				found = true
				break
			}
		}
		if !found {
			t.Errorf("Expected Go test file %s not found in results: %v", expectedFile, files)
		}
	}
}
