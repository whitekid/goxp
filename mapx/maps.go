package mapx

func SetNX[K comparable, V any](m map[K]V, k K, v V) bool {
	if _, ok := m[k]; !ok {
		m[k] = v
		return true
	}

	return false
}

func Merge[K comparable, V any](m1 map[K]V, m2 map[K]V) map[K]V {
	r := map[K]V{}

	for k, v := range m1 {
		r[k] = v
	}

	for k, v := range m2 {
		r[k] = v
	}

	return r
}
