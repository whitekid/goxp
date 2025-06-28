# errors - Enhanced Error Handling with Stack Traces

**Enhanced error handling for Go with automatic stack trace capture, error wrapping, and rich formatting capabilities.**

## Key Features

- **Stack Trace Capture**: Automatic stack trace recording for better debugging
- **Error Wrapping**: Enhanced error wrapping with context preservation
- **Rich Formatting**: Detailed error output with full stack traces
- **Standard Library Compatible**: Drop-in replacement for Go's `errors` package
- **Performance Optimized**: Efficient stack trace capture and formatting
- **Chain Support**: Full support for error chains and unwrapping

## Quick Start

```go
import "github.com/whitekid/goxp/errors"

// Create new error with stack trace
err := errors.New("something went wrong")

// Wrap existing error with context
if err != nil {
    return errors.Wrap(err, "failed to process user data")
}

// Create formatted error
if userID == 0 {
    return errors.Errorf(nil, "invalid user ID: %d", userID)
}

// Display full error with stack trace
fmt.Printf("%+v\n", err)
```

## Error Creation

### Basic Error Creation

```go
// Create new error with automatic stack trace
err := errors.New("database connection failed")

// Create formatted error
err := errors.Errorf(nil, "user %s not found", username)

// Use standard library compatibility
var ErrNotFound = errors.New("resource not found")
```

### Error Wrapping

```go
// Wrap existing error with additional context
originalErr := someFunction()
if originalErr != nil {
    return errors.Wrap(originalErr, "failed to execute operation")
}

// Wrap with formatted message
if err := validateInput(data); err != nil {
    return errors.Wrapf(err, "validation failed for user %s", userID)
}

// Create new error or wrap existing one
func processData(data []byte) error {
    if len(data) == 0 {
        return errors.Errorf(nil, "empty data provided")
    }
    
    if err := parseData(data); err != nil {
        return errors.Errorf(err, "failed to parse data: %d bytes", len(data))
    }
    
    return nil
}
```

## Stack Trace Support

### Automatic Stack Capture

```go
func deepFunction() error {
    return errors.New("deep error occurred")
}

func middleFunction() error {
    if err := deepFunction(); err != nil {
        return errors.Wrap(err, "middle layer error")
    }
    return nil
}

func topFunction() error {
    if err := middleFunction(); err != nil {
        return errors.Wrap(err, "top layer error")
    }
    return nil
}

// Usage
if err := topFunction(); err != nil {
    // Print basic error
    fmt.Printf("Error: %v\n", err)
    // Output: Error: top layer error: middle layer error: deep error occurred
    
    // Print with full stack trace
    fmt.Printf("Full trace:\n%+v\n", err)
    // Output: Full stack trace with file paths, line numbers, and function names
}
```

### Rich Formatting Options

```go
err := errors.New("example error")

// Basic string representation
fmt.Printf("%s\n", err)        // example error
fmt.Printf("%v\n", err)        // example error

// Full stack trace
fmt.Printf("%+v\n", err)       // example error with complete stack trace

// Error details in logs
log.Printf("Operation failed: %+v", err)
```

## Error Chain Operations

### Standard Library Compatibility

```go
// Error comparison
var ErrNotFound = errors.New("not found")

if errors.Is(err, ErrNotFound) {
    fmt.Println("Resource not found")
}

// Error type assertion
var pathError *os.PathError
if errors.As(err, &pathError) {
    fmt.Printf("Path error: %s\n", pathError.Path)
}

// Error unwrapping
originalErr := errors.Unwrap(wrappedErr)

// Join multiple errors
err1 := errors.New("first error")
err2 := errors.New("second error")
combined := errors.Join(err1, err2)
```

### Error Chain Inspection

