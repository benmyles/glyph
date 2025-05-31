package main

import (
	"fmt"
	"strings"
)

// FormatSymbols formats symbols for output
func FormatSymbols(symbols []Symbol, detailLevel DetailLevel) string {
	if len(symbols) == 0 {
		return "No symbols found"
	}

	var sb strings.Builder
	sb.WriteString("# Symbol Outline\n\n")

	// Group symbols by file
	fileSymbols := make(map[string][]Symbol)
	for _, sym := range symbols {
		fileSymbols[sym.FilePath] = append(fileSymbols[sym.FilePath], sym)
	}

	// Format output
	for file, syms := range fileSymbols {
		sb.WriteString(fmt.Sprintf("## %s\n\n", file))

		for _, sym := range syms {
			formatSymbol(&sb, sym, detailLevel, 0)
		}

		sb.WriteString("\n")
	}

	return sb.String()
}

func formatSymbol(sb *strings.Builder, symbol Symbol, detailLevel DetailLevel, indent int) {
	indentStr := strings.Repeat("  ", indent)

	switch detailLevel {
	case Minimal:
		sb.WriteString(fmt.Sprintf("%s- %s: %s (line %d)\n",
			indentStr, symbol.Kind, symbol.Name, symbol.StartLine))
	case Standard:
		if symbol.Signature != "" {
			// For variables and constants, show name with type/signature
			if symbol.Kind == "var" || symbol.Kind == "const" {
				// Avoid duplicate names when signature equals name
				if symbol.Signature == symbol.Name {
					sb.WriteString(fmt.Sprintf("%s- %s: %s\n",
						indentStr, symbol.Kind, symbol.Name))
				} else {
					sb.WriteString(fmt.Sprintf("%s- %s: %s %s\n",
						indentStr, symbol.Kind, symbol.Name, symbol.Signature))
				}
			} else {
				sb.WriteString(fmt.Sprintf("%s- %s: %s\n",
					indentStr, symbol.Kind, symbol.Signature))
			}
		} else {
			sb.WriteString(fmt.Sprintf("%s- %s: %s (lines %d-%d)\n",
				indentStr, symbol.Kind, symbol.Name, symbol.StartLine, symbol.EndLine))
		}
	case Full:
		sb.WriteString(fmt.Sprintf("%s- %s (lines %d-%d):\n",
			indentStr, symbol.Kind, symbol.StartLine, symbol.EndLine))
		if symbol.Signature != "" {
			sb.WriteString(fmt.Sprintf("%s  ```\n%s  %s\n%s  ```\n",
				indentStr, indentStr, symbol.Signature, indentStr))
		}
	}
}
