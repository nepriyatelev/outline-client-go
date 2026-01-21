# Go Testing Instructions for GitHub Copilot

**Purpose**: Generate high-quality, idiomatic Go tests that achieve 100% code coverage with clear, readable code.

## 🎯 AI Test Generation Priorities

**ALWAYS follow these in order:**

1. **ALWAYS** use `package mypackage_test` (black-box testing - test public API like external users)
2. **ALWAYS** use table-driven tests with `t.Run()` for multiple test cases
3. **ALWAYS** check and handle errors explicitly
4. **ALWAYS** use `t.Errorf()` for assertions (continue on failure)
5. **ALWAYS** test error paths when function returns `error` type
6. **ALWAYS** use descriptive subtest names that explain the scenario
7. **ALWAYS** use `got`/`want` naming pattern for assertions
8. **ALWAYS** use assert functions from testify/assert or standard comparisons
9. **ALWAYS** test all code paths: happy path, errors, edge cases, boundaries
10. **ALWAYS** mock interfaces using mockery, not concrete types
11. **PREFER** `t.Parallel()` for independent tests (see guidelines below)
12. **PREFER** `t.Errorf()` over `t.Fatalf()` unless setup fails
13. **INCLUDE** at least one edge case per test function
14. **INCLUDE** boundary values (zero, empty, nil, max, min)
15. **AVOID** testing unexported functions directly
16. **AVOID** asserting on exact error message strings (test error types)
17. **AVOID** global state in tests
18. **AVOID** non-deterministic tests (no random data, sleeps, or timing dependencies)

---

## Test Package Convention

### Recommended: Black-Box Testing (STANDARD)

```go
// file_test.go
package mypackage_test  // SEPARATE PACKAGE - test as external user

import (
    "testing"
    "mymodule/mypackage"
)

func TestPublicFunction(t *testing.T) {
    got := mypackage.PublicFunction()
    if got != expected {
        t.Errorf("PublicFunction() = %v; want %v", got, expected)
    }
}
```

**Advantages:**
- ✅ Tests public API like external users would use it
- ✅ Cannot accidentally access unexported symbols
- ✅ Better for integration testing
- ✅ More realistic usage patterns

**Use when:** Testing public exported functions (99% of cases)

### Alternative: White-Box Testing (RARE)

```go
// file_test.go
package mypackage  // SAME PACKAGE - can access unexported

import "testing"

func TestUnexportedHelper(t *testing.T) {
    got := unexportedHelper()  // Direct access to unexported
    if got != expected {
        t.Errorf("unexportedHelper() = %v; want %v", got, expected)
    }
}
```

**Use when:** Testing unexported helper functions (1% of cases)

---

## Official Go Testing Conventions

### Test Function Naming and Structure

Test functions must follow convention `TestXxx(*testing.T)` where Xxx starts with uppercase:

```go
package mypackage_test

import "testing"

func TestAdd(t *testing.T) {
    got := Add(2, 3)
    if got != 5 {
        t.Errorf("Add(2, 3) = %d; want 5", got)
    }
}
```

Test file must end with `_test.go`.

### Table-Driven Tests (REQUIRED)

Use subtests with `t.Run()` for organizing test cases. Field ordering matters for readability:

```go
func TestAdd(t *testing.T) {
    tests := []struct {
        name    string  // First: test name
        a, b    int     // Inputs
        want    int     // Expected output
    }{
        {"positive", 2, 3, 5},
        {"negative", -2, -3, -5},
        {"zero", 0, 0, 0},
        {"max", math.MaxInt, 1, math.MaxInt + 1}, // edge case
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got := Add(tt.a, tt.b)
            if got != tt.want {
                t.Errorf("Add(%d, %d) = %d; want %d", tt.a, tt.b, got, tt.want)
            }
        })
    }
}
```

**Table struct field ordering:**
```go
tests := []struct {
    name    string      // Always first for readability
    input   string      // All inputs
    param1  int
    param2  bool
    want    string      // Expected output
    wantErr bool        // If function returns error
}{}
```

### Error Reporting Methods

```go
t.Error(args...)          // Log error, mark failed, continue
t.Errorf(format, args)    // Formatted log, mark failed, continue (PREFERRED)
t.Fatal(args...)          // Log error, stop test immediately
t.Fatalf(format, args)    // Formatted log, stop test immediately (use for setup only)
t.Fail()                  // Mark failed, continue
t.FailNow()               // Stop immediately
```

