# validate - Struct Validation Made Easy

**Simple and powerful struct validation with comprehensive error handling and type safety.**

## Key Features

- **Simple API**: Easy-to-use validation functions
- **Struct Validation**: Comprehensive struct field validation
- **Variable Validation**: Direct variable validation
- **Rich Error Handling**: Detailed validation error information
- **Tag-Based**: Uses struct tags for validation rules

## Quick Start

```go
import "github.com/whitekid/goxp/validate"

type User struct {
    Name  string `validate:"required,min=2,max=50"`
    Email string `validate:"required,email"`
    Age   int    `validate:"min=18,max=120"`
}

user := User{
    Name:  "Alice",
    Email: "alice@example.com",
    Age:   25,
}

if err := validate.Struct(&user); err != nil {
    log.Printf("Validation failed: %v", err)
}
```

## Struct Validation

### Basic Usage
```go
type Product struct {
    Name     string  `validate:"required,min=1,max=100"`
    Price    float64 `validate:"required,gt=0"`
    Category string  `validate:"required,oneof=electronics books clothing"`
    SKU      string  `validate:"required,alphanum,len=8"`
}

product := Product{
    Name:     "Laptop",
    Price:    999.99,
    Category: "electronics",
    SKU:      "LAP12345",
}

err := validate.Struct(&product)
if err != nil {
    if validate.IsValidationError(err) {
        // Handle validation-specific errors
        fmt.Printf("Validation failed: %v\n", err)
    }
}
```

### Nested Struct Validation
```go
type Address struct {
    Street  string `validate:"required"`
    City    string `validate:"required"`
    ZipCode string `validate:"required,len=5"`
}

type Customer struct {
    Name    string   `validate:"required"`
    Email   string   `validate:"required,email"`
    Address *Address `validate:"required"`
}

customer := Customer{
    Name:  "John Doe",
    Email: "john@example.com",
    Address: &Address{
        Street:  "123 Main St",
        City:    "Anytown",
        ZipCode: "12345",
    },
}

err := validate.Struct(&customer)
```

## Variable Validation

### Single Variable
```go
// Validate individual variables
email := "invalid-email"
err := validate.Var(email, "required,email")
if err != nil {
    fmt.Printf("Email validation failed: %v\n", err)
}

// Validate numbers
age := -5
err = validate.Var(age, "required,min=0,max=150")
if err != nil {
    fmt.Printf("Age validation failed: %v\n", err)
}
```

### Multiple Variables
```go
// Validate multiple variables at once
values := map[string]interface{}{
    "username": "john_doe",
    "password": "secret123",
    "age":      25,
}

rules := map[string]string{
    "username": "required,alphanum,min=3,max=20",
    "password": "required,min=8",
    "age":      "required,min=18,max=120",
}

err := validate.Vars(values, rules)
if err != nil {
    fmt.Printf("Multi-variable validation failed: %v\n", err)
}
```

## Error Handling

### Check Validation Errors
```go
err := validate.Var(0, "required")
if err != nil {
    if validate.IsValidationError(err) {
        fmt.Println("This is a validation error")
        // Handle validation-specific logic
    } else {
        fmt.Println("This is a different type of error")
        // Handle other errors
    }
}
```

### Detailed Error Information
```go
type Config struct {
    Port     int    `validate:"required,min=1,max=65535"`
    Host     string `validate:"required,hostname"`
    Database string `validate:"required,min=1"`
}

config := Config{
    Port:     0,      // Invalid: below minimum
    Host:     "",     // Invalid: required but empty
    Database: "mydb", // Valid
}

err := validate.Struct(&config)
if err != nil {
    // Get detailed validation errors
    fmt.Printf("Validation errors:\n%v\n", err)
    
    // Example output:
    // Key: 'Config.Port' Error:Field validation for 'Port' failed on the 'min' tag
    // Key: 'Config.Host' Error:Field validation for 'Host' failed on the 'required' tag
}
```

## Common Validation Tags

### String Validation
- `required` - Field cannot be empty
- `min=N` - Minimum length
- `max=N` - Maximum length
- `len=N` - Exact length
- `email` - Valid email format
- `url` - Valid URL format
- `alpha` - Alphabetic characters only
- `alphanum` - Alphanumeric characters only
- `oneof=val1 val2` - Must be one of specified values

### Numeric Validation
- `min=N` - Minimum value
- `max=N` - Maximum value
- `gt=N` - Greater than
- `gte=N` - Greater than or equal
- `lt=N` - Less than
- `lte=N` - Less than or equal

### Advanced Usage
```go
type APIRequest struct {
    Method   string            `validate:"required,oneof=GET POST PUT DELETE"`
    URL      string            `validate:"required,url"`
    Headers  map[string]string `validate:"dive,keys,required,endkeys,required"`
    Body     []byte            `validate:"omitempty"`
    Timeout  time.Duration     `validate:"min=1s,max=5m"`
}
```

## Integration Examples

### HTTP Handler Validation
```go
func createUserHandler(w http.ResponseWriter, r *http.Request) {
    var user User
    if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
        http.Error(w, "Invalid JSON", http.StatusBadRequest)
        return
    }
    
    if err := validate.Struct(&user); err != nil {
        if validate.IsValidationError(err) {
            http.Error(w, fmt.Sprintf("Validation failed: %v", err), 
                      http.StatusBadRequest)
            return
        }
        http.Error(w, "Internal error", http.StatusInternalServerError)
        return
    }
    
    // Process valid user...
}
```

### Configuration Validation
```go
type AppConfig struct {
    Server   ServerConfig   `validate:"required"`
    Database DatabaseConfig `validate:"required"`
    Redis    RedisConfig    `validate:"required"`
}

func LoadConfig(filename string) (*AppConfig, error) {
    var config AppConfig
    
    // Load from file...
    data, err := os.ReadFile(filename)
    if err != nil {
        return nil, err
    }
    
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, err
    }
    
    // Validate configuration
    if err := validate.Struct(&config); err != nil {
        return nil, fmt.Errorf("invalid configuration: %w", err)
    }
    
    return &config, nil
}
```

---

**Note**: This package is part of the [goxp](https://github.com/whitekid/goxp) utility collection.
