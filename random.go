package goxp

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"sync"
	"time"
)

const (
	Digits          = "0123456789"
	AsciiUpperCases = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AsciiLowerCases = "abcdefghijklmnopqrstuvwxyz"
	AsciiLetters    = AsciiLowerCases + AsciiUpperCases
	specialLetters  = "!@#$%^&*()"
	randomLetters   = AsciiLetters + Digits + specialLetters
)

var (
	seed     *rand.Rand
	seedOnce sync.Once
)

// RandomString generate random string
func RandomString(size int) string {
	seedOnce.Do(func() {
		seed = rand.New(rand.NewSource(time.Now().UnixNano()))
	})

	b := make([]byte, size)
	l := len(randomLetters)

	for i := 0; i < size; i++ {
		b[i] = randomLetters[seed.Intn(l)]
	}
	return string(b)
}

// RandomStringWithCrypto generate random string securly but much slower than RandomString()
func RandomStringWithCrypto(size int) string {
	b := make([]byte, size)
	l := big.NewInt(int64(len(randomLetters)))

	for i := 0; i < size; i++ {
		n, _ := crand.Int(crand.Reader, l)
		b[i] = randomLetters[int(n.Int64())]
	}
	return string(b)
}
