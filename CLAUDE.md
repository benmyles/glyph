# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

glyph is a Model Context Protocol (MCP) server written in Go that extracts symbol outlines from codebases using Tree-sitter's declarative query language. It helps LLM coding agents understand code structure by providing clean, multi-file symbol maps with minimal token usage.

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

## Architecture

The project uses a modern, declarative approach:

- **Tree-sitter queries** (`queries.go`) - Declarative patterns for symbol extraction
- **MCP Go SDK** (`github.com/mark3labs/mcp-go`) for the Model Context Protocol server implementation
- **Modular design** with clear separation of concerns

### Key Components

- `queries.go` - Language-specific Tree-sitter query patterns
- `symbol_extractor.go` - Query-based symbol extraction engine
- `file_utils.go` - File discovery and language detection
- `types.go` - Core data structures (Symbol, DetailLevel)
- `formatter.go` - Output formatting for different detail levels
- `extract_symbols.go` - Main extraction orchestration
- `main.go` - CLI and MCP server entry points

## Development Commands

### Building
```bash
make build
```

### Running Tests
```bash
make test
```

### Installing
```bash
make install DESTDIR=/usr/local/bin
```

### Cleaning
```bash
make clean
```

### Manual commands (if needed)
```bash
go build         # Build manually
go test ./...    # Run tests manually
go mod download  # Install dependencies
go mod tidy      # Update dependencies
```

## Key Implementation Notes

### Query-Based Symbol Extraction
- Uses Tree-sitter's S-expression query syntax for precise symbol matching
- Each language has dedicated query patterns in `queries.go`
- Queries capture symbol names and types declaratively
- Much simpler and more maintainable than manual AST traversal

### Language Support
- File type detection based on extensions in `file_utils.go`
- Language-specific queries in `queries.go`
- Adding new languages requires only ~20 lines of code
- Currently supports: Go, Java, JavaScript, TypeScript, Python

### Symbol Processing
- Three detail levels: minimal, standard, full
- Configurable output formatting
- Efficient processing with parser reuse
- Graceful error handling for unsupported files

### File Pattern Handling
- Supports glob patterns with `**` for recursive matching
- All file patterns must be absolute paths (relative paths not supported)
- Robust file discovery with error tolerance

### Output Format
- Markdown-formatted symbol outlines
- Grouped by file for multi-file processing
- Token-efficient symbol type names (e.g., `func` instead of `function_declaration`)
- Configurable detail levels for different use cases

## Adding New Language Support

Thanks to the query-based architecture, adding a new language is straightforward:

1. **Add language detection** in `file_utils.go`:
   ```go
   case ".rs":
       return rust.GetLanguage(), nil
   ```

2. **Add query patterns** in `queries.go`:
   ```go
   var rustQueries = map[string]string{
       "functions": `
           (function_item
               name: (identifier) @name
           ) @function
       `,
       // ... more patterns
   }
   ```

3. **Add to query dispatcher**:
   ```go
   case ".rs":
       return &LanguageQueries{
           Language: rust.GetLanguage(),
           Queries:  rustQueries,
       }
   ```

## Testing Strategy

- **Unit tests** for each component (`*_test.go` files)
- **Integration tests** with real code samples
- **Language-specific tests** for each supported language
- **File utility tests** for pattern matching and discovery
- All tests use temporary directories for isolation

## Documentation

```
├── README.md
├── doc -- project overview
│ └── llm -- docs tailored for AI agents like yourself
│     ├── go-tree-sitter.md -- go-tree-sitter package docs
│     └── mcp-go.md -- mcp-go package docs
├── REFACTOR_SUMMARY.md -- details of the query-based refactoring
```

- If you're using a dependency that's documented in `doc/llm`: always review the docs before writing code that uses the dependency.
- The refactoring summary provides context on the architectural improvements made.

## Performance Considerations

- **Parser reuse**: Single parser instance per language for better performance
- **Streaming processing**: Files processed individually to minimize memory usage
- **Query optimization**: Tree-sitter queries are compiled once and reused
- **Error tolerance**: Failed queries don't stop processing of other symbols

## Code Quality

- **Separation of concerns**: Each file has a single, clear responsibility
- **Error handling**: Graceful degradation with informative error messages
- **Testing**: Comprehensive test coverage for all components
- **Documentation**: Self-documenting query patterns and clear function names
- **Maintainability**: Simple, declarative approach reduces complexity