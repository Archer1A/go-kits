package slice

func ToMap[T any, K comparable](slice []T, keyFunc func(item T) K) map[K]T {
	m := make(map[K]T, len(slice))
	for _, i := range slice {
		m[keyFunc(i)] = i
	}
	return m
}

func ToMapFirstWin[T any, K comparable](slice []T, keyFunc func(item T) K) map[K]T {
	m := make(map[K]T, len(slice))
	for _, i := range slice {
		if _, ok := m[keyFunc(i)]; ok {
			continue
		}
		m[keyFunc(i)] = i
	}
	return m
}
