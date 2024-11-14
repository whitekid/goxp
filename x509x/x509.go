package x509x

import (
	"bytes"
	"crypto"
	"crypto/ecdsa"
	"crypto/ed25519"
	"crypto/elliptic"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"maps"
	"math/big"
	"slices"
	"sort"

	"github.com/whitekid/goxp/errors"
	"github.com/whitekid/goxp/log"
	"github.com/whitekid/goxp/mapx"
)

const (
	CertificatePEMBlockType              = "CERTIFICATE"
	CrlPEMBlockType                      = "X509 CRL"
	CsrPEMBlockType                      = "CERTIFICATE REQUEST"
	OldCsrPEMBlockType                   = "NEW CERTIFICATE REQUEST"
	RsaPrivateKeyPEMBlockType            = "RSA PRIVATE KEY"
	EcdsaPrivateKeyPEMBlockType          = "EC PRIVATE KEY"
	Pkcs8PrivateKeyPEMBlockType          = "PRIVATE KEY"
	EncryptedPKCS8PrivateKeyPEMBLockType = "ENCRYPTED PRIVATE KEY"

	pemPrefix = "-----BEGIN "
)

var (
	pemPrefixCertificate     = []byte(pemPrefix + CertificatePEMBlockType)
	pemPrefixCSR             = []byte(pemPrefix + CsrPEMBlockType)
	pemPrefixRsaPrivateKey   = []byte(pemPrefix + RsaPrivateKeyPEMBlockType)
	pemPrefixEcdsaPrivateKey = []byte(pemPrefix + EcdsaPrivateKeyPEMBlockType)
	pemPrefixPkcs8PrivateKey = []byte(pemPrefix + Pkcs8PrivateKeyPEMBlockType)
)

// ParseCertificate parse x509 certificate PEM block or DER bytes
func ParseCertificate(certBytes []byte) (*x509.Certificate, error) {
	if bytes.HasPrefix(certBytes, pemPrefixCertificate) {
		p, _ := pem.Decode(certBytes)
		if p == nil {
			return nil, errors.New("invalid PEM")
		}

		certBytes = p.Bytes
	}

	return x509.ParseCertificate(certBytes)
}

func ParseCertificateChain(pemBytes []byte) ([]*x509.Certificate, error) {
	certs := make([]*x509.Certificate, 0)
	for {
		p, rest := pem.Decode(pemBytes)
		if p == nil {
			return certs, nil
		}

		cert, err := ParseCertificate(p.Bytes)
		if err != nil {
			return nil, fmt.Errorf("certificate parse failed: %+w", err)
		}
		certs = append(certs, cert)
		pemBytes = rest
	}
}

// ParseCSR parse x509 CSR PEM block
func ParseCSR(csrBytes []byte) (*x509.CertificateRequest, error) {
	if bytes.HasPrefix(csrBytes, pemPrefixCSR) {
		p, _ := pem.Decode(csrBytes)
		if p == nil {
			return nil, errors.New("invalid PEM")
		}

		csrBytes = p.Bytes
	}

	csr, err := x509.ParseCertificateRequest(csrBytes)
	if err != nil {
		return nil, err
	}

	return csr, nil
}

// PublicKey PrivateKey and Signer interfaces
// crypto.PrivateKey의 설명을 보면 함수가 다 있다고 함.
type PrivateKey interface {
	crypto.PrivateKey
	crypto.Signer
}

type PublicKey interface {
	Equal(x crypto.PublicKey) bool
}

// GenerateKey generate private and public key pair
func GenerateKey(algorithm x509.SignatureAlgorithm) (privateKey PrivateKey, err error) {
	switch algorithm {
	case x509.ECDSAWithSHA256:
		privateKey, err = ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
	case x509.ECDSAWithSHA384:
		privateKey, err = ecdsa.GenerateKey(elliptic.P384(), rand.Reader)
	case x509.ECDSAWithSHA512:
		privateKey, err = ecdsa.GenerateKey(elliptic.P521(), rand.Reader)
	case x509.SHA256WithRSA:
		privateKey, err = rsa.GenerateKey(rand.Reader, 256*8)
	case x509.SHA384WithRSA:
		privateKey, err = rsa.GenerateKey(rand.Reader, 384*8)
	case x509.SHA512WithRSA:
		privateKey, err = rsa.GenerateKey(rand.Reader, 512*8)
	case x509.PureEd25519:
		_, privateKey, err = ed25519.GenerateKey(rand.Reader)
	default:
		return nil, fmt.Errorf("unknown algorithm: %s", algorithm.String())
	}

	if err != nil {
		return nil, err
	}

	return
}

