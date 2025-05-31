package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestPythonSymbolExtraction(t *testing.T) {
	extractor := NewSymbolExtractor()

	tests := []struct {
		name     string
		file     string
		expected map[string][]string // symbol type -> list of expected names
	}{
		{
			name: "BasicPython",
			file: "testdata/py_basic.py.txt",
			expected: map[string][]string{
				"class": {"User", "UserRepository", "BaseService", "UserService", "DatabaseConnection"},
				"func":  {"__post_init__", "get_display_name", "is_adult", "__init__", "save", "find_by_id", "find_all", "delete", "count", "create_connection", "from_config", "process", "validate", "get_name", "create_user", "_generate_id", "retry", "decorator", "wrapper", "log_calls", "fetch_user_data", "process_users", "create_default_config", "validate_email", "calculate_age", "user_generator", "__enter__", "__exit__", "main"},
				"var":   {"VERSION", "MAX_RETRIES", "DEFAULT_TIMEOUT", "UserID", "ConfigDict", "is_active", "host", "port", "connection_string", "required_fields", "user_id", "user", "result", "pattern", "current_year", "config", "repository", "service", "user1", "user2"},
			},
		},
		{
			name: "AdvancedPython",
			file: "testdata/py_advanced.py.txt",
			expected: map[string][]string{
				"class": {"Serializable", "Comparable", "Status", "Priority", "SingletonMeta", "ConfigManager", "Cache", "AsyncTaskManager", "ValidatedProperty", "Person", "ResourceManager", "Factory"},
				"func":  {"serialize", "deserialize", "__lt__", "__eq__", "is_terminal", "__call__", "__init__", "set", "get", "clear", "add_task", "wait_all", "cancel_all", "async_database_transaction", "async_retry", "decorator", "wrapper", "measure_time", "sync_wrapper", "async_wrapper", "fetch_data", "process_batch", "async_range", "stream_data", "__set_name__", "__get__", "__set__", "is_adult", "category", "__enter__", "__exit__", "register", "create", "create_person", "create_cache", "main"},
				"var":   {"T", "K", "V", "PENDING", "PROCESSING", "COMPLETED", "FAILED", "LOW", "MEDIUM", "HIGH", "CRITICAL", "_instances", "oldest_key", "task", "results", "last_exception", "start", "result", "end", "tasks", "current", "age", "name", "creator", "manager"},
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

func TestPythonDetailLevels(t *testing.T) {
	extractor := NewSymbolExtractor()
	testFile := "testdata/py_basic.py.txt"

	// Test different detail levels
	detailTests := []struct {
		level    DetailLevel
		contains []string
	}{
		{
			level:    Minimal,
			contains: []string{"class: User", "func: main", "var: VERSION"},
		},
		{
			level:    Standard,
			contains: []string{"class User", "def main()", "VERSION (lines"},
		},
		{
			level:    Full,
			contains: []string{"class User", "```", "def main()"},
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

func TestPythonFilePatterns(t *testing.T) {
	// Test that our Python files can be found with glob patterns
	pattern := filepath.Join("testdata", "py_*.py.txt")
	files, err := FindFiles(pattern)
	if err != nil {
		t.Fatalf("Failed to find Python test files: %v", err)
	}

	expectedFiles := []string{
		"py_basic.py.txt",
		"py_advanced.py.txt",
	}

	if len(files) < len(expectedFiles) {
		t.Errorf("Expected at least %d Python test files, found %d", len(expectedFiles), len(files))
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
			t.Errorf("Expected Python test file %s not found in results: %v", expectedFile, files)
		}
	}
}
