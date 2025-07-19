package utils

func Map[T, V any](input []T, fn func(T) V) []V {
	result := make([]V, len(input))
	for i, t := range input {
		result[i] = fn(t)
	}
	return result
}
