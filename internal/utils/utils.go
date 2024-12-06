package utils

func SliceMap[K any, V any](src []K, mapFunc func(e K) V) []V {
	out := make([]V, 0, len(src))

	for _, e := range src {
		out = append(out, mapFunc(e))
	}

	return out
}
