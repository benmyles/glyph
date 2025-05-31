# Comprehensive Guide to go-tree-sitter

## Table of Contents
1. [Introduction](#introduction)
2. [Installation](#installation)
3. [Basic Concepts](#basic-concepts)
4. [Language-Specific Usage](#language-specific-usage)
    - [Go](#working-with-go)
    - [Java](#working-with-java)
    - [JavaScript](#working-with-javascript)
    - [Python](#working-with-python)
    - [TypeScript](#working-with-typescript)
5. [Working with AST Nodes](#working-with-ast-nodes)
6. [Querying and Pattern Matching](#querying-and-pattern-matching)
7. [Advanced Features](#advanced-features)
8. [Practical Examples](#practical-examples)
9. [Best Practices](#best-practices)

## Introduction

go-tree-sitter is a Go binding for the [tree-sitter](https://tree-sitter.github.io/tree-sitter/) parsing library. Tree-sitter is a parser generator tool and an incremental parsing library that builds a concrete syntax tree for a source file and efficiently updates the syntax tree as the source file is edited.

### Key Features
- Fast incremental parsing
- Robust error recovery
- Language-agnostic query system
- Zero-copy parsing
- Thread-safe parsing
- Support for multiple programming languages

## Installation

To use go-tree-sitter in your Go project:

```bash
go get github.com/smacker/go-tree-sitter
```

For specific language support:
```bash
# For Go support
go get github.com/smacker/go-tree-sitter/golang

# For Java support
go get github.com/smacker/go-tree-sitter/java

# For JavaScript support
go get github.com/smacker/go-tree-sitter/javascript

# For Python support
go get github.com/smacker/go-tree-sitter/python

# For TypeScript support
go get github.com/smacker/go-tree-sitter/typescript/typescript
go get github.com/smacker/go-tree-sitter/typescript/tsx
```

## Basic Concepts

### Parser
The parser is the main object that performs the parsing operation. It needs to be configured with a specific language grammar.

### Tree
A tree represents the syntax tree of an entire source code file. It's immutable and can be reused.

### Node
Nodes represent individual elements in the syntax tree (e.g., functions, variables, expressions).

### Query
Queries allow you to search for specific patterns in the syntax tree using S-expressions.

## Language-Specific Usage

### Working with Go

```go
package main

import (
    "context"
    "fmt"
    sitter "github.com/smacker/go-tree-sitter"
    "github.com/smacker/go-tree-sitter/golang"
)

func parseGoCode() {
    // Go source code
    sourceCode := []byte(`
package main

import "fmt"

func main() {
    message := "Hello, World!"
    fmt.Println(message)
}

func add(a, b int) int {
    return a + b
}
`)

    // Create parser
    parser := sitter.NewParser()
    parser.SetLanguage(golang.GetLanguage())

    // Parse the code
    tree, err := parser.ParseCtx(context.Background(), nil, sourceCode)
    if err != nil {
        panic(err)
    }

    // Get root node
    root := tree.RootNode()
    fmt.Printf("Root type: %s\n", root.Type()) // source_file

    // Find all function declarations
    query := `(function_declaration name: (identifier) @func_name)`
    q, _ := sitter.NewQuery([]byte(query), golang.GetLanguage())
    qc := sitter.NewQueryCursor()
    qc.Exec(q, root)

    fmt.Println("\nFunctions found:")
    for {
        m, ok := qc.NextMatch()
        if !ok {
            break
        }
        for _, c := range m.Captures {
            fmt.Printf("- %s\n", c.Node.Content(sourceCode))
        }
    }
}
```

### Working with Java

```go
package main

import (
    "context"
    "fmt"
    sitter "github.com/smacker/go-tree-sitter"
    "github.com/smacker/go-tree-sitter/java"
)

func parseJavaCode() {
    // Java source code
    sourceCode := []byte(`
import java.util.*;

public class HelloWorld {
    private String message;
    
    public HelloWorld(String message) {
        this.message = message;
    }
    
    public void printMessage() {
        System.out.println(message);
    }
    
    public static void main(String[] args) {
        HelloWorld hw = new HelloWorld("Hello, World!");
        hw.printMessage();
    }
}
`)

    parser := sitter.NewParser()
    parser.SetLanguage(java.GetLanguage())

    tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)
    root := tree.RootNode()

    // Find all method declarations
    query := `(method_declaration
        name: (identifier) @method_name
        parameters: (formal_parameters) @params
    )`
    
    q, _ := sitter.NewQuery([]byte(query), java.GetLanguage())
    qc := sitter.NewQueryCursor()
    qc.Exec(q, root)

    fmt.Println("Methods found:")
    for {
        m, ok := qc.NextMatch()
        if !ok {
            break
        }
        var methodName, params string
        for _, c := range m.Captures {
            captureName := q.CaptureNameForId(c.Index)
            if captureName == "method_name" {
                methodName = c.Node.Content(sourceCode)
            } else if captureName == "params" {
                params = c.Node.Content(sourceCode)
            }
        }
        fmt.Printf("- %s%s\n", methodName, params)
    }
}
```

### Working with JavaScript

```go
package main

import (
    "context"
    "fmt"
    sitter "github.com/smacker/go-tree-sitter"
    "github.com/smacker/go-tree-sitter/javascript"
)

func parseJavaScriptCode() {
    sourceCode := []byte(`
// ES6 class example
class Calculator {
    constructor() {
        this.result = 0;
    }
    
    add(a, b) {
        return a + b;
    }
    
    multiply(a, b) {
        return a * b;
    }
}

// Arrow functions
const square = (x) => x * x;
const greet = name => \`Hello, \${name}!\`;

// Async function
async function fetchData(url) {
    const response = await fetch(url);
    return response.json();
}
`)

    parser := sitter.NewParser()
    parser.SetLanguage(javascript.GetLanguage())

    tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)
    root := tree.RootNode()

    // Find all function types (regular, arrow, async)
    query := `[
        (function_declaration name: (identifier) @func_name)
        (variable_declarator 
            name: (identifier) @func_name
            value: (arrow_function))
        (method_definition 
            name: (property_identifier) @func_name)
    ]`
    
    q, _ := sitter.NewQuery([]byte(query), javascript.GetLanguage())
    qc := sitter.NewQueryCursor()
    qc.Exec(q, root)

    fmt.Println("Functions and methods found:")
    for {
        m, ok := qc.NextMatch()
        if !ok {
            break
        }
        for _, c := range m.Captures {
            fmt.Printf("- %s (type: %s)\n", 
                c.Node.Content(sourceCode), 
                c.Node.Parent().Type())
        }
    }
}
```

### Working with Python

```go
package main

import (
    "context"
    "fmt"
    sitter "github.com/smacker/go-tree-sitter"
    "github.com/smacker/go-tree-sitter/python"
)

func parsePythonCode() {
    sourceCode := []byte(`
import math
from typing import List, Optional

class DataProcessor:
    """A class for processing data."""
    
    def __init__(self, name: str):
        self.name = name
        self.data: List[float] = []
    
    def add_data(self, value: float) -> None:
        """Add a data point."""
        self.data.append(value)
    
    @property
    def mean(self) -> Optional[float]:
        """Calculate the mean of the data."""
        if not self.data:
            return None
        return sum(self.data) / len(self.data)
    
    @staticmethod
    def validate_input(value: float) -> bool:
        """Validate input value."""
        return not math.isnan(value) and not math.isinf(value)

def main():
    processor = DataProcessor("MyProcessor")
    processor.add_data(10.5)
    processor.add_data(20.3)
    print(f"Mean: {processor.mean}")

if __name__ == "__main__":
    main()
`)

    parser := sitter.NewParser()
    parser.SetLanguage(python.GetLanguage())

    tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)
    root := tree.RootNode()

    // Find all function definitions with decorators
    query := `(function_definition
        name: (identifier) @func_name
        parameters: (parameters) @params
        (#match? @func_name "^[a-z_]")
    )`
    
    q, _ := sitter.NewQuery([]byte(query), python.GetLanguage())
    qc := sitter.NewQueryCursor()
    qc.Exec(q, root)

    fmt.Println("Functions found:")
    for {
        m, ok := qc.NextMatch()
        if !ok {
            break
        }
        m = qc.FilterPredicates(m, sourceCode)
        for _, c := range m.Captures {
            if q.CaptureNameForId(c.Index) == "func_name" {
                fmt.Printf("- %s\n", c.Node.Content(sourceCode))
                
                // Check for decorators
                parent := c.Node.Parent()
                if parent != nil && parent.Type() == "decorated_definition" {
                    fmt.Println("  Has decorators")
                }
            }
        }
    }
}
```

### Working with TypeScript

```go
package main

import (
    "context"
    "fmt"
    sitter "github.com/smacker/go-tree-sitter"
    "github.com/smacker/go-tree-sitter/typescript/typescript"
)

func parseTypeScriptCode() {
    sourceCode := []byte(`
interface User {
    id: number;
    name: string;
    email?: string;
}

type Status = 'active' | 'inactive' | 'pending';

class UserService {
    private users: Map<number, User>;
    
    constructor() {
        this.users = new Map();
    }
    
    async addUser(user: User): Promise<void> {
        this.users.set(user.id, user);
    }
    
    getUser(id: number): User | undefined {
        return this.users.get(id);
    }
    
    getAllUsers(): User[] {
        return Array.from(this.users.values());
    }
}

// Generic function
function processArray<T>(items: T[], processor: (item: T) => void): void {
    items.forEach(processor);
}

// Type guard
function isUser(obj: any): obj is User {
    return 'id' in obj && 'name' in obj;
}
`)

    parser := sitter.NewParser()
    parser.SetLanguage(typescript.GetLanguage())

    tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)
    root := tree.RootNode()

    // Find all type definitions
    query := `[
        (interface_declaration name: (type_identifier) @type_name)
        (type_alias_declaration name: (type_identifier) @type_name)
        (class_declaration name: (type_identifier) @type_name)
    ]`
    
    q, _ := sitter.NewQuery([]byte(query), typescript.GetLanguage())
    qc := sitter.NewQueryCursor()
    qc.Exec(q, root)

    fmt.Println("Types found:")
    for {
        m, ok := qc.NextMatch()
        if !ok {
            break
        }
        for _, c := range m.Captures {
            fmt.Printf("- %s (%s)\n", 
                c.Node.Content(sourceCode),
                c.Node.Parent().Type())
        }
    }
}
```

## Working with AST Nodes

### Node Properties and Methods

```go
// Basic node information
node.Type()          // Get node type (e.g., "function_declaration")
node.StartByte()     // Starting byte position
node.EndByte()       // Ending byte position
node.StartPoint()    // Starting row/column
node.EndPoint()      // Ending row/column
node.Range()         // Complete range information

// Node relationships
node.Parent()        // Get parent node
node.Child(0)        // Get child by index
node.NamedChild(0)   // Get named child by index
node.ChildCount()    // Total number of children
node.NamedChildCount() // Number of named children
node.NextSibling()   // Next sibling node
node.PrevSibling()   // Previous sibling node

// Node properties
node.IsNamed()       // Check if node is named
node.IsMissing()     // Check if node is missing (error recovery)
node.IsExtra()       // Check if node is extra (e.g., comments)
node.IsError()       // Check if node is an error
node.HasError()      // Check if subtree contains errors
node.HasChanges()    // Check if node changed after edit

// Content extraction
node.Content(sourceCode) // Get the source code for this node

// Field access
node.ChildByFieldName("name") // Get child by field name
node.FieldNameForChild(0)     // Get field name for child at index
```

### Tree Cursor for Efficient Traversal

```go
func traverseWithCursor(root *sitter.Node) {
    cursor := sitter.NewTreeCursor(root)
    defer cursor.Close()

    visitNode := func() {
        node := cursor.CurrentNode()
        fieldName := cursor.CurrentFieldName()
        
        fmt.Printf("%s%s", 
            strings.Repeat("  ", int(node.StartPoint().Row)),
            node.Type())
        
        if fieldName != "" {
            fmt.Printf(" (%s)", fieldName)
        }
        fmt.Println()
    }

    var traverse func() bool
    traverse = func() bool {
        visitNode()
        
        if cursor.GoToFirstChild() {
            for {
                traverse()
                if !cursor.GoToNextSibling() {
                    break
                }
            }
            cursor.GoToParent()
        }
        return true
    }
    
    traverse()
}
```

## Querying and Pattern Matching

### Query Syntax

Tree-sitter uses S-expression syntax for queries:

```scheme
; Basic node matching
(identifier)

; Node with specific type
(function_declaration)

; Nested matching
(function_declaration
  name: (identifier))

; Capture nodes with @name
(function_declaration
  name: (identifier) @function_name)

; Multiple patterns
[
  (function_declaration)
  (method_definition)
]

; Wildcards
(call_expression
  function: (_) @callable)

; Field names
(function_declaration
  name: (identifier) @name
  parameters: (formal_parameters) @params
  body: (statement_block) @body)
```

### Predicates

```go
func demonstratePredicates() {
    sourceCode := []byte(`
        const CONSTANT_VALUE = 42;
        const camelCaseVar = "hello";
        let snake_case_var = true;
    `)

    parser := sitter.NewParser()
    parser.SetLanguage(javascript.GetLanguage())
    tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)
    root := tree.RootNode()

    // Match only SCREAMING_SNAKE_CASE constants
    query := `(
        (variable_declarator
            name: (identifier) @const_name)
        (#match? @const_name "^[A-Z_]+$")
    )`

    q, _ := sitter.NewQuery([]byte(query), javascript.GetLanguage())
    qc := sitter.NewQueryCursor()
    qc.Exec(q, root)

    for {
        m, ok := qc.NextMatch()
        if !ok {
            break
        }
        // Apply predicate filtering
        m = qc.FilterPredicates(m, sourceCode)
        for _, c := range m.Captures {
            fmt.Printf("Constant: %s\n", c.Node.Content(sourceCode))
        }
    }
}
```

### Available Predicates

1. **`eq?`** - Test equality
   ```scheme
   (#eq? @node "specific_text")
   (#eq? @node1 @node2)
   ```

2. **`not-eq?`** - Test inequality
   ```scheme
   (#not-eq? @node "text")
   ```

3. **`match?`** - Regex matching
   ```scheme
   (#match? @identifier "^[A-Z][a-zA-Z0-9]*$")
   ```

4. **`not-match?`** - Negative regex matching
   ```scheme
   (#not-match? @string "[0-9]")
   ```

## Advanced Features

### Incremental Parsing

```go
func incrementalParsing() {
    parser := sitter.NewParser()
    parser.SetLanguage(javascript.GetLanguage())

    // Initial parse
    oldSource := []byte("let x = 1")
    tree, _ := parser.ParseCtx(context.Background(), nil, oldSource)
    
    // Make an edit: change "1" to "42"
    newSource := []byte("let x = 42")
    
    // Notify the tree about the edit
    tree.Edit(sitter.EditInput{
        StartIndex:  8,
        OldEndIndex: 9,
        NewEndIndex: 10,
        StartPoint: sitter.Point{Row: 0, Column: 8},
        OldEndPoint: sitter.Point{Row: 0, Column: 9},
        NewEndPoint: sitter.Point{Row: 0, Column: 10},
    })

    // Re-parse incrementally
    newTree, _ := parser.ParseCtx(context.Background(), tree, newSource)
    
    // Check what changed
    root := newTree.RootNode()
    checkChanges(root)
}

func checkChanges(node *sitter.Node) {
    if node.HasChanges() {
        fmt.Printf("Node changed: %s\n", node.Type())
    }
    for i := 0; i < int(node.ChildCount()); i++ {
        checkChanges(node.Child(i))
    }
}
```

### Custom Input Reader

```go
func customInputReader() {
    // Simulate reading from a large file or network stream
    chunks := []string{
        "function hello() {",
        "  console.log('Hello');",
        "}",
    }
    
    parser := sitter.NewParser()
    parser.SetLanguage(javascript.GetLanguage())

    input := sitter.Input{
        Encoding: sitter.InputEncodingUTF8,
        Read: func(offset uint32, position sitter.Point) []byte {
            // Concatenate all chunks up to the offset
            fullText := strings.Join(chunks, "\n")
            if int(offset) >= len(fullText) {
                return nil
            }
            
            // Return remaining text from offset
            return []byte(fullText[offset:])
        },
    }

    tree, _ := parser.ParseInputCtx(context.Background(), nil, input)
    fmt.Println(tree.RootNode().String())
}
```

### Iterator Pattern

```go
func iterateNodes() {
    sourceCode := []byte(`
        function outer() {
            function inner() {
                return 42;
            }
            return inner();
        }
    `)

    root, _ := sitter.ParseCtx(context.Background(), sourceCode, javascript.GetLanguage())

    // Depth-first iteration
    fmt.Println("DFS Traversal:")
    dfsIter := sitter.NewIterator(root, sitter.DFSMode)
    dfsIter.ForEach(func(node *sitter.Node) error {
        if node.IsNamed() {
            fmt.Printf("- %s\n", node.Type())
        }
        return nil
    })

    // Breadth-first iteration
    fmt.Println("\nBFS Traversal:")
    bfsIter := sitter.NewIterator(root, sitter.BFSMode)
    for {
        node, err := bfsIter.Next()
        if err != nil {
            break
        }
        if node.IsNamed() {
            fmt.Printf("- %s (depth: %d)\n", 
                node.Type(), 
                getDepth(node))
        }
    }
}

func getDepth(node *sitter.Node) int {
    depth := 0
    current := node
    for current.Parent() != nil {
        depth++
        current = current.Parent()
    }
    return depth
}
```

## Practical Examples

### 1. Extract All Imports

```go
func extractImports(sourceCode []byte, lang *sitter.Language) []string {
    parser := sitter.NewParser()
    parser.SetLanguage(lang)
    
    tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)
    root := tree.RootNode()
    
    var query string
    switch lang {
    case golang.GetLanguage():
        query = `(import_declaration (import_spec path: (interpreted_string_literal) @import))`
    case python.GetLanguage():
        query = `[
            (import_statement (dotted_name) @import)
            (import_from_statement module_name: (dotted_name) @import)
        ]`
    case javascript.GetLanguage(), typescript.GetLanguage():
        query = `[
            (import_statement source: (string) @import)
            (import_declaration source: (string) @import)
        ]`
    case java.GetLanguage():
        query = `(import_declaration (scoped_identifier) @import)`
    }
    
    q, _ := sitter.NewQuery([]byte(query), lang)
    qc := sitter.NewQueryCursor()
    qc.Exec(q, root)
    
    var imports []string
    for {
        m, ok := qc.NextMatch()
        if !ok {
            break
        }
        for _, c := range m.Captures {
            imports = append(imports, c.Node.Content(sourceCode))
        }
    }
    
    return imports
}
```

### 2. Find All Function Calls

```go
func findFunctionCalls(sourceCode []byte, lang *sitter.Language) map[string][]string {
    parser := sitter.NewParser()
    parser.SetLanguage(lang)
    
    tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)
    root := tree.RootNode()
    
    var query string
    switch lang {
    case golang.GetLanguage():
        query = `(call_expression function: (identifier) @func_name)`
    case python.GetLanguage():
        query = `(call function: [(identifier) @func_name (attribute) @func_name])`
    case javascript.GetLanguage(), typescript.GetLanguage():
        query = `(call_expression function: [(identifier) @func_name (member_expression) @func_name])`
    case java.GetLanguage():
        query = `(method_invocation name: (identifier) @func_name)`
    }
    
    q, _ := sitter.NewQuery([]byte(query), lang)
    qc := sitter.NewQueryCursor()
    qc.Exec(q, root)
    
    calls := make(map[string][]string)
    for {
        m, ok := qc.NextMatch()
        if !ok {
            break
        }
        for _, c := range m.Captures {
            funcName := c.Node.Content(sourceCode)
            line := c.Node.StartPoint().Row + 1
            calls[funcName] = append(calls[funcName], 
                fmt.Sprintf("line %d", line))
        }
    }
    
    return calls
}
```

### 3. Code Complexity Analysis

```go
func analyzeComplexity(sourceCode []byte, lang *sitter.Language) {
    parser := sitter.NewParser()
    parser.SetLanguage(lang)
    
    tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)
    root := tree.RootNode()
    
    // Count different complexity indicators
    var functionQuery, loopQuery, conditionalQuery string
    
    switch lang {
    case golang.GetLanguage():
        functionQuery = `(function_declaration)`
        loopQuery = `[(for_statement) (range_statement)]`
        conditionalQuery = `[(if_statement) (switch_statement)]`
    case python.GetLanguage():
        functionQuery = `(function_definition)`
        loopQuery = `[(for_statement) (while_statement)]`
        conditionalQuery = `[(if_statement) (match_statement)]`
    case javascript.GetLanguage(), typescript.GetLanguage():
        functionQuery = `[(function_declaration) (arrow_function) (function)]`
        loopQuery = `[(for_statement) (while_statement) (do_statement)]`
        conditionalQuery = `[(if_statement) (switch_statement) (ternary_expression)]`
    }
    
    fmt.Println("Code Complexity Metrics:")
    fmt.Printf("Functions: %d\n", countMatches(root, functionQuery, lang))
    fmt.Printf("Loops: %d\n", countMatches(root, loopQuery, lang))
    fmt.Printf("Conditionals: %d\n", countMatches(root, conditionalQuery, lang))
}

func countMatches(root *sitter.Node, query string, lang *sitter.Language) int {
    q, _ := sitter.NewQuery([]byte(query), lang)
    qc := sitter.NewQueryCursor()
    qc.Exec(q, root)
    
    count := 0
    for {
        _, ok := qc.NextMatch()
        if !ok {
            break
        }
        count++
    }
    return count
}
```

### 4. Documentation Extractor

```go
func extractDocumentation(sourceCode []byte, lang *sitter.Language) {
    parser := sitter.NewParser()
    parser.SetLanguage(lang)
    
    tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)
    root := tree.RootNode()
    
    var query string
    switch lang {
    case golang.GetLanguage():
        // Go uses comments above declarations
        query = `(
            (comment)+ @doc
            .
            (function_declaration name: (identifier) @name)
        )`
    case python.GetLanguage():
        // Python uses docstrings
        query = `(
            function_definition
            name: (identifier) @name
            body: (block . (expression_statement (string) @doc))
        )`
    case java.GetLanguage():
        // Java uses Javadoc comments
        query = `(
            (block_comment) @doc
            .
            (method_declaration name: (identifier) @name)
            (#match? @doc "^/\\*\\*")
        )`
    }
    
    if query == "" {
        return
    }
    
    q, _ := sitter.NewQuery([]byte(query), lang)
    qc := sitter.NewQueryCursor()
    qc.Exec(q, root)
    
    fmt.Println("Documented Functions:")
    for {
        m, ok := qc.NextMatch()
        if !ok {
            break
        }
        m = qc.FilterPredicates(m, sourceCode)
        
        var name, doc string
        for _, c := range m.Captures {
            captureName := q.CaptureNameForId(c.Index)
            content := c.Node.Content(sourceCode)
            
            switch captureName {
            case "name":
                name = content
            case "doc":
                doc = content
            }
        }
        
        if name != "" && doc != "" {
            fmt.Printf("\nFunction: %s\n", name)
            fmt.Printf("Documentation: %s\n", doc)
        }
    }
}
```

## Best Practices

### 1. Parser Reuse
```go
// Good: Reuse parser instances
var globalParser = sitter.NewParser()

func parseMultipleFiles(files []string, lang *sitter.Language) {
    globalParser.SetLanguage(lang)
    
    for _, file := range files {
        content, _ := os.ReadFile(file)
        tree, _ := globalParser.ParseCtx(context.Background(), nil, content)
        // Process tree...
    }
}

// Bad: Creating new parser for each file
func badExample(files []string, lang *sitter.Language) {
    for _, file := range files {
        parser := sitter.NewParser() // Wasteful
        parser.SetLanguage(lang)
        // ...
    }
}
```

### 2. Context Usage
```go
// Always use context for cancellation support
func parseWithTimeout(sourceCode []byte, lang *sitter.Language) (*sitter.Tree, error) {
    parser := sitter.NewParser()
    parser.SetLanguage(lang)
    
    ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
    defer cancel()
    
    return parser.ParseCtx(ctx, nil, sourceCode)
}
```

### 3. Error Handling
```go
func robustParsing(sourceCode []byte, lang *sitter.Language) (*sitter.Node, error) {
    parser := sitter.NewParser()
    
    // Check language is set
    if lang == nil {
        return nil, fmt.Errorf("language not specified")
    }
    parser.SetLanguage(lang)
    
    // Set operation limit for safety
    parser.SetOperationLimit(1000000) // microseconds
    
    tree, err := parser.ParseCtx(context.Background(), nil, sourceCode)
    if err != nil {
        switch err {
        case sitter.ErrOperationLimit:
            return nil, fmt.Errorf("parsing took too long")
        case sitter.ErrNoLanguage:
            return nil, fmt.Errorf("no language set")
        default:
            return nil, fmt.Errorf("parsing failed: %w", err)
        }
    }
    
    root := tree.RootNode()
    if root.HasError() {
        // Find and report syntax errors
        var errors []string
        findErrors(root, sourceCode, &errors)
        return root, fmt.Errorf("syntax errors found: %v", errors)
    }
    
    return root, nil
}

func findErrors(node *sitter.Node, source []byte, errors *[]string) {
    if node.IsError() || node.IsMissing() {
        *errors = append(*errors, fmt.Sprintf(
            "Error at %d:%d: %s",
            node.StartPoint().Row+1,
            node.StartPoint().Column+1,
            node.String(),
        ))
    }
    
    for i := 0; i < int(node.ChildCount()); i++ {
        findErrors(node.Child(i), source, errors)
    }
}
```

### 4. Memory Management
```go
// Tree-sitter objects are automatically garbage collected
// but you can explicitly close them for immediate cleanup

func processLargeFile(filename string, lang *sitter.Language) error {
    parser := sitter.NewParser()
    defer parser.Close() // Optional but good practice
    
    parser.SetLanguage(lang)
    
    content, err := os.ReadFile(filename)
    if err != nil {
        return err
    }
    
    tree, err := parser.ParseCtx(context.Background(), nil, content)
    if err != nil {
        return err
    }
    defer tree.Close() // Optional but good practice
    
    // Process tree...
    
    return nil
}
```

### 5. Query Optimization
```go
// Cache compiled queries for reuse
var queryCache = make(map[string]*sitter.Query)
var queryCacheMu sync.RWMutex

func getCachedQuery(pattern string, lang *sitter.Language) (*sitter.Query, error) {
    key := fmt.Sprintf("%s:%p", pattern, lang)
    
    queryCacheMu.RLock()
    if q, ok := queryCache[key]; ok {
        queryCacheMu.RUnlock()
        return q, nil
    }
    queryCacheMu.RUnlock()
    
    queryCacheMu.Lock()
    defer queryCacheMu.Unlock()
    
    // Double-check after acquiring write lock
    if q, ok := queryCache[key]; ok {
        return q, nil
    }
    
    q, err := sitter.NewQuery([]byte(pattern), lang)
    if err != nil {
        return nil, err
    }
    
    queryCache[key] = q
    return q, nil
}
```

## Conclusion

go-tree-sitter provides a powerful and efficient way to parse and analyze source code across multiple programming languages. Its key strengths include:

- **Performance**: Fast parsing with incremental updates
- **Accuracy**: Precise syntax trees that handle edge cases
- **Flexibility**: Powerful query system for pattern matching
- **Language Support**: Consistent API across different languages
- **Error Recovery**: Robust handling of syntax errors

Whether you're building development tools, code analyzers, or syntax highlighters, go-tree-sitter offers the foundation for sophisticated source code processing in Go.

### Resources

- [Official tree-sitter documentation](https://tree-sitter.github.io/tree-sitter/)
- [go-tree-sitter repository](https://github.com/smacker/go-tree-sitter)
- [Tree-sitter playground](https://tree-sitter.github.io/tree-sitter/playground)
- [Query syntax reference](https://tree-sitter.github.io/tree-sitter/using-parsers#pattern-matching-with-queries)