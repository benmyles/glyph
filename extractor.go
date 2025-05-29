package main

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

type Symbol struct {
	Name       string
	Kind       string
	StartLine  uint32
	EndLine    uint32
	Signature  string
	FilePath   string
	Children   []Symbol
}

type DetailLevel int

const (
	Minimal DetailLevel = iota
	Standard
	Full
)

func parseDetailLevel(detail string) DetailLevel {
	switch strings.ToLower(detail) {
	case "minimal":
		return Minimal
	case "full":
		return Full
	default:
		return Standard
	}
}

func extractSymbols(pattern string, detail string) (string, error) {
	detailLevel := parseDetailLevel(detail)
	
	// Find files matching the pattern
	files, err := findFiles(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to find files: %w", err)
	}

	if len(files) == 0 {
		return "No files found matching pattern: " + pattern, nil
	}

	var allSymbols []Symbol
	
	for _, file := range files {
		symbols, err := extractFileSymbols(file, detailLevel)
		if err != nil {
			continue // Skip files that can't be parsed
		}
		allSymbols = append(allSymbols, symbols...)
	}

	return formatSymbols(allSymbols, detailLevel), nil
}

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
		if strings.HasPrefix(filePattern, "/") {
			filePattern = filePattern[1:]
		}
		
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
	default:
		return nil, nil
	}
	
	parser := sitter.NewParser()
	parser.SetLanguage(lang)
	return lang, parser
}

func extractSymbolsFromNode(node *sitter.Node, filePath string, content []byte, symbols *[]Symbol, detailLevel DetailLevel) {
	nodeType := node.Type()
	
	// Check if this node represents a symbol we want to extract
	if isSymbolNode(nodeType) {
		symbol := Symbol{
			Kind:      nodeType,
			StartLine: node.StartPoint().Row + 1,
			EndLine:   node.EndPoint().Row + 1,
			FilePath:  filePath,
		}
		
		// Extract name and signature based on node type
		extractSymbolDetails(node, content, &symbol, detailLevel)
		
		if symbol.Name != "" {
			*symbols = append(*symbols, symbol)
		}
	}
	
	// Recursively process children
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		extractSymbolsFromNode(child, filePath, content, symbols, detailLevel)
	}
}

func isSymbolNode(nodeType string) bool {
	symbolTypes := map[string]bool{
		"function_declaration":    true,
		"method_declaration":      true,
		"class_declaration":       true,
		"interface_declaration":   true,
		"type_alias_declaration":  true,
		"variable_declaration":    true,
		"function_definition":     true,
		"class_definition":        true,
		"function_item":           true,
		"struct_item":             true,
		"enum_item":               true,
		"trait_item":              true,
		"impl_item":               true,
		"type_spec":               true,
		"method_spec":             true,
		"field_declaration":       true,
		"const_spec":              true, // Go constants (inside const_declaration)
		"var_spec":                true, // Go variables (inside var_declaration)
	}
	
	return symbolTypes[nodeType]
}

func extractSymbolDetails(node *sitter.Node, content []byte, symbol *Symbol, detailLevel DetailLevel) {
	nodeType := node.Type()
	
	switch nodeType {
	case "function_declaration", "method_declaration", "function_definition", "function_item", "method_spec":
		extractFunctionDetails(node, content, symbol, detailLevel)
	case "class_declaration", "class_definition", "struct_item", "type_spec":
		extractClassDetails(node, content, symbol, detailLevel)
	case "interface_declaration", "trait_item":
		extractInterfaceDetails(node, content, symbol, detailLevel)
	case "variable_declaration", "const_declaration", "field_declaration", "const_spec", "var_spec":
		extractVariableDetails(node, content, symbol, detailLevel)
	}
}

func extractFunctionDetails(node *sitter.Node, content []byte, symbol *Symbol, detailLevel DetailLevel) {
	// Find name node
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "identifier" || child.Type() == "property_identifier" || child.Type() == "field_identifier" {
			symbol.Name = string(content[child.StartByte():child.EndByte()])
			break
		}
	}
	
	// Extract signature for standard and full detail levels
	if detailLevel >= Standard {
		startByte := node.StartByte()
		endByte := node.EndByte()
		
		// For full detail, include the entire function declaration
		if detailLevel == Full {
			symbol.Signature = strings.TrimSpace(string(content[startByte:endByte]))
		} else {
			// For standard detail, extract just the function signature (up to the opening brace)
			for i := startByte; i < endByte && i < uint32(len(content)); i++ {
				if content[i] == '{' {
					symbol.Signature = strings.TrimSpace(string(content[startByte:i]))
					break
				}
			}
			if symbol.Signature == "" {
				symbol.Signature = strings.TrimSpace(string(content[startByte:endByte]))
			}
		}
	}
}

