# cryptox - Encryption/Decryption Utilities

Simple encryption/decryption functions with AES support.

## Usage

```go
// Generate secure random key
key := goxp.RandomString(32)  // 32 bytes for AES-256

// Encrypt
encrypted, err := cryptox.Encrypt(key, "plaintext")
if err != nil {
    log.Fatal(err)
}

// Decrypt
decrypted, err := cryptox.Decrypt(key, encrypted)
if err != nil {
    log.Fatal(err)
}
```

[Go Playground](https://go.dev/play/p/-Rl8Ci8x0Xp)

## Security Notes

- **Use AES**: Recommended for all new applications
- **DES Deprecated**: Legacy support only, cryptographically broken
- **Key Length**: Use 32 bytes (256-bit) for AES-256
- **Key Storage**: Never hardcode keys, use environment variables or key management systems
