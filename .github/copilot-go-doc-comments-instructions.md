# Go Doc Comments Instructions for GitHub Copilot

**Purpose**: Generate high-quality, idiomatic godoc comments that appear in pkg.go.dev and `go doc` output.

## 🎯 AI Documentation Generation Priorities

**ALWAYS follow these in order:**

1. **ALWAYS** document ONLY exported (capitalized) identifiers
2. **ALWAYS** place doc comment immediately before declaration (no blank lines)
3. **ALWAYS** start with the identifier name being documented
4. **ALWAYS** use complete sentences with proper punctuation
5. **ALWAYS** use semantic line breaks (one thought per line for better diffs)
6. **ALWAYS** explain "what" and "why", not "how" (implementation details)
7. **ALWAYS** use "reports whether" for boolean-returning functions
8. **ALWAYS** document concurrency safety for types/methods
9. **ALWAYS** use doc links `[Type]`, `[Function]` for cross-references
10. **PREFER** including special cases, edge cases, error conditions
11. **INCLUDE** zero value behavior for types
12. **AVOID** blank lines between comment and code
13. **AVOID** articles (a, an, the) at the start instead of identifier name
14. **AVOID** documenting unexported identifiers (not in godoc)
15. **AVOID** comments that just repeat the code

---

## Doc Comments vs Regular Comments

### ✅ Doc Comment (appears in godoc)

```go
// Add returns the sum of a and b.
func Add(a, b int) int {
    return a + b
}
```

**Usage**: Placed immediately before exported declarations. Visible in `go doc`, `godoc`, pkg.go.dev.

### ❌ Regular Comment (NOT in godoc)

```go
// This is an internal implementation note
func Add(a, b int) int {
    return a + b
}
```

**Rule**: Only comments immediately before exported declarations become doc comments.

### ❌ Unexported (NOT documented)

```go
// This comment is ignored (function is unexported)
func add(a, b int) int {
    return a + b
}
```

**Rule**: Unexported identifiers (lowercase first letter) do NOT generate documentation.

---

## Semantic Line Breaks

Doc comments should use semantic line breaks (one complete thought per line). This makes diffs clearer:

### ❌ Without Semantic Line Breaks (hard to diff)

```go
// Package gob manages streams of gobs - binary values exchanged between an Encoder (transmitter) and a Decoder (receiver). A typical use is transporting arguments and results of method calls (an RPC system).
```

### ✅ With Semantic Line Breaks (clear diffs)

```go
// Package gob manages streams of gobs - binary values exchanged between
// an Encoder (transmitter) and a Decoder (receiver).
// A typical use is transporting arguments and results of method calls
// (an RPC system).
```

---

## Official Go Doc Comments Standards

### Packages

Every package should have a package comment introducing the package:

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

For command packages (package main):

```go
/*
Gofmt formats Go programs.

It uses tabs for indentation and blanks for alignment.
Alignment assumes that an editor is using a fixed-width font.
*/
package main
```

### Types

Document what each instance represents or provides:

```go
// A Reader serves content from a ZIP archive.
type Reader struct {
    ...
}
```

Document concurrency guarantees:

```go
// Regexp is the representation of a compiled regular expression.
// A Regexp is safe for concurrent use by multiple goroutines,
// except for configuration methods, such as Longest.
type Regexp struct {
    ...
}
```

Document zero value behavior:

```go
// A Buffer is a variable-sized buffer of bytes with Read and Write methods.
// The zero value for Buffer is an empty buffer ready to use.
type Buffer struct {
    ...
}
```

Document struct fields:

```go
// A LimitedReader reads from R but limits the amount of
// data returned to just N bytes. Each call to Read
// updates N to reflect the new amount remaining.
// Read returns EOF when N <= 0.
type LimitedReader struct {
    R Reader // underlying reader
    N int64  // max bytes remaining
}
```

### Functions and Methods

Document what the function returns or what it does:

```go
// Quote returns a double-quoted Go string literal representing s.
// The returned string uses Go escape sequences (\t, \n, \xFF, \u0100)
// for control characters and non-printable characters as defined by IsPrint.
func Quote(s string) string {
    ...
}
```

For functions with side effects:

```go
// Exit causes the current program to exit with the given status code.
// Conventionally, code zero indicates success, non-zero an error.
// The program terminates immediately; deferred functions are not run.
//
// For portability, the status code should be in the range [0, 125].
func Exit(code int) {
    ...
}
```

For boolean functions, use "reports whether":

```go
// HasPrefix reports whether the string s begins with prefix.
func HasPrefix(s, prefix string) bool
```

For multiple return values:

```go
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

For special cases:

```go
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

For performance characteristics:

```go
// Sort sorts data in ascending order as determined by the Less method.
// It makes one call to data.Len to determine n and O(n*log(n)) calls to
// data.Less and data.Swap. The sort is not guaranteed to be stable.
func Sort(data Interface) {
    ...
}
```

For concurrency-safe methods:

```go
// Close returns the connection to the connection pool.
// All operations after a Close will return with ErrConnDone.
// Close is safe to call concurrently with other operations and will
// block until all other operations finish. It may be useful to first
// cancel any used context and then call Close directly after.
func (c *Conn) Close() error {
    ...
}
```