// ParsePrivateKey parse pem formatted priate key
func ParsePrivateKey(keyPemBytes []byte) (PrivateKey, error) {
	p, _ := pem.Decode(keyPemBytes)
	if p == nil {
		return nil, errors.New("invalid PEM")
	}

	var key PrivateKey
	var err error
	switch {
	case bytes.HasPrefix(keyPemBytes, pemPrefixRsaPrivateKey):
		key, err = x509.ParsePKCS1PrivateKey(p.Bytes)

	case bytes.HasPrefix(keyPemBytes, pemPrefixEcdsaPrivateKey):
		key, err = x509.ParseECPrivateKey(p.Bytes)

	default:
		return nil, errors.New("unknown pem type")
	}

	if err != nil {
		return nil, fmt.Errorf("fail to parse private key, %+w", err)
	}
	return key, nil
}

// PrivateKeyAlgorithm return private key algorithm
func PrivateKeyAlgorithm(priv PrivateKey) x509.SignatureAlgorithm {
	switch p := priv.(type) {
	case *ecdsa.PrivateKey:
		if alg, ok := map[int]x509.SignatureAlgorithm{
			256: x509.ECDSAWithSHA256,
			384: x509.ECDSAWithSHA384,
			521: x509.ECDSAWithSHA512,
		}[p.Params().BitSize]; ok {
			return alg
		}

	case *rsa.PrivateKey:
		if alg, ok := map[int]x509.SignatureAlgorithm{
			256: x509.SHA256WithRSA,
			384: x509.SHA384WithRSA,
			512: x509.SHA512WithRSA,
		}[p.Size()]; ok {
			return alg
		}

	case ed25519.PrivateKey:
		return x509.PureEd25519
	}

	return x509.UnknownSignatureAlgorithm
}

// CreateCertificateRequest create CSR and return PEM
// privateKey: signer private key
func CreateCertificateRequest(template *x509.CertificateRequest) (csr []byte, pemBytes []byte, err error) {
	privKey, err := GenerateKey(template.SignatureAlgorithm)
	if err != nil {
		return nil, nil, err
	}

	derBytes, err := x509.CreateCertificateRequest(rand.Reader, template, privKey)
	if err != nil {
		return nil, nil, err
	}

	block := &pem.Block{
		Type:    CsrPEMBlockType,
		Headers: nil,
		Bytes:   derBytes,
	}
	return derBytes, pem.EncodeToMemory(block), nil
}

func EncodeCertificateToPEM(derBytes []byte) []byte {
	return pem.EncodeToMemory(&pem.Block{
		Type:    CertificatePEMBlockType,
		Headers: nil,
		Bytes:   derBytes,
	})
}

func EncodePrivateKeyToPEM(privateKey PrivateKey) ([]byte, error) {
	var pemType string
	var keyBytes []byte

	switch key := privateKey.(type) {
	case *rsa.PrivateKey:
		pemType = RsaPrivateKeyPEMBlockType
		keyBytes = x509.MarshalPKCS1PrivateKey(key)
	case *ecdsa.PrivateKey:
		pemType = EcdsaPrivateKeyPEMBlockType
		derBytes, err := x509.MarshalECPrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("fail to encode private key, %+w", err)
		}
		keyBytes = derBytes
	case ed25519.PrivateKey:
		pemType = EcdsaPrivateKeyPEMBlockType
		derBytes, err := x509.MarshalPKCS8PrivateKey(key)
		if err != nil {
			return nil, fmt.Errorf("fail to encode private key, %+w", err)
		}
		keyBytes = derBytes
	default:
		return nil, fmt.Errorf("unsupported private key: %T", privateKey)
	}

	return pem.EncodeToMemory(&pem.Block{
		Type:  pemType,
		Bytes: keyBytes,
	}), nil
}

