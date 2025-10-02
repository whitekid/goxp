# slug - URL-Safe Encoders

UUID and integer encoding for URL-safe short identifiers.

## UUID Encoding

```go
uid := uuid.New()
sg := slug.NewUUID()

// Encode: 0ed24fc4-d599-4c06-adae-80526c63d7a4 => DtJPxNWZTAatroBSbGPXpA
encoded := sg.Encode(uid)

// Decode: DtJPxNWZTAatroBSbGPXpA => 0ed24fc4-d599-4c06-adae-80526c63d7a4
decoded := sg.Decode(encoded)
```

## Integer Shortener

```go
encoding := "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789-_"
shortner := slug.NewShortner(encoding)

// 22775 => B9L
short := shortner.Encode(22775)
```
