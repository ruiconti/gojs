package internal

// m.Values()
func MapValues[K comparable, V any](m map[K]V) []V {
	vs := make([]V, 0, len(m))
	for _, v := range m {
		vs = append(vs, v)
	}
	return vs
}

// k := m.Keys()
// v = m.Values()
// m[v[i]] = k[i]
func MapInvert[K comparable, V comparable](m map[K]V) map[V]K {
	inv := make(map[V]K, len(m))
	for k, v := range m {
		inv[v] = k
	}
	return inv
}

// Find returns the first occurrence of x in s, or -1 if x is not present in s.
func Find[T comparable](s []T, x T, compare func(T /* arrayElem */, T /* target */) bool) *T {
	for i, v := range s {
		if compare(v, x) {
			return &s[i]
		}
	}
	return nil
}