**Guidance**: Use `t.Errorf()` for assertions. Use `t.Fatalf()` only for critical setup failures.

### Assert Functions with testify

```go
import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/require"
)

func TestWithAssert(t *testing.T) {
    result := Add(2, 3)
    
    // Use require for critical checks in setup (test stops on failure)
    require.NoError(t, err)
    require.NotNil(t, value)
    
    // Use assert for main assertions (test continues on failure)
    assert.Equal(t, 5, result, "Add(2, 3) should equal 5")
    assert.NotNil(t, result)
    assert.NoError(t, err)
    assert.Error(t, err)
    assert.True(t, condition)
    assert.False(t, condition)
    assert.Contains(t, list, item)
    assert.Empty(t, str)
    assert.NotEmpty(t, str)
}
```

**Principle:**
- Use `require` for setup (if fails, stop test immediately)
- Use `assert` for actual assertions (report all failures)

### Testing with Interfaces and Mockery

Define interfaces in consumer package:

```go
package service

type Repository interface {
    GetUser(id string) (User, error)
    SaveUser(user User) error
}

type UserService struct {
    repo Repository
}

func NewUserService(repo Repository) *UserService {
    return &UserService{repo: repo}
}

func (s *UserService) UpdateUser(id string, name string) error {
    user, err := s.repo.GetUser(id)
    if err != nil {
        return err
    }
    user.Name = name
    return s.repo.SaveUser(user)
}
```

Generate mocks with mockery:

```bash
mockery --name=Repository --output=mocks --outpkg=mocks
```

Test with mockery-generated mocks:

```go
package service_test

import (
    "testing"
    "github.com/stretchr/testify/assert"
    "github.com/stretchr/testify/mock"
    "mymodule/mocks"
    "mymodule/service"
)

func TestUpdateUser(t *testing.T) {
    tests := []struct {
        name      string
        userID    string
        newName   string
        mockSetup func(*mocks.Repository)
        wantErr   bool
    }{
        {
            name:    "successful update",
            userID:  "user123",
            newName: "John Doe",
            mockSetup: func(mockRepo *mocks.Repository) {
                mockRepo.On("GetUser", "user123").Return(
                    service.User{ID: "user123", Name: "Old Name"},
                    nil,
                )
                mockRepo.On("SaveUser", mock.MatchedBy(func(u service.User) bool {
                    return u.ID == "user123" && u.Name == "John Doe"
                })).Return(nil)
            },
            wantErr: false,
        },
        {
            name:    "get user fails",
            userID:  "user456",
            newName: "Jane Doe",
            mockSetup: func(mockRepo *mocks.Repository) {
                mockRepo.On("GetUser", "user456").Return(
                    service.User{},
                    errors.New("user not found"),
                )
            },
            wantErr: true,
        },
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            mockRepo := new(mocks.Repository)
            tt.mockSetup(mockRepo)

            svc := service.NewUserService(mockRepo)
            err := svc.UpdateUser(tt.userID, tt.newName)

            if tt.wantErr {
                assert.Error(t, err)
            } else {
                assert.NoError(t, err)
            }

            mockRepo.AssertExpectations(t)
        })
    }
}
```

### Test Helper Functions

Use `t.Helper()` to mark setup functions:

```go
func TestSomething(t *testing.T) {
    config := createTestConfig(t)
    // use config
}

func createTestConfig(t *testing.T) Config {
    t.Helper() // Attributes failures to the caller
    // setup logic
    return config
}
```

### Parallel Tests

Use `t.Parallel()` for independent tests:

```go
func TestGroupedParallel(t *testing.T) {
    for _, tc := range tests {
        // Go < 1.22: capture range variable for closure
        tc := tc
        t.Run(tc.Name, func(t *testing.T) {
            t.Parallel() // Run in parallel with other parallel tests
            // test body
        })
    }
}

// Go 1.22+: no capture needed - automatic by the language
func TestParallelGo122(t *testing.T) {
    for _, tc := range tests {
        t.Run(tc.Name, func(t *testing.T) {
            t.Parallel()  // Works correctly without capture
            // test body
        })
    }
}
```

