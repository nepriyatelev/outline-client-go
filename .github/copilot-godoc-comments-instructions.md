# Go Doc Comments Instructions

Doc comments are comments that appear immediately before top-level package, const, func, type, and var declarations with no intervening newlines. Every exported (capitalized) name should have a doc comment.

## Packages

Every package should have a package comment introducing the package. It provides information relevant to the package as a whole and generally sets expectations for the package.

### Package Comments

```go
// Package path implements utility routines for manipulating slash-separated
// paths.
//
// The path package should only be used for paths separated by forward
// slashes, such as the paths in URLs. This package does not deal with
// Windows paths with drive letters or backslashes; to manipulate
// operating system paths, use the [path/filepath] package.
package path
```

### Command Packages

A package comment for a command (package main) describes the behavior of the program:

```go
/*
Gofmt formats Go programs.

It uses tabs for indentation and blanks for alignment.
Alignment assumes that an editor is using a fixed-width font.
*/
package main
```

## Types

A type's doc comment should explain what each instance of that type represents or provides.

### Simple Type Documentation

```go
package zip

// A Reader serves content from a ZIP archive.
type Reader struct {
	...
}
```

### Concurrency Guarantees

If a type provides concurrency guarantees, document them:

```go
package regexp

// Regexp is the representation of a compiled regular expression.
// A Regexp is safe for concurrent use by multiple goroutines,
// except for configuration methods, such as Longest.
type Regexp struct {
	...
}
```

### Zero Value Documentation

```go
package bytes

// A Buffer is a variable-sized buffer of bytes with Read and Write methods.
// The zero value for Buffer is an empty buffer ready to use.
type Buffer struct {
	...
}
```

### Struct Fields

Document the meaning of each exported field:

```go
package io

// A LimitedReader reads from R but limits the amount of
// data returned to just N bytes. Each call to Read
// updates N to reflect the new amount remaining.
// Read returns EOF when N <= 0.
type LimitedReader struct {
	R Reader // underlying reader
	N int64  // max bytes remaining
}
```

## Functions and Methods

A function's doc comment should explain what the function returns or what it does.

### Basic Function Documentation

```go
package strconv

// Quote returns a double-quoted Go string literal representing s.
// The returned string uses Go escape sequences (\t, \n, \xFF, \u0100)
// for control characters and non-printable characters as defined by IsPrint.
func Quote(s string) string {
	...
}
```

### Functions with Side Effects

```go
package os

// Exit causes the current program to exit with the given status code.
// Conventionally, code zero indicates success, non-zero an error.
// The program terminates immediately; deferred functions are not run.
//
// For portability, the status code should be in the range [0, 125].
func Exit(code int) {
	...
}
```

### Boolean Return Values

Use the phrase "reports whether" for functions that return a boolean:

```go
package strings

// HasPrefix reports whether the string s begins with prefix.
func HasPrefix(s, prefix string) bool
```

### Multiple Return Values

```go
package io

// Copy copies from src to dst until either EOF is reached
// on src or an error occurs. It returns the total number of bytes
// written and the first error encountered while copying, if any.
//
// A successful Copy returns err == nil, not err == EOF.
// Because Copy is defined to read from src until EOF, it does
// not treat an EOF from Read as an error to be reported.
func Copy(dst Writer, src Reader) (n int64, err error) {
	...
}
```

### Special Cases

Document special cases that are important:

```go
package math

// Sqrt returns the square root of x.
//
// Special cases are:
//
// Sqrt(+Inf) = +Inf
// Sqrt(±0) = ±0
// Sqrt(x < 0) = NaN
// Sqrt(NaN) = NaN
func Sqrt(x float64) float64 {
	...
}
```

### Performance Characteristics

Document asymptotic bounds when important:

```go
package sort

// Sort sorts data in ascending order as determined by the Less method.
// It makes one call to data.Len to determine n and O(n*log(n)) calls to
// data.Less and data.Swap. The sort is not guaranteed to be stable.
func Sort(data Interface) {
	...
}
```

### Concurrency Safety for Methods

```go
package sql

// Close returns the connection to the connection pool.
// All operations after a Close will return with ErrConnDone.
// Close is safe to call concurrently with other operations and will
// block until all other operations finish. It may be useful to first
// cancel any used context and then call Close directly after.
func (c *Conn) Close() error {
	...
}
```

## Constants

A doc comment can introduce a group of related constants:

```go
package scanner

// The result of Scan is one of these tokens or a Unicode character.
const (
	EOF = -(iota + 1)
	Ident
	Int
	Float
	Char
	...
)
```

