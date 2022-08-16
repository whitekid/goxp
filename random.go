package goxp

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"time"

	"github.com/whitekid/goxp/fx"
)

const (
	digits          = "0123456789"
	asciiUpperCases = "ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	asciiLowerCases = "abcdefghijklmnopqrstuvwxyz"
	asciiLetters    = asciiLowerCases + asciiUpperCases
	specialLetters  = "!@#$%^&*()"
	randomLetters   = asciiLetters + digits + specialLetters
)

var (
	rnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// RandomString generate random string
func RandomString(size int) string { return RandomStringFrom(size, randomLetters) }

func RandomStringFrom(size int, source string) string {
	l := len(source)
	return string(fx.Times(size, func(i int) byte { return source[rnd.Intn(l)] }))
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

func RandomByte(size int) []byte {
	r := make([]byte, size)
	crand.Read(r)
	return r
}
