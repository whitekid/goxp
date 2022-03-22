package utils

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"time"
)

const (
	Digits          = "0123456789"
	AsciiUpperCases = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AsciiLowerCases = "abcdefghijklmnopqrstuvwxyz"
	AsciiLetters    = AsciiLowerCases + AsciiUpperCases
	randomLetters   = AsciiLetters + Digits
)

var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomString generate random string
func RandomString(size int) string {
	b := make([]byte, size)
	l := len(randomLetters)

	for i := 0; i < size; i++ {
		b[i] = randomLetters[seed.Intn(l)]
	}
	return string(b)
}

func RandomStringWithCrypto(size int) string {
	b := make([]byte, size)
	l := big.NewInt(int64(len(randomLetters)))

	for i := 0; i < size; i++ {
		n, _ := crand.Int(crand.Reader, l)
		b[i] = randomLetters[int(n.Int64())]
	}
	return string(b)
}
