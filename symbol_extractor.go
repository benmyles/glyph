package main

import (
	"context"
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
)

// SymbolExtractor handles symbol extraction using Tree-sitter queries
type SymbolExtractor struct {
	parser *sitter.Parser
}

// NewSymbolExtractor creates a new symbol extractor
func NewSymbolExtractor() *SymbolExtractor {
	return &SymbolExtractor{
		parser: sitter.NewParser(),
	}
}

// ExtractFromFile extracts symbols from a single file
func (e *SymbolExtractor) ExtractFromFile(filePath string, detailLevel DetailLevel) ([]Symbol, error) {
	content, err := ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	langQueries := GetLanguageQueriesForFile(filePath)
	if langQueries == nil {
		return nil, fmt.Errorf("unsupported file type: %s", filePath)
	}

	e.parser.SetLanguage(langQueries.Language)
	tree, err := e.parser.ParseCtx(context.Background(), nil, content)
	if err != nil {
		return nil, err
	}

	return e.extractSymbolsFromTree(tree, content, filePath, langQueries, detailLevel)
}

// extractSymbolsFromTree extracts symbols using Tree-sitter queries
func (e *SymbolExtractor) extractSymbolsFromTree(tree *sitter.Tree, content []byte, filePath string, langQueries *LanguageQueries, detailLevel DetailLevel) ([]Symbol, error) {
	var allSymbols []Symbol
	root := tree.RootNode()

	// Execute each query for this language
	for symbolType, queryStr := range langQueries.Queries {
		symbols, err := e.executeQuery(root, content, filePath, queryStr, symbolType, detailLevel, langQueries.Language)
		if err != nil {
			// Skip queries that fail to compile or execute
			continue
		}
		allSymbols = append(allSymbols, symbols...)
	}

	return allSymbols, nil
}

// executeQuery runs a single Tree-sitter query and extracts symbols
func (e *SymbolExtractor) executeQuery(root *sitter.Node, content []byte, filePath, queryStr, symbolType string, detailLevel DetailLevel, lang *sitter.Language) ([]Symbol, error) {
	query, err := sitter.NewQuery([]byte(queryStr), lang)
	if err != nil {
		return nil, fmt.Errorf("failed to create query for %s: %w", symbolType, err)
	}

	cursor := sitter.NewQueryCursor()
	cursor.Exec(query, root)

	var symbols []Symbol

	for {
		match, ok := cursor.NextMatch()
		if !ok {
			break
		}

		symbol := e.extractSymbolFromMatch(match, query, content, filePath, symbolType, detailLevel)
		if symbol.Name != "" {
			symbols = append(symbols, symbol)
		}
	}

	return symbols, nil
}

// extractSymbolFromMatch creates a Symbol from a query match
func (e *SymbolExtractor) extractSymbolFromMatch(match *sitter.QueryMatch, query *sitter.Query, content []byte, filePath, symbolType string, detailLevel DetailLevel) Symbol {
	symbol := Symbol{
		Kind:     mapSymbolKind(symbolType),
		FilePath: filePath,
	}

	var mainNode *sitter.Node
	var nameNode *sitter.Node

	// Extract information from captures
	for _, capture := range match.Captures {
		captureName := query.CaptureNameForId(capture.Index)
		node := capture.Node

		switch captureName {
		case "name":
			nameNode = node
			symbol.Name = string(content[node.StartByte():node.EndByte()])
		case "function", "method", "class", "interface", "type", "const", "var", "struct", "enum", "record", "annotation", "constructor", "field":
			mainNode = node
			symbol.StartLine = node.StartPoint().Row + 1
			symbol.EndLine = node.EndPoint().Row + 1
		}
	}

	// If we have a main node, extract signature based on detail level
	if mainNode != nil && detailLevel >= Standard {
		symbol.Signature = e.extractSignature(mainNode, content, detailLevel)
	}

	// If we don't have a main node but have a name node, use that for position
	if mainNode == nil && nameNode != nil {
		symbol.StartLine = nameNode.StartPoint().Row + 1
		symbol.EndLine = nameNode.EndPoint().Row + 1
	}

	return symbol
}

// extractSignature extracts the signature based on detail level
func (e *SymbolExtractor) extractSignature(node *sitter.Node, content []byte, detailLevel DetailLevel) string {
	if detailLevel == Full {
		// For full detail, include the entire node content
		return strings.TrimSpace(string(content[node.StartByte():node.EndByte()]))
	}

	// For standard detail, try to extract just the declaration part
	return e.extractDeclarationSignature(node, content)
}

// extractDeclarationSignature extracts just the declaration part (before the body)
func (e *SymbolExtractor) extractDeclarationSignature(node *sitter.Node, content []byte) string {
	startByte := node.StartByte()
	endByte := node.EndByte()

	// Look for common body indicators to stop before the implementation
	bodyIndicators := []string{"{", ":", "="}

	for i := startByte; i < endByte && i < uint32(len(content)); i++ {
		char := string(content[i])
		for _, indicator := range bodyIndicators {
			if char == indicator {
				// Found body start, return everything up to this point
				signature := strings.TrimSpace(string(content[startByte:i]))
				if signature != "" {
					return signature
				}
			}
		}
	}

	// If no body indicator found, return the whole content
	return strings.TrimSpace(string(content[startByte:endByte]))
}

// mapSymbolKind maps query symbol types to display kinds
func mapSymbolKind(symbolType string) string {
	kindMap := map[string]string{
		"functions":            "func",
		"generator_functions":  "func",
		"arrow_functions":      "func",
		"function_expressions": "func",
		"methods":              "method",
		"classes":              "class",
		"interfaces":           "interface",
		"types":                "type",
		"constants":            "const",
		"variables":            "var",
		"structs":              "struct",
		"enums":                "enum",
		"records":              "record",
		"annotations":          "annotation",
		"constructors":         "constructor",
		"fields":               "field",
		"interface_constants":  "field",
		"annotation_methods":   "method",
		"async_functions":      "func",
		"decorated_functions":  "func",
		"decorated_classes":    "class",
		"assignments":          "var",
		"type_aliases":         "type",
		"properties":           "property",
	}

	if mapped, ok := kindMap[symbolType]; ok {
		return mapped
	}
	return symbolType
}
