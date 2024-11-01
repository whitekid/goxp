package goxp

import (
	"crypto/rand"
	"math/big"
)

var (
	digits       = []rune("0123456789")
	upperCases   = []rune("ABCDEFGHIJKLMNOPQRSTUVWXYZ")
	lowerCases   = []rune("abcdefghijklmnopqrstuvwxyz")
	letters      = append(lowerCases, upperCases...)
	specialChars = []rune("!@#$%^&*()")
	randomChars  = append(append(letters, digits...), specialChars...)
)

// RandomString generate random string
func RandomString(size int) string { return RandomStringWith(size, randomChars) }

func RandomStringWith(size int, source []rune) string {
	if size < 0 {
		panic("size must be greater than 0")
	}

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

func RandomByte(size int) []byte {
	if size < 0 {
		panic("size must be greater than 0")
	}

	r := make([]byte, size)
	rand.Read(r)
	return r
}
