# requests - High-Performance HTTP Client for Go

**Performance-optimized HTTP client with connection pooling, buffer reuse, and modern Go features.**

## Key Features

- **High Performance**: Connection pooling, buffer pooling, optimized timeouts
- **Simple API**: Fluent interface for easy request building
- **Automatic Compression**: Built-in support for gzip, brotli, zstd, deflate
- **Smart Redirects**: Configurable redirect handling with state isolation
- **Type Safety**: Modern Go with generics support
- **Well Tested**: Comprehensive test coverage

## Quick Start

```go
import "github.com/whitekid/goxp/requests"

// Simple GET request
resp, err := requests.Get("https://api.github.com").Do(context.Background())
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

// JSON response parsing
var result map[string]interface{}
err = resp.JSON(&result)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("URL: %s\n", result["hub_url"])
```

## Advanced Usage

### JSON Requests with Headers
```go
payload := map[string]string{"key": "value"}

resp, err := requests.Post("https://api.example.com/data").
    Header("Authorization", "Bearer your-token").
    Header("Content-Type", "application/json").
    JSON(payload).
    Do(ctx)
```

### Form Data and Query Parameters
```go
resp, err := requests.Post("https://httpbin.org/post").
    Query("param1", "value1").
    Query("param2", "value2").
    Form("field1", "data1").
    Form("field2", "data2").
    Do(ctx)
```

### Custom Client with Timeouts
```go
client := &http.Client{
    Timeout: 10 * time.Second,
}

resp, err := requests.Get("https://slow-api.com").
    WithClient(client).
    Do(ctx)
```

### Authentication
```go
// Bearer token
resp, err := requests.Get("https://api.example.com").
    AuthBearer("your-jwt-token").
    Do(ctx)

// Basic auth
resp, err := requests.Get("https://api.example.com").
    AuthBasic("username", "password").
    Do(ctx)

// Custom token
resp, err := requests.Get("https://api.example.com").
    AuthToken("your-api-token").
    Do(ctx)
```

### Redirect Control
```go
// Disable redirects
resp, err := requests.Get("https://httpbin.org/redirect/3").
    FollowRedirect(false).
    Do(ctx)

// Custom redirect handling (default: follow redirects)
resp, err := requests.Get("https://httpbin.org/redirect/3").
    FollowRedirect(true).
    Do(ctx)
```

## Performance Features

### Connection Pooling
The library uses an optimized default HTTP client with:
- **MaxIdleConns**: 100 connections
- **MaxIdleConnsPerHost**: 10 connections per host
- **IdleConnTimeout**: 90 seconds
- **TLSHandshakeTimeout**: 5 seconds

### Buffer Pooling
JSON encoding uses a buffer pool to reduce memory allocations and GC pressure.

### Compression Support
Automatic support for:
- **gzip** - Standard compression
- **brotli** - Google's compression algorithm
- **zstd** - Facebook's Zstandard
- **deflate** - Standard deflate compression

## HTTP Methods

All standard HTTP methods are supported:

```go
requests.Get(url)
requests.Post(url)
requests.Put(url)
requests.Patch(url)
requests.Delete(url)
requests.Head(url)
requests.Options(url)
requests.Verb("CUSTOM", url)  // Custom HTTP verbs
```

## Error Handling

```go
resp, err := requests.Get("https://api.example.com").Do(ctx)
if err != nil {
    // Network or request building error
    return err
}

if err := resp.Success(); err != nil {
    // HTTP status code indicates failure (not 2xx)
    return err
}
```

## Response Utilities

```go
// Get response as string
text := resp.String()

// Parse JSON
var data MyStruct
err := resp.JSON(&data)

// Access raw response
statusCode := resp.StatusCode
headers := resp.Header
body := resp.Body  // io.ReadCloser
```

## Testing

For comprehensive examples and edge cases, refer to [test cases](request_test.go).

## Migration from Other Libraries

Coming from other HTTP libraries? The requests package provides a familiar, chainable API that's both powerful and easy to use.

---

**Note**: This package is part of the [goxp](https://github.com/whitekid/goxp) utility collection.
