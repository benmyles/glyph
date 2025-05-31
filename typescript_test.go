package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestTypeScriptSymbolExtraction(t *testing.T) {
	extractor := NewSymbolExtractor()

	tests := []struct {
		name     string
		file     string
		expected map[string][]string // symbol type -> list of expected names
	}{
		{
			name: "BasicTS",
			file: "testdata/ts_basic.ts.txt",
			expected: map[string][]string{
				"type":      {"UserID", "Status"},
				"interface": {"User", "UserRepository", "Repository", "Comparable", "Logger"},
				"class":     {"UserService", "GenericRepository", "ConsoleLogger"},
				"func":      {"createUser", "processUsers", "identity", "mapArray", "generateId", "validateEmail", "sortUsers"},
				"method":    {"constructor", "findById", "save", "delete", "getUserCount", "getId", "getCreatedAt", "findAll", "count", "log", "error"},
				"property":  {"id", "name", "email", "status", "createdAt"},
				"var":       {"deleted", "validateEmail", "emailRegex", "sortUsers", "aVal", "bVal"},
			},
		},
		{
			name: "AdvancedTS",
			file: "testdata/ts_advanced.ts.txt",
			expected: map[string][]string{
				"interface": {"Config", "EventMap", "Lengthwise", "Window"},
				"class":     {"HttpClient", "Calculator", "InMemoryRepository"},
				"func":      {"createClient", "logged", "validate", "logLength", "processValue"},
				"method":    {"constructor", "get", "post", "add", "divide", "multiply", "save", "findById", "findAll", "count", "deleteById"},
				"type":      {"EventType", "EventHandler", "ApiResponse", "Partial", "Required", "HttpMethod", "ApiEndpoint", "HttpUrl", "Pick", "Omit"},
				"property":  {"apiUrl", "timeout", "click", "hover", "focus", "blur", "message", "count", "data", "length", "id", "myApp", "version", "config"},
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

func TestTypeScriptDetailLevels(t *testing.T) {
	extractor := NewSymbolExtractor()
	testFile := "testdata/ts_basic.ts.txt"

	// Test different detail levels
	detailTests := []struct {
		level    DetailLevel
		contains []string
	}{
		{
			level:    Minimal,
			contains: []string{"interface: User", "class: UserService", "func: createUser"},
		},
		{
			level:    Standard,
			contains: []string{"interface User", "class UserService", "function createUser"},
		},
		{
			level:    Full,
			contains: []string{"interface User", "```", "class UserService"},
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

func TestTypeScriptFilePatterns(t *testing.T) {
	// Test that our TypeScript files can be found with glob patterns
	pattern := filepath.Join("testdata", "ts_*.ts.txt")
	files, err := FindFiles(pattern)
	if err != nil {
		t.Fatalf("Failed to find TypeScript test files: %v", err)
	}

	expectedFiles := []string{
		"ts_basic.ts.txt",
		"ts_advanced.ts.txt",
	}

	if len(files) < len(expectedFiles) {
		t.Errorf("Expected at least %d TypeScript test files, found %d", len(expectedFiles), len(files))
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
			t.Errorf("Expected TypeScript test file %s not found in results: %v", expectedFile, files)
		}
	}
}
