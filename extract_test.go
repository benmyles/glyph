package main

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

func TestParseDetailLevel(t *testing.T) {
	tests := []struct {
		input    string
		expected DetailLevel
	}{
		{"minimal", Minimal},
		{"standard", Standard},
		{"full", Full},
		{"", Standard},
		{"invalid", Standard},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			result := parseDetailLevel(tt.input)
			if result != tt.expected {
				t.Errorf("parseDetailLevel(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

func TestExtractSymbols(t *testing.T) {
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
	
	// Test different detail levels
	tests := []struct {
		detail   string
		contains []string
	}{
		{
			detail:   "minimal",
			contains: []string{"type: Server", "method: Start", "func: main"},
		},
		{
			detail:   "standard",
			contains: []string{"type: Server struct", "func (s *Server) Start() error", "func main()"},
		},
		{
			detail:   "full",
			contains: []string{"Server struct", "func (s *Server) Start() error", "```"},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.detail, func(t *testing.T) {
			result, err := extractSymbols(testFile, tt.detail)
			if err != nil {
				t.Fatalf("extractSymbols error = %v", err)
			}
			
			// Check if we have expected content
			if result == "No symbols found" {
				t.Errorf("No symbols were extracted")
			}
			
			// For debugging
			t.Logf("Result for %s:\n%s", tt.detail, result)
			
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Result does not contain expected string %q", expected)
				}
			}
		})
	}
}

func TestJavaAnnotatedClass(t *testing.T) {
	// Test case for Java class with annotations (like DataLabCli)
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "DataLabCli.java")
	
	testCode := `package com.example.cli;

import picocli.CommandLine.Command;
import picocli.CommandLine.Option;

@Command(name = "datalab", 
         mixinStandardHelpOptions = true, 
         version = "1.0",
         description = "DataLab API Client for PDF to Markdown conversion",
         subcommands = {DataLabCli.SubmitCommand.class, DataLabCli.GetCommand.class})
public class DataLabCli implements Callable<Integer> {
    
    @Option(names = {"-u", "--url"}, 
            description = "DataLab API base URL", 
            defaultValue = "https://api.datalab.to")
    private String baseUrl;
    
    @Command(name = "submit", description = "Submit a PDF file for marker processing")
    static class SubmitCommand implements Callable<Integer> {
        
        @Option(names = {"--max-pages"}, description = "Maximum number of pages to process")
        private Integer maxPages;
        
        public Integer call() throws Exception {
            return 0;
        }
    }
}
`

	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		t.Fatal(err)
	}
	
	// Test different detail levels
	tests := []struct {
		detail   string
		contains []string
	}{
		{
			detail:   "minimal",
			contains: []string{
				"class: DataLabCli",
				"field: baseUrl",
				"class: SubmitCommand",
				"field: maxPages",
			},
		},
		{
			detail:   "standard",
			contains: []string{
				"public class DataLabCli implements Callable<Integer>",
				"private String baseUrl",
				"static class SubmitCommand implements Callable<Integer>",
				"private Integer maxPages",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.detail, func(t *testing.T) {
			result, err := extractSymbols(testFile, tt.detail)
			if err != nil {
				t.Fatalf("extractSymbols error = %v", err)
			}
			
			// For debugging
			t.Logf("Result for %s:\n%s", tt.detail, result)
			
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Result does not contain expected string %q", expected)
				}
			}
		})
	}
}