**When to use `t.Parallel()`:**
- ✅ Tests are independent (no shared state)
- ✅ No TestMain setup/teardown
- ✅ No global variables modified
- ✅ No race conditions possible
- ✅ Each test uses t.TempDir() or private state

**When NOT to use `t.Parallel()`:**
- ❌ Tests depend on TestMain setup
- ❌ Shared mutable state
- ❌ Order of execution matters
- ❌ Concurrent execution would cause issues

### TestMain Pattern

Global setup/teardown for entire package:

```go
func TestMain(m *testing.M) {
    // Global setup
    setup()
    
    // Run tests
    code := m.Run()
    
    // Global teardown
    teardown()
    
    os.Exit(code)
}

func setup() {
    // Initialize database, create temp files, etc.
}

func teardown() {
    // Clean up resources
}
```

**Use when:**
- Tests need shared setup/teardown
- Database initialization required
- Temp directories needed for all tests
- Environmental setup needed

### Skipping Tests

Mark expensive tests to skip with `-short` flag:

```go
func TestExpensive(t *testing.T) {
    if testing.Short() {
        t.Skip("skipping expensive test in short mode")
    }
    // expensive test body
}

// Run with: go test -short
```

### Resource Management

```go
func TestWithTempDir(t *testing.T) {
    dir := t.TempDir()  // Auto-cleaned after test
    // use dir

    t.Cleanup(func() {
        // cleanup code
        // Called in LIFO order (last registered, first called)
    })
}
```

---

## 100% Code Coverage Checklist

Ensure tests cover:

- [ ] **Happy path** - normal successful execution
- [ ] **All error paths** - every `return err` statement
- [ ] **Edge cases** - zero, empty, nil, max/min values
- [ ] **All if/else branches** - both true and false paths
- [ ] **All loops** - empty, single item, multiple items
- [ ] **Boundary conditions** - first/last elements
- [ ] **Concurrent access** (if applicable)
- [ ] **Timeout/context deadlines** (if applicable)
- [ ] **Interface implementations** - all methods tested
- [ ] **Error types** - specific error conditions
- [ ] **Special cases** documented in function

---

## Error Handling in Tests

### Correct: Check error existence

```go
// ✅ GOOD - check if error occurred
if err != nil {
    t.Errorf("unexpected error: %v", err)
}
```

### Correct: Check error type

```go
// ✅ GOOD - check specific error type
if !errors.Is(err, ErrNotFound) {
    t.Errorf("expected ErrNotFound, got %v", err)
}

var validErr *ValidationError
if !errors.As(err, &validErr) {
    t.Error("expected ValidationError type")
}
```

### Incorrect: Check error message

```go
// ❌ BAD - brittle, message may change
if err.Error() != "specific message" {
    t.Error("wrong error message")
}
```

---

## ❌ Common Testing Mistakes (Anti-Patterns)

### Mistake 1: Testing Unexported Functions Directly

```go
// ❌ BAD - white-box testing of internal helpers
package mypackage  // Same package

func Test_internalHelper(t *testing.T) {
    result := internalHelper()  // Testing unexported function
    assert.Equal(t, expected, result)
}

// ✅ GOOD - black-box testing of public API
package mypackage_test  // Separate package

func TestPublicFunction(t *testing.T) {
    result := mypackage.PublicFunction()  // Uses internalHelper internally
    assert.Equal(t, expected, result)
}
```

### Mistake 2: Using Concrete Types Instead of Interfaces

```go
// ❌ BAD - hard to test, cannot mock
type Service struct {
    db *sql.DB  // Concrete type
}

func NewService(db *sql.DB) *Service {
    return &Service{db: db}
}

// ✅ GOOD - easy to mock
type Database interface {
    Query(query string) ([]Row, error)
    Exec(query string) error
}

type Service struct {
    db Database  // Interface
}

func NewService(db Database) *Service {
    return &Service{db: db}
}
```

### Mistake 3: Non-Deterministic Tests

```go
// ❌ BAD - uses time.Sleep, flaky
func TestProcess(t *testing.T) {
    go processor.Start()
    time.Sleep(100 * time.Millisecond)  // Race condition
    assert.True(t, processor.IsRunning())
}

// ✅ GOOD - uses synchronization
func TestProcess(t *testing.T) {
    done := make(chan bool)
    go func() {
        processor.Start()
        done <- true
    }()
    
    select {
    case <-done:
        assert.True(t, processor.IsRunning())
    case <-time.After(1 * time.Second):
        t.Fatal("timeout waiting for processor")
    }
}
```

