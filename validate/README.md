# Validator

## `Struct()` - validate struct

```go
v := &struct {
    RequiredInt int `validate:"required"`
}{
    RequiredInt: 0,
}
require.Error(t, Struct(v))
// err: Key: 'RequiredInt' Error:Field validation for 'RequiredInt' failed on the 'required' tag
```

## `Var(), Vars()` - validate variable

```go
require.Error(t, Var(0, "required"))
// err: Key: '' Error:Field validation for '' failed on the 'required' tag
```

## `IsValidationError()` - check if error is validation error

```go
err := Var(0, "required")
require.Error(t, err)
require.True(t, IsValidationError(err))
```
