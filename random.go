package goxp

import (
	"crypto/rand"
)

var (
	digits       = []rune("0123456789")
	upperCases   = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lowerCases   = []rune("abcdefghijklmnopqrstuvwxyz")
	letters      = append(lowerCases, upperCases...)
	specialChars = []rune("!@#$%^&*()")
	randomChars  = append(append(letters, digits...), specialChars...)
)

// RandomString generate cryptographically secure random string
func RandomString(size int) string { return RandomStringWith(size, randomChars) }

// RandomStringWith generates a cryptographically secure random string from the given source runes.
// It uses crypto/rand for security-critical applications.
// For better performance in non-security contexts, consider using math/rand/v2.
func RandomStringWith(size int, source []rune) string {
	if size < 0 {
		return ""
	}
	// Prevent excessive memory allocation - limit to 1MB of runes
	const maxSize = 1024 * 1024 / 4 // 4 bytes per rune
	if size > maxSize {
		panic("size too large: maximum allowed is 262144")
	}
	if len(source) == 0 {
		panic("source cannot be empty")
	}

	r := make([]rune, size)
	bytes := make([]byte, size)

	// Batch read random bytes for better performance
	_, err := rand.Read(bytes)
	Must(err)

	// Map bytes to source runes
	sourceLen := len(source)
	for i, b := range bytes {
		r[i] = source[int(b)%sourceLen]
	}

	return string(r)
}

func RandomByte(size int) []byte {
	if size < 0 {
		panic("size must be greater than 0")
	}

	r := make([]byte, size)
	rand.Read(r)
	return r
}