### Mistake 4: Ignoring Errors in Test Setup

```go
// ❌ BAD - ignores setup errors
func TestRead(t *testing.T) {
    os.WriteFile("test.txt", []byte("data"), 0644)  // Error ignored
    data, _ := os.ReadFile("test.txt")  // Error ignored
    assert.Equal(t, "data", string(data))
}

// ✅ GOOD - checks all errors
func TestRead(t *testing.T) {
    if err := os.WriteFile("test.txt", []byte("data"), 0644); err != nil {
        t.Fatalf("setup failed: %v", err)
    }
    
    data, err := os.ReadFile("test.txt")
    if err != nil {
        t.Fatalf("ReadFile failed: %v", err)
    }
    
    assert.Equal(t, "data", string(data))
}
```

### Mistake 5: Global State Between Tests

```go
// ❌ BAD - shared global state
var counter int

func TestIncrement(t *testing.T) {
    counter++
    assert.Equal(t, 1, counter)  // Fails if tests run in parallel
}

// ✅ GOOD - isolated state
func TestIncrement(t *testing.T) {
    counter := 0
    counter++
    assert.Equal(t, 1, counter)
}
```

### Mistake 6: Not Using Table-Driven Tests

```go
// ❌ BAD - repeated code
func TestAdd_Positive(t *testing.T) {
    assert.Equal(t, 5, Add(2, 3))
}

func TestAdd_Negative(t *testing.T) {
    assert.Equal(t, -5, Add(-2, -3))
}

func TestAdd_Zero(t *testing.T) {
    assert.Equal(t, 0, Add(0, 0))
}

// ✅ GOOD - table-driven
func TestAdd(t *testing.T) {
    tests := []struct {
        name string
        a, b int
        want int
    }{
        {"positive", 2, 3, 5},
        {"negative", -2, -3, -5},
        {"zero", 0, 0, 0},
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            assert.Equal(t, tt.want, Add(tt.a, tt.b))
        })
    }
}
```

### Mistake 7: Testing with Random Data

```go
// ❌ BAD - non-deterministic
func TestGenerate(t *testing.T) {
    id := GenerateRandomID()
    assert.NotEmpty(t, id)  // Different every run
}

// ✅ GOOD - deterministic with seed
func TestGenerate(t *testing.T) {
    gen := NewGenerator(WithSeed(42))
    id := gen.GenerateID()
    assert.Equal(t, "expected_id_with_seed_42", id)
}
```

### Mistake 8: Not Capturing Loop Variables (Go < 1.22)

```go
// ❌ BAD - loop variable not captured (Go < 1.22)
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        t.Parallel()
        // tt may have wrong value
        assert.Equal(t, tt.want, Process(tt.input))
    })
}

// ✅ GOOD - variable captured
for _, tt := range tests {
    tt := tt  // Capture for Go < 1.22
    t.Run(tt.name, func(t *testing.T) {
        t.Parallel()
        assert.Equal(t, tt.want, Process(tt.input))
    })
}

// ✅ Go 1.22+ - automatic capture
for _, tt := range tests {
    t.Run(tt.name, func(t *testing.T) {
        t.Parallel()  // Works without capture
        assert.Equal(t, tt.want, Process(tt.input))
    })
}
```

---

## Common Test Patterns

### Pattern 1: Simple Function Testing

```go
func TestParse(t *testing.T) {
    tests := []struct {
        name    string
        input   string
        want    int
        wantErr bool
    }{
        {"valid", "123", 123, false},
        {"invalid", "abc", 0, true},
        {"empty", "", 0, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            got, err := Parse(tt.input)
            
            if (err != nil) != tt.wantErr {
                t.Errorf("Parse(%q) error = %v; wantErr %v", tt.input, err, tt.wantErr)
                return
            }
            
            if got != tt.want {
                t.Errorf("Parse(%q) = %d; want %d", tt.input, got, tt.want)
            }
        })
    }
}
```

### Pattern 2: Interface Mocking with mockery

