package fx

import (
	"math/rand"
	"time"
)

var rnd *rand.Rand

func init() {
	rnd = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func ForEach[T any](collection []T, fx func(int, T)) {
	for i, e := range collection {
		fx(i, e)
	}
}

func Filter[T any](collection []T, fx func(T) bool) []T {
	var result []T

	ForEach(collection, func(i int, e T) {
		if fx(e) {
			result = append(result, e)
		}
	})

	return result
}

func Map[T any, R any](collection []T, fx func(T) R) []R {
	result := make([]R, len(collection))
	ForEach(collection, func(i int, e T) { result[i] = fx(e) })

	return result
}

func Reduce[T any, R any](collection []T, fx func(r R, e T) R) R {
	var agg R

	ForEach(collection, func(i int, e T) { agg = fx(agg, collection[i]) })

	return agg
}

func Times[T any](count int, fx func(int) T) []T {
	result := make([]T, count)

	ForEach(result, func(i int, e T) { result[i] = fx(i) })

	return result
}

// Shuffle return shuffled slice
func Shuffle[T any](x []T) []T {

	sf := make([]T, len(x))
	copy(sf, x)

	rnd.Shuffle(len(sf), func(i, j int) {
		sf[i], sf[j] = sf[j], sf[i]
	})

	return sf
}

func Distinct[T comparable](collection []T) []T {
	set := make(map[T]struct{})
	var result []T

	ForEach(collection, func(_ int, e T) {
		if _, ok := set[e]; !ok {
			set[e] = struct{}{}
			result = append(result, e)
		}
	})

	return result
}

func Contains[T comparable](x []T, e T) bool {
	for _, el := range x {
		if e == el {
			return true
		}
	}
	return false
}

func Index[T comparable](collection []T, e T) int {
	for i := range collection {
		if collection[i] == e {
			return i
		}
	}

	return -1
}

// Every return true if y is subset x
func Every[T comparable](x, y []T) bool {
	for _, ey := range y {
		if !Contains(x, ey) {
			return false
		}
	}

	return true
}

func Sample[T any](collection []T) T { return collection[rnd.Intn(len(collection))] }

func Samples[T any](collection []T, count int) []T {
	return Map(make([]T, count), func(e T) T { return Sample(collection) })
}

func Keys[K comparable, V any](m map[K]V) []K {
	var result []K
	for k := range m {
		result = append(result, k)
	}
	return result
}

func Values[K comparable, V any](m map[K]V) []V {
	var result []V
	for _, v := range m {
		result = append(result, v)
	}
	return result
}

func ForEachMap[K comparable, V any](collection map[K]V, fx func(k K, v V)) {
	for k, v := range collection {
		fx(k, v)
	}
}

func MapMap[K comparable, V any, U any](collection map[K]V, fx func(K) U) map[K]U {
	result := make(map[K]U)
	for k := range collection {
		result[k] = fx(k)
	}
	return result
}

func MapValues[K comparable, V any](m map[K]V, fx func(x V) V) map[K]V {
	var result map[K]V
	for k, v := range m {
		result[k] = fx(v)
	}
	return result
}

func Zip[K comparable, V any](keys []K, values []V) map[K]V {
	var result map[K]V

	ForEach(keys, func(i int, k K) { result[k] = values[i] })

	return result
}

type ifElse[T any] struct {
	cond  func() bool
	value T
}

func If[T any](cond func() bool, value T) *ifElse[T] {
	return &ifElse[T]{
		value: value,
		cond:  cond,
	}
}

func (i *ifElse[T]) Else(elseValue T) T {
	if i.cond() {
		return i.value
	}

	return elseValue
}

func Ternary[T any](cond bool, trueValue T, falseValue T) T {
	if cond {
		return trueValue
	}

	return falseValue
}

func TernaryF[T any](cond func() bool, trueValue T, falseValue T) T {
	return Ternary(cond(), trueValue, falseValue)
}

type Number interface {
	uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int64 | float32 | float64
}

func Sum[T Number](collection []T) T { return Reduce(collection, func(x T, y T) T { return x + y }) }

func Max[T Number](collection []T) T {
	return Reduce(collection, func(x T, y T) T { return Ternary(x > y, x, y) })
}

func Min[T Number](collection []T) T {
	return Reduce(collection, func(x T, y T) T { return Ternary(x > y, y, x) })
}
