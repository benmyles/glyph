# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Project Overview

glyph is a Model Context Protocol (MCP) server written in Go that extracts symbol outlines from codebases using tree-sitter parsing. It helps LLM coding agents understand code structure by providing clean, multi-file symbol maps.

## Architecture

The project uses:
- **tree-sitter** (`github.com/smacker/go-tree-sitter`) for AST parsing and symbol extraction
- **MCP Go SDK** (`github.com/mark3labs/mcp-go`) for the Model Context Protocol server implementation

## Development Commands

### Building
```bash
go build
```

### Running Tests
```bash
go test ./...
```

### Installing Dependencies
```bash
go mod download
```

### Updating Dependencies
```bash
go mod tidy
```

## Key Implementation Notes

- The server processes file paths using glob patterns
- Symbol extraction is performed via tree-sitter AST parsing
- Output format should be optimized for LLM context windows
- Detail levels should be configurable (from high-level structure to full signatures)

## Documentation

```
├── README.md
├── doc -- project overview
│ └── llm -- docs tailored for AI agents like yourself
│     ├── go-tree-sitter.md -- go-tree-sitter package docs
│     └── mcp-go.md -- mcp-go package docs
```

- If a p