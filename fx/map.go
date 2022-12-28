package fx

import (
	"math/rand"

	"golang.org/x/exp/maps"
)

// Keys returns key slice
func Keys[M ~map[K]V, K comparable, V any](m M) []K { return maps.Keys(m) }

// Values return values slice
func Values[M ~map[K]V, K comparable, V any](m M) []V { return maps.Values(m) }

func EqualMap[M1, M2 ~map[K]V, K, V comparable](m1 M1, m2 M2) bool { return maps.Equal(m1, m2) }
func EqualMapFunc[M1 ~map[K]V1, M2 ~map[K]V2, K comparable, V1, V2 any](m1 M1, m2 M2, eq func(V1, V2) bool) bool {
	return maps.EqualFunc(m1, m2, eq)
}

func ClearMap[M ~map[K]V, K comparable, V any](m M)                         { maps.Clear(m) }
func CloneMap[M ~map[K]V, K comparable, V any](m M) M                       { return maps.Clone(m) }
func CopyMap[M1 ~map[K]V, M2 ~map[K]V, K comparable, V any](dst M1, src M2) { maps.Copy(dst, src) }
func DeleteMapFunc[M ~map[K]V, K comparable, V any](m M, del func(K, V) bool) {
	maps.DeleteFunc(m, del)
}

type Pair[K comparable, V any] struct {
	Key   K
	Value V
}

func Items[M ~map[K]V, K comparable, V any](m M) []Pair[K, V] {
	r := make([]Pair[K, V], 0, len(m))
	for k, v := range m {
		r = append(r, Pair[K, V]{Key: k, Value: v})
	}

	return r
}

func FromPair[K comparable, V any](p []Pair[K, V]) map[K]V {
	r := make(map[K]V, len(p))
	for i := 0; i < len(p); i++ {
		r[p[i].Key] = p[i].Value
	}

	return r
}

// FilterMap ...
func FilterMap[M ~map[K]V, K comparable, V any](m M, fx func(K, V) bool) M {
	r := make(M)

	for k, v := range m {
		if fx(k, v) {
			r[k] = v
		}
	}

	return r
}

// ForEachMap iterate map and apply function
func ForEachMap[M ~map[K]V, K comparable, V any](m M, fx func(k K, v V)) {
	for k, v := range m {
		fx(k, v)
	}
}

// ForEachMapE stop for each if error
func ForEachMapE[M ~map[K]V, K comparable, V any](m M, fx func(k K, v V) error) error {
	for k, v := range m {
		if err := fx(k, v); err != nil {
			return err
		}
	}

	return nil
}

func MapKeys[M ~map[K]V, K comparable, V any, U comparable](m M, fx func(K, V) U) map[U]V {
	if m == nil {
		return nil
	}

	r := make(map[U]V, len(m))

	for k, v := range m {
		r[fx(k, v)] = v
	}

	return r
}

// MapValues map mappings
func MapValues[M ~map[K]V, K comparable, V any, U any](m M, fx func(K, V) U) map[K]U {
	if m == nil {
		return nil
	}

	r := make(map[K]U, len(m))

	for k, v := range m {
		r[k] = fx(k, v)
	}

	return r
}

func MapItems[M ~map[K1]V1, K1 comparable, V1 any, K2 comparable, V2 any](m M, f func(K1, V1) (K2, V2)) map[K2]V2 {
	r := make(map[K2]V2, len(m))
	for k, v := range m {
		k2, v2 := f(k, v)
		r[k2] = v2
	}
	return r
}

func MapToSlice[M ~map[K]V, K comparable, V any, T any](m M, f func(K, V) T) []T {
	r := make([]T, 0, len(m))
	for k, v := range m {
		r = append(r, f(k, v))
	}
	return r
}

func MergeMap[M ~map[K]V, K comparable, V any](m ...M) M {
	if m == nil {
		return nil
	}

	r := make(M)

	for _, m := range m {
		ForEachMap(m, func(k K, v V) { r[k] = v })
	}

	return r
}

func SampleMap[M ~map[K]V, K comparable, V any](m M) (rk K, rv V) {
	n := rand.Intn(len(m))

	i := 0

	for k, v := range m {
		if i == n {
			return k, v
		}
		i++
	}

	return
}
