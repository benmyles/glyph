package main

import "strings"

// Symbol represents a code symbol with its metadata
type Symbol struct {
	Name      string
	Kind      string
	StartLine uint32
	EndLine   uint32
	Signature string
	FilePath  string
}

// DetailLevel controls how much information to include in symbol extraction
type DetailLevel int

const (
	Minimal DetailLevel = iota
	Standard
	Full
)

// ParseDetailLevel converts a string to DetailLevel
func ParseDetailLevel(detail string) DetailLevel {
	switch strings.ToLower(detail) {
	case "minimal":
		return Minimal
	case "full":
		return Full
	default:
		return Standard
	}
}
