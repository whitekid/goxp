// functional and generics...
package goxp

import (
	"math/rand"
	"sync"
	"time"
)

var (
	shuffleOnce sync.Once
	shuffleRand *rand.Rand
)

// Shuffle return shuffled slice
func Shuffle[T any](x []T) []T {
	shuffleOnce.Do(func() {
		shuffleRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	})

	sf := make([]T, len(x))
	copy(sf, x)

	shuffleRand.Shuffle(len(sf), func(i, j int) {
		sf[i], sf[j] = sf[j], sf[i]
	})

	return sf
}

// ShuffleInline return shuffled slice
func ShuffleInline[T any](x []T) []T {
	shuffleOnce.Do(func() {
		shuffleRand = rand.New(rand.NewSource(time.Now().UnixNano()))
	})

	sf := make([]T, len(x))
	copy(sf, x)

	shuffleRand.Shuffle(len(sf), func(i, j int) {
		sf[i], sf[j] = sf[j], sf[i]
	})

	return sf
}
