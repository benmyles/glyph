package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Create MCP server
	mcpServer := server.NewMCPServer(
		"glyph",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Register tools
	extractSymbolsTool := mcp.NewTool(
		"extract_symbols",
		mcp.WithDescription("Extract symbol outlines from source code files using tree-sitter parsing"),
		mcp.WithString("pattern", mcp.Required(), mcp.Description("Glob pattern to match files (e.g., '**/*.go', 'src/**/*.js')")),
		mcp.WithString("detail", mcp.Description("Level of detail: 'minimal', 'standard', 'full' (default: 'standard')")),
	)

	mcpServer.AddTool(extractSymbolsTool, extractSymbolsHandler)

	// Start server
	if err := server.ServeStdio(mcpServer); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func extractSymbolsHandler(_ context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	pattern, err := request.RequireString("pattern")
	if err != nil {
		return mcp.NewToolResultError("pattern argument is required"), nil
	}

	detail := "standard"
	if d := request.GetString("detail", ""); d != "" {
		detail = d
	}

	// Extract symbols from files matching the pattern
	result, err := extractSymbols(pattern, detail)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to extract symbols: %v", err)), nil
	}

	return mcp.NewToolResultText(result), nil
}