package main

import (
	"path/filepath"
	"strings"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/golang"
	"github.com/smacker/go-tree-sitter/java"
	"github.com/smacker/go-tree-sitter/javascript"
	"github.com/smacker/go-tree-sitter/python"
	"github.com/smacker/go-tree-sitter/typescript/typescript"
)

// LanguageQueries holds the Tree-sitter queries for a specific language
type LanguageQueries struct {
	Language *sitter.Language
	Queries  map[string]string
}

// GetLanguageQueries returns the appropriate queries for a given file path
func GetLanguageQueriesForFile(filePath string) *LanguageQueries {
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
					return &LanguageQueries{
						Language: java.GetLanguage(),
						Queries:  javaQueries,
					}
				case "go":
					return &LanguageQueries{
						Language: golang.GetLanguage(),
						Queries:  goQueries,
					}
				case "js", "javascript":
					return &LanguageQueries{
						Language: javascript.GetLanguage(),
						Queries:  javascriptQueries,
					}
				case "ts", "typescript":
					return &LanguageQueries{
						Language: typescript.GetLanguage(),
						Queries:  typescriptQueries,
					}
				case "py", "python":
					return &LanguageQueries{
						Language: python.GetLanguage(),
						Queries:  pythonQueries,
					}
				}
			}
		}
		// Also check for patterns like "something.java.txt"
		if strings.Contains(filename, ".java.txt") {
			return &LanguageQueries{
				Language: java.GetLanguage(),
				Queries:  javaQueries,
			}
		}
		if strings.Contains(filename, ".go.txt") {
			return &LanguageQueries{
				Language: golang.GetLanguage(),
				Queries:  goQueries,
			}
		}
		if strings.Contains(filename, ".js.txt") || strings.Contains(filename, ".jsx.txt") {
			return &LanguageQueries{
				Language: javascript.GetLanguage(),
				Queries:  javascriptQueries,
			}
		}
		if strings.Contains(filename, ".ts.txt") || strings.Contains(filename, ".tsx.txt") {
			return &LanguageQueries{
				Language: typescript.GetLanguage(),
				Queries:  typescriptQueries,
			}
		}
		if strings.Contains(filename, ".py.txt") {
			return &LanguageQueries{
				Language: python.GetLanguage(),
				Queries:  pythonQueries,
			}
		}
	}

	ext := strings.ToLower(filepath.Ext(filePath))

	switch ext {
	case ".go":
		return &LanguageQueries{
			Language: golang.GetLanguage(),
			Queries:  goQueries,
		}
	case ".java":
		return &LanguageQueries{
			Language: java.GetLanguage(),
			Queries:  javaQueries,
		}
	case ".js", ".jsx":
		return &LanguageQueries{
			Language: javascript.GetLanguage(),
			Queries:  javascriptQueries,
		}
	case ".py":
		return &LanguageQueries{
			Language: python.GetLanguage(),
			Queries:  pythonQueries,
		}
	case ".ts", ".tsx":
		return &LanguageQueries{
			Language: typescript.GetLanguage(),
			Queries:  typescriptQueries,
		}
	default:
		return nil
	}
}

// GetLanguageQueries returns the appropriate queries for a given language
func GetLanguageQueries(lang *sitter.Language) *LanguageQueries {
	// This is a fallback method - prefer GetLanguageQueriesForFile when possible
	switch lang {
	case golang.GetLanguage():
		return &LanguageQueries{
			Language: lang,
			Queries:  goQueries,
		}
	case java.GetLanguage():
		return &LanguageQueries{
			Language: lang,
			Queries:  javaQueries,
		}
	case javascript.GetLanguage():
		return &LanguageQueries{
			Language: lang,
			Queries:  javascriptQueries,
		}
	case python.GetLanguage():
		return &LanguageQueries{
			Language: lang,
			Queries:  pythonQueries,
		}
	case typescript.GetLanguage():
		return &LanguageQueries{
			Language: lang,
			Queries:  typescriptQueries,
		}
	default:
		return nil
	}
}

