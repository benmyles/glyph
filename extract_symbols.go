package main

import (
	"fmt"
)

// ExtractSymbols extracts symbols from files matching a pattern
func ExtractSymbols(pattern string, detail string) (string, error) {
	detailLevel := ParseDetailLevel(detail)

	// Find files matching the pattern
	files, err := FindFiles(pattern)
	if err != nil {
		return "", fmt.Errorf("failed to find files: %w", err)
	}

	if len(files) == 0 {
		return "No files found matching pattern: " + pattern, nil
	}

	var allSymbols []Symbol
	extractor := NewSymbolExtractor()

	for _, file := range files {
		symbols, err := extractor.ExtractFromFile(file, detailLevel)
		if err != nil {
			continue // Skip files that can't be parsed
		}
		allSymbols = append(allSymbols, symbols...)
	}

	if len(allSymbols) == 0 {
		return "No symbols found", nil
	}

	return FormatSymbols(allSymbols, detailLevel), nil
}
