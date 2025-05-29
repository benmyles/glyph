=== README.md ===

# go tree-sitter

[![Build Status](https://github.com/smacker/go-tree-sitter/workflows/Test/badge.svg?branch=master)](https://github.com/smacker/go-tree-sitter/actions/workflows/test.yml?query=branch%3Amaster)
[![GoDoc](https://godoc.org/github.com/smacker/go-tree-sitter?status.svg)](https://godoc.org/github.com/smacker/go-tree-sitter)

Golang bindings for [tree-sitter](https://github.com/tree-sitter/tree-sitter)

## Usage

Create a parser with a grammar:

```go
import (
	"context"
	"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

parser := sitter.NewParser()
parser.SetLanguage(javascript.GetLanguage())
```

Parse some code:

```go
sourceCode := []byte("let a = 1")
tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)
```

Inspect the syntax tree:

```go
n := tree.RootNode()

fmt.Println(n) // (program (lexical_declaration (variable_declarator (identifier) (number))))

child := n.NamedChild(0)
fmt.Println(child.Type()) // lexical_declaration
fmt.Println(child.StartByte()) // 0
fmt.Println(child.EndByte()) // 9
```

### Custom grammars

This repository provides grammars for many common languages out of the box.

But if you need support for any other language you can keep it inside your own project or publish it as a separate repository to share with the community.

See explanation on how to create a grammar for go-tree-sitter [here](https://github.com/smacker/go-tree-sitter/issues/57).

Known external grammars:

- [Salesforce grammars](https://github.com/aheber/tree-sitter-sfapex) - including Apex, SOQL, and SOSL languages.
- [Ruby](https://github.com/shagabutdinov/go-tree-sitter-ruby) - Deprecated, grammar is provided by main repo instead
- [Go Template](https://github.com/mrjosh/helm-ls/tree/master/internal/tree-sitter/gotemplate) - Used for helm

### Editing

If your source code changes, you can update the syntax tree. This will take less time than the first parse.

```go
// change 1 -> true
newText := []byte("let a = true")
tree.Edit(sitter.EditInput{
    StartIndex:  8,
    OldEndIndex: 9,
    NewEndIndex: 12,
    StartPoint: sitter.Point{
        Row:    0,
        Column: 8,
    },
    OldEndPoint: sitter.Point{
        Row:    0,
        Column: 9,
    },
    NewEndPoint: sitter.Point{
        Row:    0,
        Column: 12,
    },
})

// check that it changed tree
assert.True(n.HasChanges())
assert.True(n.Child(0).HasChanges())
assert.False(n.Child(0).Child(0).HasChanges()) // left side of the tree didn't change
assert.True(n.Child(0).Child(1).HasChanges())

// generate new tree
newTree := parser.Parse(tree, newText)
```

### Predicates

You can filter AST by using [predicate](https://tree-sitter.github.io/tree-sitter/using-parsers#predicates) S-expressions.

Similar to [Rust](https://github.com/tree-sitter/tree-sitter/tree/master/lib/binding_rust) or [WebAssembly](https://github.com/tree-sitter/tree-sitter/blob/master/lib/binding_web) bindings we support filtering on a few common predicates:
- `eq?`, `not-eq?`
- `match?`, `not-match?`

Usage [example](./_examples/predicates/main.go):

```go
func main() {
	// Javascript code
	sourceCode := []byte(`
		const camelCaseConst = 1;
		const SCREAMING_SNAKE_CASE_CONST = 2;
		const lower_snake_case_const = 3;`)
	// Query with predicates
	screamingSnakeCasePattern := `(
		(identifier) @constant
		(#match? @constant "^[A-Z][A-Z_]+")
	)`

	// Parse source code
	lang := javascript.GetLanguage()
	n, _ := sitter.ParseCtx(context.Background(), sourceCode, lang)
	// Execute the query
	q, _ := sitter.NewQuery([]byte(screamingSnakeCasePattern), lang)
	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)
	// Iterate over query results
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		// Apply predicates filtering
		m = qc.FilterPredicates(m, sourceCode)
		for _, c := range m.Captures {
			fmt.Println(c.Node.Content(sourceCode))
		}
	}
}

// Output of this program:
// SCREAMING_SNAKE_CASE_CONST
```

## Development

### Updating a grammar

Check if any updates for vendored files are available:

```
go run _automation/main.go check-updates
```

Update vendor files:

- open `_automation/grammars.json`
- modify `reference` (for tagged grammars) or `revision` (for grammars from a branch)
- run `go run _automation/main.go update <grammar-name>`

It is also possible to update all grammars in one go using

```
go run _automation/main.go update-all
```

=== _examples/main.go ===

package main

import (
"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

func main() {
input := []byte("function hello() { console.log('hello') }; function goodbye(){}")

	parser := sitter.NewParser()
	parser.SetLanguage(javascript.GetLanguage())

	tree := parser.Parse(nil, input)

	n := tree.RootNode()

	fmt.Println("AST:", n)
	fmt.Println("Root type:", n.Type())
	fmt.Println("Root children:", n.ChildCount())

	fmt.Println("\nFunctions in input:")
	q, _ := sitter.NewQuery([]byte("(function_declaration) @func"), javascript.GetLanguage())
	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)

	var funcs []*sitter.Node
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}

		for _, c := range m.Captures {
			funcs = append(funcs, c.Node)
			fmt.Println("-", funcName(input, c.Node))
		}
	}

	fmt.Println("\nEdit input")
	input = []byte("function hello() { console.log('hello') }; function goodbye(){ console.log('goodbye') }")
	// reuse tree
	tree.Edit(sitter.EditInput{
		StartIndex:  62,
		OldEndIndex: 63,
		NewEndIndex: 87,
		StartPoint: sitter.Point{
			Row:    0,
			Column: 62,
		},
		OldEndPoint: sitter.Point{
			Row:    0,
			Column: 63,
		},
		NewEndPoint: sitter.Point{
			Row:    0,
			Column: 87,
		},
	})

	for _, f := range funcs {
		var textChange string
		if f.HasChanges() {
			textChange = "has change"
		} else {
			textChange = "no changes"
		}
		fmt.Println("-", funcName(input, f), ">", textChange)
	}

	newTree := parser.Parse(tree, input)
	n = newTree.RootNode()
	fmt.Println("\nNew AST:", n)
}

func funcName(content []byte, n *sitter.Node) string {
if n == nil {
return ""
}

	if n.Type() != "function_declaration" {
		return ""
	}

	return n.ChildByFieldName("name").Content(content)
}

=== _examples/predicates/main.go ===

package main

import (
"context"
"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

func main() {
// Javascript code
sourceCode := []byte(`
	const camelCaseConst = 1;
	const SCREAMING_SNAKE_CASE_CONST = 2;
	const lower_snake_case_const = 3;`)
// Query with predicates
screamingSnakeCasePattern := `(
	(identifier) @constant
	(#match? @constant "^[A-Z][A-Z_]+")
)`

	// Parse source code
	lang := javascript.GetLanguage()
	n, _ := sitter.ParseCtx(context.Background(), sourceCode, lang)
	// Execute the query
	q, _ := sitter.NewQuery([]byte(screamingSnakeCasePattern), lang)
	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)
	// Iterate over query results
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		// Apply predicates filtering
		m = qc.FilterPredicates(m, sourceCode)
		for _, c := range m.Captures {
			fmt.Println(c.Node.Content(sourceCode))
		}
	}
}

=== godoc ===

README ¶
go tree-sitter
Build Status GoDoc

Golang bindings for tree-sitter

Usage
Create a parser with a grammar:

import (
"context"
"fmt"

	sitter "github.com/smacker/go-tree-sitter"
	"github.com/smacker/go-tree-sitter/javascript"
)

parser := sitter.NewParser()
parser.SetLanguage(javascript.GetLanguage())
Parse some code:

sourceCode := []byte("let a = 1")
tree, _ := parser.ParseCtx(context.Background(), nil, sourceCode)
Inspect the syntax tree:

n := tree.RootNode()

fmt.Println(n) // (program (lexical_declaration (variable_declarator (identifier) (number))))

child := n.NamedChild(0)
fmt.Println(child.Type()) // lexical_declaration
fmt.Println(child.StartByte()) // 0
fmt.Println(child.EndByte()) // 9
Custom grammars
This repository provides grammars for many common languages out of the box.

But if you need support for any other language you can keep it inside your own project or publish it as a separate repository to share with the community.

See explanation on how to create a grammar for go-tree-sitter here.

Known external grammars:

Salesforce grammars - including Apex, SOQL, and SOSL languages.
Ruby - Deprecated, grammar is provided by main repo instead
Go Template - Used for helm
Editing
If your source code changes, you can update the syntax tree. This will take less time than the first parse.

// change 1 -> true
newText := []byte("let a = true")
tree.Edit(sitter.EditInput{
StartIndex:  8,
OldEndIndex: 9,
NewEndIndex: 12,
StartPoint: sitter.Point{
Row:    0,
Column: 8,
},
OldEndPoint: sitter.Point{
Row:    0,
Column: 9,
},
NewEndPoint: sitter.Point{
Row:    0,
Column: 12,
},
})

// check that it changed tree
assert.True(n.HasChanges())
assert.True(n.Child(0).HasChanges())
assert.False(n.Child(0).Child(0).HasChanges()) // left side of the tree didn't change
assert.True(n.Child(0).Child(1).HasChanges())

// generate new tree
newTree := parser.Parse(tree, newText)
Predicates
You can filter AST by using predicate S-expressions.

Similar to Rust or WebAssembly bindings we support filtering on a few common predicates:

eq?, not-eq?
match?, not-match?
Usage example:

func main() {
// Javascript code
sourceCode := []byte(`
		const camelCaseConst = 1;
		const SCREAMING_SNAKE_CASE_CONST = 2;
		const lower_snake_case_const = 3;`)
// Query with predicates
screamingSnakeCasePattern := `(
		(identifier) @constant
		(#match? @constant "^[A-Z][A-Z_]+")
	)`

	// Parse source code
	lang := javascript.GetLanguage()
	n, _ := sitter.ParseCtx(context.Background(), sourceCode, lang)
	// Execute the query
	q, _ := sitter.NewQuery([]byte(screamingSnakeCasePattern), lang)
	qc := sitter.NewQueryCursor()
	qc.Exec(q, n)
	// Iterate over query results
	for {
		m, ok := qc.NextMatch()
		if !ok {
			break
		}
		// Apply predicates filtering
		m = qc.FilterPredicates(m, sourceCode)
		for _, c := range m.Captures {
			fmt.Println(c.Node.Content(sourceCode))
		}
	}
}

// Output of this program:
// SCREAMING_SNAKE_CASE_CONST
Development
Updating a grammar
Check if any updates for vendored files are available:

go run _automation/main.go check-updates
Update vendor files:

open _automation/grammars.json
modify reference (for tagged grammars) or revision (for grammars from a branch)
run go run _automation/main.go update <grammar-name>
It is also possible to update all grammars in one go using

go run _automation/main.go update-all
Collapse ▴
Documentation ¶
Overview ¶
Code generated by test_grammar_generate.sh; DO NOT EDIT.

Index ¶
Constants
Variables
func QueryErrorTypeToString(errorType QueryErrorType) string
type BaseTree
func (t *BaseTree) Close()
type EditInput
type Input
type InputEncoding
type IterMode
type Iterator
func NewIterator(n *Node, mode IterMode) *Iterator
func NewNamedIterator(n *Node, mode IterMode) *Iterator
func (iter *Iterator) ForEach(fn func(*Node) error) error
func (iter *Iterator) Next() (*Node, error)
type Language
func NewLanguage(ptr unsafe.Pointer) *Language
func (l *Language) FieldName(idx int) string
func (l *Language) SymbolCount() uint32
func (l *Language) SymbolName(s Symbol) string
func (l *Language) SymbolType(s Symbol) SymbolType
type Node
func Parse(content []byte, lang *Language) *Nodedeprecated
func ParseCtx(ctx context.Context, content []byte, lang *Language) (*Node, error)
func (n Node) Child(idx int) *Node
func (n Node) ChildByFieldName(name string) *Node
func (n Node) ChildCount() uint32
func (n Node) Content(input []byte) string
func (n Node) Edit(i EditInput)
func (n Node) EndByte() uint32
func (n Node) EndPoint() Point
func (n Node) Equal(other *Node) bool
func (n Node) FieldNameForChild(idx int) string
func (n Node) HasChanges() bool
func (n Node) HasError() bool
func (n Node) ID() uintptr
func (n Node) IsError() bool
func (n Node) IsExtra() bool
func (n Node) IsMissing() bool
func (n Node) IsNamed() bool
func (n Node) IsNull() bool
func (n Node) NamedChild(idx int) *Node
func (n Node) NamedChildCount() uint32
func (n Node) NamedDescendantForPointRange(start Point, end Point) *Node
func (n Node) NextNamedSibling() *Node
func (n Node) NextSibling() *Node
func (n Node) Parent() *Node
func (n Node) PrevNamedSibling() *Node
func (n Node) PrevSibling() *Node
func (n Node) Range() Range
func (n Node) StartByte() uint32
func (n Node) StartPoint() Point
func (n Node) String() string
func (n Node) Symbol() Symbol
func (n Node) Type() string
type Parser
func NewParser() *Parser
func (p *Parser) Close()
func (p *Parser) Debug()
func (p *Parser) OperationLimit() int
func (p *Parser) Parse(oldTree *Tree, content []byte) *Treedeprecated
func (p *Parser) ParseCtx(ctx context.Context, oldTree *Tree, content []byte) (*Tree, error)
func (p *Parser) ParseInput(oldTree *Tree, input Input) *Tree
func (p *Parser) ParseInputCtx(ctx context.Context, oldTree *Tree, input Input) (*Tree, error)
func (p *Parser) Reset()
func (p *Parser) SetIncludedRanges(ranges []Range)
func (p *Parser) SetLanguage(lang *Language)
func (p *Parser) SetOperationLimit(limit int)
type Point
type Quantifier
type Query
func NewQuery(pattern []byte, lang *Language) (*Query, error)
func (q *Query) CaptureCount() uint32
func (q *Query) CaptureNameForId(id uint32) string
func (q *Query) CaptureQuantifierForId(id uint32, captureId uint32) Quantifier
func (q *Query) Close()
func (q *Query) PatternCount() uint32
func (q *Query) PredicatesForPattern(patternIndex uint32) [][]QueryPredicateStep
func (q *Query) StringCount() uint32
func (q *Query) StringValueForId(id uint32) string
type QueryCapture
type QueryCursor
func NewQueryCursor() *QueryCursor
func (qc *QueryCursor) Close()
func (qc *QueryCursor) Exec(q *Query, n *Node)
func (qc *QueryCursor) FilterPredicates(m *QueryMatch, input []byte) *QueryMatch
func (qc *QueryCursor) NextCapture() (*QueryMatch, uint32, bool)
func (qc *QueryCursor) NextMatch() (*QueryMatch, bool)
func (qc *QueryCursor) SetPointRange(startPoint Point, endPoint Point)
type QueryError
func (qe *QueryError) Error() string
type QueryErrorType
type QueryMatch
type QueryPredicateStep
type QueryPredicateStepType
type Range
type ReadFunc
type Symbol
type SymbolType
func (t SymbolType) String() string
type Tree
func (t *Tree) Copy() *Tree
func (t *Tree) Edit(i EditInput)
func (t *Tree) RootNode() *Node
type TreeCursor
func NewTreeCursor(n *Node) *TreeCursor
func (c *TreeCursor) Close()
func (c *TreeCursor) CurrentFieldName() string
func (c *TreeCursor) CurrentNode() *Node
func (c *TreeCursor) GoToFirstChild() bool
func (c *TreeCursor) GoToFirstChildForByte(b uint32) int64
func (c *TreeCursor) GoToNextSibling() bool
func (c *TreeCursor) GoToParent() bool
func (c *TreeCursor) Reset(n *Node)
Constants ¶
View Source
const (
QuantifierZero = iota
QuantifierZeroOrOne
QuantifierZeroOrMore
QuantifierOne
QuantifierOneOrMore
)
Variables ¶
View Source
var (
ErrOperationLimit = errors.New("operation limit was hit")
ErrNoLanguage     = errors.New("cannot parse without language")
)
Functions ¶
func QueryErrorTypeToString ¶
func QueryErrorTypeToString(errorType QueryErrorType) string
Types ¶
type BaseTree ¶
type BaseTree struct {
// contains filtered or unexported fields
}
we use cache for nodes on normal tree object it prevent run of SetFinalizer as it introduces cycle we can workaround it using separate object for details see: https://github.com/golang/go/issues/7358#issuecomment-66091558

func (*BaseTree) Close ¶
func (t *BaseTree) Close()
Close should be called to ensure that all the memory used by the tree is freed.

As the constructor in go-tree-sitter would set this func call through runtime.SetFinalizer, parser.Close() will be called by Go's garbage collector and users would not have to call this manually.

type EditInput ¶
type EditInput struct {
StartIndex  uint32
OldEndIndex uint32
NewEndIndex uint32
StartPoint  Point
OldEndPoint Point
NewEndPoint Point
}
type Input ¶
type Input struct {
Read     ReadFunc
Encoding InputEncoding
}
Input defines parameters for parse method

type InputEncoding ¶
type InputEncoding int
InputEncoding is a encoding of the text to parse

const (
InputEncodingUTF8 InputEncoding = iota
InputEncodingUTF16
)
type IterMode ¶
type IterMode int
const (
DFSMode IterMode = iota
BFSMode
)
type Iterator ¶
type Iterator struct {
// contains filtered or unexported fields
}
Iterator for a tree of nodes

func NewIterator ¶
func NewIterator(n *Node, mode IterMode) *Iterator
NewIterator takes a node and mode (DFS/BFS) and returns iterator over children of the node

func NewNamedIterator ¶
func NewNamedIterator(n *Node, mode IterMode) *Iterator
NewNamedIterator takes a node and mode (DFS/BFS) and returns iterator over named children of the node

func (*Iterator) ForEach ¶
func (iter *Iterator) ForEach(fn func(*Node) error) error
func (*Iterator) Next ¶
func (iter *Iterator) Next() (*Node, error)
type Language ¶
type Language struct {
// contains filtered or unexported fields
}
Language defines how to parse a particular programming language

func NewLanguage ¶
func NewLanguage(ptr unsafe.Pointer) *Language
NewLanguage creates new Language from c pointer

func (*Language) FieldName ¶
func (l *Language) FieldName(idx int) string
func (*Language) SymbolCount ¶
func (l *Language) SymbolCount() uint32
SymbolCount returns the number of distinct field names in the language.

func (*Language) SymbolName ¶
func (l *Language) SymbolName(s Symbol) string
SymbolName returns a node type string for the given Symbol.

func (*Language) SymbolType ¶
func (l *Language) SymbolType(s Symbol) SymbolType
SymbolType returns named, anonymous, or a hidden type for a Symbol.

type Node ¶
type Node struct {
// contains filtered or unexported fields
}
Node represents a single node in the syntax tree It tracks its start and end positions in the source code, as well as its relation to other nodes like its parent, siblings and children.

func
Parse
deprecated
func ParseCtx ¶
func ParseCtx(ctx context.Context, content []byte, lang *Language) (*Node, error)
ParseCtx is a shortcut for parsing bytes of source code, returns root node

func (Node) Child ¶
func (n Node) Child(idx int) *Node
Child returns the node's child at the given index, where zero represents the first child.

func (Node) ChildByFieldName ¶
func (n Node) ChildByFieldName(name string) *Node
ChildByFieldName returns the node's child with the given field name.

func (Node) ChildCount ¶
func (n Node) ChildCount() uint32
ChildCount returns the node's number of children.

func (Node) Content ¶
func (n Node) Content(input []byte) string
Content returns node's source code from input as a string

func (Node) Edit ¶
func (n Node) Edit(i EditInput)
Edit the node to keep it in-sync with source code that has been edited.

func (Node) EndByte ¶
func (n Node) EndByte() uint32
EndByte returns the node's end byte.

func (Node) EndPoint ¶
func (n Node) EndPoint() Point
EndPoint returns the node's end position in terms of rows and columns.

func (Node) Equal ¶
func (n Node) Equal(other *Node) bool
Equal checks if two nodes are identical.

func (Node) FieldNameForChild ¶
func (n Node) FieldNameForChild(idx int) string
FieldNameForChild returns the field name of the child at the given index, or "" if not named.

func (Node) HasChanges ¶
func (n Node) HasChanges() bool
HasChanges checks if a syntax node has been edited.

func (Node) HasError ¶
func (n Node) HasError() bool
HasError check if the node is a syntax error or contains any syntax errors.

func (Node) ID ¶
func (n Node) ID() uintptr
func (Node) IsError ¶
func (n Node) IsError() bool
IsError checks if the node is a syntax error. Syntax errors represent parts of the code that could not be incorporated into a valid syntax tree.

func (Node) IsExtra ¶
func (n Node) IsExtra() bool
IsExtra checks if the node is *extra*. Extra nodes represent things like comments, which are not required the grammar, but can appear anywhere.

func (Node) IsMissing ¶
func (n Node) IsMissing() bool
IsMissing checks if the node is *missing*. Missing nodes are inserted by the parser in order to recover from certain kinds of syntax errors.

func (Node) IsNamed ¶
func (n Node) IsNamed() bool
IsNamed checks if the node is *named*. Named nodes correspond to named rules in the grammar, whereas *anonymous* nodes correspond to string literals in the grammar.

func (Node) IsNull ¶
func (n Node) IsNull() bool
IsNull checks if the node is null.

func (Node) NamedChild ¶
func (n Node) NamedChild(idx int) *Node
NamedChild returns the node's *named* child at the given index.

func (Node) NamedChildCount ¶
func (n Node) NamedChildCount() uint32
NamedChildCount returns the node's number of *named* children.

func (Node) NamedDescendantForPointRange ¶
func (n Node) NamedDescendantForPointRange(start Point, end Point) *Node
func (Node) NextNamedSibling ¶
func (n Node) NextNamedSibling() *Node
NextNamedSibling returns the node's next *named* sibling.

func (Node) NextSibling ¶
func (n Node) NextSibling() *Node
NextSibling returns the node's next sibling.

func (Node) Parent ¶
func (n Node) Parent() *Node
Parent returns the node's immediate parent.

func (Node) PrevNamedSibling ¶
func (n Node) PrevNamedSibling() *Node
PrevNamedSibling returns the node's previous *named* sibling.

func (Node) PrevSibling ¶
func (n Node) PrevSibling() *Node
PrevSibling returns the node's previous sibling.

func (Node) Range ¶
func (n Node) Range() Range
func (Node) StartByte ¶
func (n Node) StartByte() uint32
StartByte returns the node's start byte.

func (Node) StartPoint ¶
func (n Node) StartPoint() Point
StartPoint returns the node's start position in terms of rows and columns.

func (Node) String ¶
func (n Node) String() string
String returns an S-expression representing the node as a string.

func (Node) Symbol ¶
func (n Node) Symbol() Symbol
Symbol returns the node's type as a Symbol.

func (Node) Type ¶
func (n Node) Type() string
Type returns the node's type as a string.

type Parser ¶
type Parser struct {
// contains filtered or unexported fields
}
Parser produces concrete syntax tree based on source code using Language

func NewParser ¶
func NewParser() *Parser
NewParser creates new Parser

func (*Parser) Close ¶
func (p *Parser) Close()
Close should be called to ensure that all the memory used by the parse is freed.

As the constructor in go-tree-sitter would set this func call through runtime.SetFinalizer, parser.Close() will be called by Go's garbage collector and users would not have to call this manually.

func (*Parser) Debug ¶
func (p *Parser) Debug()
Debug enables debug output to stderr

func (*Parser) OperationLimit ¶
func (p *Parser) OperationLimit() int
OperationLimit returns the duration in microseconds that parsing is allowed to take

func (*Parser)
Parse
deprecated
func (*Parser) ParseCtx ¶
func (p *Parser) ParseCtx(ctx context.Context, oldTree *Tree, content []byte) (*Tree, error)
ParseCtx produces new Tree from content using old tree

func (*Parser) ParseInput ¶
func (p *Parser) ParseInput(oldTree *Tree, input Input) *Tree
ParseInput produces new Tree by reading from a callback defined in input it is useful if your data is stored in specialized data structure as it will avoid copying the data into []bytes and faster access to edited part of the data

func (*Parser) ParseInputCtx ¶
func (p *Parser) ParseInputCtx(ctx context.Context, oldTree *Tree, input Input) (*Tree, error)
ParseInputCtx produces new Tree by reading from a callback defined in input it is useful if your data is stored in specialized data structure as it will avoid copying the data into []bytes and faster access to edited part of the data

func (*Parser) Reset ¶
func (p *Parser) Reset()
Reset causes the parser to parse from scratch on the next call to parse, instead of resuming so that it sees the changes to the beginning of the source code.

func (*Parser) SetIncludedRanges ¶
func (p *Parser) SetIncludedRanges(ranges []Range)
SetIncludedRanges sets text ranges of a file

func (*Parser) SetLanguage ¶
func (p *Parser) SetLanguage(lang *Language)
SetLanguage assignes Language to a parser

func (*Parser) SetOperationLimit ¶
func (p *Parser) SetOperationLimit(limit int)
SetOperationLimit limits the maximum duration in microseconds that parsing should be allowed to take before halting

type Point ¶
type Point struct {
Row    uint32
Column uint32
}
type Quantifier ¶
type Quantifier int
type Query ¶
type Query struct {
// contains filtered or unexported fields
}
Query API

func NewQuery ¶
func NewQuery(pattern []byte, lang *Language) (*Query, error)
NewQuery creates a query by specifying a string containing one or more patterns. In case of error returns QueryError.

func (*Query) CaptureCount ¶
func (q *Query) CaptureCount() uint32
func (*Query) CaptureNameForId ¶
func (q *Query) CaptureNameForId(id uint32) string
func (*Query) CaptureQuantifierForId ¶
func (q *Query) CaptureQuantifierForId(id uint32, captureId uint32) Quantifier
func (*Query) Close ¶
func (q *Query) Close()
Close should be called to ensure that all the memory used by the query is freed.

As the constructor in go-tree-sitter would set this func call through runtime.SetFinalizer, parser.Close() will be called by Go's garbage collector and users would not have to call this manually.

func (*Query) PatternCount ¶
func (q *Query) PatternCount() uint32
func (*Query) PredicatesForPattern ¶
func (q *Query) PredicatesForPattern(patternIndex uint32) [][]QueryPredicateStep
func (*Query) StringCount ¶
func (q *Query) StringCount() uint32
func (*Query) StringValueForId ¶
func (q *Query) StringValueForId(id uint32) string
type QueryCapture ¶
type QueryCapture struct {
Index uint32
Node  *Node
}
QueryCapture is a captured node by a query with an index

type QueryCursor ¶
type QueryCursor struct {
// contains filtered or unexported fields
}
QueryCursor carries the state needed for processing the queries.

func NewQueryCursor ¶
func NewQueryCursor() *QueryCursor
NewQueryCursor creates a query cursor.

func (*QueryCursor) Close ¶
func (qc *QueryCursor) Close()
Close should be called to ensure that all the memory used by the query cursor is freed.

As the constructor in go-tree-sitter would set this func call through runtime.SetFinalizer, parser.Close() will be called by Go's garbage collector and users would not have to call this manually.

func (*QueryCursor) Exec ¶
func (qc *QueryCursor) Exec(q *Query, n *Node)
Exec executes the query on a given syntax node.

func (*QueryCursor) FilterPredicates ¶
func (qc *QueryCursor) FilterPredicates(m *QueryMatch, input []byte) *QueryMatch
func (*QueryCursor) NextCapture ¶
func (qc *QueryCursor) NextCapture() (*QueryMatch, uint32, bool)
func (*QueryCursor) NextMatch ¶
func (qc *QueryCursor) NextMatch() (*QueryMatch, bool)
NextMatch iterates over matches. This function will return (nil, false) when there are no more matches. Otherwise, it will populate the QueryMatch with data about which pattern matched and which nodes were captured.

func (*QueryCursor) SetPointRange ¶
func (qc *QueryCursor) SetPointRange(startPoint Point, endPoint Point)
type QueryError ¶
type QueryError struct {
Offset  uint32
Type    QueryErrorType
Message string
}
QueryError - if there is an error in the query, then the Offset argument will be set to the byte offset of the error, and the Type argument will be set to a value that indicates the type of error.

func (*QueryError) Error ¶
func (qe *QueryError) Error() string
type QueryErrorType ¶
type QueryErrorType int
QueryErrorType - value that indicates the type of QueryError.

const (
QueryErrorNone QueryErrorType = iota
QueryErrorSyntax
QueryErrorNodeType
QueryErrorField
QueryErrorCapture
QueryErrorStructure
QueryErrorLanguage
)
type QueryMatch ¶
type QueryMatch struct {
ID           uint32
PatternIndex uint16
Captures     []QueryCapture
}
QueryMatch - you can then iterate over the matches.

type QueryPredicateStep ¶
type QueryPredicateStep struct {
Type    QueryPredicateStepType
ValueId uint32
}
type QueryPredicateStepType ¶
type QueryPredicateStepType int
const (
QueryPredicateStepTypeDone QueryPredicateStepType = iota
QueryPredicateStepTypeCapture
QueryPredicateStepTypeString
)
type Range ¶
type Range struct {
StartPoint Point
EndPoint   Point
StartByte  uint32
EndByte    uint32
}
type ReadFunc ¶
type ReadFunc func(offset uint32, position Point) []byte
ReadFunc is a function to retrieve a chunk of text at a given byte offset and (row, column) position it should return nil to indicate the end of the document

type Symbol ¶
type Symbol = C.TSSymbol
type SymbolType ¶
type SymbolType int
const (
SymbolTypeRegular SymbolType = iota
SymbolTypeAnonymous
SymbolTypeAuxiliary
)
func (SymbolType) String ¶
func (t SymbolType) String() string
type Tree ¶
type Tree struct {
*BaseTree
// contains filtered or unexported fields
}
Tree represents the syntax tree of an entire source code file Note: Tree instances are not thread safe; you must copy a tree if you want to use it on multiple threads simultaneously.

func (*Tree) Copy ¶
func (t *Tree) Copy() *Tree
Copy returns a new copy of a tree

func (*Tree) Edit ¶
func (t *Tree) Edit(i EditInput)
Edit the syntax tree to keep it in sync with source code that has been edited.

func (*Tree) RootNode ¶
func (t *Tree) RootNode() *Node
RootNode returns root node of a tree

type TreeCursor ¶
type TreeCursor struct {
// contains filtered or unexported fields
}
TreeCursor allows you to walk a syntax tree more efficiently than is possible using the `Node` functions. It is a mutable object that is always on a certain syntax node, and can be moved imperatively to different nodes.

func NewTreeCursor ¶
func NewTreeCursor(n *Node) *TreeCursor
NewTreeCursor creates a new tree cursor starting from the given node.

func (*TreeCursor) Close ¶
func (c *TreeCursor) Close()
Close should be called to ensure that all the memory used by the tree cursor is freed.

As the constructor in go-tree-sitter would set this func call through runtime.SetFinalizer, parser.Close() will be called by Go's garbage collector and users would not have to call this manually.

func (*TreeCursor) CurrentFieldName ¶
func (c *TreeCursor) CurrentFieldName() string
CurrentFieldName gets the field name of the tree cursor's current node.

This returns empty string if the current node doesn't have a field.

func (*TreeCursor) CurrentNode ¶
func (c *TreeCursor) CurrentNode() *Node
CurrentNode of the tree cursor.

func (*TreeCursor) GoToFirstChild ¶
func (c *TreeCursor) GoToFirstChild() bool
GoToFirstChild moves the cursor to the first child of its current node.

This returns `true` if the cursor successfully moved, and returns `false` if there were no children.

func (*TreeCursor) GoToFirstChildForByte ¶
func (c *TreeCursor) GoToFirstChildForByte(b uint32) int64
GoToFirstChildForByte moves the cursor to the first child of its current node that extends beyond the given byte offset.

This returns the index of the child node if one was found, and returns -1 if no such child was found.

func (*TreeCursor) GoToNextSibling ¶
func (c *TreeCursor) GoToNextSibling() bool
GoToNextSibling moves the cursor to the next sibling of its current node.

This returns `true` if the cursor successfully moved, and returns `false` if there was no next sibling node.

func (*TreeCursor) GoToParent ¶
func (c *TreeCursor) GoToParent() bool
GoToParent moves the cursor to the parent of its current node.

This returns `true` if the cursor successfully moved, and returns `false` if there was no parent node (the cursor was already on the root node).

func (*TreeCursor) Reset ¶
func (c *TreeCursor) Reset(n *Node)
Reset re-initializes a tree cursor to start at a different node.