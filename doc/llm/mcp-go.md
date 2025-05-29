=== README.md ===

<!-- omit in toc -->
<div align="center">
<img src="./logo.png" alt="MCP Go Logo">

[![Build](https://github.com/mark3labs/mcp-go/actions/workflows/ci.yml/badge.svg?branch=main)](https://github.com/mark3labs/mcp-go/actions/workflows/ci.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/mark3labs/mcp-go?cache)](https://goreportcard.com/report/github.com/mark3labs/mcp-go)
[![GoDoc](https://pkg.go.dev/badge/github.com/mark3labs/mcp-go.svg)](https://pkg.go.dev/github.com/mark3labs/mcp-go)

<strong>A Go implementation of the Model Context Protocol (MCP), enabling seamless integration between LLM applications and external data sources and tools.</strong>

<br>

[![Tutorial](http://img.youtube.com/vi/qoaeYMrXJH0/0.jpg)](http://www.youtube.com/watch?v=qoaeYMrXJH0 "Tutorial")

<br>

Discuss the SDK on [Discord](https://discord.gg/RqSS2NQVsY)

</div>


```go
package main

import (
    "context"
    "errors"
    "fmt"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    // Create a new MCP server
    s := server.NewMCPServer(
        "Demo 🚀",
        "1.0.0",
        server.WithToolCapabilities(false),
    )

    // Add tool
    tool := mcp.NewTool("hello_world",
        mcp.WithDescription("Say hello to someone"),
        mcp.WithString("name",
            mcp.Required(),
            mcp.Description("Name of the person to greet"),
        ),
    )

    // Add tool handler
    s.AddTool(tool, helloHandler)

    // Start the stdio server
    if err := server.ServeStdio(s); err != nil {
        fmt.Printf("Server error: %v\n", err)
    }
}

func helloHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    name, err := request.RequireString("name")
    if err != nil {
        return mcp.NewToolResultError(err.Error()), nil
    }

    return mcp.NewToolResultText(fmt.Sprintf("Hello, %s!", name)), nil
}
```

That's it!

MCP Go handles all the complex protocol details and server management, so you can focus on building great tools. It aims to be high-level and easy to use.

### Key features:
* **Fast**: High-level interface means less code and faster development
* **Simple**: Build MCP servers with minimal boilerplate
* **Complete***: MCP Go aims to provide a full implementation of the core MCP specification

(\*emphasis on *aims*)

🚨 🚧 🏗️ *MCP Go is under active development, as is the MCP specification itself. Core features are working but some advanced capabilities are still in progress.*


<!-- omit in toc -->
## Table of Contents

- [Installation](#installation)
- [Quickstart](#quickstart)
- [What is MCP?](#what-is-mcp)
- [Core Concepts](#core-concepts)
    - [Server](#server)
    - [Resources](#resources)
    - [Tools](#tools)
    - [Prompts](#prompts)
- [Examples](#examples)
- [Extras](#extras)
    - [Transports](#transports)
    - [Session Management](#session-management)
        - [Basic Session Handling](#basic-session-handling)
        - [Per-Session Tools](#per-session-tools)
        - [Tool Filtering](#tool-filtering)
        - [Working with Context](#working-with-context)
    - [Request Hooks](#request-hooks)
    - [Tool Handler Middleware](#tool-handler-middleware)
    - [Regenerating Server Code](#regenerating-server-code)

## Installation

```bash
go get github.com/mark3labs/mcp-go
```

## Quickstart

Let's create a simple MCP server that exposes a calculator tool and some data:

```go
package main

import (
    "context"
    "fmt"

    "github.com/mark3labs/mcp-go/mcp"
    "github.com/mark3labs/mcp-go/server"
)

func main() {
    // Create a new MCP server
    s := server.NewMCPServer(
        "Calculator Demo",
        "1.0.0",
        server.WithToolCapabilities(false),
        server.WithRecovery(),
    )

    // Add a calculator tool
    calculatorTool := mcp.NewTool("calculate",
        mcp.WithDescription("Perform basic arithmetic operations"),
        mcp.WithString("operation",
            mcp.Required(),
            mcp.Description("The operation to perform (add, subtract, multiply, divide)"),
            mcp.Enum("add", "subtract", "multiply", "divide"),
        ),
        mcp.WithNumber("x",
            mcp.Required(),
            mcp.Description("First number"),
        ),
        mcp.WithNumber("y",
            mcp.Required(),
            mcp.Description("Second number"),
        ),
    )

    // Add the calculator handler
    s.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // Using helper functions for type-safe argument access
        op, err := request.RequireString("operation")
        if err != nil {
            return mcp.NewToolResultError(err.Error()), nil
        }
        
        x, err := request.RequireFloat("x")
        if err != nil {
            return mcp.NewToolResultError(err.Error()), nil
        }
        
        y, err := request.RequireFloat("y")
        if err != nil {
            return mcp.NewToolResultError(err.Error()), nil
        }

        var result float64
        switch op {
        case "add":
            result = x + y
        case "subtract":
            result = x - y
        case "multiply":
            result = x * y
        case "divide":
            if y == 0 {
                return mcp.NewToolResultError("cannot divide by zero"), nil
            }
            result = x / y
        }

        return mcp.NewToolResultText(fmt.Sprintf("%.2f", result)), nil
    })

    // Start the server
    if err := server.ServeStdio(s); err != nil {
        fmt.Printf("Server error: %v\n", err)
    }
}
```

## What is MCP?

The [Model Context Protocol (MCP)](https://modelcontextprotocol.io) lets you build servers that expose data and functionality to LLM applications in a secure, standardized way. Think of it like a web API, but specifically designed for LLM interactions. MCP servers can:

- Expose data through **Resources** (think of these sort of like GET endpoints; they are used to load information into the LLM's context)
- Provide functionality through **Tools** (sort of like POST endpoints; they are used to execute code or otherwise produce a side effect)
- Define interaction patterns through **Prompts** (reusable templates for LLM interactions)
- And more!


## Core Concepts


### Server

<details>
<summary>Show Server Examples</summary>

The server is your core interface to the MCP protocol. It handles connection management, protocol compliance, and message routing:

```go
// Create a basic server
s := server.NewMCPServer(
    "My Server",  // Server name
    "1.0.0",     // Version
)

// Start the server using stdio
if err := server.ServeStdio(s); err != nil {
    log.Fatalf("Server error: %v", err)
}
```

</details>

### Resources

<details>
<summary>Show Resource Examples</summary>
Resources are how you expose data to LLMs. They can be anything - files, API responses, database queries, system information, etc. Resources can be:

- Static (fixed URI)
- Dynamic (using URI templates)

Here's a simple example of a static resource:

```go
// Static resource example - exposing a README file
resource := mcp.NewResource(
    "docs://readme",
    "Project README",
    mcp.WithResourceDescription("The project's README file"), 
    mcp.WithMIMEType("text/markdown"),
)

// Add resource with its handler
s.AddResource(resource, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
    content, err := os.ReadFile("README.md")
    if err != nil {
        return nil, err
    }
    
    return []mcp.ResourceContents{
        mcp.TextResourceContents{
            URI:      "docs://readme",
            MIMEType: "text/markdown",
            Text:     string(content),
        },
    }, nil
})
```

And here's an example of a dynamic resource using a template:

```go
// Dynamic resource example - user profiles by ID
template := mcp.NewResourceTemplate(
    "users://{id}/profile",
    "User Profile",
    mcp.WithTemplateDescription("Returns user profile information"),
    mcp.WithTemplateMIMEType("application/json"),
)

// Add template with its handler
s.AddResourceTemplate(template, func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error) {
    // Extract ID from the URI using regex matching
    // The server automatically matches URIs to templates
    userID := extractIDFromURI(request.Params.URI)
    
    profile, err := getUserProfile(userID)  // Your DB/API call here
    if err != nil {
        return nil, err
    }
    
    return []mcp.ResourceContents{
        mcp.TextResourceContents{
            URI:      request.Params.URI,
            MIMEType: "application/json",
            Text:     profile,
        },
    }, nil
})
```

The examples are simple but demonstrate the core concepts. Resources can be much more sophisticated - serving multiple contents, integrating with databases or external APIs, etc.
</details>

### Tools

<details>
<summary>Show Tool Examples</summary>

Tools let LLMs take actions through your server. Unlike resources, tools are expected to perform computation and have side effects. They're similar to POST endpoints in a REST API.

Simple calculation example:
```go
calculatorTool := mcp.NewTool("calculate",
    mcp.WithDescription("Perform basic arithmetic calculations"),
    mcp.WithString("operation",
        mcp.Required(),
        mcp.Description("The arithmetic operation to perform"),
        mcp.Enum("add", "subtract", "multiply", "divide"),
    ),
    mcp.WithNumber("x",
        mcp.Required(),
        mcp.Description("First number"),
    ),
    mcp.WithNumber("y",
        mcp.Required(),
        mcp.Description("Second number"),
    ),
)

s.AddTool(calculatorTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    args := request.GetArguments()
    op := args["operation"].(string)
    x := args["x"].(float64)
    y := args["y"].(float64)

    var result float64
    switch op {
    case "add":
        result = x + y
    case "subtract":
        result = x - y
    case "multiply":
        result = x * y
    case "divide":
        if y == 0 {
            return mcp.NewToolResultError("cannot divide by zero"), nil
        }
        result = x / y
    }
    
    return mcp.FormatNumberResult(result), nil
})
```

HTTP request example:
```go
httpTool := mcp.NewTool("http_request",
    mcp.WithDescription("Make HTTP requests to external APIs"),
    mcp.WithString("method",
        mcp.Required(),
        mcp.Description("HTTP method to use"),
        mcp.Enum("GET", "POST", "PUT", "DELETE"),
    ),
    mcp.WithString("url",
        mcp.Required(),
        mcp.Description("URL to send the request to"),
        mcp.Pattern("^https?://.*"),
    ),
    mcp.WithString("body",
        mcp.Description("Request body (for POST/PUT)"),
    ),
)

s.AddTool(httpTool, func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    args := request.GetArguments()
    method := args["method"].(string)
    url := args["url"].(string)
    body := ""
    if b, ok := args["body"].(string); ok {
        body = b
    }

    // Create and send request
    var req *http.Request
    var err error
    if body != "" {
        req, err = http.NewRequest(method, url, strings.NewReader(body))
    } else {
        req, err = http.NewRequest(method, url, nil)
    }
    if err != nil {
        return mcp.NewToolResultErrorFromErr("unable to create request", err), nil
    }

    client := &http.Client{}
    resp, err := client.Do(req)
    if err != nil {
        return mcp.NewToolResultErrorFromErr("unable to execute request", err), nil
    }
    defer resp.Body.Close()

    // Return response
    respBody, err := io.ReadAll(resp.Body)
    if err != nil {
        return mcp.NewToolResultErrorFromErr("unable to read request response", err), nil
    }

    return mcp.NewToolResultText(fmt.Sprintf("Status: %d\nBody: %s", resp.StatusCode, string(respBody))), nil
})
```

Tools can be used for any kind of computation or side effect:
- Database queries
- File operations
- External API calls
- Calculations
- System operations

Each tool should:
- Have a clear description
- Validate inputs
- Handle errors gracefully
- Return structured responses
- Use appropriate result types

</details>

### Prompts

<details>
<summary>Show Prompt Examples</summary>

Prompts are reusable templates that help LLMs interact with your server effectively. They're like "best practices" encoded into your server. Here are some examples:

```go
// Simple greeting prompt
s.AddPrompt(mcp.NewPrompt("greeting",
    mcp.WithPromptDescription("A friendly greeting prompt"),
    mcp.WithArgument("name",
        mcp.ArgumentDescription("Name of the person to greet"),
    ),
), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    name := request.Params.Arguments["name"]
    if name == "" {
        name = "friend"
    }
    
    return mcp.NewGetPromptResult(
        "A friendly greeting",
        []mcp.PromptMessage{
            mcp.NewPromptMessage(
                mcp.RoleAssistant,
                mcp.NewTextContent(fmt.Sprintf("Hello, %s! How can I help you today?", name)),
            ),
        },
    ), nil
})

// Code review prompt with embedded resource
s.AddPrompt(mcp.NewPrompt("code_review",
    mcp.WithPromptDescription("Code review assistance"),
    mcp.WithArgument("pr_number",
        mcp.ArgumentDescription("Pull request number to review"),
        mcp.RequiredArgument(),
    ),
), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    prNumber := request.Params.Arguments["pr_number"]
    if prNumber == "" {
        return nil, fmt.Errorf("pr_number is required")
    }
    
    return mcp.NewGetPromptResult(
        "Code review assistance",
        []mcp.PromptMessage{
            mcp.NewPromptMessage(
                mcp.RoleUser,
                mcp.NewTextContent("Review the changes and provide constructive feedback."),
            ),
            mcp.NewPromptMessage(
                mcp.RoleAssistant,
                mcp.NewEmbeddedResource(mcp.ResourceContents{
                    URI: fmt.Sprintf("git://pulls/%s/diff", prNumber),
                    MIMEType: "text/x-diff",
                }),
            ),
        },
    ), nil
})

// Database query builder prompt
s.AddPrompt(mcp.NewPrompt("query_builder",
    mcp.WithPromptDescription("SQL query builder assistance"),
    mcp.WithArgument("table",
        mcp.ArgumentDescription("Name of the table to query"),
        mcp.RequiredArgument(),
    ),
), func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error) {
    tableName := request.Params.Arguments["table"]
    if tableName == "" {
        return nil, fmt.Errorf("table name is required")
    }
    
    return mcp.NewGetPromptResult(
        "SQL query builder assistance",
        []mcp.PromptMessage{
            mcp.NewPromptMessage(
                mcp.RoleUser,
                mcp.NewTextContent("Help construct efficient and safe queries for the provided schema."),
            ),
            mcp.NewPromptMessage(
                mcp.RoleUser,
                mcp.NewEmbeddedResource(mcp.ResourceContents{
                    URI: fmt.Sprintf("db://schema/%s", tableName),
                    MIMEType: "application/json",
                }),
            ),
        },
    ), nil
})
```

Prompts can include:
- System instructions
- Required arguments
- Embedded resources
- Multiple messages
- Different content types (text, images, etc.)
- Custom URI schemes

</details>

## Examples

For examples, see the [`examples/`](examples/) directory.

## Extras

### Transports

MCP-Go supports stdio, SSE and streamable-HTTP transport layers.

### Session Management

MCP-Go provides a robust session management system that allows you to:
- Maintain separate state for each connected client
- Register and track client sessions
- Send notifications to specific clients
- Provide per-session tool customization

<details>
<summary>Show Session Management Examples</summary>

#### Basic Session Handling

```go
// Create a server with session capabilities
s := server.NewMCPServer(
    "Session Demo",
    "1.0.0",
    server.WithToolCapabilities(true),
)

// Implement your own ClientSession
type MySession struct {
    id           string
    notifChannel chan mcp.JSONRPCNotification
    isInitialized bool
    // Add custom fields for your application
}

// Implement the ClientSession interface
func (s *MySession) SessionID() string {
    return s.id
}

func (s *MySession) NotificationChannel() chan<- mcp.JSONRPCNotification {
    return s.notifChannel
}

func (s *MySession) Initialize() {
    s.isInitialized = true
}

func (s *MySession) Initialized() bool {
    return s.isInitialized
}

// Register a session
session := &MySession{
    id:           "user-123",
    notifChannel: make(chan mcp.JSONRPCNotification, 10),
}
if err := s.RegisterSession(context.Background(), session); err != nil {
    log.Printf("Failed to register session: %v", err)
}

// Send notification to a specific client
err := s.SendNotificationToSpecificClient(
    session.SessionID(),
    "notification/update",
    map[string]any{"message": "New data available!"},
)
if err != nil {
    log.Printf("Failed to send notification: %v", err)
}

// Unregister session when done
s.UnregisterSession(context.Background(), session.SessionID())
```

#### Per-Session Tools

For more advanced use cases, you can implement the `SessionWithTools` interface to support per-session tool customization:

```go
// Implement SessionWithTools interface for per-session tools
type MyAdvancedSession struct {
    MySession  // Embed the basic session
    sessionTools map[string]server.ServerTool
}

// Implement additional methods for SessionWithTools
func (s *MyAdvancedSession) GetSessionTools() map[string]server.ServerTool {
    return s.sessionTools
}

func (s *MyAdvancedSession) SetSessionTools(tools map[string]server.ServerTool) {
    s.sessionTools = tools
}

// Create and register a session with tools support
advSession := &MyAdvancedSession{
    MySession: MySession{
        id:           "user-456",
        notifChannel: make(chan mcp.JSONRPCNotification, 10),
    },
    sessionTools: make(map[string]server.ServerTool),
}
if err := s.RegisterSession(context.Background(), advSession); err != nil {
    log.Printf("Failed to register session: %v", err)
}

// Add session-specific tools
userSpecificTool := mcp.NewTool(
    "user_data",
    mcp.WithDescription("Access user-specific data"),
)
// You can use AddSessionTool (similar to AddTool)
err := s.AddSessionTool(
    advSession.SessionID(),
    userSpecificTool,
    func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
        // This handler is only available to this specific session
        return mcp.NewToolResultText("User-specific data for " + advSession.SessionID()), nil
    },
)
if err != nil {
    log.Printf("Failed to add session tool: %v", err)
}

// Or use AddSessionTools directly with ServerTool
/*
err := s.AddSessionTools(
    advSession.SessionID(),
    server.ServerTool{
        Tool: userSpecificTool,
        Handler: func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
            // This handler is only available to this specific session
            return mcp.NewToolResultText("User-specific data for " + advSession.SessionID()), nil
        },
    },
)
if err != nil {
    log.Printf("Failed to add session tool: %v", err)
}
*/

// Delete session-specific tools when no longer needed
err = s.DeleteSessionTools(advSession.SessionID(), "user_data")
if err != nil {
    log.Printf("Failed to delete session tool: %v", err)
}
```

#### Tool Filtering

You can also apply filters to control which tools are available to certain sessions:

```go
// Add a tool filter that only shows tools with certain prefixes
s := server.NewMCPServer(
    "Tool Filtering Demo",
    "1.0.0",
    server.WithToolCapabilities(true),
    server.WithToolFilter(func(ctx context.Context, tools []mcp.Tool) []mcp.Tool {
        // Get session from context
        session := server.ClientSessionFromContext(ctx)
        if session == nil {
            return tools // Return all tools if no session
        }
        
        // Example: filter tools based on session ID prefix
        if strings.HasPrefix(session.SessionID(), "admin-") {
            // Admin users get all tools
            return tools
        } else {
            // Regular users only get tools with "public-" prefix
            var filteredTools []mcp.Tool
            for _, tool := range tools {
                if strings.HasPrefix(tool.Name, "public-") {
                    filteredTools = append(filteredTools, tool)
                }
            }
            return filteredTools
        }
    }),
)
```

#### Working with Context

The session context is automatically passed to tool and resource handlers:

```go
s.AddTool(mcp.NewTool("session_aware"), func(ctx context.Context, req mcp.CallToolRequest) (*mcp.CallToolResult, error) {
    // Get the current session from context
    session := server.ClientSessionFromContext(ctx)
    if session == nil {
        return mcp.NewToolResultError("No active session"), nil
    }
    
    return mcp.NewToolResultText("Hello, session " + session.SessionID()), nil
})

// When using handlers in HTTP/SSE servers, you need to pass the context with the session
httpHandler := func(w http.ResponseWriter, r *http.Request) {
    // Get session from somewhere (like a cookie or header)
    session := getSessionFromRequest(r)
    
    // Add session to context
    ctx := s.WithContext(r.Context(), session)
    
    // Use this context when handling requests
    // ...
}
```

</details>

### Request Hooks

Hook into the request lifecycle by creating a `Hooks` object with your
selection among the possible callbacks.  This enables telemetry across all
functionality, and observability of various facts, for example the ability
to count improperly-formatted requests, or to log the agent identity during
initialization.

Add the `Hooks` to the server at the time of creation using the
`server.WithHooks` option.

### Tool Handler Middleware

Add middleware to tool call handlers using the `server.WithToolHandlerMiddleware` option. Middlewares can be registered on server creation and are applied on every tool call.

A recovery middleware option is available to recover from panics in a tool call and can be added to the server with the `server.WithRecovery` option.

### Regenerating Server Code

Server hooks and request handlers are generated. Regenerate them by running:

```bash
go generate ./...
```

You need `go` installed and the `goimports` tool available. The generator runs
`goimports` automatically to format and fix imports.

=== examples/typed_tools/main.go ===

```go
package main

import (
	"context"
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

// Define a struct for our typed arguments
type GreetingArgs struct {
	Name      string   `json:"name"`
	Age       int      `json:"age"`
	IsVIP     bool     `json:"is_vip"`
	Languages []string `json:"languages"`
	Metadata  struct {
		Location string `json:"location"`
		Timezone string `json:"timezone"`
	} `json:"metadata"`
}

func main() {
	// Create a new MCP server
	s := server.NewMCPServer(
		"Typed Tools Demo 🚀",
		"1.0.0",
		server.WithToolCapabilities(false),
	)

	// Add tool with complex schema
	tool := mcp.NewTool("greeting",
		mcp.WithDescription("Generate a personalized greeting"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("Name of the person to greet"),
		),
		mcp.WithNumber("age",
			mcp.Description("Age of the person"),
			mcp.Min(0),
			mcp.Max(150),
		),
		mcp.WithBoolean("is_vip",
			mcp.Description("Whether the person is a VIP"),
			mcp.DefaultBool(false),
		),
		mcp.WithArray("languages",
			mcp.Description("Languages the person speaks"),
			mcp.Items(map[string]any{"type": "string"}),
		),
		mcp.WithObject("metadata",
			mcp.Description("Additional information about the person"),
			mcp.Properties(map[string]any{
				"location": map[string]any{
					"type":        "string",
					"description": "Current location",
				},
				"timezone": map[string]any{
					"type":        "string",
					"description": "Timezone",
				},
			}),
		),
	)

	// Add tool handler using the typed handler
	s.AddTool(tool, mcp.NewTypedToolHandler(typedGreetingHandler))

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

// Our typed handler function that receives strongly-typed arguments
func typedGreetingHandler(ctx context.Context, request mcp.CallToolRequest, args GreetingArgs) (*mcp.CallToolResult, error) {
	if args.Name == "" {
		return mcp.NewToolResultError("name is required"), nil
	}

	// Build a personalized greeting based on the complex arguments
	greeting := fmt.Sprintf("Hello, %s!", args.Name)

	if args.Age > 0 {
		greeting += fmt.Sprintf(" You are %d years old.", args.Age)
	}

	if args.IsVIP {
		greeting += " Welcome back, valued VIP customer!"
	}

	if len(args.Languages) > 0 {
		greeting += fmt.Sprintf(" You speak %d languages: %v.", len(args.Languages), args.Languages)
	}

	if args.Metadata.Location != "" {
		greeting += fmt.Sprintf(" I see you're from %s.", args.Metadata.Location)

		if args.Metadata.Timezone != "" {
			greeting += fmt.Sprintf(" Your timezone is %s.", args.Metadata.Timezone)
		}
	}

	return mcp.NewToolResultText(greeting), nil
}
```

=== godoc ===

Documentation ¶
Overview ¶
Code generated by `go generate`. DO NOT EDIT. source: server/internal/gen/hooks.go.tmpl

Code generated by `go generate`. DO NOT EDIT. source: server/internal/gen/request_handler.go.tmpl

Package server provides MCP (Model Context Protocol) server implementations.

Index ¶
Variables
func NewTestServer(server *MCPServer, opts ...SSEOption) *httptest.Server
func NewTestStreamableHTTPServer(server *MCPServer, opts ...StreamableHTTPOption) *httptest.Server
func ServeStdio(server *MCPServer, opts ...StdioOption) error
type BeforeAnyHookFunc
type ClientSession
func ClientSessionFromContext(ctx context.Context) ClientSession
type DynamicBasePathFunc
type ErrDynamicPathConfig
func (e *ErrDynamicPathConfig) Error() string
type HTTPContextFunc
type Hooks
func (c *Hooks) AddAfterCallTool(hook OnAfterCallToolFunc)
func (c *Hooks) AddAfterGetPrompt(hook OnAfterGetPromptFunc)
func (c *Hooks) AddAfterInitialize(hook OnAfterInitializeFunc)
func (c *Hooks) AddAfterListPrompts(hook OnAfterListPromptsFunc)
func (c *Hooks) AddAfterListResourceTemplates(hook OnAfterListResourceTemplatesFunc)
func (c *Hooks) AddAfterListResources(hook OnAfterListResourcesFunc)
func (c *Hooks) AddAfterListTools(hook OnAfterListToolsFunc)
func (c *Hooks) AddAfterPing(hook OnAfterPingFunc)
func (c *Hooks) AddAfterReadResource(hook OnAfterReadResourceFunc)
func (c *Hooks) AddAfterSetLevel(hook OnAfterSetLevelFunc)
func (c *Hooks) AddBeforeAny(hook BeforeAnyHookFunc)
func (c *Hooks) AddBeforeCallTool(hook OnBeforeCallToolFunc)
func (c *Hooks) AddBeforeGetPrompt(hook OnBeforeGetPromptFunc)
func (c *Hooks) AddBeforeInitialize(hook OnBeforeInitializeFunc)
func (c *Hooks) AddBeforeListPrompts(hook OnBeforeListPromptsFunc)
func (c *Hooks) AddBeforeListResourceTemplates(hook OnBeforeListResourceTemplatesFunc)
func (c *Hooks) AddBeforeListResources(hook OnBeforeListResourcesFunc)
func (c *Hooks) AddBeforeListTools(hook OnBeforeListToolsFunc)
func (c *Hooks) AddBeforePing(hook OnBeforePingFunc)
func (c *Hooks) AddBeforeReadResource(hook OnBeforeReadResourceFunc)
func (c *Hooks) AddBeforeSetLevel(hook OnBeforeSetLevelFunc)
func (c *Hooks) AddOnError(hook OnErrorHookFunc)
func (c *Hooks) AddOnRegisterSession(hook OnRegisterSessionHookFunc)
func (c *Hooks) AddOnRequestInitialization(hook OnRequestInitializationFunc)
func (c *Hooks) AddOnSuccess(hook OnSuccessHookFunc)
func (c *Hooks) AddOnUnregisterSession(hook OnUnregisterSessionHookFunc)
func (c *Hooks) RegisterSession(ctx context.Context, session ClientSession)
func (c *Hooks) UnregisterSession(ctx context.Context, session ClientSession)
type InsecureStatefulSessionIdManager
func (s *InsecureStatefulSessionIdManager) Generate() string
func (s *InsecureStatefulSessionIdManager) Terminate(sessionID string) (isNotAllowed bool, err error)
func (s *InsecureStatefulSessionIdManager) Validate(sessionID string) (isTerminated bool, err error)
type MCPServer
func NewMCPServer(name, version string, opts ...ServerOption) *MCPServer
func ServerFromContext(ctx context.Context) *MCPServer
func (s *MCPServer) AddNotificationHandler(method string, handler NotificationHandlerFunc)
func (s *MCPServer) AddPrompt(prompt mcp.Prompt, handler PromptHandlerFunc)
func (s *MCPServer) AddResource(resource mcp.Resource, handler ResourceHandlerFunc)
func (s *MCPServer) AddResourceTemplate(template mcp.ResourceTemplate, handler ResourceTemplateHandlerFunc)
func (s *MCPServer) AddSessionTool(sessionID string, tool mcp.Tool, handler ToolHandlerFunc) error
func (s *MCPServer) AddSessionTools(sessionID string, tools ...ServerTool) error
func (s *MCPServer) AddTool(tool mcp.Tool, handler ToolHandlerFunc)
func (s *MCPServer) AddTools(tools ...ServerTool)
func (s *MCPServer) DeletePrompts(names ...string)
func (s *MCPServer) DeleteSessionTools(sessionID string, names ...string) error
func (s *MCPServer) DeleteTools(names ...string)
func (s *MCPServer) HandleMessage(ctx context.Context, message json.RawMessage) mcp.JSONRPCMessage
func (s *MCPServer) RegisterSession(ctx context.Context, session ClientSession) error
func (s *MCPServer) RemoveResource(uri string)
func (s *MCPServer) SendNotificationToAllClients(method string, params map[string]any)
func (s *MCPServer) SendNotificationToClient(ctx context.Context, method string, params map[string]any) error
func (s *MCPServer) SendNotificationToSpecificClient(sessionID string, method string, params map[string]any) error
func (s *MCPServer) SetTools(tools ...ServerTool)
func (s *MCPServer) UnregisterSession(ctx context.Context, sessionID string)
func (s *MCPServer) WithContext(ctx context.Context, session ClientSession) context.Context
type NotificationHandlerFunc
type OnAfterCallToolFunc
type OnAfterGetPromptFunc
type OnAfterInitializeFunc
type OnAfterListPromptsFunc
type OnAfterListResourceTemplatesFunc
type OnAfterListResourcesFunc
type OnAfterListToolsFunc
type OnAfterPingFunc
type OnAfterReadResourceFunc
type OnAfterSetLevelFunc
type OnBeforeCallToolFunc
type OnBeforeGetPromptFunc
type OnBeforeInitializeFunc
type OnBeforeListPromptsFunc
type OnBeforeListResourceTemplatesFunc
type OnBeforeListResourcesFunc
type OnBeforeListToolsFunc
type OnBeforePingFunc
type OnBeforeReadResourceFunc
type OnBeforeSetLevelFunc
type OnErrorHookFunc
type OnRegisterSessionHookFunc
type OnRequestInitializationFunc
type OnSuccessHookFunc
type OnUnregisterSessionHookFunc
type PromptHandlerFunc
type ResourceHandlerFunc
type ResourceTemplateHandlerFunc
type SSEContextFunc
type SSEOption
func WithAppendQueryToMessageEndpoint() SSEOption
func WithBasePath(basePath string) SSEOptiondeprecated
func WithBaseURL(baseURL string) SSEOption
func WithDynamicBasePath(fn DynamicBasePathFunc) SSEOption
func WithHTTPServer(srv *http.Server) SSEOption
func WithKeepAlive(keepAlive bool) SSEOption
func WithKeepAliveInterval(keepAliveInterval time.Duration) SSEOption
func WithMessageEndpoint(endpoint string) SSEOption
func WithSSEContextFunc(fn SSEContextFunc) SSEOption
func WithSSEEndpoint(endpoint string) SSEOption
func WithStaticBasePath(basePath string) SSEOption
func WithUseFullURLForMessageEndpoint(useFullURLForMessageEndpoint bool) SSEOption
type SSEServer
func NewSSEServer(server *MCPServer, opts ...SSEOption) *SSEServer
func (s *SSEServer) CompleteMessageEndpoint() (string, error)
func (s *SSEServer) CompleteMessagePath() string
func (s *SSEServer) CompleteSseEndpoint() (string, error)
func (s *SSEServer) CompleteSsePath() string
func (s *SSEServer) GetMessageEndpointForClient(r *http.Request, sessionID string) string
func (s *SSEServer) GetUrlPath(input string) (string, error)
func (s *SSEServer) MessageHandler() http.Handler
func (s *SSEServer) SSEHandler() http.Handler
func (s *SSEServer) SendEventToSession(sessionID string, event any) error
func (s *SSEServer) ServeHTTP(w http.ResponseWriter, r *http.Request)
func (s *SSEServer) Shutdown(ctx context.Context) error
func (s *SSEServer) Start(addr string) error
type ServerOption
func WithHooks(hooks *Hooks) ServerOption
func WithInstructions(instructions string) ServerOption
func WithLogging() ServerOption
func WithPaginationLimit(limit int) ServerOption
func WithPromptCapabilities(listChanged bool) ServerOption
func WithRecovery() ServerOption
func WithResourceCapabilities(subscribe, listChanged bool) ServerOption
func WithToolCapabilities(listChanged bool) ServerOption
func WithToolFilter(toolFilter ToolFilterFunc) ServerOption
func WithToolHandlerMiddleware(toolHandlerMiddleware ToolHandlerMiddleware) ServerOption
type ServerTool
type SessionIdManager
type SessionWithClientInfo
type SessionWithLogging
type SessionWithTools
type StatelessSessionIdManager
func (s *StatelessSessionIdManager) Generate() string
func (s *StatelessSessionIdManager) Terminate(sessionID string) (isNotAllowed bool, err error)
func (s *StatelessSessionIdManager) Validate(sessionID string) (isTerminated bool, err error)
type StdioContextFunc
type StdioOption
func WithErrorLogger(logger *log.Logger) StdioOption
func WithStdioContextFunc(fn StdioContextFunc) StdioOption
type StdioServer
func NewStdioServer(server *MCPServer) *StdioServer
func (s *StdioServer) Listen(ctx context.Context, stdin io.Reader, stdout io.Writer) error
func (s *StdioServer) SetContextFunc(fn StdioContextFunc)
func (s *StdioServer) SetErrorLogger(logger *log.Logger)
type StreamableHTTPOption
func WithEndpointPath(endpointPath string) StreamableHTTPOption
func WithHTTPContextFunc(fn HTTPContextFunc) StreamableHTTPOption
func WithHeartbeatInterval(interval time.Duration) StreamableHTTPOption
func WithLogger(logger util.Logger) StreamableHTTPOption
func WithSessionIdManager(manager SessionIdManager) StreamableHTTPOption
func WithStateLess(stateLess bool) StreamableHTTPOption
type StreamableHTTPServer
func NewStreamableHTTPServer(server *MCPServer, opts ...StreamableHTTPOption) *StreamableHTTPServer
func (s *StreamableHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request)
func (s *StreamableHTTPServer) Shutdown(ctx context.Context) error
func (s *StreamableHTTPServer) Start(addr string) error
type ToolFilterFunc
type ToolHandlerFunc
type ToolHandlerMiddleware
type UnparsableMessageError
func (e *UnparsableMessageError) Error() string
func (e *UnparsableMessageError) GetMessage() json.RawMessage
func (e *UnparsableMessageError) GetMethod() mcp.MCPMethod
func (e *UnparsableMessageError) Unwrap() error
Constants ¶
This section is empty.

Variables ¶
View Source
var (
// Common server errors
ErrUnsupported      = errors.New("not supported")
ErrResourceNotFound = errors.New("resource not found")
ErrPromptNotFound   = errors.New("prompt not found")
ErrToolNotFound     = errors.New("tool not found")

	// Session-related errors
	ErrSessionNotFound              = errors.New("session not found")
	ErrSessionExists                = errors.New("session already exists")
	ErrSessionNotInitialized        = errors.New("session not properly initialized")
	ErrSessionDoesNotSupportTools   = errors.New("session does not support per-session tools")
	ErrSessionDoesNotSupportLogging = errors.New("session does not support setting logging level")

	// Notification-related errors
	ErrNotificationNotInitialized = errors.New("notification channel not initialized")
	ErrNotificationChannelBlocked = errors.New("notification channel full or blocked")
)
Functions ¶
func NewTestServer ¶
added in v0.2.0
func NewTestServer(server *MCPServer, opts ...SSEOption) *httptest.Server
NewTestServer creates a test server for testing purposes

func NewTestStreamableHTTPServer ¶
added in v0.30.0
func NewTestStreamableHTTPServer(server *MCPServer, opts ...StreamableHTTPOption) *httptest.Server
NewTestStreamableHTTPServer creates a test server for testing purposes

func ServeStdio ¶
func ServeStdio(server *MCPServer, opts ...StdioOption) error
ServeStdio is a convenience function that creates and starts a StdioServer with os.Stdin and os.Stdout. It sets up signal handling for graceful shutdown on SIGTERM and SIGINT. Returns an error if the server encounters any issues during operation.

Types ¶
type BeforeAnyHookFunc ¶
added in v0.16.0
type BeforeAnyHookFunc func(ctx context.Context, id any, method mcp.MCPMethod, message any)
BeforeAnyHookFunc is a function that is called after the request is parsed but before the method is called.

type ClientSession ¶
added in v0.14.0
type ClientSession interface {
// Initialize marks session as fully initialized and ready for notifications
Initialize()
// Initialized returns if session is ready to accept notifications
Initialized() bool
// NotificationChannel provides a channel suitable for sending notifications to client.
NotificationChannel() chan<- mcp.JSONRPCNotification
// SessionID is a unique identifier used to track user session.
SessionID() string
}
ClientSession represents an active session that can be used by MCPServer to interact with client.

func ClientSessionFromContext ¶
added in v0.14.0
func ClientSessionFromContext(ctx context.Context) ClientSession
ClientSessionFromContext retrieves current client notification context from context.

type DynamicBasePathFunc ¶
added in v0.25.0
type DynamicBasePathFunc func(r *http.Request, sessionID string) string
DynamicBasePathFunc allows the user to provide a function to generate the base path for a given request and sessionID. This is useful for cases where the base path is not known at the time of SSE server creation, such as when using a reverse proxy or when the base path is dynamically generated. The function should return the base path (e.g., "/mcp/tenant123").

type ErrDynamicPathConfig ¶
added in v0.25.0
type ErrDynamicPathConfig struct {
Method string
}
ErrDynamicPathConfig is returned when attempting to use static path methods with dynamic path configuration

func (*ErrDynamicPathConfig) Error ¶
added in v0.25.0
func (e *ErrDynamicPathConfig) Error() string
type HTTPContextFunc ¶
added in v0.27.0
type HTTPContextFunc func(ctx context.Context, r *http.Request) context.Context
HTTPContextFunc is a function that takes an existing context and the current request and returns a potentially modified context based on the request content. This can be used to inject context values from headers, for example.

type Hooks ¶
added in v0.16.0
type Hooks struct {
OnRegisterSession             []OnRegisterSessionHookFunc
OnUnregisterSession           []OnUnregisterSessionHookFunc
OnBeforeAny                   []BeforeAnyHookFunc
OnSuccess                     []OnSuccessHookFunc
OnError                       []OnErrorHookFunc
OnRequestInitialization       []OnRequestInitializationFunc
OnBeforeInitialize            []OnBeforeInitializeFunc
OnAfterInitialize             []OnAfterInitializeFunc
OnBeforePing                  []OnBeforePingFunc
OnAfterPing                   []OnAfterPingFunc
OnBeforeSetLevel              []OnBeforeSetLevelFunc
OnAfterSetLevel               []OnAfterSetLevelFunc
OnBeforeListResources         []OnBeforeListResourcesFunc
OnAfterListResources          []OnAfterListResourcesFunc
OnBeforeListResourceTemplates []OnBeforeListResourceTemplatesFunc
OnAfterListResourceTemplates  []OnAfterListResourceTemplatesFunc
OnBeforeReadResource          []OnBeforeReadResourceFunc
OnAfterReadResource           []OnAfterReadResourceFunc
OnBeforeListPrompts           []OnBeforeListPromptsFunc
OnAfterListPrompts            []OnAfterListPromptsFunc
OnBeforeGetPrompt             []OnBeforeGetPromptFunc
OnAfterGetPrompt              []OnAfterGetPromptFunc
OnBeforeListTools             []OnBeforeListToolsFunc
OnAfterListTools              []OnAfterListToolsFunc
OnBeforeCallTool              []OnBeforeCallToolFunc
OnAfterCallTool               []OnAfterCallToolFunc
}
func (*Hooks) AddAfterCallTool ¶
added in v0.16.0
func (c *Hooks) AddAfterCallTool(hook OnAfterCallToolFunc)
func (*Hooks) AddAfterGetPrompt ¶
added in v0.16.0
func (c *Hooks) AddAfterGetPrompt(hook OnAfterGetPromptFunc)
func (*Hooks) AddAfterInitialize ¶
added in v0.16.0
func (c *Hooks) AddAfterInitialize(hook OnAfterInitializeFunc)
func (*Hooks) AddAfterListPrompts ¶
added in v0.16.0
func (c *Hooks) AddAfterListPrompts(hook OnAfterListPromptsFunc)
func (*Hooks) AddAfterListResourceTemplates ¶
added in v0.16.0
func (c *Hooks) AddAfterListResourceTemplates(hook OnAfterListResourceTemplatesFunc)
func (*Hooks) AddAfterListResources ¶
added in v0.16.0
func (c *Hooks) AddAfterListResources(hook OnAfterListResourcesFunc)
func (*Hooks) AddAfterListTools ¶
added in v0.16.0
func (c *Hooks) AddAfterListTools(hook OnAfterListToolsFunc)
func (*Hooks) AddAfterPing ¶
added in v0.16.0
func (c *Hooks) AddAfterPing(hook OnAfterPingFunc)
func (*Hooks) AddAfterReadResource ¶
added in v0.16.0
func (c *Hooks) AddAfterReadResource(hook OnAfterReadResourceFunc)
func (*Hooks) AddAfterSetLevel ¶
added in v0.28.0
func (c *Hooks) AddAfterSetLevel(hook OnAfterSetLevelFunc)
func (*Hooks) AddBeforeAny ¶
added in v0.16.0
func (c *Hooks) AddBeforeAny(hook BeforeAnyHookFunc)
func (*Hooks) AddBeforeCallTool ¶
added in v0.16.0
func (c *Hooks) AddBeforeCallTool(hook OnBeforeCallToolFunc)
func (*Hooks) AddBeforeGetPrompt ¶
added in v0.16.0
func (c *Hooks) AddBeforeGetPrompt(hook OnBeforeGetPromptFunc)
func (*Hooks) AddBeforeInitialize ¶
added in v0.16.0
func (c *Hooks) AddBeforeInitialize(hook OnBeforeInitializeFunc)
func (*Hooks) AddBeforeListPrompts ¶
added in v0.16.0
func (c *Hooks) AddBeforeListPrompts(hook OnBeforeListPromptsFunc)
func (*Hooks) AddBeforeListResourceTemplates ¶
added in v0.16.0
func (c *Hooks) AddBeforeListResourceTemplates(hook OnBeforeListResourceTemplatesFunc)
func (*Hooks) AddBeforeListResources ¶
added in v0.16.0
func (c *Hooks) AddBeforeListResources(hook OnBeforeListResourcesFunc)
func (*Hooks) AddBeforeListTools ¶
added in v0.16.0
func (c *Hooks) AddBeforeListTools(hook OnBeforeListToolsFunc)
func (*Hooks) AddBeforePing ¶
added in v0.16.0
func (c *Hooks) AddBeforePing(hook OnBeforePingFunc)
func (*Hooks) AddBeforeReadResource ¶
added in v0.16.0
func (c *Hooks) AddBeforeReadResource(hook OnBeforeReadResourceFunc)
func (*Hooks) AddBeforeSetLevel ¶
added in v0.28.0
func (c *Hooks) AddBeforeSetLevel(hook OnBeforeSetLevelFunc)
func (*Hooks) AddOnError ¶
added in v0.16.0
func (c *Hooks) AddOnError(hook OnErrorHookFunc)
AddOnError registers a hook function that will be called when an error occurs. The error parameter contains the actual error object, which can be interrogated using Go's error handling patterns like errors.Is and errors.As.

Example: ``` // Create a channel to receive errors for testing errChan := make(chan error, 1)

// Register hook to capture and inspect errors hooks := &Hooks{}

hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
// For capability-related errors
if errors.Is(err, ErrUnsupported) {
// Handle capability not supported
errChan <- err
return
}

    // For parsing errors
    var parseErr = &UnparsableMessageError{}
    if errors.As(err, &parseErr) {
        // Handle unparsable message errors
        fmt.Printf("Failed to parse %s request: %v\n",
                   parseErr.GetMethod(), parseErr.Unwrap())
        errChan <- parseErr
        return
    }

    // For resource/prompt/tool not found errors
    if errors.Is(err, ErrResourceNotFound) ||
       errors.Is(err, ErrPromptNotFound) ||
       errors.Is(err, ErrToolNotFound) {
        // Handle not found errors
        errChan <- err
        return
    }

    // For other errors
    errChan <- err
})
server := NewMCPServer("test-server", "1.0.0", WithHooks(hooks)) ```

func (*Hooks) AddOnRegisterSession ¶
added in v0.18.0
func (c *Hooks) AddOnRegisterSession(hook OnRegisterSessionHookFunc)
func (*Hooks) AddOnRequestInitialization ¶
added in v0.24.0
func (c *Hooks) AddOnRequestInitialization(hook OnRequestInitializationFunc)
func (*Hooks) AddOnSuccess ¶
added in v0.16.0
func (c *Hooks) AddOnSuccess(hook OnSuccessHookFunc)
func (*Hooks) AddOnUnregisterSession ¶
added in v0.23.0
func (c *Hooks) AddOnUnregisterSession(hook OnUnregisterSessionHookFunc)
func (*Hooks) RegisterSession ¶
added in v0.18.0
func (c *Hooks) RegisterSession(ctx context.Context, session ClientSession)
func (*Hooks) UnregisterSession ¶
added in v0.23.0
func (c *Hooks) UnregisterSession(ctx context.Context, session ClientSession)
type InsecureStatefulSessionIdManager ¶
added in v0.30.0
type InsecureStatefulSessionIdManager struct{}
InsecureStatefulSessionIdManager generate id with uuid It won't validate the id indeed, so it could be fake. For more secure session id, use a more complex generator, like a JWT.

func (*InsecureStatefulSessionIdManager) Generate ¶
added in v0.30.0
func (s *InsecureStatefulSessionIdManager) Generate() string
func (*InsecureStatefulSessionIdManager) Terminate ¶
added in v0.30.0
func (s *InsecureStatefulSessionIdManager) Terminate(sessionID string) (isNotAllowed bool, err error)
func (*InsecureStatefulSessionIdManager) Validate ¶
added in v0.30.0
func (s *InsecureStatefulSessionIdManager) Validate(sessionID string) (isTerminated bool, err error)
type MCPServer ¶
type MCPServer struct {
// contains filtered or unexported fields
}
MCPServer implements a Model Context Protocol server that can handle various types of requests including resources, prompts, and tools.

func NewMCPServer ¶
added in v0.5.0
func NewMCPServer(
name, version string,
opts ...ServerOption,
) *MCPServer
NewMCPServer creates a new MCP server instance with the given name, version and options

func ServerFromContext ¶
added in v0.8.0
func ServerFromContext(ctx context.Context) *MCPServer
ServerFromContext retrieves the MCPServer instance from a context

func (*MCPServer) AddNotificationHandler ¶
added in v0.5.0
func (s *MCPServer) AddNotificationHandler(
method string,
handler NotificationHandlerFunc,
)
AddNotificationHandler registers a new handler for incoming notifications

func (*MCPServer) AddPrompt ¶
added in v0.5.0
func (s *MCPServer) AddPrompt(prompt mcp.Prompt, handler PromptHandlerFunc)
AddPrompt registers a new prompt handler with the given name

func (*MCPServer) AddResource ¶
added in v0.5.0
func (s *MCPServer) AddResource(
resource mcp.Resource,
handler ResourceHandlerFunc,
)
AddResource registers a new resource and its handler

func (*MCPServer) AddResourceTemplate ¶
added in v0.5.0
func (s *MCPServer) AddResourceTemplate(
template mcp.ResourceTemplate,
handler ResourceTemplateHandlerFunc,
)
AddResourceTemplate registers a new resource template and its handler

func (*MCPServer) AddSessionTool ¶
added in v0.24.0
func (s *MCPServer) AddSessionTool(sessionID string, tool mcp.Tool, handler ToolHandlerFunc) error
AddSessionTool adds a tool for a specific session

func (*MCPServer) AddSessionTools ¶
added in v0.24.0
func (s *MCPServer) AddSessionTools(sessionID string, tools ...ServerTool) error
AddSessionTools adds tools for a specific session

func (*MCPServer) AddTool ¶
added in v0.5.0
func (s *MCPServer) AddTool(tool mcp.Tool, handler ToolHandlerFunc)
AddTool registers a new tool and its handler

func (*MCPServer) AddTools ¶
added in v0.8.5
func (s *MCPServer) AddTools(tools ...ServerTool)
AddTools registers multiple tools at once

func (*MCPServer) DeletePrompts ¶
added in v0.30.0
func (s *MCPServer) DeletePrompts(names ...string)
DeletePrompts removes prompts from the server

func (*MCPServer) DeleteSessionTools ¶
added in v0.24.0
func (s *MCPServer) DeleteSessionTools(sessionID string, names ...string) error
DeleteSessionTools removes tools from a specific session

func (*MCPServer) DeleteTools ¶
added in v0.8.5
func (s *MCPServer) DeleteTools(names ...string)
DeleteTools removes tools from the server

func (*MCPServer) HandleMessage ¶
added in v0.5.0
func (s *MCPServer) HandleMessage(
ctx context.Context,
message json.RawMessage,
) mcp.JSONRPCMessage
HandleMessage processes an incoming JSON-RPC message and returns an appropriate response

func (*MCPServer) RegisterSession ¶
added in v0.14.0
func (s *MCPServer) RegisterSession(
ctx context.Context,
session ClientSession,
) error
RegisterSession saves session that should be notified in case if some server attributes changed.

func (*MCPServer) RemoveResource ¶
added in v0.22.0
func (s *MCPServer) RemoveResource(uri string)
RemoveResource removes a resource from the server

func (*MCPServer) SendNotificationToAllClients ¶
added in v0.23.0
func (s *MCPServer) SendNotificationToAllClients(
method string,
params map[string]any,
)
SendNotificationToAllClients sends a notification to all the currently active clients.

func (*MCPServer) SendNotificationToClient ¶
added in v0.8.0
func (s *MCPServer) SendNotificationToClient(
ctx context.Context,
method string,
params map[string]any,
) error
SendNotificationToClient sends a notification to the current client

func (*MCPServer) SendNotificationToSpecificClient ¶
added in v0.24.0
func (s *MCPServer) SendNotificationToSpecificClient(
sessionID string,
method string,
params map[string]any,
) error
SendNotificationToSpecificClient sends a notification to a specific client by session ID

func (*MCPServer) SetTools ¶
added in v0.8.5
func (s *MCPServer) SetTools(tools ...ServerTool)
SetTools replaces all existing tools with the provided list

func (*MCPServer) UnregisterSession ¶
added in v0.14.0
func (s *MCPServer) UnregisterSession(
ctx context.Context,
sessionID string,
)
UnregisterSession removes from storage session that is shut down.

func (*MCPServer) WithContext ¶
added in v0.8.0
func (s *MCPServer) WithContext(
ctx context.Context,
session ClientSession,
) context.Context
WithContext sets the current client session and returns the provided context

type NotificationHandlerFunc ¶
added in v0.5.0
type NotificationHandlerFunc func(ctx context.Context, notification mcp.JSONRPCNotification)
NotificationHandlerFunc handles incoming notifications.

type OnAfterCallToolFunc ¶
added in v0.16.0
type OnAfterCallToolFunc func(ctx context.Context, id any, message *mcp.CallToolRequest, result *mcp.CallToolResult)
type OnAfterGetPromptFunc ¶
added in v0.16.0
type OnAfterGetPromptFunc func(ctx context.Context, id any, message *mcp.GetPromptRequest, result *mcp.GetPromptResult)
type OnAfterInitializeFunc ¶
added in v0.16.0
type OnAfterInitializeFunc func(ctx context.Context, id any, message *mcp.InitializeRequest, result *mcp.InitializeResult)
type OnAfterListPromptsFunc ¶
added in v0.16.0
type OnAfterListPromptsFunc func(ctx context.Context, id any, message *mcp.ListPromptsRequest, result *mcp.ListPromptsResult)
type OnAfterListResourceTemplatesFunc ¶
added in v0.16.0
type OnAfterListResourceTemplatesFunc func(ctx context.Context, id any, message *mcp.ListResourceTemplatesRequest, result *mcp.ListResourceTemplatesResult)
type OnAfterListResourcesFunc ¶
added in v0.16.0
type OnAfterListResourcesFunc func(ctx context.Context, id any, message *mcp.ListResourcesRequest, result *mcp.ListResourcesResult)
type OnAfterListToolsFunc ¶
added in v0.16.0
type OnAfterListToolsFunc func(ctx context.Context, id any, message *mcp.ListToolsRequest, result *mcp.ListToolsResult)
type OnAfterPingFunc ¶
added in v0.16.0
type OnAfterPingFunc func(ctx context.Context, id any, message *mcp.PingRequest, result *mcp.EmptyResult)
type OnAfterReadResourceFunc ¶
added in v0.16.0
type OnAfterReadResourceFunc func(ctx context.Context, id any, message *mcp.ReadResourceRequest, result *mcp.ReadResourceResult)
type OnAfterSetLevelFunc ¶
added in v0.28.0
type OnAfterSetLevelFunc func(ctx context.Context, id any, message *mcp.SetLevelRequest, result *mcp.EmptyResult)
type OnBeforeCallToolFunc ¶
added in v0.16.0
type OnBeforeCallToolFunc func(ctx context.Context, id any, message *mcp.CallToolRequest)
type OnBeforeGetPromptFunc ¶
added in v0.16.0
type OnBeforeGetPromptFunc func(ctx context.Context, id any, message *mcp.GetPromptRequest)
type OnBeforeInitializeFunc ¶
added in v0.16.0
type OnBeforeInitializeFunc func(ctx context.Context, id any, message *mcp.InitializeRequest)
type OnBeforeListPromptsFunc ¶
added in v0.16.0
type OnBeforeListPromptsFunc func(ctx context.Context, id any, message *mcp.ListPromptsRequest)
type OnBeforeListResourceTemplatesFunc ¶
added in v0.16.0
type OnBeforeListResourceTemplatesFunc func(ctx context.Context, id any, message *mcp.ListResourceTemplatesRequest)
type OnBeforeListResourcesFunc ¶
added in v0.16.0
type OnBeforeListResourcesFunc func(ctx context.Context, id any, message *mcp.ListResourcesRequest)
type OnBeforeListToolsFunc ¶
added in v0.16.0
type OnBeforeListToolsFunc func(ctx context.Context, id any, message *mcp.ListToolsRequest)
type OnBeforePingFunc ¶
added in v0.16.0
type OnBeforePingFunc func(ctx context.Context, id any, message *mcp.PingRequest)
type OnBeforeReadResourceFunc ¶
added in v0.16.0
type OnBeforeReadResourceFunc func(ctx context.Context, id any, message *mcp.ReadResourceRequest)
type OnBeforeSetLevelFunc ¶
added in v0.28.0
type OnBeforeSetLevelFunc func(ctx context.Context, id any, message *mcp.SetLevelRequest)
type OnErrorHookFunc ¶
added in v0.16.0
type OnErrorHookFunc func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error)
OnErrorHookFunc is a hook that will be called when an error occurs, either during the request parsing or the method execution.

Example usage: ```

hooks.AddOnError(func(ctx context.Context, id any, method mcp.MCPMethod, message any, err error) {
// Check for specific error types using errors.Is
if errors.Is(err, ErrUnsupported) {
// Handle capability not supported errors
log.Printf("Capability not supported: %v", err)
}

// Use errors.As to get specific error types
var parseErr = &UnparsableMessageError{}
if errors.As(err, &parseErr) {
// Access specific methods/fields of the error type
log.Printf("Failed to parse message for method %s: %v",
parseErr.GetMethod(), parseErr.Unwrap())
// Access the raw message that failed to parse
rawMsg := parseErr.GetMessage()
}

// Check for specific resource/prompt/tool errors
switch {
case errors.Is(err, ErrResourceNotFound):
log.Printf("Resource not found: %v", err)
case errors.Is(err, ErrPromptNotFound):
log.Printf("Prompt not found: %v", err)
case errors.Is(err, ErrToolNotFound):
log.Printf("Tool not found: %v", err)
}
})
type OnRegisterSessionHookFunc ¶
added in v0.18.0
type OnRegisterSessionHookFunc func(ctx context.Context, session ClientSession)
OnRegisterSessionHookFunc is a hook that will be called when a new session is registered.

type OnRequestInitializationFunc ¶
added in v0.24.0
type OnRequestInitializationFunc func(ctx context.Context, id any, message any) error
OnRequestInitializationFunc is a function that called before handle diff request method Should any errors arise during func execution, the service will promptly return the corresponding error message.

type OnSuccessHookFunc ¶
added in v0.16.0
type OnSuccessHookFunc func(ctx context.Context, id any, method mcp.MCPMethod, message any, result any)
OnSuccessHookFunc is a hook that will be called after the request successfully generates a result, but before the result is sent to the client.

type OnUnregisterSessionHookFunc ¶
added in v0.23.0
type OnUnregisterSessionHookFunc func(ctx context.Context, session ClientSession)
OnUnregisterSessionHookFunc is a hook that will be called when a session is being unregistered.

type PromptHandlerFunc ¶
added in v0.5.0
type PromptHandlerFunc func(ctx context.Context, request mcp.GetPromptRequest) (*mcp.GetPromptResult, error)
PromptHandlerFunc handles prompt requests with given arguments.

type ResourceHandlerFunc ¶
added in v0.5.0
type ResourceHandlerFunc func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)
ResourceHandlerFunc is a function that returns resource contents.

type ResourceTemplateHandlerFunc ¶
added in v0.5.0
type ResourceTemplateHandlerFunc func(ctx context.Context, request mcp.ReadResourceRequest) ([]mcp.ResourceContents, error)
ResourceTemplateHandlerFunc is a function that returns a resource template.

type SSEContextFunc ¶
added in v0.13.0
type SSEContextFunc func(ctx context.Context, r *http.Request) context.Context
SSEContextFunc is a function that takes an existing context and the current request and returns a potentially modified context based on the request content. This can be used to inject context values from headers, for example.

type SSEOption ¶
added in v0.13.0
type SSEOption func(*SSEServer)
SSEOption defines a function type for configuring SSEServer

func WithAppendQueryToMessageEndpoint ¶
added in v0.25.0
func WithAppendQueryToMessageEndpoint() SSEOption
WithAppendQueryToMessageEndpoint configures the SSE server to append the original request's query parameters to the message endpoint URL that is sent to clients during the SSE connection initialization. This is useful when you need to preserve query parameters from the initial SSE connection request and carry them over to subsequent message requests, maintaining context or authentication details across the communication channel.

func
WithBasePath
deprecated
added in v0.12.0
func WithBaseURL ¶
added in v0.12.0
func WithBaseURL(baseURL string) SSEOption
WithBaseURL sets the base URL for the SSE server

func WithDynamicBasePath ¶
added in v0.25.0
func WithDynamicBasePath(fn DynamicBasePathFunc) SSEOption
WithDynamicBasePath accepts a function for generating the base path. This is useful for cases where the base path is not known at the time of SSE server creation, such as when using a reverse proxy or when the server is mounted at a dynamic path.

func WithHTTPServer ¶
added in v0.12.0
func WithHTTPServer(srv *http.Server) SSEOption
WithHTTPServer sets the HTTP server instance

func WithKeepAlive ¶
added in v0.20.0
func WithKeepAlive(keepAlive bool) SSEOption
func WithKeepAliveInterval ¶
added in v0.20.0
func WithKeepAliveInterval(keepAliveInterval time.Duration) SSEOption
func WithMessageEndpoint ¶
added in v0.12.0
func WithMessageEndpoint(endpoint string) SSEOption
WithMessageEndpoint sets the message endpoint path

func WithSSEContextFunc ¶
added in v0.13.0
func WithSSEContextFunc(fn SSEContextFunc) SSEOption
WithSSEContextFunc sets a function that will be called to customise the context to the server using the incoming request.

func WithSSEEndpoint ¶
added in v0.12.0
func WithSSEEndpoint(endpoint string) SSEOption
WithSSEEndpoint sets the SSE endpoint path

func WithStaticBasePath ¶
added in v0.26.0
func WithStaticBasePath(basePath string) SSEOption
WithStaticBasePath adds a new option for setting a static base path

func WithUseFullURLForMessageEndpoint ¶
added in v0.18.0
func WithUseFullURLForMessageEndpoint(useFullURLForMessageEndpoint bool) SSEOption
WithUseFullURLForMessageEndpoint controls whether the SSE server returns a complete URL (including baseURL) or just the path portion for the message endpoint. Set to false when clients will concatenate the baseURL themselves to avoid malformed URLs like "http://localhost/mcphttp://localhost/mcp/message".

type SSEServer ¶
added in v0.2.0
type SSEServer struct {
// contains filtered or unexported fields
}
SSEServer implements a Server-Sent Events (SSE) based MCP server. It provides real-time communication capabilities over HTTP using the SSE protocol.

func NewSSEServer ¶
added in v0.2.0
func NewSSEServer(server *MCPServer, opts ...SSEOption) *SSEServer
NewSSEServer creates a new SSE server instance with the given MCP server and options.

func (*SSEServer) CompleteMessageEndpoint ¶
added in v0.16.0
func (s *SSEServer) CompleteMessageEndpoint() (string, error)
func (*SSEServer) CompleteMessagePath ¶
added in v0.16.0
func (s *SSEServer) CompleteMessagePath() string
func (*SSEServer) CompleteSseEndpoint ¶
added in v0.16.0
func (s *SSEServer) CompleteSseEndpoint() (string, error)
func (*SSEServer) CompleteSsePath ¶
added in v0.16.0
func (s *SSEServer) CompleteSsePath() string
func (*SSEServer) GetMessageEndpointForClient ¶
added in v0.18.0
func (s *SSEServer) GetMessageEndpointForClient(r *http.Request, sessionID string) string
GetMessageEndpointForClient returns the appropriate message endpoint URL with session ID for the given request. This is the canonical way to compute the message endpoint for a client. It handles both dynamic and static path modes, and honors the WithUseFullURLForMessageEndpoint flag.

func (*SSEServer) GetUrlPath ¶
added in v0.16.0
func (s *SSEServer) GetUrlPath(input string) (string, error)
func (*SSEServer) MessageHandler ¶
added in v0.25.0
func (s *SSEServer) MessageHandler() http.Handler
MessageHandler returns an http.Handler for the message endpoint.

This method allows you to mount the message handler at any arbitrary path using your own router (e.g. net/http, gorilla/mux, chi, etc.). It is intended for advanced scenarios where you want to control the routing or support dynamic segments.

IMPORTANT: When using this handler in advanced/dynamic mounting scenarios, you must use the WithDynamicBasePath option to ensure the correct base path is communicated to clients.

Example usage:

// Advanced/dynamic:
sseServer := NewSSEServer(mcpServer,
WithDynamicBasePath(func(r *http.Request, sessionID string) string {
tenant := r.PathValue("tenant")
return "/mcp/" + tenant
}),
WithBaseURL("http://localhost:8080")
)
mux.Handle("/mcp/{tenant}/sse", sseServer.SSEHandler())
mux.Handle("/mcp/{tenant}/message", sseServer.MessageHandler())
For non-dynamic cases, use ServeHTTP method instead.

func (*SSEServer) SSEHandler ¶
added in v0.25.0
func (s *SSEServer) SSEHandler() http.Handler
SSEHandler returns an http.Handler for the SSE endpoint.

This method allows you to mount the SSE handler at any arbitrary path using your own router (e.g. net/http, gorilla/mux, chi, etc.). It is intended for advanced scenarios where you want to control the routing or support dynamic segments.

IMPORTANT: When using this handler in advanced/dynamic mounting scenarios, you must use the WithDynamicBasePath option to ensure the correct base path is communicated to clients.

Example usage:

// Advanced/dynamic:
sseServer := NewSSEServer(mcpServer,
WithDynamicBasePath(func(r *http.Request, sessionID string) string {
tenant := r.PathValue("tenant")
return "/mcp/" + tenant
}),
WithBaseURL("http://localhost:8080")
)
mux.Handle("/mcp/{tenant}/sse", sseServer.SSEHandler())
mux.Handle("/mcp/{tenant}/message", sseServer.MessageHandler())
For non-dynamic cases, use ServeHTTP method instead.

func (*SSEServer) SendEventToSession ¶
added in v0.2.0
func (s *SSEServer) SendEventToSession(
sessionID string,
event any,
) error
SendEventToSession sends an event to a specific SSE session identified by sessionID. Returns an error if the session is not found or closed.

func (*SSEServer) ServeHTTP ¶
added in v0.10.0
func (s *SSEServer) ServeHTTP(w http.ResponseWriter, r *http.Request)
ServeHTTP implements the http.Handler interface.

func (*SSEServer) Shutdown ¶
added in v0.2.0
func (s *SSEServer) Shutdown(ctx context.Context) error
Shutdown gracefully stops the SSE server, closing all active sessions and shutting down the HTTP server.

func (*SSEServer) Start ¶
added in v0.2.0
func (s *SSEServer) Start(addr string) error
Start begins serving SSE connections on the specified address. It sets up HTTP handlers for SSE and message endpoints.

type ServerOption ¶
added in v0.5.0
type ServerOption func(*MCPServer)
ServerOption is a function that configures an MCPServer.

func WithHooks ¶
added in v0.16.0
func WithHooks(hooks *Hooks) ServerOption
WithHooks allows adding hooks that will be called before or after either [all] requests or before / after specific request methods, or else prior to returning an error to the client.

func WithInstructions ¶
added in v0.12.0
func WithInstructions(instructions string) ServerOption
WithInstructions sets the server instructions for the client returned in the initialize response

func WithLogging ¶
added in v0.5.0
func WithLogging() ServerOption
WithLogging enables logging capabilities for the server

func WithPaginationLimit ¶
added in v0.19.0
func WithPaginationLimit(limit int) ServerOption
WithPaginationLimit sets the pagination limit for the server.

func WithPromptCapabilities ¶
added in v0.5.0
func WithPromptCapabilities(listChanged bool) ServerOption
WithPromptCapabilities configures prompt-related server capabilities

func WithRecovery ¶
added in v0.20.0
func WithRecovery() ServerOption
WithRecovery adds a middleware that recovers from panics in tool handlers.

func WithResourceCapabilities ¶
added in v0.5.0
func WithResourceCapabilities(subscribe, listChanged bool) ServerOption
WithResourceCapabilities configures resource-related server capabilities

func WithToolCapabilities ¶
added in v0.5.0
func WithToolCapabilities(listChanged bool) ServerOption
WithToolCapabilities configures tool-related server capabilities

func WithToolFilter ¶
added in v0.24.0
func WithToolFilter(
toolFilter ToolFilterFunc,
) ServerOption
WithToolFilter adds a filter function that will be applied to tools before they are returned in list_tools

func WithToolHandlerMiddleware ¶
added in v0.20.0
func WithToolHandlerMiddleware(
toolHandlerMiddleware ToolHandlerMiddleware,
) ServerOption
WithToolHandlerMiddleware allows adding a middleware for the tool handler call chain.

type ServerTool ¶
added in v0.8.5
type ServerTool struct {
Tool    mcp.Tool
Handler ToolHandlerFunc
}
ServerTool combines a Tool with its ToolHandlerFunc.

type SessionIdManager ¶
added in v0.30.0
type SessionIdManager interface {
Generate() string
// Validate checks if a session ID is valid and not terminated.
// Returns isTerminated=true if the ID is valid but belongs to a terminated session.
// Returns err!=nil if the ID format is invalid or lookup failed.
Validate(sessionID string) (isTerminated bool, err error)
// Terminate marks a session ID as terminated.
// Returns isNotAllowed=true if the server policy prevents client termination.
// Returns err!=nil if the ID is invalid or termination failed.
Terminate(sessionID string) (isNotAllowed bool, err error)
}
type SessionWithClientInfo ¶
added in v0.30.0
type SessionWithClientInfo interface {
ClientSession
// GetClientInfo returns the client information for this session
GetClientInfo() mcp.Implementation
// SetClientInfo sets the client information for this session
SetClientInfo(clientInfo mcp.Implementation)
}
SessionWithClientInfo is an extension of ClientSession that can store client info

type SessionWithLogging ¶
added in v0.28.0
type SessionWithLogging interface {
ClientSession
// SetLogLevel sets the minimum log level
SetLogLevel(level mcp.LoggingLevel)
// GetLogLevel retrieves the minimum log level
GetLogLevel() mcp.LoggingLevel
}
SessionWithLogging is an extension of ClientSession that can receive log message notifications and set log level

type SessionWithTools ¶
added in v0.24.0
type SessionWithTools interface {
ClientSession
// GetSessionTools returns the tools specific to this session, if any
// This method must be thread-safe for concurrent access
GetSessionTools() map[string]ServerTool
// SetSessionTools sets tools specific to this session
// This method must be thread-safe for concurrent access
SetSessionTools(tools map[string]ServerTool)
}
SessionWithTools is an extension of ClientSession that can store session-specific tool data

type StatelessSessionIdManager ¶
added in v0.30.0
type StatelessSessionIdManager struct{}
StatelessSessionIdManager does nothing, which means it has no session management, which is stateless.

func (*StatelessSessionIdManager) Generate ¶
added in v0.30.0
func (s *StatelessSessionIdManager) Generate() string
func (*StatelessSessionIdManager) Terminate ¶
added in v0.30.0
func (s *StatelessSessionIdManager) Terminate(sessionID string) (isNotAllowed bool, err error)
func (*StatelessSessionIdManager) Validate ¶
added in v0.30.0
func (s *StatelessSessionIdManager) Validate(sessionID string) (isTerminated bool, err error)
type StdioContextFunc ¶
added in v0.13.0
type StdioContextFunc func(ctx context.Context) context.Context
StdioContextFunc is a function that takes an existing context and returns a potentially modified context. This can be used to inject context values from environment variables, for example.

type StdioOption ¶
added in v0.13.0
type StdioOption func(*StdioServer)
StdioOption defines a function type for configuring StdioServer

func WithErrorLogger ¶
added in v0.13.0
func WithErrorLogger(logger *log.Logger) StdioOption
WithErrorLogger sets the error logger for the server

func WithStdioContextFunc ¶
added in v0.13.0
func WithStdioContextFunc(fn StdioContextFunc) StdioOption
WithStdioContextFunc sets a function that will be called to customise the context to the server. Note that the stdio server uses the same context for all requests, so this function will only be called once per server instance.

type StdioServer ¶
type StdioServer struct {
// contains filtered or unexported fields
}
StdioServer wraps a MCPServer and handles stdio communication. It provides a simple way to create command-line MCP servers that communicate via standard input/output streams using JSON-RPC messages.

func NewStdioServer ¶
added in v0.5.5
func NewStdioServer(server *MCPServer) *StdioServer
NewStdioServer creates a new stdio server wrapper around an MCPServer. It initializes the server with a default error logger that discards all output.

func (*StdioServer) Listen ¶
added in v0.5.5
func (s *StdioServer) Listen(
ctx context.Context,
stdin io.Reader,
stdout io.Writer,
) error
Listen starts listening for JSON-RPC messages on the provided input and writes responses to the provided output. It runs until the context is cancelled or an error occurs. Returns an error if there are issues with reading input or writing output.

func (*StdioServer) SetContextFunc ¶
added in v0.13.0
func (s *StdioServer) SetContextFunc(fn StdioContextFunc)
SetContextFunc sets a function that will be called to customise the context to the server. Note that the stdio server uses the same context for all requests, so this function will only be called once per server instance.

func (*StdioServer) SetErrorLogger ¶
added in v0.5.5
func (s *StdioServer) SetErrorLogger(logger *log.Logger)
SetErrorLogger configures where error messages from the StdioServer are logged. The provided logger will receive all error messages generated during server operation.

type StreamableHTTPOption ¶
added in v0.27.0
type StreamableHTTPOption func(*StreamableHTTPServer)
StreamableHTTPOption defines a function type for configuring StreamableHTTPServer

func WithEndpointPath ¶
added in v0.30.0
func WithEndpointPath(endpointPath string) StreamableHTTPOption
WithEndpointPath sets the endpoint path for the server. The default is "/mcp". It's only works for `Start` method. When used as a http.Handler, it has no effect.

func WithHTTPContextFunc ¶
added in v0.27.0
func WithHTTPContextFunc(fn HTTPContextFunc) StreamableHTTPOption
WithHTTPContextFunc sets a function that will be called to customise the context to the server using the incoming request. This can be used to inject context values from headers, for example.

func WithHeartbeatInterval ¶
added in v0.30.0
func WithHeartbeatInterval(interval time.Duration) StreamableHTTPOption
WithHeartbeatInterval sets the heartbeat interval. Positive interval means the server will send a heartbeat to the client through the GET connection, to keep the connection alive from being closed by the network infrastructure (e.g. gateways). If the client does not establish a GET connection, it has no effect. The default is not to send heartbeats.

func WithLogger ¶
added in v0.30.0
func WithLogger(logger util.Logger) StreamableHTTPOption
WithLogger sets the logger for the server

func WithSessionIdManager ¶
added in v0.30.0
func WithSessionIdManager(manager SessionIdManager) StreamableHTTPOption
WithSessionIdManager sets a custom session id generator for the server. By default, the server will use SimpleStatefulSessionIdGenerator, which generates session ids with uuid, and it's insecure. Notice: it will override the WithStateLess option.

func WithStateLess ¶
added in v0.30.0
func WithStateLess(stateLess bool) StreamableHTTPOption
WithStateLess sets the server to stateless mode. If true, the server will manage no session information. Every request will be treated as a new session. No session id returned to the client. The default is false.

Notice: This is a convenience method. It's identical to set WithSessionIdManager option to StatelessSessionIdManager.

type StreamableHTTPServer ¶
added in v0.27.0
type StreamableHTTPServer struct {
// contains filtered or unexported fields
}
StreamableHTTPServer implements a Streamable-http based MCP server. It communicates with clients over HTTP protocol, supporting both direct HTTP responses, and SSE streams. https://modelcontextprotocol.io/specification/2025-03-26/basic/transports#streamable-http

Usage:

server := NewStreamableHTTPServer(mcpServer)
server.Start(":8080") // The final url for client is http://xxxx:8080/mcp by default
or the server itself can be used as a http.Handler, which is convenient to integrate with existing http servers, or advanced usage:

handler := NewStreamableHTTPServer(mcpServer)
http.Handle("/streamable-http", handler)
http.ListenAndServe(":8080", nil)
Notice: Except for the GET handlers(listening), the POST handlers(request/notification) will not trigger the session registration. So the methods like `SendNotificationToSpecificClient` or `hooks.onRegisterSession` will not be triggered for POST messages.

The current implementation does not support the following features from the specification:

Batching of requests/notifications/responses in arrays.
Stream Resumability
func NewStreamableHTTPServer ¶
added in v0.30.0
func NewStreamableHTTPServer(server *MCPServer, opts ...StreamableHTTPOption) *StreamableHTTPServer
NewStreamableHTTPServer creates a new streamable-http server instance

func (*StreamableHTTPServer) ServeHTTP ¶
added in v0.30.0
func (s *StreamableHTTPServer) ServeHTTP(w http.ResponseWriter, r *http.Request)
ServeHTTP implements the http.Handler interface.

func (*StreamableHTTPServer) Shutdown ¶
added in v0.30.0
func (s *StreamableHTTPServer) Shutdown(ctx context.Context) error
Shutdown gracefully stops the server, closing all active sessions and shutting down the HTTP server.

func (*StreamableHTTPServer) Start ¶
added in v0.30.0
func (s *StreamableHTTPServer) Start(addr string) error
Start begins serving the http server on the specified address and path (endpointPath). like:

s.Start(":8080")
type ToolFilterFunc ¶
added in v0.24.0
type ToolFilterFunc func(ctx context.Context, tools []mcp.Tool) []mcp.Tool
ToolFilterFunc is a function that filters tools based on context, typically using session information.

type ToolHandlerFunc ¶
added in v0.5.0
type ToolHandlerFunc func(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error)
ToolHandlerFunc handles tool calls with given arguments.

type ToolHandlerMiddleware ¶
added in v0.20.0
type ToolHandlerMiddleware func(ToolHandlerFunc) ToolHandlerFunc
ToolHandlerMiddleware is a middleware function that wraps a ToolHandlerFunc.

type UnparsableMessageError ¶
added in v0.23.0
type UnparsableMessageError struct {
// contains filtered or unexported fields
}
UnparsableMessageError is attached to the RequestError when json.Unmarshal fails on the request.

func (*UnparsableMessageError) Error ¶
added in v0.23.0
func (e *UnparsableMessageError) Error() string
func (*UnparsableMessageError) GetMessage ¶
added in v0.23.0
func (e *UnparsableMessageError) GetMessage() json.RawMessage
func (*UnparsableMessageError) GetMethod ¶
added in v0.23.0
func (e *UnparsableMessageError) GetMethod() mcp.MCPMethod
func (*UnparsableMessageError) Unwrap ¶
added in v0.23.0
func (e *UnparsableMessageError) Unwrap() error