### Interfaces

Document what the interface represents:

```go
// Reader is the interface that wraps the basic Read method.
//
// Read reads up to len(p) bytes into p. It returns the number
// of bytes read (0 <= n <= len(p)) and any error encountered.
// Even if Read returns n < len(p), it may use all of p as scratch
// space during the call.
type Reader interface {
    Read(p []byte) (n int, err error)
}
```

### Constants

For grouped constants:

```go
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

For ungrouped constants:

```go
// Version is the Unicode edition from which the tables are derived.
const Version = "13.0.0"
```

### Variables

Same conventions as constants:

```go
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

### Sentinel Errors

```go
// Common errors returned by the package.
var (
    // ErrNotFound is returned when the requested resource does not exist.
    ErrNotFound = errors.New("resource not found")

    // ErrInvalidInput is returned when input validation fails.
    ErrInvalidInput = errors.New("invalid input")

    // ErrTimeout is returned when an operation times out.
    ErrTimeout = errors.New("operation timed out")
)
```

---

## Comment Syntax

### Paragraphs

A paragraph is unindented non-blank lines. Blank lines separate paragraphs:

```go
// First paragraph.
//
// Second paragraph.
```

### Headings

Headings begin with `#` and a space. Must be unindented and set off by blank lines:

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

Doc links reference exported identifiers in current or other packages:

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

Preformatted text is indented:

```go
// On the wire, the JSON will look something like this:
//
// {
//     "kind":"MyAPIObject",
//     "apiVersion":"v1",
//     "myPlugin": {
//         "kind":"PluginA",
//         "aOption":"foo",
//     },
// }
```

### Notes

Use form `MARKER(uid): body`. Common markers: `TODO`, `BUG`, `FIXME`:

```go
// TODO(user1): refactor to use standard library context
// BUG(user2): not cleaned up
var ctx context.Context
```

### Deprecations

Must start with "Deprecated: " (note the capital D):

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

### Example Functions

Example functions generate documentation:

```go
package mypackage_test

import "fmt"

func ExampleAdd() {
    fmt.Println(Add(2, 3))
    // Output: 5
}

func ExampleSort() {
    fmt.Println(Sort([]int{3, 1, 2}))
    // Output: [1 2 3]
}

// Unordered output
func ExampleRange() {
    for _, v := range []int{1, 2, 3} {
        fmt.Println(v)
    }
    // Unordered output:
    // 1
    // 2
    // 3
}
```

---

## ❌ Common Documentation Mistakes

### Anti-Pattern 1: Blank Line Between Comment and Declaration

```go
// ❌ BAD - blank line breaks the doc comment association
// Add returns the sum of a and b.

func Add(a, b int) int {
    return a + b
}

// ✅ GOOD - no blank line
// Add returns the sum of a and b.
func Add(a, b int) int {
    return a + b
}
```

### Anti-Pattern 2: Not Starting with Identifier Name

```go
// ❌ BAD - doesn't start with function name
// Returns the sum of two integers.
func Add(a, b int) int {
    return a + b
}

// ✅ GOOD - starts with function name
// Add returns the sum of two integers.
func Add(a, b int) int {
    return a + b
}
```

### Anti-Pattern 3: Starting with Article Instead of Name

```go
// ❌ BAD - starts with "The" instead of identifier
// The Add function returns the sum.
func Add(a, b int) int {
    return a + b
}

// ✅ GOOD - starts with identifier
// Add returns the sum of a and b.
func Add(a, b int) int {
    return a + b
}
```

### Anti-Pattern 4: Incomplete Sentences

```go
// ❌ BAD - not a complete sentence, no period
// Add two numbers
func Add(a, b int) int {
    return a + b
}

// ✅ GOOD - complete sentence with period
// Add returns the sum of a and b.
func Add(a, b int) int {
    return a + b
}
```

### Anti-Pattern 5: Explaining Implementation Instead of Purpose

```go
// ❌ BAD - explains how (implementation details)
// ParseJSON loops through bytes and uses json.Unmarshal to convert them.
func ParseJSON(data []byte) (map[string]interface{}, error) {
    var result map[string]interface{}
    err := json.Unmarshal(data, &result)
    return result, err
}

// ✅ GOOD - explains what and why
// ParseJSON decodes JSON-encoded data and returns a map representation.
// It returns an error if the data is not valid JSON.
func ParseJSON(data []byte) (map[string]interface{}, error) {
    var result map[string]interface{}
    err := json.Unmarshal(data, &result)
    return result, err
}
```

### Anti-Pattern 6: Missing Error Documentation

```go
// ❌ BAD - doesn't mention error return
// ReadConfig reads configuration from file.
func ReadConfig(path string) (*Config, error) {
    // ...
}

// ✅ GOOD - documents error conditions
// ReadConfig reads and parses configuration from the specified file.
// It returns an error if the file doesn't exist or contains invalid data.
func ReadConfig(path string) (*Config, error) {
    // ...
}
```

### Anti-Pattern 7: Wrong Phrasing for Boolean Returns

