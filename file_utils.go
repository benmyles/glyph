package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

// ReadFile reads the content of a file
func ReadFile(filePath string) ([]byte, error) {
	return os.ReadFile(filePath)
}

// GetLanguageForFile determines the Tree-sitter language for a file
func GetLanguageForFile(filePath string) (*sitter.Language, error) {
	// For test files with .txt extension, check the filename pattern
	if strings.HasSuffix(filePath, ".txt") {
		filename := filepath.Base(filePath)
		if strings.Contains(filename, "_") {
			// Extract language from filename pattern like "java_basic_class.java.txt"
			parts := strings.Split(filename, "_")
			if len(parts) > 0 {
				lang := parts[0]
				switch lang {
				case "java":
					return java.GetLanguage(), nil
				case "go":
					return golang.GetLanguage(), nil
				case "js", "javascript":
					return javascript.GetLanguage(), nil
				case "ts", "typescript":
					return typescript.GetLanguage(), nil
				case "py", "python":
					return python.GetLanguage(), nil
				}
			}
		}
		// Also check for patterns like "something.java.txt"
		if strings.Contains(filename, ".java.txt") {
			return java.GetLanguage(), nil
		}
		if strings.Contains(filename, ".go.txt") {
			return golang.GetLanguage(), nil
		}
		if strings.Contains(filename, ".js.txt") || strings.Contains(filename, ".jsx.txt") {
			return javascript.GetLanguage(), nil
		}
		if strings.Contains(filename, ".ts.txt") || strings.Contains(filename, ".tsx.txt") {
			return typescript.GetLanguage(), nil
		}
		if strings.Contains(filename, ".py.txt") {
			return python.GetLanguage(), nil
		}
	}

	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".go":
		return golang.GetLanguage(), nil
	case ".js", ".jsx":
		return javascript.GetLanguage(), nil
	case ".ts", ".tsx":
		return typescript.GetLanguage(), nil
	case ".py":
		return python.GetLanguage(), nil
	case ".java":
		return java.GetLanguage(), nil
	default:
		return nil, fmt.Errorf("unsupported file type: %s", filePath)
	}
}

// FindFiles finds files matching a glob pattern
func FindFiles(pattern string) ([]string, error) {
	// If pattern contains **, use filepath.Walk for recursive matching
	if strings.Contains(pattern, "**") {
		var files []string

		// Split pattern at **
		parts := strings.Split(pattern, "**")
		if len(parts) != 2 {
			return nil, fmt.Errorf("invalid pattern with **: %s", pattern)
		}

		baseDir := parts[0]
		if baseDir == "" {
			baseDir = "."
		} else {
			// Remove trailing slash
			baseDir = strings.TrimSuffix(baseDir, "/")
		}

		// Get the file pattern after **
		filePattern := parts[1]
		filePattern = strings.TrimPrefix(filePattern, "/")

		err := filepath.Walk(baseDir, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return nil // Skip errors
			}

			if !info.IsDir() {
				// Check if the filename matches the pattern
				matched, _ := filepath.Match(filePattern, filepath.Base(path))
				if matched {
					files = append(files, path)
				}
			}
			return nil
		})

		if err != nil {
			return nil, err
		}

		return files, nil
	}

	// For patterns without **, use standard glob
	matches, err := filepath.Glob(pattern)
	if err != nil {
		return nil, err
	}

	return matches, nil
}
