# x509x - X.509 Certificate Utilities

**Enhanced X.509 certificate utilities for Go with support for certificate parsing, key generation, CSR handling, and cryptographic operations.**

## Key Features

- **Certificate Operations**: Parse and encode X.509 certificates from PEM/DER formats
- **Key Generation**: Generate RSA, ECDSA, and Ed25519 key pairs with secure defaults
- **CSR Handling**: Create and parse Certificate Signing Requests
- **Security Validation**: Algorithm validation and signature verification
- **Format Support**: Full PEM and DER format compatibility
- **Modern Algorithms**: Support for current cryptographic standards

## Quick Start

```go
import "github.com/whitekid/goxp/x509x"

// Generate a new key pair
privKey, pubKey, err := x509x.GenerateKey("ecdsa", 256)
if err != nil {
    log.Fatal(err)
}

// Create a Certificate Signing Request
template := &x509.CertificateRequest{
    Subject: pkix.Name{
        CommonName:   "example.com",
        Organization: []string{"Example Corp"},
    },
    DNSNames: []string{"example.com", "www.example.com"},
}

csrDER, err := x509x.CreateCertificateRequest(template, privKey)
if err != nil {
    log.Fatal(err)
}

// Parse certificates from PEM data
certPEM := `-----BEGIN CERTIFICATE-----
...
-----END CERTIFICATE-----`

cert, err := x509x.ParseCertificate([]byte(certPEM))
if err != nil {
    log.Fatal(err)
}
```

## Certificate Operations

### Parsing Certificates

```go
// Parse single certificate from PEM or DER
cert, err := x509x.ParseCertificate(certData)
if err != nil {
    log.Fatal(err)
}

// Parse certificate chain from PEM
certs, err := x509x.ParseCertificateChain(chainPEM)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("Parsed %d certificates\n", len(certs))
for i, cert := range certs {
    fmt.Printf("Certificate %d: %s\n", i+1, cert.Subject.CommonName)
}
```

### Certificate Encoding

```go
// Convert certificate to PEM format
pemData, err := x509x.EncodeCertificateToPEM(cert)
if err != nil {
    log.Fatal(err)
}

// Write to file
err = os.WriteFile("certificate.pem", pemData, 0644)
if err != nil {
    log.Fatal(err)
}
```

## Key Generation and Management

### Generate Key Pairs

```go
// Generate RSA key (2048, 3072, or 4096 bits)
rsaPriv, rsaPub, err := x509x.GenerateKey("rsa", 2048)
if err != nil {
    log.Fatal(err)
}

// Generate ECDSA key (P-256, P-384, P-521 curves)
ecdsaPriv, ecdsaPub, err := x509x.GenerateKey("ecdsa", 256) // P-256
if err != nil {
    log.Fatal(err)
}

// Generate Ed25519 key
ed25519Priv, ed25519Pub, err := x509x.GenerateKey("ed25519", 256)
if err != nil {
    log.Fatal(err)
}
```

### Key Operations

```go
// Parse private key from PEM
keyPEM := `-----BEGIN PRIVATE KEY-----
...
-----END PRIVATE KEY-----`

privKey, err := x509x.ParsePrivateKey([]byte(keyPEM))
if err != nil {
    log.Fatal(err)
}

// Encode private key to PEM
pemData, err := x509x.EncodePrivateKeyToPEM(privKey)
if err != nil {
    log.Fatal(err)
}

// Determine key algorithm
algorithm := x509x.PrivateKeyAlgorithm(privKey)
fmt.Printf("Key algorithm: %s\n", algorithm)
```

## Certificate Signing Requests (CSR)

### Create CSR