func extractClassDetails(node *sitter.Node, content []byte, symbol *Symbol, detailLevel DetailLevel) {
	// Find name node
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "identifier" || child.Type() == "type_identifier" {
			symbol.Name = string(content[child.StartByte():child.EndByte()])
			break
		}
	}
	
	if detailLevel >= Standard {
		// For type_spec in Go, we want to include the type keyword
		if node.Type() == "type_spec" {
			// Extract everything up to the opening brace
			startByte := node.StartByte()
			for i := startByte; i < node.EndByte() && i < uint32(len(content)); i++ {
				if content[i] == '{' {
					symbol.Signature = strings.TrimSpace(string(content[startByte:i+1]))
					break
				}
			}
			// If no brace found, it might be a type alias
			if symbol.Signature == "" {
				symbol.Signature = strings.TrimSpace(string(content[node.StartByte():node.EndByte()]))
			}
		} else {
			// Extract class/struct header
			startByte := node.StartByte()
			for i := startByte; i < node.EndByte() && i < uint32(len(content)); i++ {
				if content[i] == '{' {
					symbol.Signature = strings.TrimSpace(string(content[startByte:i]))
					break
				}
			}
		}
	}
}

func extractInterfaceDetails(node *sitter.Node, content []byte, symbol *Symbol, detailLevel DetailLevel) {
	// Similar to class extraction
	extractClassDetails(node, content, symbol, detailLevel)
}

func extractVariableDetails(node *sitter.Node, content []byte, symbol *Symbol, detailLevel DetailLevel) {
	// Find variable name
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "variable_declarator" || child.Type() == "identifier" {
			if child.Type() == "variable_declarator" && child.ChildCount() > 0 {
				nameNode := child.Child(0)
				if nameNode.Type() == "identifier" {
					symbol.Name = string(content[nameNode.StartByte():nameNode.EndByte()])
				}
			} else if child.Type() == "identifier" {
				symbol.Name = string(content[child.StartByte():child.EndByte()])
			}
			break
		} else if child.Type() == "const_spec" || child.Type() == "var_spec" {
			// For Go const/var declarations, find the identifier within the spec
			if child.ChildCount() > 0 {
				firstChild := child.Child(0)
				if firstChild.Type() == "identifier" {
					symbol.Name = string(content[firstChild.StartByte():firstChild.EndByte()])
				}
			}
		}
	}
	
	if detailLevel >= Standard {
		// For const_spec and var_spec, use the spec content directly
		if node.Type() == "const_spec" || node.Type() == "var_spec" {
			symbol.Signature = strings.TrimSpace(string(content[node.StartByte():node.EndByte()]))
		} else {
			// For const_declaration and var_declaration, use the whole declaration
			symbol.Signature = strings.TrimSpace(string(content[node.StartByte():node.EndByte()]))
		}
	}
}

func formatSymbols(symbols []Symbol, detailLevel DetailLevel) string {
	if len(symbols) == 0 {
		return "No symbols found"
	}
	
	var sb strings.Builder
	sb.WriteString("# Symbol Outline\n\n")
	
	// Group symbols by file
	fileSymbols := make(map[string][]Symbol)
	for _, sym := range symbols {
		fileSymbols[sym.FilePath] = append(fileSymbols[sym.FilePath], sym)
	}
	
	// Format output
	for file, syms := range fileSymbols {
		sb.WriteString(fmt.Sprintf("## %s\n\n", file))
		
		for _, sym := range syms {
			formatSymbol(&sb, sym, detailLevel, 0)
		}
		
		sb.WriteString("\n")
	}
	
	return sb.String()
}

func formatSymbol(sb *strings.Builder, symbol Symbol, detailLevel DetailLevel, indent int) {
	indentStr := strings.Repeat("  ", indent)
	
	switch detailLevel {
	case Minimal:
		sb.WriteString(fmt.Sprintf("%s- %s: %s (line %d)\n", 
			indentStr, symbol.Kind, symbol.Name, symbol.StartLine))
	case Standard:
		if symbol.Signature != "" {
			sb.WriteString(fmt.Sprintf("%s- %s: %s\n", 
				indentStr, symbol.Kind, symbol.Signature))
		} else {
			sb.WriteString(fmt.Sprintf("%s- %s: %s (lines %d-%d)\n", 
				indentStr, symbol.Kind, symbol.Name, symbol.StartLine, symbol.EndLine))
		}
	case Full:
		sb.WriteString(fmt.Sprintf("%s- %s (lines %d-%d):\n", 
			indentStr, symbol.Kind, symbol.StartLine, symbol.EndLine))
		if symbol.Signature != "" {
			sb.WriteString(fmt.Sprintf("%s  ```\n%s  %s\n%s  ```\n", 
				indentStr, indentStr, symbol.Signature, indentStr))
		}
	}
	
	// Format children if any
	for _, child := range symbol.Children {
		formatSymbol(sb, child, detailLevel, indent+1)
	}
}