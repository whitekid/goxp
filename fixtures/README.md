# fixtures - Test Utilities and Environment Management

**Comprehensive test fixtures for managing environment variables, temporary files, and test state isolation.**

## Key Features

- **Environment Management**: Safe environment variable manipulation for tests
- **File Fixtures**: Temporary file and directory management
- **State Isolation**: Ensure tests don't interfere with each other
- **Auto Cleanup**: Automatic resource cleanup with defer patterns
- **Simple API**: Easy-to-use testing utilities

## Environment Variable Fixtures

### Basic Usage
```go
import (
    "os"
    "testing"
    "github.com/whitekid/goxp/fixtures"
    "github.com/stretchr/testify/require"
)

func TestEnvironmentVariable(t *testing.T) {
    // Set environment variable for test duration
    defer fixtures.Env("API_KEY", "test-key-123")()
    
    // Variable is available during test
    require.Equal(t, "test-key-123", os.Getenv("API_KEY"))
    
    // Automatically restored when test completes
}
```

### Multiple Environment Variables
```go
func TestMultipleEnvVars(t *testing.T) {
    // Set multiple environment variables
    defer fixtures.Env("DB_HOST", "localhost")()
    defer fixtures.Env("DB_PORT", "5432")()
    defer fixtures.Env("DB_NAME", "testdb")()
    
    // All variables are available
    require.Equal(t, "localhost", os.Getenv("DB_HOST"))
    require.Equal(t, "5432", os.Getenv("DB_PORT"))
    require.Equal(t, "testdb", os.Getenv("DB_NAME"))
}
```

### Environment Isolation
```go
func TestEnvironmentIsolation(t *testing.T) {
    // Original value
    originalPath := os.Getenv("PATH")
    
    // Temporarily override
    defer fixtures.Env("PATH", "/custom/path")()
    require.Equal(t, "/custom/path", os.Getenv("PATH"))
    
    // After test, original value is restored
    // (verified by test framework cleanup)
}
```

## File and Directory Fixtures

### Temporary Files
```go
func TestWithTempFile(t *testing.T) {
    // Create temporary file
    tempFile := fixtures.TempFile(t, "test-data", "Hello, World!")
    defer os.Remove(tempFile) // Cleanup
    
    // File exists and contains expected data
    data, err := os.ReadFile(tempFile)
    require.NoError(t, err)
    require.Equal(t, "Hello, World!", string(data))
}
```

### Temporary Directories
```go
func TestWithTempDir(t *testing.T) {
    // Create temporary directory
    tempDir := fixtures.TempDir(t, "test-workspace")
    defer os.RemoveAll(tempDir) // Cleanup
    
    // Create files in temp directory
    testFile := filepath.Join(tempDir, "test.txt")
    err := os.WriteFile(testFile, []byte("test data"), 0644)
    require.NoError(t, err)
    
    // Verify file exists
    require.FileExists(t, testFile)
}
```

## Configuration Testing

### Database Configuration
```go
func TestDatabaseConnection(t *testing.T) {
    // Set test database configuration
    defer fixtures.Env("DB_HOST", "localhost")()
    defer fixtures.Env("DB_PORT", "5432")()
    defer fixtures.Env("DB_NAME", "testdb")()
    defer fixtures.Env("DB_USER", "testuser")()
    defer fixtures.Env("DB_PASS", "testpass")()
    
    // Test database configuration loading
    config := loadDatabaseConfig()
    require.Equal(t, "localhost", config.Host)
    require.Equal(t, 5432, config.Port)
    require.Equal(t, "testdb", config.Database)
}
```

### API Configuration
```go
func TestAPIConfiguration(t *testing.T) {
    // Set API configuration
    defer fixtures.Env("API_BASE_URL", "https://api.test.com")()
    defer fixtures.Env("API_TOKEN", "test-token-123")()
    defer fixtures.Env("API_TIMEOUT", "30s")()
    
    client := NewAPIClient()
    require.Equal(t, "https://api.test.com", client.BaseURL)
    require.Equal(t, "test-token-123", client.Token)
    require.Equal(t, 30*time.Second, client.Timeout)
}
```

