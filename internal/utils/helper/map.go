package helper

func GetMapCopy[K comparable, V any](original map[K]V) map[K]V {
	copy := make(map[K]V)
	for k, v := range original {
		copy[k] = v
	}
	return copy
}
