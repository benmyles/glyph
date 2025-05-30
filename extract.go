package main

import (
	"fmt"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
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
		// General
		"function_declaration":    true,
		"method_declaration":      true,
		"class_declaration":       true,
		"interface_declaration":   true,
		"type_alias_declaration":  true,
		"variable_declaration":    true,
		"function_definition":     true,
		"class_definition":        true,
		"field_declaration":       true,
		
		// Go specific
		"function_item":           true,
		"struct_item":             true,
		"type_spec":               true,
		"method_spec":             true,
		"const_spec":              true,
		"var_spec":                true,
		
		// JavaScript/TypeScript specific
		"lexical_declaration":     true,
		
		// Python specific
		"decorated_definition":    true,
		
		// Rust specific
		"enum_item":               true,
		"trait_item":              true,
		"impl_item":               true,
		
		// Java specific
		"constructor_declaration": true,
		"record_declaration":      true,
		"enum_declaration":        true,
		"annotation_type_declaration": true,
		"constant_declaration":    true,
		"static_initializer":      true,
		"instance_initializer":    true,
		"compact_constructor_declaration": true,
	}
	
	return symbolTypes[nodeType]
}

func extractSymbolDetails(node *sitter.Node, content []byte, symbol *Symbol, detailLevel DetailLevel) {
	nodeType := node.Type()
	
	switch nodeType {
	case "function_declaration", "method_declaration", "function_definition", "function_item", "method_spec":
		extractFunctionDetails(node, content, symbol, detailLevel)
	case "constructor_declaration", "compact_constructor_declaration":
		extractConstructorDetails(node, content, symbol, detailLevel)
	case "class_declaration", "class_definition", "struct_item", "type_spec", "record_declaration", "enum_declaration":
		extractClassDetails(node, content, symbol, detailLevel)
	case "interface_declaration", "trait_item", "annotation_type_declaration":
		extractInterfaceDetails(node, content, symbol, detailLevel)
	case "variable_declaration", "const_declaration", "field_declaration", "const_spec", "var_spec", "constant_declaration":
		extractVariableDetails(node, content, symbol, detailLevel)
	case "static_initializer", "instance_initializer":
		extractInitializerDetails(node, content, symbol, detailLevel)
	case "lexical_declaration":
		extractVariableDetails(node, content, symbol, detailLevel)
	case "decorated_definition":
		extractDecoratedDefinition(node, content, symbol, detailLevel)
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
		// For Java classes, we need to find where the actual class declaration starts
		// (after annotations) by looking for the class/interface/enum/record keyword
		if node.Type() == "class_declaration" || node.Type() == "interface_declaration" || 
		   node.Type() == "enum_declaration" || node.Type() == "record_declaration" {
			// Find where to start the signature (skip annotations)
			signatureStart := node.StartByte()
			foundNonAnnotationModifier := false
			
			// Look for modifiers node
			for i := 0; i < int(node.ChildCount()); i++ {
				child := node.Child(i)
				if child.Type() == "modifiers" {
					// Find first non-annotation modifier
					for j := 0; j < int(child.ChildCount()); j++ {
						modChild := child.Child(j)
						if modChild.Type() != "annotation" && modChild.Type() != "marker_annotation" {
							signatureStart = modChild.StartByte()
							foundNonAnnotationModifier = true
							break
						}
					}
				}
				// If we found the class/interface/enum/record keyword and haven't found non-annotation modifiers
				if (child.Type() == "class" || child.Type() == "interface" || 
				    child.Type() == "enum" || child.Type() == "record") && !foundNonAnnotationModifier {
					signatureStart = child.StartByte()
					break
				}
			}
			
			// Find the opening brace
			bracePos := node.EndByte()
			for i := 0; i < int(node.ChildCount()); i++ {
				child := node.Child(i)
				if child.Type() == "class_body" || child.Type() == "interface_body" || 
				   child.Type() == "enum_body" || child.Type() == "record_body" {
					bracePos = child.StartByte()
					break
				}
			}
			
			// Extract the signature
			if bracePos > signatureStart {
				symbol.Signature = strings.TrimSpace(string(content[signatureStart:bracePos]))
			}
		} else if node.Type() == "type_spec" {
			// For Go type_spec, include up to and including the opening brace
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
			// Default behavior for other types
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

func extractConstructorDetails(node *sitter.Node, content []byte, symbol *Symbol, detailLevel DetailLevel) {
	// For constructors, look for the identifier (constructor name)
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "identifier" {
			symbol.Name = string(content[child.StartByte():child.EndByte()])
			break
		}
	}
	
	// Extract signature for standard and full detail levels
	if detailLevel >= Standard {
		startByte := node.StartByte()
		endByte := node.EndByte()
		
		if detailLevel == Full {
			symbol.Signature = strings.TrimSpace(string(content[startByte:endByte]))
		} else {
			// For standard detail, extract just the constructor signature
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

func extractInitializerDetails(node *sitter.Node, content []byte, symbol *Symbol, detailLevel DetailLevel) {
	// For static/instance initializers, use the node type as the name
	symbol.Name = strings.Replace(node.Type(), "_", " ", -1)
	
	if detailLevel >= Standard {
		symbol.Signature = strings.TrimSpace(string(content[node.StartByte():node.EndByte()]))
	}
}

func extractDecoratedDefinition(node *sitter.Node, content []byte, symbol *Symbol, detailLevel DetailLevel) {
	// For Python decorated definitions, extract the actual definition
	for i := 0; i < int(node.ChildCount()); i++ {
		child := node.Child(i)
		if child.Type() == "function_definition" || child.Type() == "class_definition" {
			extractSymbolDetails(child, content, symbol, detailLevel)
			return
		}
	}
}
