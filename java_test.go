package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestJavaSymbolExtraction(t *testing.T) {
	extractor := NewSymbolExtractor()

	tests := []struct {
		name     string
		file     string
		expected map[string][]string // symbol type -> list of expected names
	}{
		{
			name: "BasicClass",
			file: "testdata/java_basic_class.java.txt",
			expected: map[string][]string{
				"class":       {"BasicExample"},
				"field":       {"VERSION", "MAX_SIZE", "name", "items", "count"},
				"constructor": {"BasicExample", "BasicExample"}, // two constructors
				"method":      {"getName", "setName", "addItem", "printVersion", "processFile", "main"},
			},
		},
		{
			name: "Interface",
			file: "testdata/java_interface.java.txt",
			expected: map[string][]string{
				"interface": {"DataProcessor"},
				"field":     {"DEFAULT_NAME", "MAX_ITEMS"},
				"method":    {"process", "processAll", "processWithValidation", "initialize", "isValid", "printInfo", "createDefault"},
			},
		},
		{
			name: "Enum",
			file: "testdata/java_enum.java.txt",
			expected: map[string][]string{
				"enum":        {"Status"},
				"field":       {"displayName", "code", "active"},
				"constructor": {"Status"},
				"method":      {"getDisplayName", "getCode", "isActive", "getDescription", "fromCode", "toString"},
			},
		},
		{
			name: "Record",
			file: "testdata/java_record.java.txt",
			expected: map[string][]string{
				"record":      {"Person"},
				"field":       {"MIN_AGE", "MAX_AGE"},
				"constructor": {"Person"},
				"method":      {"isAdult", "isMinor", "getDisplayName", "createChild", "isValidAge", "toString"},
			},
		},
		{
			name: "Annotation",
			file: "testdata/java_annotation.java.txt",
			expected: map[string][]string{
				"annotation": {"Benchmark"},
				"enum":       {"TimeUnit"},
				"method":     {"value", "description", "iterations", "enabled", "tags", "unit"},
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

func TestJavaDetailLevels(t *testing.T) {
	extractor := NewSymbolExtractor()
	testFile := "testdata/java_basic_class.java.txt"

	// Test different detail levels
	detailTests := []struct {
		level    DetailLevel
		contains []string
	}{
		{
			level:    Minimal,
			contains: []string{"class: BasicExample", "method: getName", "field: name"},
		},
		{
			level:    Standard,
			contains: []string{"public class BasicExample", "public String getName()", "private String name"},
		},
		{
			level:    Full,
			contains: []string{"public class BasicExample", "```", "public String getName()"},
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

func TestJavaFilePatterns(t *testing.T) {
	// Test that our Java files can be found with glob patterns
	pattern := filepath.Join("testdata", "java_*.java.txt")
	files, err := FindFiles(pattern)
	if err != nil {
		t.Fatalf("Failed to find Java test files: %v", err)
	}

	expectedFiles := []string{
		"java_basic_class.java.txt",
		"java_interface.java.txt",
		"java_enum.java.txt",
		"java_record.java.txt",
		"java_annotation.java.txt",
	}

	if len(files) < len(expectedFiles) {
		t.Errorf("Expected at least %d Java test files, found %d", len(expectedFiles), len(files))
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
			t.Errorf("Expected Java test file %s not found in results: %v", expectedFile, files)
		}
	}
}

// Helper function to check if a slice contains a string
func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}
