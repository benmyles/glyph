# glyph

A Model Context Protocol (MCP) server that extracts symbol outlines from your codebase, giving LLM coding agents the context they need. Can also be used as a standalone CLI tool for direct symbol extraction.

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

glyph takes a file path glob, recursively discovers matching code files, and parses them using tree-sitter to generate accurate symbol outlines. You control the level of detail—from high-level structure to complete function signatures with visibility modifiers.

## Why glyph?

LLM coding agents work best when they understand your code's structure. glyph bridges that gap by providing clean, multi-file symbol maps that serve as efficient context without overwhelming token limits.

## Key Features

- **Glob-based file discovery** - Point it at your code with familiar patterns
- **Tree-sitter parsing** - Language-aware AST parsing for accurate symbol extraction
- **Configurable detail levels** - Choose how much context you need
- **Multi-file support** - Get a unified view across your entire project
- **MCP-native** - Built for seamless integration with AI coding workflows

## Supported Languages

- Go
- JavaScript/TypeScript
- Python
- Java (including modern features like records, sealed classes, and lambdas)
- Easy to add any language with a tree-sitter parser

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
- `trait` - Traits (Rust)
- `impl` - Implementations (Rust)
- `record` - Records (Java)
- `annotation` - Annotations (Java)
- `static_init` - Static initializers
- `init` - Instance initializers
- `decorated` - Decorated definitions (Python)

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
- `-detail`: Level of detail (`minimal` or `standard`). Default is `standard`.

Note: All file patterns must be absolute paths.

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