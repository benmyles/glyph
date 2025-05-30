package main

import (
	"context"
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


func findFiles(pattern string) ([]string, error) {
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

func extractFileSymbols(filePath string, detailLevel DetailLevel) ([]Symbol, error) {
	content, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	// Determine language and parser
	_, parser := getParserForFile(filePath)
	if parser == nil {
		return nil, fmt.Errorf("unsupported file type: %s", filePath)
	}

	// Parse the file
	tree, err := parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, err
	}

	// Extract symbols from AST
	var symbols []Symbol
	extractSymbolsFromNode(tree.RootNode(), filePath, content, &symbols, detailLevel)
	
	return symbols, nil
}

func getParserForFile(filePath string) (*sitter.Language, *sitter.Parser) {
	ext := strings.ToLower(filepath.Ext(filePath))
	
	var lang *sitter.Language
	switch ext {
	case ".go":
		lang = golang.GetLanguage()
	case ".js", ".jsx":
		lang = javascript.GetLanguage()
	case ".ts", ".tsx":
		lang = typescript.GetLanguage()
	case ".py":
		lang = python.GetLanguage()
	case ".java":
		lang = java.GetLanguage()
	default:
		return nil, nil
	}
	
	parser := sitter.NewParser()
	parser.SetLanguage(lang)
	return lang, parser
}
