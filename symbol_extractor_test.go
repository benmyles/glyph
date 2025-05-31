package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestSymbolExtractor_ExtractFromFile_Go(t *testing.T) {
	// Create a test Go file
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "test.go")

	testCode := `package main

import "fmt"

type Server struct {
	host string
	port int
}

func (s *Server) Start() error {
	fmt.Printf("Starting server on %s:%d\n", s.host, s.port)
	return nil
}

func main() {
	server := &Server{
		host: "localhost",
		port: 8080,
	}
	server.Start()
}

const Version = "1.0.0"

var Config = map[string]string{
	"env": "production",
}
`

	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		t.Fatal(err)
	}

	extractor := NewSymbolExtractor()

	// Test different detail levels
	tests := []struct {
		detail   DetailLevel
		contains []string
	}{
		{
			detail:   Minimal,
			contains: []string{"struct: Server", "method: Start", "func: main"},
		},
		{
			detail:   Standard,
			contains: []string{"Server", "Start", "main"},
		},
		{
			detail:   Full,
			contains: []string{"Server", "Start", "main"},
		},
	}

	for _, tt := range tests {
		t.Run(tt.detail.String(), func(t *testing.T) {
			symbols, err := extractor.ExtractFromFile(testFile, tt.detail)
			if err != nil {
				t.Fatalf("ExtractFromFile error = %v", err)
			}

			if len(symbols) == 0 {
				t.Errorf("No symbols were extracted")
			}

			// Convert symbols to string for easier checking
			result := FormatSymbols(symbols, tt.detail)
			t.Logf("Result for %v:\n%s", tt.detail, result)

			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Result does not contain expected string %q", expected)
				}
			}
		})
	}
}

func TestSymbolExtractor_ExtractFromFile_Java(t *testing.T) {
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "Example.java")

	testCode := `package com.example;

import java.util.*;

public class Example {
    private String name;
    
    public Example(String name) {
        this.name = name;
    }
    
    public String getName() {
        return name;
    }
    
    public static void main(String[] args) {
        Example ex = new Example("test");
        System.out.println(ex.getName());
    }
}
`

	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		t.Fatal(err)
	}

	extractor := NewSymbolExtractor()
	symbols, err := extractor.ExtractFromFile(testFile, Standard)
	if err != nil {
		t.Fatalf("ExtractFromFile error = %v", err)
	}

	if len(symbols) == 0 {
		t.Errorf("No symbols were extracted")
	}

	result := FormatSymbols(symbols, Standard)
	t.Logf("Result:\n%s", result)

	expectedSymbols := []string{"class", "field", "constructor", "method"}
	for _, expected := range expectedSymbols {
		if !strings.Contains(result, expected) {
			t.Errorf("Result does not contain expected symbol type %q", expected)
		}
	}
}

func TestSymbolExtractor_ExtractFromFile_JavaScript(t *testing.T) {
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "example.js")

	testCode := `class Calculator {
    constructor() {
        this.result = 0;
    }
    
    add(a, b) {
        return a + b;
    }
}

const square = (x) => x * x;

async function fetchData(url) {
    const response = await fetch(url);
    return response.json();
}

function greet(name) {
    return "Hello, " + name;
}
`

	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		t.Fatal(err)
	}

	extractor := NewSymbolExtractor()
	symbols, err := extractor.ExtractFromFile(testFile, Standard)
	if err != nil {
		t.Fatalf("ExtractFromFile error = %v", err)
	}

	if len(symbols) == 0 {
		t.Errorf("No symbols were extracted")
	}

	result := FormatSymbols(symbols, Standard)
	t.Logf("Result:\n%s", result)

	expectedSymbols := []string{"class", "method", "func"}
	for _, expected := range expectedSymbols {
		if !strings.Contains(result, expected) {
			t.Errorf("Result does not contain expected symbol type %q", expected)
		}
	}
}

func TestSymbolExtractor_ExtractFromFile_Python(t *testing.T) {
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "example.py")

	testCode := `import math

class DataProcessor:
    def __init__(self, name):
        self.name = name
        self.data = []
    
    def add_data(self, value):
        self.data.append(value)
    
    @property
    def mean(self):
        if not self.data:
            return None
        return sum(self.data) / len(self.data)

def main():
    processor = DataProcessor("test")
    processor.add_data(10)
    print(processor.mean)

if __name__ == "__main__":
    main()
`

	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		t.Fatal(err)
	}

	extractor := NewSymbolExtractor()
	symbols, err := extractor.ExtractFromFile(testFile, Standard)
	if err != nil {
		t.Fatalf("ExtractFromFile error = %v", err)
	}

	if len(symbols) == 0 {
		t.Errorf("No symbols were extracted")
	}

	result := FormatSymbols(symbols, Standard)
	t.Logf("Result:\n%s", result)

	expectedSymbols := []string{"class", "func"}
	for _, expected := range expectedSymbols {
		if !strings.Contains(result, expected) {
			t.Errorf("Result does not contain expected symbol type %q", expected)
		}
	}
}

// Helper method for DetailLevel
func (d DetailLevel) String() string {
	switch d {
	case Minimal:
		return "minimal"
	case Standard:
		return "standard"
	case Full:
		return "full"
	default:
		return "unknown"
	}
}
