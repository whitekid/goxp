package fx

import (
	"math/rand"
	"time"
)

var rnd *rand.Rand = rand.New(rand.NewSource(time.Now().UnixNano()))

// ForEach iteraterate slice and apply function
func ForEach[T any](collection []T, fx func(int, T)) {
	for i, e := range collection {
		fx(i, e)
	}
}

// ForEachE stop foreach if error
func ForEachE[T any](collection []T, fx func(int, T) error) error {
	for i, e := range collection {
		if err := fx(i, e); err != nil {
			return err
		}
	}
	return nil
}

// Filter return filtered slice
func Filter[T any](collection []T, fx func(T) bool) []T {
	if collection == nil {
		return nil
	}

	result := make([]T, len(collection))

	j := 0
	ForEach(collection, func(i int, e T) {
		if !fx(e) {
			return
		}
		result[j] = e
		j++
	})

	return result[:j]
}

// Map map element and return new type
func Map[T any, R any](collection []T, fx func(T) R) []R {
	if collection == nil {
		return nil
	}

	result := make([]R, len(collection))
	ForEach(collection, func(i int, e T) { result[i] = fx(e) })

	return result
}

func Reduce[T any](collection []T, fx func(r T, e T) T) T {
	if len(collection) == 1 {
		return collection[0]
	}

	agg := collection[0]

	ForEach(collection[1:], func(i int, e T) { agg = fx(agg, collection[i+1]) })

	return agg
}

// Times repeat count times
func Times[T any](count int, fx func(int) T) []T {
	result := make([]T, count)

	ForEach(result, func(i int, e T) { result[i] = fx(i) })

	return result
}

// Shuffle return shuffled slice
func Shuffle[T any](collection []T) []T {
	if collection == nil {
		return nil
	}

	sf := make([]T, len(collection))
	copy(sf, collection)

	rnd.Shuffle(len(sf), func(i, j int) { sf[i], sf[j] = sf[j], sf[i] })

	return sf
}

// Distinct return distinct slice
func Distinct[T comparable](collection []T) []T {
	if collection == nil {
		return nil
	}

	set := NewSet[T]()
	set.Append(collection...)
	return set.Slice()
}

// Contains return true if e in collection
func Contains[T comparable](collection []T, e T) bool {
	for _, el := range collection {
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

func Find[T any](collection []T, fx func(T) bool) (T, bool) {
	for _, e := range collection {
		if fx(e) {
			return e, true
		}
	}

	var result T
	return result, false
}

// Every return true if y is subset x
func Every[T comparable](collection, subset []T) bool {
	for _, e := range subset {
		if !Contains(collection, e) {
			return false
		}
	}

	return true
}

func Sample[T any](collection []T) T { return collection[rnd.Intn(len(collection))] }

func Samples[T any](collection []T, count int) []T {
	if collection == nil {
		return nil
	}

	return Map(make([]T, count), func(e T) T { return Sample(collection) })
}

// Keys returns key slice
func Keys[K comparable, V any](mapping map[K]V) []K {
	if mapping == nil {
		return nil
	}

	result := make([]K, len(mapping))

	i := 0
	for k := range mapping {
		result[i] = k
		i++
	}

	return result
}

// Values return values slice
func Values[K comparable, V any](mapping map[K]V) []V {
	if mapping == nil {
		return nil
	}

	result := make([]V, len(mapping))

	i := 0
	for _, v := range mapping {
		result[i] = v
		i++
	}

	return result
}

// ForEachMap iterate map and apply function
func ForEachMap[K comparable, V any](mapping map[K]V, fx func(k K, v V)) {
	for k, v := range mapping {
		fx(k, v)
	}
}

// ForEachMapE stop for each if error
func ForEachMapE[K comparable, V any](mapping map[K]V, fx func(k K, v V) error) error {
	for k, v := range mapping {
		if err := fx(k, v); err != nil {
			return err
		}
	}

	return nil
}

// MapValues map mappings
func MapValues[K comparable, V any, U any](mapping map[K]V, fx func(K) U) map[K]U {
	if mapping == nil {
		return nil
	}

	result := make(map[K]U)

	for k := range mapping {
		result[k] = fx(k)
	}
	return result
}

func MergeMap[K comparable, V any](mapping ...map[K]V) map[K]V {
	if mapping == nil {
		return nil
	}

	result := map[K]V{}

	for _, m := range mapping {
		ForEachMap(m, func(k K, v V) { result[k] = v })
	}

	return result
}

func SampleMap[K comparable, V any](mapping map[K]V) (K, V) {
	n := rand.Intn(len(mapping))

	i := 0
	var rk K
	var rv V
	for k, v := range mapping {
		if i == n {
			rk, rv = k, v
			break
		}
		i++
	}

	return rk, rv
}

// Zip zip slice pair to mapping
// (key1, key2, key3), (values1, value2, values3) --> (key1: value1), (key2: value2), (key3: value3)
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

func TernaryCF[T any](cond bool, trueFn func() T, falseFn func() T) T {
	if cond {
		return trueFn()
	}
	return falseFn()
}

type Int interface {
	uint | uint8 | uint16 | uint32 | uint64 | int | int8 | int16 | int32 | int64
}

type Number interface {
	Int | float64
}

type Ordered interface {
	Number | rune | byte | uintptr
}

func Sum[T Ordered](collection []T) T { return Reduce(collection, func(x T, y T) T { return x + y }) }

func Max[T Ordered](col []T) T {
	switch len(col) {
	case 0:
		panic("empty collection")
	case 1:
		return col[0]
	case 2:
		if col[0] > col[1] {
			return col[0]
		}
		return col[1]
	default:
		return Reduce(col, func(x T, y T) T { return Ternary(x > y, x, y) })
	}
}

func Min[T Ordered](col []T) T {
	switch len(col) {
	case 0:
		panic("empty collection")
	case 1:
		return col[0]
	case 2:
		if col[0] > col[1] {
			return col[1]
		}
		return col[1]
	default:
		return Reduce(col, func(x T, y T) T { return Ternary(x > y, y, x) })
	}
}