var (
	keyUsageToStr = map[x509.KeyUsage]string{
		x509.KeyUsageDigitalSignature:  "Digital Signature",
		x509.KeyUsageContentCommitment: "Non Repudiation",
		x509.KeyUsageKeyEncipherment:   "Key Encipherment",
		x509.KeyUsageDataEncipherment:  "Data Encipherment",
		x509.KeyUsageKeyAgreement:      "Key Agreement",
		x509.KeyUsageCertSign:          "Certificate Sign",
		x509.KeyUsageCRLSign:           "CRL Sign",
		x509.KeyUsageEncipherOnly:      "Encipher Only",
		x509.KeyUsageDecipherOnly:      "Decipher Only",
	}
	strToKeyUsage = make(map[string]x509.KeyUsage, len(keyUsageToStr))

	extKeyUsageToStr = map[x509.ExtKeyUsage]string{
		x509.ExtKeyUsageAny:                            "UsageA ny",
		x509.ExtKeyUsageServerAuth:                     "TLS Web Server Authentication",
		x509.ExtKeyUsageClientAuth:                     "TLS Web Client Authentication",
		x509.ExtKeyUsageCodeSigning:                    "Code Signing",
		x509.ExtKeyUsageEmailProtection:                "Email Protection",
		x509.ExtKeyUsageIPSECEndSystem:                 "IPSEC End System",
		x509.ExtKeyUsageIPSECTunnel:                    "IPSEC Tunnel",
		x509.ExtKeyUsageIPSECUser:                      "IPSEC User",
		x509.ExtKeyUsageTimeStamping:                   "Time Stamping",
		x509.ExtKeyUsageOCSPSigning:                    "OCSP Signing",
		x509.ExtKeyUsageMicrosoftServerGatedCrypto:     "Microsoft Server Gated Crypto",
		x509.ExtKeyUsageNetscapeServerGatedCrypto:      "Netscape Server Gated Crypto",
		x509.ExtKeyUsageMicrosoftCommercialCodeSigning: "Microsoft Commercial Code Signing",
		x509.ExtKeyUsageMicrosoftKernelCodeSigning:     "Microsoft Kernel Code Signing",
	}
	strToExtKeyUsage = make(map[string]x509.ExtKeyUsage, len(extKeyUsageToStr))
)

func init() {
	keyUsages := slices.Collect(maps.Keys(keyUsageToStr))
	sort.Slice(keyUsages, func(i, j int) bool { return int(keyUsages[i]) < int(keyUsages[j]) })
	mapx.Each(keyUsageToStr, func(k x509.KeyUsage, v string) { strToKeyUsage[v] = k })
	extKeyUsages := slices.Collect(maps.Keys(extKeyUsageToStr))
	sort.Slice(extKeyUsages, func(i, j int) bool { return int(extKeyUsages[i]) < int(extKeyUsages[j]) })
	mapx.Each(extKeyUsageToStr, func(k x509.ExtKeyUsage, v string) { strToExtKeyUsage[v] = k })
}

func KeyUsageToStr(keyUsage x509.KeyUsage) string       { return keyUsageToStr[keyUsage] }
func StrToKeyUsage(s string) x509.KeyUsage              { return strToKeyUsage[s] }
func ExtKeyUsageToStr(keyUsage x509.ExtKeyUsage) string { return extKeyUsageToStr[keyUsage] }
func StrToExtKeyUsage(s string) x509.ExtKeyUsage        { return strToExtKeyUsage[s] }

func RandomSerial() *big.Int {
	s, _ := rand.Int(rand.Reader, new(big.Int).Lsh(big.NewInt(1), 128))
	return s
}

func VerifySignature(pub PublicKey, hash []byte, signature []byte) bool {
	switch p := pub.(type) {
	case *ecdsa.PublicKey:
		return ecdsa.VerifyASN1(p, hash, signature)
	default:
		log.Fatalf("unsupported public key: %T", pub)
		return false
	}
}

var (
	unsupportedAlgorithm = []x509.SignatureAlgorithm{x509.MD2WithRSA, x509.MD5WithRSA, x509.SHA1WithRSA}
	leafAlgorithms       = []x509.SignatureAlgorithm{x509.SHA256WithRSA, x509.SHA384WithRSA, x509.SHA512WithRSA, x509.ECDSAWithSHA256, x509.ECDSAWithSHA384, x509.ECDSAWithSHA512, x509.PureEd25519}
)

// ValidCertificateAlgorithm check
func ValidCertificateAlgorithm(isCA bool, keyAlgorithm x509.SignatureAlgorithm, signatureAlgorithm x509.SignatureAlgorithm) error {
	// unsupported
	if slices.Contains(unsupportedAlgorithm, keyAlgorithm) || slices.Contains(unsupportedAlgorithm, signatureAlgorithm) {
		return fmt.Errorf("invalid algorithm: %s", keyAlgorithm.String())
	}

	if !isCA {
		if !slices.Contains(leafAlgorithms, keyAlgorithm) {
			return fmt.Errorf("invalid key algorithm: %s", keyAlgorithm.String())
		}

		if signatureAlgorithm != x509.UnknownSignatureAlgorithm && !slices.Contains(leafAlgorithms, keyAlgorithm) {
			return fmt.Errorf("invalid signature algorithm: %s", signatureAlgorithm.String())
		}
	}

	return nil
}
