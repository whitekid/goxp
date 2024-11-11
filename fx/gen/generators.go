package gen

import (
	"crypto/rand"
	"math"
	"math/big"

	"golang.org/x/exp/constraints"
)

// Serial int generator
func Serial() Gen[uint64] { return IntN[uint64](math.MaxUint64) }

// IntN generate int to max
func IntN[T constraints.Integer](max T) Gen[T] {
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
func RandInt(max ...int) Gen[int] {
	_max := big.NewInt(math.MaxInt)
	if len(max) > 0 {
		_max = big.NewInt(int64(max[0]))
	}

	return func() (int, bool) {
		v, err := rand.Int(rand.Reader, _max)
		return int(v.Int64()), err != nil
	}
}

// Byte generate readom buffer
func Byte(size int) Gen[[]byte] {
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

func Cycle[T any](seed []T) Gen[T] {
	if len(seed) == 0 {
		return stop[T]()
	}

	i := 0

	return func() (T, bool) {
		j := i % len(seed)
		i++
		return seed[j], true
	}
}

func stop[T any]() func() (T, bool) {
	return func() (T, bool) {
		var v T
		return v, false
	}
}

func Sample[T any](s []T) Gen[T] {
	if len(s) == 0 {
		return stop[T]()
	}

	next := IntN(len(s))
	return func() (T, bool) {
		i, _ := next()
		return s[i], true
	}
}