Ungrouped constants warrant a full doc comment:

```go
package unicode

// Version is the Unicode edition from which the tables are derived.
const Version = "13.0.0"
```

## Variables

The conventions for variables are the same as those for constants:

```go
package fs

// Generic file system errors.
// Errors returned by file systems can be tested against these errors
// using errors.Is.
var (
	ErrInvalid   = errInvalid()   // "invalid argument"
	ErrPermission = errPermission() // "permission denied"
	ErrExist     = errExist()     // "file already exists"
	ErrNotExist  = errNotExist()  // "file does not exist"
	ErrClosed    = errClosed()    // "file already closed"
)
```

## Comment Syntax

### Paragraphs

A paragraph is a span of unindented non-blank lines. Blank lines separate paragraphs.

### Headings

A heading begins with `#` and then a space and the heading text. Headings must be unindented and set off from adjacent text by blank lines:

```go
// Package strconv implements conversions to and from string representations
// of basic data types.
//
// # Numeric Conversions
//
// The most common numeric conversions are [Atoi] (string to int) and [Itoa] (int to string).
```

### Links

Text links in square brackets can reference URLs:

```go
// Package json implements encoding and decoding of JSON as defined in
// [RFC 7159]. The mapping between JSON and Go values is described
// in the documentation for the Marshal and Unmarshal functions.
//
// For an introduction to this package, see the article
// "[JSON and Go]".
//
// [RFC 7159]: https://tools.ietf.org/html/rfc7159
// [JSON and Go]: https://golang.org/doc/articles/json_and_go.html
```

### Doc Links

Doc links reference exported identifiers in the current package or other packages:

```go
// ReadFrom reads data from r until EOF and appends it to the buffer,
// growing the buffer as needed. The return value n is the number of bytes read.
// Any error except [io.EOF] encountered during the read is also returned.
// If the buffer becomes too large, ReadFrom will panic with [ErrTooLarge].
func (b *Buffer) ReadFrom(r io.Reader) (n int64, err error) {
	...
}
```

For other packages, use `[pkg.Name]` or `[pkg.Name.Method]`.

### Lists

Bullet lists:

```go
// PublicSuffixList provides the public suffix of a domain. For example:
// - the public suffix of "example.com" is "com",
// - the public suffix of "foo1.foo2.foo3.co.uk" is "co.uk", and
// - the public suffix of "bar.pvt.k12.ma.us" is "pvt.k12.ma.us".
```

Numbered lists:

```go
// Clean returns the shortest path name equivalent to path
// by purely lexical processing. It applies the following rules
// iteratively until no further processing can be done:
//
// 1. Replace multiple slashes with a single slash.
// 2. Eliminate each . path name element (the current directory).
// 3. Eliminate each inner .. path name element (the parent directory)
// along with the non-.. element that precedes it.
// 4. Eliminate .. elements that begin a rooted path.
func Clean(path string) string {
	...
}
```

### Code Blocks

Preformatted text (code blocks) is indented. Indented lines that don't start a list are rendered as code:

```go
// On the wire, the JSON will look something like this:
//
// {
// 	"kind":"MyAPIObject",
// 	"apiVersion":"v1",
// 	"myPlugin": {
// 		"kind":"PluginA",
// 		"aOption":"foo",
// 	},
// }
```

### Notes

Notes of the form `MARKER(uid): body`. Common markers are `TODO`, `BUG`, `FIXME`:

```go
// TODO(user1): refactor to use standard library context
// BUG(user2): not cleaned up
var ctx context.Context
```

### Deprecations

Paragraphs starting with "Deprecated: " are treated as deprecation notices:

```go
// Package rc4 implements the RC4 stream cipher.
//
// Deprecated: RC4 is cryptographically broken and should not be used
// except for compatibility with legacy systems.
//
// This package is frozen and no new functionality will be added.
package rc4

// Reset zeros the key data and makes the Cipher unusable.
//
// Deprecated: Reset can't guarantee that the key will be entirely removed from
// the process's memory.
func (c *Cipher) Reset()
```

## Rules

- Doc comments appear immediately before declarations with no intervening blank lines
- Comments must be complete sentences, starting with the identifier name
- Comments should explain "why" and "what", not implementation details
- Every exported name should have a doc comment
- For multi-file packages, doc comment should be in one file only
- Use proper sentence structure and punctuation
- Use code blocks for examples
- Use doc links `[Package]`, `[Type]`, `[Function]` for references
- Mark deprecated code with "Deprecated: ..."
- For functions: explain return values and side effects
- For types: explain what instances represent and concurrency guarantees
- For constants/vars: explain the purpose
