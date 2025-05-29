# glyph

A Model Context Protocol (MCP) server that extracts symbol outlines from your codebase, giving LLM coding agents the context they need.

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