package gen

import (
	"crypto/rand"
	"math"
	"math/big"

	"golang.org/x/exp/constraints"
)

// Serial int generator
func Serial() Generator[int] {
	i := 0

	return func() (int, bool) {
		defer func() { i++ }()

		return i, true
	}
}

// IntN generate int to max
func IntN[T constraints.Integer](max T) Generator[T] {
	var i T

	return func() (T, bool) {
		if i >= max {
			return 0, false
		}
		
		defer func() { i++ }()
		return i, true
	}
}

// RandInt generate random int
func RandInt(max ...int) Generator[int] {
	genMax := big.NewInt(math.MaxInt)
	if len(max) > 0 {
		genMax = big.NewInt(int64(max[0]))
	}

	return func() (int, bool) {
		v, err := rand.Int(rand.Reader, genMax)
		return int(v.Int64()), err != nil
	}
}

// Byte generate readom buffer
func Byte(size int) Generator[[]byte] {
	return func() ([]byte, bool) {
		buf := make([]byte, size)

		if _, err := rand.Read(buf); err != nil {
			return nil, false
		}

		return buf, true
	}
}

func Cycle[T any](seed []T) Generator[T] {
	i := 0

	return func() (T, bool) {
		defer func() {
			i++
			i = i % len(seed)
		}()

		return seed[i], true
	}
}
