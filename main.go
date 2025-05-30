package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"path/filepath"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	// Check if running with subcommands
	if len(os.Args) < 2 {
		printUsage()
		os.Exit(1)
	}

	switch os.Args[1] {
	case "mcp":
		runMCPServer(os.Args[2:])
	case "cli":
		runCLI(os.Args[2:])
	default:
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Fprintf(os.Stderr, "Usage: %s [mcp|cli] [options]\n", os.Args[0])
	fmt.Fprintf(os.Stderr, "  mcp  - Run as MCP server (default)\n")
	fmt.Fprintf(os.Stderr, "  cli  - Run in CLI mode\n")
}

func validateAbsolutePath(pattern string) error {
	if !filepath.IsAbs(pattern) {
		return fmt.Errorf("pattern must be an absolute path, got: %s", pattern)
	}
	return nil
}

func runCLI(args []string) {
	// Set up CLI flags
	cliFlags := flag.NewFlagSet("cli", flag.ExitOnError)
	detail := cliFlags.String("detail", "standard", "Level of detail: minimal or standard")

	cliFlags.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s cli [options] <pattern>\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "\nOptions:\n")
		cliFlags.PrintDefaults()
		fmt.Fprintf(os.Stderr, "\nExamples:\n")
		fmt.Fprintf(os.Stderr, "  %s cli '/path/to/project/*.go'                    # Extract symbols from all .go files\n", os.Args[0])
		fmt.Fprintf(os.Stderr, "  %s cli -detail=minimal '/path/to/project/**/*.js' # Extract minimal symbols from all .js files\n", os.Args[0])
	}

	if err := cliFlags.Parse(args); err != nil {
		os.Exit(1)
	}

	// Check for pattern argument
	if cliFlags.NArg() < 1 {
		cliFlags.Usage()
		os.Exit(1)
	}

	pattern := cliFlags.Arg(0)
	if err := validateAbsolutePath(pattern); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Extract symbols
	result, err := extractSymbols(pattern, *detail)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}

	// Print results to stdout
	fmt.Print(result)
}

func runMCPServer(args []string) {
	// Set up MCP flags
	mcpFlags := flag.NewFlagSet("mcp", flag.ExitOnError)

	if err := mcpFlags.Parse(args); err != nil {
		os.Exit(1)
	}

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
		mcp.WithString("pattern", mcp.Required(), mcp.Description("Absolute path glob pattern to match files (e.g., '/path/to/project/**/*.go', '/home/user/src/**/*.js')")),
		mcp.WithString("detail", mcp.Description("Level of detail: 'minimal', 'standard' (default: 'standard')")),
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

	if err := validateAbsolutePath(pattern); err != nil {
		return mcp.NewToolResultError(err.Error()), nil
	}

	// Extract symbols from files matching the pattern
	result, err := extractSymbols(pattern, detail)
	if err != nil {
		return mcp.NewToolResultError(fmt.Sprintf("failed to extract symbols: %v", err)), nil
	}

	return mcp.NewToolResultText(result), nil
}
