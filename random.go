package goxp

import (
	"crypto/rand"
	"math/big"

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

var ()

// RandomString generate random string
func RandomString(size int) string                    { return RandomStringWith(size, randomChars) }
func RandomStringWith(size int, source []rune) string { return randomStringWithRand(size, source) }

func randomStringWithRand(size int, source []rune) string {
	l := int64(len(source))
	r := make([]rune, size)

	max := big.NewInt(l)
	for i := 0; i < size; i++ {
		v, err := rand.Int(rand.Reader, max)
		if err != nil {
			panic(err)
		}

		r[i] = source[v.Int64()%l]
	}

	return string(r)
}

// randomStringWithCrypto generate random string securly but much slower than RandomString()
func randomStringWithCrypto(size int, source []rune) string {
	b := make([]rune, fx.Abs(size))
	l := big.NewInt(int64(len(source)))

	for i := 0; i < size; i++ {
		n, _ := rand.Int(rand.Reader, l)
		b[i] = source[int(n.Int64())]
	}
	return string(b)
}

func RandomByte(size int) []byte { return randomByteWithRand(size) }

func randomByteWithRand(size int) []byte {
	r := make([]byte, fx.Abs(size))
	rand.Read(r)
	return r
}

func randomByteWithCrypto(size int) []byte {
	r := make([]byte, fx.Abs(size))
	rand.Read(r)
	return r
}