```go
func analyzeError(err error) {
    // Walk through error chain
    current := err
    level := 0
    
    for current != nil {
        fmt.Printf("Level %d: %v\n", level, current)
        current = errors.Unwrap(current)
        level++
    }
}

// Usage
wrappedErr := errors.Wrap(
    errors.Wrap(
        errors.New("root cause"), 
        "middle layer"),
    "top layer")

analyzeError(wrappedErr)
// Level 0: top layer: middle layer: root cause
// Level 1: middle layer: root cause
// Level 2: root cause
```

## Advanced Usage

### Custom Error Types with Stack Traces

```go
type ValidationError struct {
    Field string
    Value interface{}
    err   error
}

func (e *ValidationError) Error() string {
    return fmt.Sprintf("validation failed for field %s: %v", e.Field, e.Value)
}

func (e *ValidationError) Unwrap() error {
    return e.err
}

func NewValidationError(field string, value interface{}, cause error) error {
    verr := &ValidationError{
        Field: field,
        Value: value,
        err:   cause,
    }
    
    // Wrap with stack trace
    return errors.Wrap(verr, "validation error")
}

// Usage
if age < 0 {
    return NewValidationError("age", age, errors.New("must be positive"))
}
```

### Error Context Preservation

```go
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return errors.Wrapf(err, "failed to open file %s", filename)
    }
    defer file.Close()
    
    data, err := io.ReadAll(file)
    if err != nil {
        return errors.Wrapf(err, "failed to read file %s", filename)
    }
    
    if err := processData(data); err != nil {
        return errors.Wrapf(err, "failed to process content of %s", filename)
    }
    
    return nil
}

// Usage provides clear error context
if err := processFile("config.json"); err != nil {
    log.Printf("File processing failed: %+v", err)
    // Shows complete chain: file operation -> read operation -> processing
}
```

## Integration Examples

### HTTP Handler Error Handling

```go
func userHandler(w http.ResponseWriter, r *http.Request) {
    userID := r.URL.Query().Get("id")
    if userID == "" {
        err := errors.New("missing user ID parameter")
        http.Error(w, err.Error(), http.StatusBadRequest)
        log.Printf("Bad request: %+v", err)
        return
    }
    
    user, err := getUserByID(userID)
    if err != nil {
        wrappedErr := errors.Wrapf(err, "failed to get user %s", userID)
        
        if errors.Is(err, ErrUserNotFound) {
            http.Error(w, "User not found", http.StatusNotFound)
        } else {
            http.Error(w, "Internal server error", http.StatusInternalServerError)
        }
        
        log.Printf("User handler error: %+v", wrappedErr)
        return
    }
    
    // Success case...
}
```

### Database Operation Error Handling

```go
func CreateUser(db *sql.DB, user *User) error {
    tx, err := db.Begin()
    if err != nil {
        return errors.Wrap(err, "failed to begin transaction")
    }
    defer tx.Rollback()
    
    _, err = tx.Exec("INSERT INTO users (name, email) VALUES (?, ?)", 
                     user.Name, user.Email)
    if err != nil {
        return errors.Wrapf(err, "failed to insert user %s", user.Email)
    }
    
    if err = tx.Commit(); err != nil {
        return errors.Wrap(err, "failed to commit user creation")
    }
    
    return nil
}

// Usage
if err := CreateUser(db, newUser); err != nil {
    log.Printf("User creation failed: %+v", err)
    // Provides full context: transaction -> insert -> commit
}
```

### Service Layer Error Propagation

```go
type UserService struct {
    repo UserRepository
}

func (s *UserService) UpdateUser(id string, updates UserUpdates) error {
    // Validate input
    if err := updates.Validate(); err != nil {
        return errors.Wrapf(err, "invalid user updates for user %s", id)
    }
    
    // Check if user exists
    existing, err := s.repo.GetByID(id)
    if err != nil {
        return errors.Wrapf(err, "failed to retrieve user %s", id)
    }
    
    if existing == nil {
        return errors.Errorf(nil, "user %s not found", id)
    }
    
    // Apply updates
    if err := s.repo.Update(id, updates); err != nil {
        return errors.Wrapf(err, "failed to update user %s", id)
    }
    
    return nil
}
```

