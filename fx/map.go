package fx

import "math/rand"

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

// FilterMap ...
func FilterMap[K comparable, V any](m map[K]V, fx func(K, V) bool) map[K]V {
	r := make(map[K]V)

	for k, v := range m {
		if fx(k, v) {
			r[k] = v
		}
	}

	return r
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

func MapKeys[K comparable, V any, U comparable](mapping map[K]V, fx func(K, V) U) map[U]V {
	if mapping == nil {
		return nil
	}

	result := make(map[U]V)

	for k, v := range mapping {
		result[fx(k, v)] = v
	}
	return result
}

// MapValues map mappings
func MapValues[K comparable, V any, U any](mapping map[K]V, fx func(K, V) U) map[K]U {
	if mapping == nil {
		return nil
	}

	result := make(map[K]U)

	for k, v := range mapping {
		result[k] = fx(k, v)
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

func SampleMap[K comparable, V any](mapping map[K]V) (rk K, rv V) {
	n := rand.Intn(len(mapping))

	i := 0

	for k, v := range mapping {
		if i == n {
			return k, v
		}
		i++
	}

	return
}