func TestExtractJavaSymbols(t *testing.T) {
	// Create a comprehensive test Java file
	testDir := t.TempDir()
	testFile := filepath.Join(testDir, "Example.java")
	
	testCode := `package com.example.test;

import java.util.*;
import java.util.stream.Collectors;

/**
 * Example class demonstrating modern Java features
 */
@Deprecated
public class Example<T extends Comparable<T>> {
    // Constants
    public static final String VERSION = "1.0.0";
    private static final int MAX_SIZE = 100;
    
    // Fields
    private String name;
    protected List<T> items = new ArrayList<>();
    
    // Static initializer
    static {
        System.out.println("Class loaded");
    }
    
    // Instance initializer
    {
        items = new ArrayList<>();
    }
    
    // Constructor
    public Example(String name) {
        this.name = name;
    }
    
    // Compact constructor for records (if this was a record)
    // public Example { ... }
    
    // Regular method
    public String getName() {
        return name;
    }
    
    // Generic method
    public <U> List<U> convert(Function<T, U> converter) {
        return items.stream()
            .map(converter)
            .collect(Collectors.toList());
    }
    
    // Static method
    public static void printInfo() {
        System.out.println("Version: " + VERSION);
    }
    
    // Method with lambda
    public void processItems() {
        items.forEach(item -> {
            System.out.println(item);
        });
    }
    
    // Method with method reference
    public void printAll() {
        items.forEach(System.out::println);
    }
    
    // Inner class
    public static class Builder {
        private String name;
        
        public Builder withName(String name) {
            this.name = name;
            return this;
        }
        
        public Example build() {
            return new Example(name);
        }
    }
    
    // Inner interface
    public interface Processor<T> {
        void process(T item);
        
        default void preprocessor(T item) {
            // Default method
        }
    }
    
    // Enum
    public enum Status {
        ACTIVE("Active"),
        INACTIVE("Inactive");
        
        private final String displayName;
        
        Status(String displayName) {
            this.displayName = displayName;
        }
        
        public String getDisplayName() {
            return displayName;
        }
    }
    
    // Record (Java 14+)
    public record Person(String name, int age) {
        // Compact constructor
        public Person {
            if (age < 0) {
                throw new IllegalArgumentException("Age cannot be negative");
            }
        }
        
        // Additional method
        public boolean isAdult() {
            return age >= 18;
        }
    }
    
    // Sealed class (Java 15+)
    public sealed interface Shape 
        permits Circle, Rectangle, Triangle {
        double area();
    }
    
    // Final class implementing sealed interface
    public final class Circle implements Shape {
        private final double radius;
        
        public Circle(double radius) {
            this.radius = radius;
        }
        
        @Override
        public double area() {
            return Math.PI * radius * radius;
        }
    }
    
    // Annotation type
    @Retention(RetentionPolicy.RUNTIME)
    @Target(ElementType.METHOD)
    public @interface Benchmark {
        String value() default "";
    }
}
`

	if err := os.WriteFile(testFile, []byte(testCode), 0644); err != nil {
		t.Fatal(err)
	}
	
	// Test different detail levels
	tests := []struct {
		detail   string
		contains []string
		notContains []string
	}{
		{
			detail:   "minimal",
			contains: []string{
				"class: Example",
				"field: VERSION",
				"field: name",
				"constructor: Example",
				"method: getName",
				"class: Builder",
				"interface: Processor",
				"enum: Status",
				"record: Person",
				"interface: Shape",
				"class: Circle",
				"annotation: Benchmark",
				"static_init: static block",
			},
			notContains: []string{
				"lambda_expression",
				"method_reference",
				"variable_declarator",
			},
		},
		{
			detail:   "standard",
			contains: []string{
				"public static final String VERSION",
				"public Example(String name)",
				"public String getName()",
				"public record Person(String name, int age)",
				"public sealed interface Shape",
			},
		},
		{
			detail:   "full",
			contains: []string{
				"```",
				"public static final String VERSION = \"1.0.0\"",
				"static {",
				"items.forEach(item ->",
			},
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.detail, func(t *testing.T) {
			result, err := extractSymbols(testFile, tt.detail)
			if err != nil {
				t.Fatalf("extractSymbols error = %v", err)
			}
			
			// Check if we have expected content
			if result == "No symbols found" {
				t.Errorf("No symbols were extracted")
			}
			
			// For debugging
			t.Logf("Result for %s:\n%s", tt.detail, result)
			
			for _, expected := range tt.contains {
				if !strings.Contains(result, expected) {
					t.Errorf("Result does not contain expected string %q", expected)
				}
			}
			
			for _, notExpected := range tt.notContains {
				if strings.Contains(result, notExpected) {
					t.Errorf("Result should not contain string %q", notExpected)
				}
			}
		})
	}
}
