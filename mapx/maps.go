package mapx

import (
	"iter"
	"maps"
	"slices"

	"github.com/whitekid/goxp/slicex"
)

type MapX[M ~map[K]V, K comparable, V any] map[K]V

func Of[K comparable, V any](s iter.Seq2[K, V]) MapX[map[K]V, K, V] { return maps.Collect(s) }

func (m MapX[M, K, V]) All() iter.Seq2[K, V] { return maps.All(m) }

// Map iterate map and apply function
func (m MapX[M, K, V]) Each(fx func(K, V)) { Each(m, fx) }
func Each[M ~map[K]V, K comparable, V any](m M, fx func(k K, v V)) {
	for k, v := range m {
		fx(k, v)
	}
}

func (m MapX[M, K, V]) Clone() MapX[M, K, V]    { return maps.Clone(m) }
func (m MapX[M, K, V]) Copy(dest MapX[M, K, V]) { maps.Copy(m, dest) }

// Filter ...
func (m MapX[M, K, V]) Filter(fx func(K, V) bool) MapX[M, K, V] { return Filter(m, fx) }
func Filter[M ~map[K]V, K comparable, V any](m M, fx func(K, V) bool) M {
	r := make(M)

	for k, v := range m {
		if fx(k, v) {
			r[k] = v
		}
	}

	return r
}

func (m MapX[M, K, V]) Insert(s iter.Seq2[K, V]) MapX[M, K, V] {
	maps.Insert(m, s)
	return m
}

func (m MapX[M, K, V]) Keys() iter.Seq[K] { return maps.Keys(m) }

func Map[M1 ~map[K1]V1, K1 comparable, V1 any, M2 ~map[K2]V2, K2 comparable, V2 any](m M1, f func(K1, V1) (K2, V2)) M2 {
	r := make(map[K2]V2, len(m))

	for k, v := range m {
		k2, v2 := f(k, v)
		r[k2] = v2
	}

	return r
}

func MapKey[K1 comparable, V any, K2 comparable](m map[K1]V, fx func(K1, V) K2) map[K2]V {
	if m == nil {
		return nil
	}

	r := make(map[K2]V, len(m))

	for k, v := range m {
		r[fx(k, v)] = v
	}

	return r
}

// MapValue map mappings
func MapValue[K comparable, V1 any, V2 any](m map[K]V1, fx func(K, V1) V2) map[K]V2 {
	if m == nil {
		return nil
	}

	r := make(map[K]V2, len(m))

	for k, v := range m {
		r[k] = fx(k, v)
	}

	return r
}

// Merge merge to map
func (m MapX[M, K, V]) Merge(m2 ...M) MapX[M, K, V] { return MapX[M, K, V](Merge(m, m2...)) }
func Merge[M ~map[K]V, K comparable, V any](m1 map[K]V, m2 ...M) M {
	r := map[K]V{}

	for k, v := range m1 {
		r[k] = v
	}

	for _, mm := range m2 {
		for k, v := range mm {
			r[k] = v
		}
	}

	return r
}

func (m MapX[M, K, V]) Sample() (K, V) { return Sample(m) }
func Sample[M ~map[K]V, K comparable, V any](m M) (K, V) {
	k := slicex.Sample(slices.Collect(maps.Keys(m)))
	return k, m[k]
}

func (m MapX[M, K, V]) SetNx(k K, v V) bool { return SetNX[M](m, k, v) }
func SetNX[M ~map[K]V, K comparable, V any](m map[K]V, k K, v V) bool {
	if _, ok := m[k]; !ok {
		m[k] = v
		return true
	}

	return false
}

// Slice convert map to slice
func Slice[M ~map[K]V, K comparable, V any, T any](m M, f func(K, V) T) []T {
	r := make([]T, 0, len(m))
	for k, v := range m {
		r = append(r, f(k, v))
	}
	return slices.Clip(r)
}

func (m MapX[M, K, V]) Values() iter.Seq[V] { return maps.Values(m) }
