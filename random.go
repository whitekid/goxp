package goxp

import (
	crand "crypto/rand"
	"math/big"
	"math/rand"
	"time"

	"github.com/whitekid/goxp/fx"
)

var (
	digits       = []rune("0123456789")
	upperCases   = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lowerCases   = []rune("abcdefghijklmnopqrstuvwxyz")
	letters      = append(lowerCases, upperCases...)
	specialChars = []rune("!@#$%^&*()")
	randomChars  = append(append(letters, digits...), specialChars...)
)

var (
	rnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))
)

// RandomString generate random string
func RandomString(size int) string { return RandomStringWith(size, randomChars) }

func RandomStringWith(size int, source []rune) string {
	l := len(source)
	return string(fx.Times(size, func(i int) rune { return source[rnd.Intn(l)] }))
}

// RandomStringWithCrypto generate random string securly but much slower than RandomString()
func RandomStringWithCrypto(size int) string {
	b := make([]rune, size)
	l := big.NewInt(int64(len(randomChars)))

	for i := 0; i < size; i++ {
		n, _ := crand.Int(crand.Reader, l)
		b[i] = randomChars[int(n.Int64())]
	}
	return string(b)
}

func RandomByte(size int) []byte {
	r := make([]byte, size)
	crand.Read(r)
	return r
}