## Advanced Fixtures

### Custom Fixture Functions
```go
// Create custom fixture for complex setup
func withTestServer(t *testing.T, handler http.Handler) (string, func()) {
    server := httptest.NewServer(handler)
    
    // Set environment to point to test server
    cleanup := fixtures.Env("API_URL", server.URL)
    
    return server.URL, func() {
        cleanup()      // Restore environment
        server.Close() // Close test server
    }
}

func TestWithCustomFixture(t *testing.T) {
    handler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.WriteHeader(http.StatusOK)
        w.Write([]byte(`{"status": "ok"}`))
    })
    
    url, cleanup := withTestServer(t, handler)
    defer cleanup()
    
    // Test against the fixture server
    resp, err := http.Get(url + "/health")
    require.NoError(t, err)
    require.Equal(t, http.StatusOK, resp.StatusCode)
}
```

### Fixture Composition
```go
func TestComplexScenario(t *testing.T) {
    // Combine multiple fixtures
    tempDir := fixtures.TempDir(t, "complex-test")
    defer os.RemoveAll(tempDir)
    
    defer fixtures.Env("DATA_DIR", tempDir)()
    defer fixtures.Env("LOG_LEVEL", "debug")()
    defer fixtures.Env("FEATURE_FLAG_X", "true")()
    
    // Create test data files
    configFile := filepath.Join(tempDir, "config.json")
    config := map[string]interface{}{
        "database": map[string]string{
            "host": "localhost",
            "port": "5432",
        },
    }
    
    configData, _ := json.Marshal(config)
    err := os.WriteFile(configFile, configData, 0644)
    require.NoError(t, err)
    
    // Run complex test scenario...
    app := NewApplication()
    err = app.Initialize()
    require.NoError(t, err)
}
```

## Integration with Test Suites

### Table-Driven Tests
```go
func TestAPIEndpoints(t *testing.T) {
    tests := []struct {
        name   string
        envVars map[string]string
        expect  string
    }{
        {
            name: "production environment",
            envVars: map[string]string{
                "ENV": "production",
                "API_URL": "https://api.prod.com",
            },
            expect: "https://api.prod.com",
        },
        {
            name: "development environment",
            envVars: map[string]string{
                "ENV": "development",
                "API_URL": "http://localhost:8080",
            },
            expect: "http://localhost:8080",
        },
    }
    
    for _, tt := range tests {
        t.Run(tt.name, func(t *testing.T) {
            // Set up environment for this test case
            var cleanups []func()
            for key, value := range tt.envVars {
                cleanups = append(cleanups, fixtures.Env(key, value))
            }
            defer func() {
                for _, cleanup := range cleanups {
                    cleanup()
                }
            }()
            
            // Run test with environment
            config := loadConfig()
            require.Equal(t, tt.expect, config.APIURL)
        })
    }
}
```

### Parallel Test Safety
```go
func TestParallelEnvironmentSafety(t *testing.T) {
    t.Parallel() // This test can run in parallel
    
    // Each parallel test gets its own environment isolation
    defer fixtures.Env("TEST_ID", t.Name())()
    
    // Test-specific logic...
    require.Equal(t, t.Name(), os.Getenv("TEST_ID"))
}
```

## Best Practices

### Always Use Defer
```go
// ✅ GOOD: Always use defer for cleanup
func TestGoodPractice(t *testing.T) {
    defer fixtures.Env("KEY", "value")()
    // Test logic...
}

// ❌ BAD: Manual cleanup is error-prone
func TestBadPractice(t *testing.T) {
    cleanup := fixtures.Env("KEY", "value")
    // Test logic...
    cleanup() // Might be forgotten or skipped on early return
}
```

### Environment Variable Validation
```go
func TestWithValidation(t *testing.T) {
    defer fixtures.Env("REQUIRED_VAR", "test-value")()
    
    // Verify environment is set correctly
    value := os.Getenv("REQUIRED_VAR")
    require.NotEmpty(t, value, "REQUIRED_VAR must be set for this test")
    
    // Continue with test...
}
```

---

**Note**: This package is part of the [goxp](https://github.com/whitekid/goxp) utility collection.