```go
func TestServiceWithMocks(t *testing.T) {
    t.Run("successful operation", func(t *testing.T) {
        mockDB := new(mocks.Database)
        mockDB.On("Query", "SELECT * FROM users").Return([]User{}, nil)
        
        svc := NewService(mockDB)
        result, err := svc.GetUsers()
        
        assert.NoError(t, err)
        assert.NotNil(t, result)
        mockDB.AssertExpectations(t)
    })
}
```

### Pattern 3: Context and Timeout Testing

```go
func TestWithContext(t *testing.T) {
    tests := []struct {
        name    string
        timeout time.Duration
        wantErr bool
    }{
        {"completes", 100 * time.Millisecond, false},
        {"timeout", 1 * time.Millisecond, true},
    }

    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            ctx, cancel := context.WithTimeout(context.Background(), tt.timeout)
            defer cancel()

            err := SlowOperation(ctx)
            if (err != nil) != tt.wantErr {
                t.Errorf("error = %v; wantErr %v", err, tt.wantErr)
            }
        })
    }
}
```

### Pattern 4: Temporary Files and Directories

```go
func TestFileOperations(t *testing.T) {
    dir := t.TempDir()

    file := filepath.Join(dir, "test.txt")
    if err := os.WriteFile(file, []byte("test data"), 0644); err != nil {
        t.Fatalf("setup failed: %v", err)
    }

    data, err := os.ReadFile(file)
    if err != nil {
        t.Errorf("ReadFile failed: %v", err)
    }
    
    assert.Equal(t, "test data", string(data))
}
```

### Pattern 5: Custom Assertion Helpers

```go
// Define helpers for complex assertions
func assertUserEqual(t *testing.T, expected, actual User) {
    t.Helper()  // Attribute failures correctly
    assert.Equal(t, expected.ID, actual.ID, "user ID mismatch")
    assert.Equal(t, expected.Name, actual.Name, "user name mismatch")
    assert.Equal(t, expected.Email, actual.Email, "user email mismatch")
}

// Use in tests
func TestUpdateUser(t *testing.T) {
    expected := User{ID: "123", Name: "John", Email: "john@example.com"}
    actual := updateUser(user)
    assertUserEqual(t, expected, actual)
}
```

### Pattern 6: Benchmark Testing

```go
func BenchmarkSort(b *testing.B) {
    data := generateTestData(1000)
    b.ResetTimer()  // Don't count setup time

    for i := 0; i < b.N; i++ {
        Sort(data)
    }
}

// Benchmark with sub-benchmarks for different sizes
func BenchmarkSort_Sizes(b *testing.B) {
    sizes := []int{10, 100, 1000, 10000}
    
    for _, size := range sizes {
        b.Run(fmt.Sprintf("size=%d", size), func(b *testing.B) {
            data := generateTestData(size)
            b.ResetTimer()
            
            for i := 0; i < b.N; i++ {
                Sort(data)
            }
        })
    }
}

// Parallel benchmark
func BenchmarkParallel(b *testing.B) {
    b.RunParallel(func(pb *testing.PB) {
        for pb.Next() {
            Process()
        }
    })
}
```

### Pattern 7: Concurrency and Race Condition Testing

```go
func TestConcurrentAccess(t *testing.T) {
    t.Run("concurrent writes are safe", func(t *testing.T) {
        cache := NewSafeCache()
        const numGoroutines = 100
        
        var wg sync.WaitGroup
        wg.Add(numGoroutines)
        
        for i := 0; i < numGoroutines; i++ {
            go func(n int) {
                defer wg.Done()
                cache.Set(fmt.Sprintf("key%d", n), n)
            }(i)
        }
        
        wg.Wait()
        
        assert.Equal(t, numGoroutines, cache.Len())
    })
}

// Run with: go test -race
```

### Pattern 8: Fuzz Testing (Go 1.18+)

```go
func FuzzParse(f *testing.F) {
    // Add seed corpus
    f.Add("123")
    f.Add("-456")
    f.Add("0")
    
    f.Fuzz(func(t *testing.T, input string) {
        result, err := Parse(input)
        
        if err == nil {
            // Verify result can be stringified back
            str := fmt.Sprintf("%d", result)
            result2, err2 := Parse(str)
            assert.NoError(t, err2)
            assert.Equal(t, result, result2)
        }
    })
}

// Run with: go test -fuzz=FuzzParse -fuzztime=30s
```

