# validate - Struct Validation

Simple struct and variable validation using tags.

## Struct Validation

```go
type User struct {
    Name  string `validate:"required,min=2,max=50"`
    Email string `validate:"required,email"`
    Age   int    `validate:"min=18,max=120"`
}

user := User{Name: "Alice", Email: "alice@example.com", Age: 25}

if err := validate.Struct(&user); err != nil {
    log.Printf("Validation failed: %v", err)
}
```

## Variable Validation

```go
// Single variable
email := "invalid-email"
err := validate.Var(email, "required,email")

// Multiple variables
values := map[string]interface{}{
    "username": "john_doe",
    "password": "secret123",
}
rules := map[string]string{
    "username": "required,alphanum,min=3",
    "password": "required,min=8",
}
err := validate.Vars(values, rules)
```

## Common Tags

**String**: `required`, `min=N`, `max=N`, `email`, `url`, `alpha`, `alphanum`, `oneof=val1 val2`

**Numeric**: `min=N`, `max=N`, `gt=N`, `gte=N`, `lt=N`, `lte=N`
