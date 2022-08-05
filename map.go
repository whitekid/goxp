package goxp

func SetNX[K comparable, V any](m map[K]V, k K, v V) bool {
	if _, ok := m[k]; !ok {
		m[k] = v
		return true
	}

	return false
}