---

## 🤖 AI Test Generation Decision Tree

When generating tests, follow this decision tree:

```
START: Analyze function signature
│
├─ Is function exported (capitalized)?
│  ├─ NO  → SKIP (don't test unexported functions)
│  └─ YES → Continue
│
├─ Does function return error?
│  ├─ YES → MUST include error test cases
│  │        - Test happy path (no error)
│  │        - Test each error path
│  │        - Use errors.Is() or errors.As()
│  └─ NO  → Focus on return value validation
│
├─ Does function accept interface parameter?
│  ├─ YES → MUST use mockery for mocks
│  │        - Generate mock: mockery --name=InterfaceName
│  │        - Use mock.On() for expectations
│  │        - Call mock.AssertExpectations(t)
│  └─ NO  → Use concrete test values
│
├─ Does function accept context.Context?
│  ├─ YES → MUST test timeout/cancellation
│  │        - Test with valid context
│  │        - Test with canceled context
│  │        - Test with deadline exceeded
│  └─ NO  → Skip context tests
│
├─ Does function have side effects?
│  ├─ YES → MUST verify side effects
│  │        - Check state changes
│  │        - Verify mock calls
│  │        - Check file/database modifications
│  └─ NO  → Test return values only
│
├─ Is function concurrent (uses goroutines)?
│  ├─ YES → MUST test concurrency
│  │        - Use sync.WaitGroup
│  │        - Test with race detector
│  │        - Verify thread safety
│  └─ NO  → Skip concurrency tests
│
└─ GENERATE table-driven test with:
    ├─ Happy path (success case)
    ├─ Error paths (all error conditions)
    ├─ Edge cases (empty, nil, zero)
    ├─ Boundary values (min, max, n±1)
    └─ Special cases (documented in comments)
```

## 💡 AI Test Case Generation Hints

### Hint 1: Automatic Test Case Coverage

When generating table-driven tests, **ALWAYS** include these test cases:

**For numeric functions:**
```go
{name: "zero", input: 0, want: ...},
{name: "positive", input: 42, want: ...},
{name: "negative", input: -42, want: ...},
{name: "max int", input: math.MaxInt, want: ...},
{name: "min int", input: math.MinInt, want: ...},
```

**For string functions:**
```go
{name: "empty string", input: "", want: ...},
{name: "single char", input: "a", want: ...},
{name: "multiple chars", input: "hello", want: ...},
{name: "unicode", input: "こんにちは", want: ...},
{name: "special chars", input: "!@#$%", want: ...},
{name: "whitespace", input: "  \n\t  ", want: ...},
```

**For slice/array functions:**
```go
{name: "nil slice", input: nil, want: ...},
{name: "empty slice", input: []T{}, want: ...},
{name: "single element", input: []T{item}, want: ...},
{name: "multiple elements", input: []T{item1, item2, item3}, want: ...},
{name: "large slice", input: make([]T, 1000), want: ...},
```

**For pointer functions:**
```go
{name: "nil pointer", input: nil, wantErr: true},
{name: "valid pointer", input: &value, want: ...},
```

**For boolean functions:**
```go
{name: "returns true", input: validInput, want: true},
{name: "returns false", input: invalidInput, want: false},
```

### Hint 2: Error Path Enumeration

For functions returning `error`, enumerate ALL error paths:

```go
func Example(input string) (Result, error) {
    if input == "" {           // Error path 1
        return Result{}, ErrEmpty
    }
    if !isValid(input) {       // Error path 2
        return Result{}, ErrInvalid
    }
    result, err := process(input)
    if err != nil {            // Error path 3
        return Result{}, fmt.Errorf("process failed: %w", err)
    }
    return result, nil         // Success path
}

// Required test cases:
// 1. input="" → ErrEmpty
// 2. input=invalid → ErrInvalid
// 3. process() fails → wrapped error
// 4. process() succeeds → result, nil
```

### Hint 3: Mock Setup Patterns

When using mockery mocks, follow these patterns:

**Simple mock:**
```go
mockRepo.On("GetUser", "123").Return(user, nil)
```

**Mock with any arguments:**
```go
mockRepo.On("GetUser", mock.Anything).Return(user, nil)
```

