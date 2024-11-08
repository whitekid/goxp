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
		current := i
		i++

		return current, true
	}
}

// IntN generate int to max
func IntN[T constraints.Integer](max T) Generator[T] {
	var i T

	return func() (T, bool) {
		if i >= max {
			return 0, false
		}

		current := i
		i++
		return current, true
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
	if size <= 0 {
		return func() ([]byte, bool) { return nil, false }
	}

	return func() ([]byte, bool) {
		buf := make([]byte, size)

		if _, err := rand.Read(buf); err != nil {
			return nil, false
		}

		return buf, true
	}
}

func Cycle[T any](seed []T) Generator[T] {
	if len(seed) == 0 {
		return func() (T, bool) { return *new(T), false }
	}

	i := 0

	return func() (T, bool) {
		j := i % len(seed)
		i++
		return seed[j], true
	}
}