```go
// Create CSR template
template := &x509.CertificateRequest{
    Subject: pkix.Name{
        Country:            []string{"US"},
        Province:           []string{"California"},
        Locality:           []string{"San Francisco"},
        Organization:       []string{"Example Corp"},
        OrganizationalUnit: []string{"IT Department"},
        CommonName:         "api.example.com",
    },
    EmailAddresses: []string{"admin@example.com"},
    DNSNames:       []string{"api.example.com", "www.api.example.com"},
    IPAddresses:    []net.IP{net.ParseIP("192.168.1.100")},
}

// Generate key and create CSR
privKey, _, err := x509x.GenerateKey("ecdsa", 256)
if err != nil {
    log.Fatal(err)
}

csrDER, err := x509x.CreateCertificateRequest(template, privKey)
if err != nil {
    log.Fatal(err)
}

// Convert to PEM for transmission
csrPEM := pem.EncodeToMemory(&pem.Block{
    Type:  "CERTIFICATE REQUEST",
    Bytes: csrDER,
})
```

### Parse CSR

```go
csrPEM := `-----BEGIN CERTIFICATE REQUEST-----
...
-----END CERTIFICATE REQUEST-----`

csr, err := x509x.ParseCSR([]byte(csrPEM))
if err != nil {
    log.Fatal(err)
}

fmt.Printf("CSR Subject: %s\n", csr.Subject.CommonName)
fmt.Printf("DNS Names: %v\n", csr.DNSNames)
fmt.Printf("Signature Algorithm: %s\n", csr.SignatureAlgorithm)
```

## Security and Validation

### Algorithm Validation

```go
// Validate certificate algorithm combination
cert, err := x509x.ParseCertificate(certData)
if err != nil {
    log.Fatal(err)
}

if !x509x.ValidCertificateAlgorithm(cert.SignatureAlgorithm) {
    log.Println("Warning: Certificate uses weak or deprecated algorithm")
}
```

### Signature Verification

```go
// Verify digital signatures (ECDSA supported)
pubKey := cert.PublicKey
message := []byte("data to verify")
signature := []byte{...} // signature bytes

valid, err := x509x.VerifySignature(pubKey, message, signature)
if err != nil {
    log.Fatal(err)
}

if valid {
    fmt.Println("Signature is valid")
} else {
    fmt.Println("Signature verification failed")
}
```

## Supported Algorithms

### Signature Algorithms
- **RSA**: SHA256WithRSA, SHA384WithRSA, SHA512WithRSA
- **ECDSA**: ECDSAWithSHA256, ECDSAWithSHA384, ECDSAWithSHA512
- **Ed25519**: PureEd25519

### Key Types and Sizes
- **RSA**: 2048, 3072, 4096 bits (minimum 2048 for security)
- **ECDSA**: P-256, P-384, P-521 curves
- **Ed25519**: 256 bits

### Rejected Algorithms (Security)
- MD2, MD5, SHA1 with RSA (cryptographically broken)
- RSA keys smaller than 2048 bits
- Weak elliptic curves

## Advanced Usage

### Certificate Chain Validation

```go
// Parse and validate certificate chain
chainPEM := `-----BEGIN CERTIFICATE-----
... (end entity certificate)
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
... (intermediate CA)
-----END CERTIFICATE-----
-----BEGIN CERTIFICATE-----
... (root CA)
-----END CERTIFICATE-----`

certs, err := x509x.ParseCertificateChain([]byte(chainPEM))
if err != nil {
    log.Fatal(err)
}

// Build certificate chain for validation
roots := x509.NewCertPool()
intermediates := x509.NewCertPool()

// Add root CA
roots.AddCert(certs[len(certs)-1])

// Add intermediate CAs
for i := 1; i < len(certs)-1; i++ {
    intermediates.AddCert(certs[i])
}

// Verify end entity certificate
opts := x509.VerifyOptions{
    Roots:         roots,
    Intermediates: intermediates,
}

chains, err := certs[0].Verify(opts)
if err != nil {
    log.Printf("Certificate verification failed: %v", err)
} else {
    fmt.Printf("Certificate chain verified: %d chains found\n", len(chains))
}
```

### Utility Functions

```go
// Generate cryptographically secure serial number
serial, err := x509x.RandomSerial()
if err != nil {
    log.Fatal(err)
}

// Use in certificate template
template := &x509.Certificate{
    SerialNumber: serial,
    // ... other fields
}
```

