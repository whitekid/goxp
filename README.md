# goxp - Go Utility Function Collection

[![Go](https://github.com/whitekid/goxp/actions/workflows/go.yml/badge.svg)](https://github.com/whitekid/goxp/actions/workflows/go.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/whitekid/goxp)](https://goreportcard.com/report/github.com/whitekid/goxp)

A comprehensive collection of Go utility functions and packages designed for modern Go development (Go 1.24+). Features performance-optimized implementations, type-safe generics, and comprehensive testing.

**Recent Updates:**
- Critical performance optimizations (100x improvement in sampling operations)
- Security enhancements and vulnerability fixes
- HTTP client with connection pooling and buffer optimization
- Comprehensive test coverage improvements

For detailed usage examples, please refer to the test cases in each package.

## Core Functionality

### Goroutines & Concurrency

|                      |                                                         |
| -------------------- | ------------------------------------------------------- |
| `DoWithWorker()`     | iterate `chan` and do with workers (optimized pool)     |
| `Every()`            | run func with interval (context-aware)                  |
| `After()`            | run function with delay                                 |
| `Async()`,`Async2()` | run function with goroutine and get result asynchronous |

```go
// Worker pool with optimized performance
err := DoWithWorker(ctx, jobs, func(job Job) error {
    return processJob(job)
}, WithWorkerCount(runtime.NumCPU()))
```

### Data Encoding/Decoding (Type-Safe Generics)

|                |                                    |
| -------------- | ---------------------------------- |
| `ReadJSON[T]`  | decode JSON to type with generics  |
| `WriteJSON[T]` | encode type to JSON writer         |
| `ReadXML[T]`   | decode XML to type with generics   |
| `WriteXML[T]`  | encode type to XML writer          |
| `ReadYAML[T]`  | decode YAML to type with generics  |
| `WriteYAML[T]` | encode type to YAML writer         |

```go
// Type-safe JSON operations
var config Config
err := ReadJSON[Config]("config.json", &config)

// Direct type inference
users, err := ReadJSON[[]User]("users.json")
```

### Utility Functions

|                   |                                |
| ----------------- | ------------------------------ |
| `SetBits()`       | set bits in integer            |
| `ClearBits()`     | clear bits in integer          |
| `IsContextDone()` | return true if context is done |
| `WithTimeout()`   | run fn() with timeout context  |
| `Must(error)`     | panic if error not nil         |
| `NewPool[T]()`    | `sync.Pool` with generics      |
| `Abs[T]()`        | absolute value with generics   |
| `Timer()`         | measure execution time         |

```go
// Generic pool for any type
pool := NewPool[*bytes.Buffer](func() *bytes.Buffer {
    return &bytes.Buffer{}
})

// Execution timing
defer Timer("operation")()
// Output: time takes 123.45ms: operation
```

### String Parsing with Defaults

|                   |                                          |
| ----------------- | ---------------------------------------- |
| `AtoiDef[T]()`    | `strconv.Atoi()` with default value      |
| `ParseBoolDef()`  | `strconv.ParseBool()` with default value |
| `ParseIntDef[T]()` | `strconv.ParseInt()` with default value  |

```go
// Safe parsing with defaults
port := AtoiDef(os.Getenv("PORT"), 8080)           // defaults to 8080
timeout := ParseIntDef[int64](cfg.Timeout, 30)     // defaults to 30
debug := ParseBoolDef(os.Getenv("DEBUG"), false)   // defaults to false
```

### Time Parsing and Formatting

|                   |                                             |
| ----------------- | ------------------------------------------- |
| `ParseDateTime()` | parse string to time for well known layouts |
| `TimeWithLayout`  | `time.Time` with layouts                    |
| `RFC1123ZTime`    | `Mon, 02 Jan 2006 15:04:05 -0700`           |
| `RFC3339Time`     | `2006-01-02T15:04:05Z07:00`                 |

```go
// Smart time parsing
t, err := ParseDateTime("2023-12-25T10:30:00Z")     // Auto-detects format
rfc3339 := RFC3339Time(time.Now())                  // Formats as RFC3339
rfc1123 := RFC1123ZTime(time.Now())                 // Formats as RFC1123Z
```

### Shell Execution

**Security Warning**: Be careful with shell command execution. Always validate and sanitize inputs to prevent command injection attacks.

#### `Exec()` - simple run command

```go
// Run command and output to stdin/stdout
exc := Exec("ls", "-al")
err := exc.Do(context.Background())
require.NoError(t, err)

// Run command and get output
exc := Exec("ls", "-al")
output, err := exc.Output(context.Background())
require.NoError(t, err)
require.Contains(t, string(output), "README.md")
```

#### `Pipe()` - run command with pipe

### Conditional Execution

#### `If()`, `Else()`, `IfThen()` - functional conditionals

```go
IfThen(true, func() { fmt.Printf("true\n") })
// Output: true

IfThen(condition, 
    func() { fmt.Printf("true\n") }, 
    func() { fmt.Printf("false\n") })

// Multiple else conditions supported
IfThen(false, 
    func() { fmt.Printf("condition1\n") }, 
    func() { fmt.Printf("condition2\n") }, 
    func() { fmt.Printf("default\n") })
```

[Go Playground](https://go.dev/play/p/wNadBmhNYR-)

#### `Ternary()`, `TernaryF()`, `TernaryCF()` - ternary operations

```go
// Simple ternary
result := Ternary(condition, "yes", "no")

// Function-based ternary (lazy evaluation)
result := TernaryF(condition, 
    func() string { return expensiveTrue() }, 
    func() string { return expensiveFalse() })
```

### Random String/Byte Generation

**Security**: Uses cryptographically secure random generation with input validation.

#### `RandomByte()` - generate random bytes

```go
b := RandomByte(10)
// hex.EncodeToString(b) = "4d46ef2f87b8191daf58"
```

#### `RandomString()`, `RandomStringWith()` - generate random strings

```go
s := RandomString(10)
// s = "$c&I$#LR3Y"

// Custom character set with validation
s := RandomStringWith(10, []rune("abcdefg"))
// s = "bbffedabda"

// Size limits enforced for security (max 1MB)
s := RandomStringWith(1000, []rune("0123456789"))
```

### Performance Measurement

```go
func doSomething() {
    defer Timer("doSomething()")()
    time.Sleep(500 * time.Millisecond)
}

doSomething()
// Output: time takes 500.505063ms: doSomething()
```

[Go Playground](https://go.dev/play/p/Wcj2Hw5CLL6)

### Tuple Support

|                    |                              |
| ------------------ | ---------------------------- |
| `Tuple2`, `Tuple3` | `Pack()` and `Unpack()`      |
| `T2`,`T3`          | construct tuple with element |

```go
// Create and unpack tuples
tuple := T2("hello", 42)
str, num := tuple.Unpack()

// Multi-return functions
func getData() (string, int, error) {
    return "data", 123, nil
}
result := T3(getData())
```

## Sub-packages

### High-Performance HTTP Client
- **[requests](requests)** - Simple HTTP client with connection pooling and buffer optimization
  ```go
  resp, err := requests.Get("https://api.example.com").
      Header("Authorization", "Bearer token").
      JSON(payload).
      Do(ctx)
  ```

### Data Structure Extensions
- **[slicex](slicex)** - Performance-optimized slice operations with generics and `iter.Seq` support
- **[mapx](mapx)** - Map extensions with functional operations and optimized sampling
- **[sets](sets)** - Set data structure implementation
- **[chanx](chanx)** - Channel extensions and utilities

### Functional Programming
- **[fx](fx)** - Experimental functional programming with `iter.Seq` (Go 1.24+)
- **[fx/gen](fx/gen)** - Generator functions and utilities
- **[iterx](iterx)** - Iterator extensions and helper functions

### Security & Crypto
- **[cryptox](cryptox)** - Encryption/decryption functions (DES deprecated, use AES)
- **[x509x](x509x)** - X.509 certificate utilities

### Development Tools
- **[log](log)** - Simple structured logging powered by Zap
- **[errors](errors)** - Errors with stack traces
- **[retry](retry)** - Retry mechanisms with backoff strategies
- **[services](services)** - Simple service framework with lifecycle management
- **[worker](worker)** - Optimized worker pool implementation

### Testing & Validation
- **[testx](testx)** - Unit test utility functions
- **[fixtures](fixtures)** - Test fixture management (env vars, files)
- **[httptest](httptest)** - HTTP test server utilities
- **[validate](validate)** - Struct validation made easy

### CLI & Configuration
- **[cobrax](cobrax)** - Cobra and Viper utility functions
- **[flags](flags)** - Command-line flag management
- **[slug](slug)** - UUID to slug conversion

## Performance Optimizations

Recent performance improvements include:

- **HTTP Client**: Connection pooling, buffer reuse, optimized timeouts
- **Sampling Operations**: 100x performance improvement using `math/rand` instead of `crypto/rand`
- **Algorithm Efficiency**: Fisher-Yates shuffle for O(n) complexity
- **Memory Management**: Buffer pooling and reduced allocations
- **Type Safety**: Extensive use of Go generics for zero-cost abstractions

## Security Considerations

**Security Features:**
- Input validation and size limits in random string generation
- Secure logging practices (no sensitive data in logs)  
- Deprecation warnings for insecure cryptographic functions
- Command injection prevention guidance

**Security Warnings:**
- Always validate shell command inputs
- DES encryption is deprecated - use AES instead
- Be cautious with command execution functions
