package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestJavaScriptSymbolExtraction(t *testing.T) {
	extractor := NewSymbolExtractor()

	tests := []struct {
		name     string
		file     string
		expected map[string][]string // symbol type -> list of expected names
	}{
		{
			name: "BasicJS",
			file: "testdata/js_basic.js.txt",
			expected: map[string][]string{
				"func":   {"greet", "calculateTotal", "fetchUser", "saveUser", "numberGenerator", "add", "multiply"}, // regular functions, async functions, generator functions, arrow functions
				"class":  {"User", "AdminUser"},
				"method": {"constructor", "getName", "setName", "getEmail", "toString", "fromJSON", "hasPermission", "addPermission", "createSuperAdmin", "addUser", "findUser", "removeUser"},
				"var":    {"API_URL", "currentUser", "isDebug", "userService", "processData"}, // variables and function expressions (for now)
			},
		},
		{
			name: "ModernJS",
			file: "testdata/js_modern.js.txt",
			expected: map[string][]string{
				"func":   {"highlight", "createReactiveObject"}, // regular functions that are detected correctly
				"class":  {"DataStore", "ApiClient", "SecureUser", "StateMachine"},
				"method": {"constructor", "get", "set", "subscribe", "create", "request", "authenticate", "changePassword", "addTransition", "transition"},
				"var":    {"config", "message", "privateData", "INTERNAL_STATE", "createCounter"}, // variables and complex expressions
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

func TestJavaScriptDetailLevels(t *testing.T) {
	extractor := NewSymbolExtractor()
	testFile := "testdata/js_basic.js.txt"

	// Test different detail levels
	detailTests := []struct {
		level    DetailLevel
		contains []string
	}{
		{
			level:    Minimal,
			contains: []string{"func: greet", "class: User", "method: getName"},
		},
		{
			level:    Standard,
			contains: []string{"function greet(name)", "class User", "getName()"},
		},
		{
			level:    Full,
			contains: []string{"function greet(name)", "```", "class User"},
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

func TestJavaScriptFilePatterns(t *testing.T) {
	// Test that our JavaScript files can be found with glob patterns
	pattern := filepath.Join("testdata", "js_*.js.txt")
	files, err := FindFiles(pattern)
	if err != nil {
		t.Fatalf("Failed to find JavaScript test files: %v", err)
	}

	expectedFiles := []string{
		"js_basic.js.txt",
		"js_modern.js.txt",
	}

	if len(files) < len(expectedFiles) {
		t.Errorf("Expected at least %d JavaScript test files, found %d", len(expectedFiles), len(files))
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
			t.Errorf("Expected JavaScript test file %s not found in results: %v", expectedFile, files)
		}
	}
}
