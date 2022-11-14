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

// Zip zip slice pair to mapping
// (key1, key2, key3), (values1, value2, values3) --> (key1: value1), (key2: value2), (key3: value3)
func Zip[K comparable, V any](keys []K, values []V) (r map[K]V) {
	r = make(map[K]V)

	ForEach(keys, func(i int, k K) { r[k] = values[i] })

	return r
}