// Go language queries
var goQueries = map[string]string{
	"functions": `
		(function_declaration
			name: (identifier) @name
			parameters: (parameter_list) @params
			result: (_)? @return_type
		) @function
	`,
	"methods": `
		(method_declaration
			receiver: (parameter_list) @receiver
			name: (field_identifier) @name
			parameters: (parameter_list) @params
			result: (_)? @return_type
		) @method
	`,
	"types": `
		(type_spec
			name: (type_identifier) @name
			type: (_) @type_def
		) @type
	`,
	"constants": `
		(const_spec
			name: (identifier) @name
			type: (_)? @type
			value: (_)? @value
		) @const
	`,
	"variables": `
		(var_spec
			name: (identifier) @name
			type: (_)? @type
			value: (_)? @value
		) @var
	`,
	"interfaces": `
		(type_spec
			name: (type_identifier) @name
			type: (interface_type) @interface_body
		) @interface
	`,
	"structs": `
		(type_spec
			name: (type_identifier) @name
			type: (struct_type) @struct_body
		) @struct
	`,
}

// Java language queries
var javaQueries = map[string]string{
	"classes": `
		(class_declaration
			name: (identifier) @name
		) @class
	`,
	"interfaces": `
		(interface_declaration
			name: (identifier) @name
		) @interface
	`,
	"methods": `
		(method_declaration
			name: (identifier) @name
		) @method
	`,
	"constructors": `
		(constructor_declaration
			name: (identifier) @name
		) @constructor
	`,
	"fields": `
		(field_declaration
			declarator: (variable_declarator
				name: (identifier) @name
			)
		) @field
	`,
	"interface_constants": `
		(interface_declaration
			body: (interface_body
				(constant_declaration
					declarator: (variable_declarator
						name: (identifier) @name
					)
				) @field
			)
		)
	`,
	"annotation_methods": `
		(annotation_type_declaration
			body: (annotation_type_body
				(annotation_type_element_declaration
					name: (identifier) @name
				) @method
			)
		)
	`,
	"enums": `
		(enum_declaration
			name: (identifier) @name
		) @enum
	`,
	"records": `
		(record_declaration
			name: (identifier) @name
		) @record
	`,
	"annotations": `
		(annotation_type_declaration
			name: (identifier) @name
		) @annotation
	`,
}

// JavaScript language queries
var javascriptQueries = map[string]string{
	"functions": `
		(function_declaration
			name: (identifier) @name
		) @function
	`,
	"generator_functions": `
		(generator_function_declaration
			name: (identifier) @name
		) @function
	`,
	"arrow_functions": `
		(variable_declarator
			name: (identifier) @name
			value: (arrow_function) @arrow_func
		) @function
	`,
	"function_expressions": `
		(variable_declarator
			name: (identifier) @name
			value: (function) @func_expr
		) @function
	`,
	"classes": `
		(class_declaration
			name: (identifier) @name
		) @class
	`,
	"methods": `
		(method_definition
			name: (property_identifier) @name
		) @method
	`,
	"variables": `
		(variable_declarator
			name: (identifier) @name
		) @variable
	`,
}

// Python language queries
var pythonQueries = map[string]string{
	"functions": `
		(function_definition
			name: (identifier) @name
			parameters: (parameters) @params
			return_type: (_)? @return_type
		) @function
	`,
	"classes": `
		(class_definition
			name: (identifier) @name
			superclasses: (argument_list)? @bases
			body: (block) @body
		) @class
	`,
	"decorated_functions": `
		(decorated_definition
			(decorator)+ @decorators
			definition: (function_definition
				name: (identifier) @name
				parameters: (parameters) @params
			) @function
		) @decorated_function
	`,
	"decorated_classes": `
		(decorated_definition
			(decorator)+ @decorators
			definition: (class_definition
				name: (identifier) @name
			) @class
		) @decorated_class
	`,
	"assignments": `
		(assignment
			left: (identifier) @name
			right: (_) @value
		) @assignment
	`,
}

// TypeScript language queries (extends JavaScript)
var typescriptQueries = map[string]string{
	"functions": `
		(function_declaration
			name: (identifier) @name
		) @function
	`,
	"interfaces": `
		(interface_declaration
			name: (type_identifier) @name
		) @interface
	`,
	"type_aliases": `
		(type_alias_declaration
			name: (type_identifier) @name
		) @type
	`,
	"classes": `
		(class_declaration
			name: (type_identifier) @name
		) @class
	`,
	"methods": `
		(method_definition
			name: (property_identifier) @name
		) @method
	`,
	"properties": `
		(property_signature
			name: (property_identifier) @name
		) @property
	`,
	"variables": `
		(variable_declarator
			name: (identifier) @name
		) @variable
	`,
	"arrow_functions": `
		(variable_declarator
			name: (identifier) @name
			value: (arrow_function)
		) @function
	`,
	"namespaces": `
		(module_declaration
			name: (identifier) @name
		) @namespace
	`,
}
