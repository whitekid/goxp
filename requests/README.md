# requests - HTTP Client

High-performance HTTP client with connection pooling and compression support.

## Basic Usage

```go
// Simple GET request
resp, err := requests.Get("https://api.github.com").Do(ctx)
if err != nil {
    log.Fatal(err)
}
defer resp.Body.Close()

// Parse JSON response
var result map[string]interface{}
resp.JSON(&result)
```

## HTTP Methods

```go
requests.Get(url)
requests.Post(url)
requests.Put(url)
requests.Patch(url)
requests.Delete(url)
requests.Head(url)
requests.Options(url)
requests.Verb("CUSTOM", url)
```

## Request Building

```go
// JSON request with headers
payload := map[string]string{"key": "value"}
resp, err := requests.Post("https://api.example.com").
    Header("Authorization", "Bearer token").
    JSON(payload).
    Do(ctx)

// Form data with query parameters
resp, err := requests.Post("https://httpbin.org/post").
    Query("param1", "value1").
    Form("field1", "data1").
    Do(ctx)

// Custom client
client := &http.Client{Timeout: 10 * time.Second}
resp, err := requests.Get(url).WithClient(client).Do(ctx)
```

## Authentication

```go
// Bearer token
requests.Get(url).AuthBearer("jwt-token").Do(ctx)

// Basic auth
requests.Get(url).AuthBasic("user", "pass").Do(ctx)

// API token
requests.Get(url).AuthToken("api-token").Do(ctx)
```

## Features

- **Connection Pooling**: MaxIdleConns=100, MaxIdleConnsPerHost=10
- **Buffer Pooling**: Reduces memory allocations
- **Compression**: gzip, brotli, zstd, deflate support
- **Redirect Control**: `FollowRedirect(true/false)`
