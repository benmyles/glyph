# glyph

A Model Context Protocol (MCP) server that extracts symbol outlines from your codebase, giving LLM coding agents the context they need. Can also be used as a standalone CLI tool for direct symbol extraction.

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

## Usage

### MCP Server Mode (default)

Run glyph as an MCP server for integration with LLM coding agents:

```bash
glyph mcp
```

### CLI Mode

Use glyph directly from the command line to extract symbols:

```bash
# Extract symbols from Go files
glyph cli '/path/to/project/*.go'

# Extract minimal symbols from all JavaScript files recursively
glyph cli -detail=minimal '/path/to/project/**/*.js'
```

Options:
- `-detail`: Level of detail (`minimal` or `standard`). Default is `standard`.

Note: All file patterns must be absolute paths.

## Installation

### Building from source

```bash
git clone https://github.com/benmyles/glyph
cd glyph
make build
```

Or manually:
```bash
go build -o glyph
```

### Installing to system

```bash
# Build and install to a directory in your PATH
make build
make install DESTDIR=/usr/local/bin

# Or install to a custom location
make install DESTDIR=/path/to/install/dir
```

### Development

```bash
# Run tests
make test

# Clean build artifacts
make clean
```

### Integration with AI Coding Assistants

#### Claude Code

Add glyph to Claude Code using the MCP command:

```bash
# Add glyph as a local MCP server (for your current project)
claude mcp add glyph /path/to/glyph mcp
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
      "command": "/path/to/glyph",
      "args": ["mcp"]
    }
  }
}
```

After configuration, you can use glyph through the MCP interface to extract symbol outlines from your codebase, helping the AI understand your code structure better.
