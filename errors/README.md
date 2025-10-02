# errors - Error Handling with Stack Traces

Enhanced error handling with automatic stack trace capture.

## Basic Usage

```go
// Create error with stack trace
err := errors.New("something went wrong")

// Wrap existing error
if err != nil {
    return errors.Wrap(err, "failed to process data")
}

// Formatted error
err := errors.Errorf(nil, "invalid user ID: %d", userID)

// Display with stack trace
fmt.Printf("%+v\n", err)
```

## Error Wrapping

```go
func processFile(filename string) error {
    file, err := os.Open(filename)
    if err != nil {
        return errors.Wrapf(err, "failed to open %s", filename)
    }
    defer file.Close()

    data, err := io.ReadAll(file)
    if err != nil {
        return errors.Wrapf(err, "failed to read %s", filename)
    }

    return processData(data)
}
```

## Stack Trace Formatting

```go
err := errors.New("example error")

fmt.Printf("%s\n", err)   // example error
fmt.Printf("%v\n", err)   // example error
fmt.Printf("%+v\n", err)  // example error with full stack trace
```

## Standard Library Compatible

```go
// Error comparison
var ErrNotFound = errors.New("not found")
if errors.Is(err, ErrNotFound) { ... }

// Type assertion
var pathError *os.PathError
if errors.As(err, &pathError) { ... }

// Unwrap
originalErr := errors.Unwrap(wrappedErr)

// Join multiple errors
combined := errors.Join(err1, err2)
```