**Mock with argument matcher:**
```go
mockRepo.On("SaveUser", mock.MatchedBy(func(u User) bool {
    return u.ID == "123" && u.Name != ""
})).Return(nil)
```

**Mock that returns different values on subsequent calls:**
```go
mockRepo.On("GetUser", "123").Return(user, nil).Once()
mockRepo.On("GetUser", "123").Return(User{}, ErrNotFound).Once()
```

**Mock with callback:**
```go
mockRepo.On("Process").Run(func(args mock.Arguments) {
    // Custom logic
}).Return(nil)
```

### Hint 4: Coverage Formula

To achieve 100% coverage:

```
Total Coverage = (Covered Statements / Total Statements) × 100

Required test cases = 
    Number of branches +
    Number of error paths +
    Number of edge cases +
    Number of special cases

Example:
func Parse(s string) (int, error) {
    if s == "" {          // Branch 1 → test case 1
        return 0, ErrEmpty
    }
    if len(s) > 100 {     // Branch 2 → test case 2
        return 0, ErrTooLong
    }
    n, err := strconv.Atoi(s)
    if err != nil {       // Branch 3 → test case 3
        return 0, ErrInvalid
    }
    return n, nil         // Branch 4 → test case 4
}

Minimum test cases needed: 4
Edge cases to add: empty, too long, invalid format, valid number
Total test cases: 4+ edge cases
```

---

## Running Tests

```bash
go test ./...                    # Test all packages
go test -v ./...                 # Verbose output
go test -run "TestName"          # Run specific test
go test -run "TestName/subname"  # Run specific subtest
go test -short                   # Skip long-running tests
go test -count=1                 # Run without caching
go test -cover                   # Show coverage
go test -coverprofile=c.out ./...  # Generate coverage
go tool cover -html=c.out        # View coverage in browser
go test -race ./...              # Detect race conditions
go test -timeout 30s ./...       # Set timeout
go test -parallel 4 ./...        # Limit parallelism
go test -bench . ./...           # Run benchmarks
go test -bench . -benchmem ./... # Benchmarks with memory
```

---

## Best Practices Summary

1. **Always** use `package mypackage_test` (black-box testing)
2. **Always** use table-driven tests with `t.Run()`
3. **Always** test error cases when function returns `error`
4. **Always** check and handle errors explicitly
5. **Always** use descriptive subtest names
6. **Always** test all code paths (happy, error, edge cases)
7. **Always** use mockery for interface mocking
8. Use `t.Errorf()` for most assertions
9. Use `t.Fatalf()` only for critical setup failures
10. Use assert functions from testify (assert.Equal, assert.NoError, etc.)
11. Mock interfaces with mockery, not concrete types
12. Use `t.Parallel()` only for independent tests (see guidelines above)
13. Use `t.Helper()` for test setup functions
14. Use `t.TempDir()` for temporary files
15. Use `t.Cleanup()` for resource cleanup
16. Never assert on exact error message strings (use errors.Is/As)
17. Never depend on global state in tests
18. Never use sleep/timing in tests (use channels/sync primitives)
19. Never test unexported functions directly
20. Never ignore errors in test setup
21. Aim for 100% coverage of function code paths
22. Test edge cases and boundary conditions (zero, empty, nil, max, min)
23. Keep test code simple and readable
24. Run tests with `-race` to detect race conditions
25. Use fuzzing for input validation functions (Go 1.18+)

---

## Troubleshooting

**Test fails with "undefined: mypackage.Function"**
- ✅ Import the package: `import "mymodule/mypackage"`

**Tests run in wrong order or state leaks**
- ✅ Check for global variables
- ✅ Use `t.Cleanup()` for teardown
- ✅ Use `t.TempDir()` for files
- ✅ Don't use `t.Parallel()` if tests share state

**Mock expectations not matched**
- ✅ Verify mock setup matches actual calls
- ✅ Call `mockRepo.AssertExpectations(t)`
- ✅ Use `mock.Anything` for flexible matching

**Race condition detected**
- ✅ Run with `go test -race`
- ✅ Use `sync.Mutex` to protect shared state
- ✅ Avoid shared state in tests
- ✅ Don't use `t.Parallel()` for dependent tests

---

**Last Updated**: January 2026  
**Go Version Compatibility**: Go 1.18+ (with fuzzing and generics support), Go 1.22+ (improved loop variables)
