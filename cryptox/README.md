# cryptox - Secure Encryption/Decryption for Go

**Simple and secure encryption utilities with security best practices and modern cryptographic standards.**

## Key Features

- **Security Hardened**: No sensitive data logging, secure key handling
- **Simple API**: Easy-to-use encrypt/decrypt functions
- **Multiple Algorithms**: Support for various encryption methods
- **Security Warnings**: Clear guidance on algorithm choices
- **Well Tested**: Comprehensive test coverage

## Quick Start

```go
import (
    "github.com/whitekid/goxp"
    "github.com/whitekid/goxp/cryptox"
)

// Generate a secure random key
key := goxp.RandomString(32)  // Use appropriate key length

// Encrypt data
plaintext := "동해 물과 백두산이 마르고 닳도록"
encrypted, err := cryptox.Encrypt(key, plaintext)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Encrypted: %v\n", encrypted)

// Decrypt data  
decrypted, err := cryptox.Decrypt(key, encrypted)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Decrypted: %v\n", decrypted)
```

[Go Playground](https://go.dev/play/p/-Rl8Ci8x0Xp)

## Security Considerations

### Important Security Warnings

1. **DES Algorithm Deprecated**: 
   ```go
   // ❌ AVOID: DES is cryptographically broken
   // DES encryption is deprecated and insecure
   // Use AES instead for new applications
   ```

2. **Key Management**:
   - Use cryptographically secure random keys
   - Store keys securely (environment variables, key management systems)
   - Never hardcode keys in source code
   - Use appropriate key lengths (AES-256 recommended)

3. **Secure Logging**:
   - The library has been hardened to prevent sensitive data leakage in logs
   - No encryption keys or plaintext data will be logged

### Recommended Usage

```go
// ✅ RECOMMENDED: Use strong keys
key := goxp.RandomString(32)  // 256-bit key for AES-256

// ✅ RECOMMENDED: Handle errors properly
encrypted, err := cryptox.Encrypt(key, sensitiveData)
if err != nil {
    // Log error without exposing sensitive data
    log.Printf("Encryption failed: %v", err)
    return err
}

// ✅ RECOMMENDED: Clear sensitive data from memory when done
defer func() {
    // Clear sensitive variables
    for i := range key {
        key[i] = 0
    }
}()
```

## Algorithm Support

### Modern Encryption (Recommended)
- **AES**: Advanced Encryption Standard (recommended)
- **ChaCha20**: Modern stream cipher

### Legacy Support (Deprecated)
- **DES**: **Deprecated** - Use only for legacy compatibility
  ```go
  // This will show a deprecation warning
  result, err := cryptox.EncryptDES(key, data)
  // Warning: DES encryption is deprecated and insecure
  ```

## Best Practices

### Key Generation
```go
// Generate secure keys
func generateSecureKey(length int) string {
    return goxp.RandomString(length)
}

// For AES-256
key := generateSecureKey(32)  // 32 bytes = 256 bits
```

### Error Handling
```go
func secureEncrypt(key, data string) (string, error) {
    encrypted, err := cryptox.Encrypt(key, data)
    if err != nil {
        // Don't log sensitive data
        return "", fmt.Errorf("encryption failed: %w", err)
    }
    return encrypted, nil
}
```

### Environment-Based Configuration
```go
func getEncryptionKey() string {
    key := os.Getenv("ENCRYPTION_KEY")
    if key == "" {
        log.Fatal("ENCRYPTION_KEY environment variable not set")
    }
    return key
}
```

## Security Improvements

Recent security enhancements include:

- **Secure Logging**: Removed sensitive data from error logs
- **Input Validation**: Enhanced validation for encryption parameters
- **Deprecation Warnings**: Clear warnings for insecure algorithms
- **Key Handling**: Improved secure key management practices

## Testing and Validation

The cryptox package includes comprehensive tests for:
- Encryption/decryption round trips
- Error handling and edge cases
- Security boundary testing
- Algorithm-specific validations

## Migration Guide

### From Insecure Patterns
```go
// ❌ OLD (insecure)
key := "hardcoded-key"
encrypted, _ := cryptox.EncryptDES(key, data)

// ✅ NEW (secure)
key := goxp.RandomString(32)
encrypted, err := cryptox.Encrypt(key, data)
if err != nil {
    return fmt.Errorf("encryption failed: %w", err)
}
```

### Upgrading from DES
```go
// Replace DES usage
// OLD: cryptox.EncryptDES(key, data)
// NEW: cryptox.Encrypt(key, data)  // Uses secure algorithm by default
```

---

**Security Notice**: This package follows security best practices and provides warnings for deprecated algorithms. Always use the latest encryption standards and proper key management.

**Note**: This package is part of the [goxp](https://github.com/whitekid/goxp) utility collection.
