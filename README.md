# glyph

A Model Context Protocol (MCP) server that extracts symbol outlines from your codebase using Tree-sitter's declarative query language. Gives LLM coding agents the context they need with clean, efficient symbol extraction. Can also be used as a standalone CLI tool.

## Installation

### macOS

Install Go:
```shell
$ brew install go
```

Install the latest version of glyph:
```shell
$ GOBIN=/usr/local/bin go install "github.com/benmyles/glyph@latest"
```

## What it does

glyph takes a file path glob, recursively discovers matching code files, and parses them using Tree-sitter's powerful query language to generate accurate symbol outlines. You control the level of detail—from high-level structure to complete function signatures with visibility modifiers.

## Why glyph?

LLM coding agents work best when they understand your code's structure. glyph bridges that gap by providing clean, multi-file symbol maps that serve as efficient context without overwhelming token limits.

## Key Features

- **Declarative Query-Based Extraction** - Uses Tree-sitter's query language for precise, maintainable symbol extraction
- **Glob-based file discovery** - Point it at your code with familiar patterns
- **Language-agnostic architecture** - Easy to add new languages with just query patterns
- **Configurable detail levels** - Choose how much context you need
- **Multi-file support** - Get a unified view across your entire project
- **MCP-native** - Built for seamless integration with AI coding workflows
- **High performance** - Optimized Tree-sitter queries for fast extraction

## Supported Languages

- **Go** - Functions, methods, types, structs, interfaces, constants, variables
- **Java** - Classes, interfaces, methods, constructors, fields, enums, records, annotations
- **JavaScript/TypeScript** - Functions, classes, methods, arrow functions, variables, interfaces, type aliases
- **Python** - Functions, classes, decorated definitions, assignments
- **Easy to extend** - Adding new languages requires only ~20 lines of query patterns

## Architecture

glyph uses a modern, declarative approach:

- **Tree-sitter queries** replace complex AST traversal logic
- **Language-specific query patterns** define what symbols to extract
- **Unified extraction engine** works across all languages
- **Modular design** with separate concerns for queries, extraction, formatting, and file handling

## Symbol Types

glyph uses concise symbol type names to minimize token usage:

- `func` - Functions and function-like declarations
- `method` - Class/struct methods
- `class` - Classes
- `interface` - Interfaces
- `struct` - Structs (Go)
- `type` - Type declarations (Go)
- `const` - Constants
- `var` - Variables
- `field` - Class/struct fields
- `constructor` - Constructors
- `enum` - Enumerations
- `record` - Records (Java)
- `annotation` - Annotations (Java)
- `property` - Properties (TypeScript)

## Usage

### Integration with AI Coding Assistants

#### Claude Code

Add glyph to Claude Code using the MCP command:

```bash
# Add glyph as a local MCP server (for your current project)
claude mcp add glyph /usr/local/bin/glyph mcp
```

To verify the installation:
```bash
# List all configured servers
claude mcp list

# Check glyph server details
claude mcp get glyph
```

You can also check the server status anytime within Claude Code using the `/mcp` command.

#### Cursor

Add glyph to your Cursor MCP configuration:

1. Open Cursor settings (Cmd/Ctrl + ,)
2. Search for "MCP" or navigate to Extensions → MCP
3. Add glyph to your MCP servers:

```json
{
  "mcpServers": {
    "glyph": {
      "command": "/usr/local/bin/glyph",
      "args": ["mcp"]
    }
  }
}
```

After configuration, you can use glyph through the MCP interface to extract symbol outlines from your codebase, helping the AI understand your code structure better.

### MCP Server Mode (default)

Run glyph as an MCP server for integration with LLM coding agents:

```bash
$ glyph mcp
```

### CLI Mode

Use glyph directly from the command line to extract symbols:

```bash
$ glyph cli '/path/to/project/*.go'
```

```bash
$ glyph cli -detail=minimal '/path/to/project/**/*.js'
```

Options:
- `-detail`: Level of detail (`minimal`, `standard`, or `full`). Default is `standard`.

Note: All file patterns must be absolute paths.

## Detail Levels

### Minimal
Shows just symbol names and types with line numbers:
```
- func: main (line 15)
- struct: Server (line 5)
- method: Start (line 10)
```

### Standard (default)
Shows signatures and declarations:
```
- func: func main()
- struct: type Server struct
- method: func (s *Server) Start() error
```

### Full
Shows complete symbol definitions with code blocks:
```
- func (lines 15-20):
  ```
  func main() {
      server := &Server{}
      server.Start()
  }
  ```
```

## Development

```bash
$ git clone https://github.com/benmyles/glyph
$ cd glyph
```

```shell
$ make clean
```

```shell
$ make test
```

```shell
$ make build
```

```bash
$ make install DESTDIR=/usr/local/bin
```

## Adding New Languages

Thanks to the query-based architecture, adding support for a new language is straightforward:

1. **Add language detection** in `file_utils.go` (~2 lines)
2. **Add query patterns** in `queries.go` (~10-20 lines)
3. **Add to query dispatcher** in `queries.go` (~5 lines)

Example for adding Rust support:

```go
// In file_utils.go
case ".rs":
    return rust.GetLanguage(), nil

// In queries.go
var rustQueries = map[string]string{
    "functions": `
        (function_item
            name: (identifier) @name
        ) @function
    `,
    "structs": `
        (struct_item
            name: (type_identifier) @name
        ) @struct
    `,
    // ... more patterns
}
```

## Performance

- **Fast parsing** with Tree-sitter's incremental parsing
- **Optimized queries** for efficient symbol extraction
- **Minimal memory usage** with streaming file processing
- **Parser reuse** for better performance across multiple files