## Performance Considerations

### Stack Trace Overhead

```go
// Stack traces are captured only when errors are created
// No performance impact during normal operation

// For high-performance scenarios, consider:
func fastPath() error {
    // Use standard library errors for hot paths where stack traces aren't needed
    if criticalCondition {
        return fmt.Errorf("fast error: %v", condition)
    }
    
    // Use enhanced errors for debugging scenarios
    if debugMode {
        return errors.New("detailed error with stack trace")
    }
    
    return nil
}
```

### Memory Usage

```go
// Stack traces consume memory proportional to call depth
// Typical overhead: 8-24 bytes per frame + function names

// For long-running applications, consider error lifecycle:
func processLongRunningTask() {
    var lastError error
    
    for item := range workQueue {
        if err := processItem(item); err != nil {
            // Store only the most recent error to limit memory usage
            lastError = errors.Wrapf(err, "failed to process item %v", item.ID)
            continue
        }
    }
    
    if lastError != nil {
        return lastError
    }
    return nil
}
```

## Best Practices

### Error Message Guidelines

```go
// Good: Specific, actionable error messages
return errors.Errorf(nil, "failed to connect to database %s:%d", host, port)

// Good: Include relevant context
return errors.Wrapf(err, "user validation failed for email %s", user.Email)

// Avoid: Vague error messages
return errors.New("something went wrong")

// Avoid: Redundant wrapping
if err != nil {
    return errors.Wrap(err, "error occurred") // Not helpful
}
```

### Error Handling Patterns

```go
// Pattern 1: Immediate wrapping with context
func ProcessOrder(orderID string) error {
    order, err := getOrder(orderID)
    if err != nil {
        return errors.Wrapf(err, "failed to retrieve order %s", orderID)
    }
    
    if err := validateOrder(order); err != nil {
        return errors.Wrapf(err, "order %s validation failed", orderID)
    }
    
    return nil
}

// Pattern 2: Sentinel errors for expected conditions
var (
    ErrOrderNotFound = errors.New("order not found")
    ErrInvalidOrder  = errors.New("invalid order")
)

func ValidateOrder(order *Order) error {
    if order == nil {
        return ErrOrderNotFound
    }
    
    if order.Amount <= 0 {
        return errors.Wrapf(ErrInvalidOrder, "order amount must be positive: %v", order.Amount)
    }
    
    return nil
}
```

### Logging Integration

```go
func logError(err error, context string) {
    if err == nil {
        return
    }
    
    // Log basic error for normal operation
    log.WithField("context", context).Error(err.Error())
    
    // Log full stack trace for debugging (structured logging)
    log.WithFields(log.Fields{
        "context":    context,
        "stackTrace": fmt.Sprintf("%+v", err),
    }).Debug("Detailed error information")
}
```

## Function Reference

| Function | Description |
|----------|-------------|
| `New(message string) error` | Create new error with stack trace |
| `Errorf(err error, format string, args ...any) error` | Create/wrap error with formatting |
| `Wrap(err error, message string) error` | Wrap existing error with message |
| `Wrapf(err error, format string, args ...any) error` | Wrap existing error with formatted message |
| `Is(err, target error) bool` | Compare errors (standard library) |
| `As(err error, target any) bool` | Type assertion (standard library) |
| `Join(errs ...error) error` | Join multiple errors (standard library) |
| `Unwrap(err error) error` | Unwrap error chain (standard library) |

## Migration from Standard Library

```go
// Before: standard library
import "errors"
import "fmt"

err := errors.New("something failed")
wrappedErr := fmt.Errorf("operation failed: %w", err)

// After: enhanced errors
import "github.com/whitekid/goxp/errors"

err := errors.New("something failed")           // Now includes stack trace
wrappedErr := errors.Wrap(err, "operation failed") // Better wrapping

// Print with stack trace
fmt.Printf("%+v\n", wrappedErr)
```

---

**Note**: This package is part of the [goxp](https://github.com/whitekid/goxp) utility collection.