## Error Handling

```go
// The package provides detailed error information
cert, err := x509x.ParseCertificate(invalidData)
if err != nil {
    switch {
    case strings.Contains(err.Error(), "no certificate found"):
        log.Println("No valid certificate in provided data")
    case strings.Contains(err.Error(), "failed to parse"):
        log.Println("Certificate format is invalid")
    default:
        log.Printf("Unexpected error: %v", err)
    }
}
```

## Security Considerations

### Best Practices

1. **Key Generation**: Always use recommended key sizes
   ```go
   // Good: Use strong key sizes
   privKey, _, err := x509x.GenerateKey("ecdsa", 256)  // P-256
   privKey, _, err := x509x.GenerateKey("rsa", 3072)   // 3072-bit RSA
   ```

2. **Algorithm Selection**: Prefer modern algorithms
   ```go
   // Preferred order: Ed25519 > ECDSA > RSA
   if x509x.ValidCertificateAlgorithm(x509.PureEd25519) {
       // Use Ed25519 when supported
   }
   ```

3. **Certificate Validation**: Always validate certificates
   ```go
   if !x509x.ValidCertificateAlgorithm(cert.SignatureAlgorithm) {
       return errors.New("certificate uses weak algorithm")
   }
   ```

### Key Storage Security

```go
// Generate key with secure random source
privKey, _, err := x509x.GenerateKey("ecdsa", 384)
if err != nil {
    log.Fatal(err)
}

// Encode for secure storage
pemData, err := x509x.EncodePrivateKeyToPEM(privKey)
if err != nil {
    log.Fatal(err)
}

// Store with appropriate file permissions
err = os.WriteFile("private.key", pemData, 0600) // Owner read/write only
if err != nil {
    log.Fatal(err)
}
```

## Integration Examples

### TLS Certificate Management

```go
// Generate key and create self-signed certificate
privKey, pubKey, err := x509x.GenerateKey("ecdsa", 256)
if err != nil {
    log.Fatal(err)
}

template := &x509.Certificate{
    SerialNumber: big.NewInt(1),
    Subject: pkix.Name{
        CommonName: "localhost",
    },
    DNSNames:     []string{"localhost"},
    IPAddresses:  []net.IP{net.ParseIP("127.0.0.1")},
    NotBefore:    time.Now(),
    NotAfter:     time.Now().Add(365 * 24 * time.Hour),
    KeyUsage:     x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
    ExtKeyUsage:  []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
}

certDER, err := x509.CreateCertificate(rand.Reader, template, template, pubKey, privKey)
if err != nil {
    log.Fatal(err)
}

// Convert to PEM for use with TLS
certPEM, err := x509x.EncodeCertificateToPEM(&x509.Certificate{Raw: certDER})
if err != nil {
    log.Fatal(err)
}

keyPEM, err := x509x.EncodePrivateKeyToPEM(privKey)
if err != nil {
    log.Fatal(err)
}
```

## Function Reference

|                               |                                         |
| ----------------------------- | --------------------------------------- |
| `ParseCertificate()`          | parse certificate from PEM or DER       |
| `ParseCertificateChain()`     | parse certificate chain from PEM        |
| `ParseCSR()`                  | parse CSR(Certificate Request) from PEM |
| `GenerateKey()`               | generate private, public key pair       |
| `ParsePrivateKey()`           | parse private key from PEM               |
| `CreateCertificateRequest()`  | create CSR with automatic key handling  |
| `EncodeCertificateToPEM()`    | convert certificate to PEM format       |
| `EncodePrivateKeyToPEM()`     | convert private key to PEM format       |
| `VerifySignature()`           | verify digital signatures               |
| `ValidCertificateAlgorithm()` | validate algorithm security compliance  |
| `PrivateKeyAlgorithm()`       | determine private key algorithm         |
| `RandomSerial()`              | generate secure certificate serial      |

---

**Note**: This package is part of the [goxp](https://github.com/whitekid/goxp) utility collection.