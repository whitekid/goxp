package utils

import (
	"math/rand"
	"time"
)

const (
	Digits          = "0123456789"
	AsciiUpperCases = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	AsciiLowerCases = "abcdefghijklmnopqrstuvwxyz"
	AsciiLetters    = AsciiLowerCases + AsciiUpperCases
)

var seed = rand.New(rand.NewSource(time.Now().UnixNano()))

// RandomString generate random string
func RandomString(size int) string {
	const charset = AsciiLetters + Digits

	b := make([]byte, size)
	l := len(charset)

	for i := 0; i < size; i++ {
		b[i] = charset[seed.Intn(l)]
	}
	return string(b)
}