```go
// ❌ BAD - doesn't use "reports whether"
// IsValid checks if the string is valid.
func IsValid(s string) bool {
    return len(s) > 0
}

// ✅ GOOD - uses "reports whether"
// IsValid reports whether s is a valid non-empty string.
func IsValid(s string) bool {
    return len(s) > 0
}
```

### Anti-Pattern 8: Documenting Unexported Identifiers

```go
// ❌ BAD - unexported, not in godoc (wasted effort)
// unexportedHelper does something.
func unexportedHelper() {}

// ✅ GOOD - only document exported
// ExportedHelper does something.
func ExportedHelper() {}
```

### Anti-Pattern 9: Missing Concurrency Information

```go
// ❌ BAD - doesn't mention thread safety
// Cache stores key-value pairs in memory.
type Cache struct {
    sync.RWMutex
    data map[string]interface{}
}

// ✅ GOOD - documents concurrency safety
// Cache stores key-value pairs in memory.
// Cache is safe for concurrent use by multiple goroutines.
type Cache struct {
    sync.RWMutex
    data map[string]interface{}
}
```

### Anti-Pattern 10: Not Using Doc Links

```go
// ❌ BAD - references without links
// Parse parses data using the rules from ParseConfig.
// Returns ParseError if parsing fails.
func Parse(data []byte) error {
    // ...
}

// ✅ GOOD - uses doc links
// Parse parses data using the rules from [ParseConfig].
// Returns [ParseError] if parsing fails.
func Parse(data []byte) error {
    // ...
}
```

---

## Common Verb Patterns by Function Type

### Constructor Functions (New*, Make*, Create*)

```go
// NewClient creates and returns a new HTTP client.
// NewCache creates a cache with the specified capacity.
// CreateDatabase creates a new database connection.
```

### Getters (Get*)

```go
// GetUser retrieves the user with the specified ID.
// GetConfig returns the current configuration.
```

### Setters (Set*, Update*, Modify*)

```go
// SetTimeout sets the timeout duration for requests.
// UpdateConfig updates the configuration with new values.
```

### Boolean Functions (Is*, Has*, Can*, Should*, Reports*)

```go
// IsValid reports whether the input is valid.
// HasPrefix reports whether s begins with prefix.
// CanExecute reports whether the operation can be executed.
```

### Action Functions (Execute*, Process*, Handle*, Parse*, Format*)

```go
// Execute runs the specified command.
// Process processes the incoming data.
// HandleRequest handles the HTTP request.
// ParseJSON parses JSON data into a map.
// FormatDate formats the date according to the layout.
```

### Conversion Functions (To*, From*, As*)

```go
// ToString converts the value to a string representation.
// FromJSON creates an object from JSON data.
// AsBytes returns the data as a byte slice.
```

### Validation Functions (Validate*, Check*, Verify*)

```go
// Validate validates the input data and returns an error if invalid.
// CheckPermissions checks if the user has required permissions.
// VerifySignature verifies the cryptographic signature.
```

### Registration Functions (Register*, Add*, Remove*)

```go
// Register registers a new handler for the given pattern.
// AddMiddleware adds middleware to the processing chain.
// RemoveListener removes the specified event listener.
```

---

## gofmt and Doc Comments

`gofmt` automatically reformats doc comments to canonical form:

```bash
# Reformat all doc comments in a file
gofmt -w file.go

# View formatted doc comments
gofmt -d file.go
```

**What gofmt does:**
- Normalizes spacing and indentation
- Ensures proper line wrapping
- Standardizes bullet list formatting
- Validates comment structure

---

## ✅ Quick Documentation Checklist

Before committing code, verify:

- [ ] Every exported identifier has a doc comment
- [ ] No unexported identifiers documented (wasted effort)
- [ ] Comment starts immediately before declaration (no blank line)
- [ ] Comment starts with the identifier name
- [ ] Complete sentences with proper punctuation
- [ ] No articles (a, an, the) at the start
- [ ] Error conditions are documented (if returns error)
- [ ] Boolean functions use "reports whether"
- [ ] Concurrency safety documented for types/methods
- [ ] Special cases and edge cases mentioned
- [ ] Doc links used for cross-references `[Type]`, `[Function]`
- [ ] Zero value behavior documented (for types)
- [ ] Deprecated items marked with "Deprecated: " (capital D)
- [ ] Semantic line breaks used (one thought per line)
- [ ] Examples provided in Example* functions where helpful

---

## Rules Summary

- Doc comments appear immediately before exported declarations, no blank lines
- Comments must be complete sentences starting with the identifier name
- Comments should explain "why" and "what", not implementation details
- Every exported name should have a doc comment
- Unexported identifiers should NOT be documented
- Use proper sentence structure and punctuation
- Use code blocks for examples
- Use doc links `[Package]`, `[Type]`, `[Function]` for references
- Mark deprecated code with "Deprecated: ..." (capital D)
- For functions: explain return values and side effects
- For types: explain what instances represent and concurrency guarantees
- For constants/vars: explain the purpose
- Use semantic line breaks for better diffs

---
**Last Updated**: January 2026