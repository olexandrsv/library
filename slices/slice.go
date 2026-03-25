package slices

import "library/errors"

func Convert[T, K any](slice []T, f func(t T) (K, error)) ([]K, error) {
	convertedSlice := make([]K, 0, len(slice))
	for _, el := range slice {
		converted, err := f(el)
		if err != nil {
			return nil, err
		}
		convertedSlice = append(convertedSlice, converted)
	}
	return convertedSlice, nil
}

func Map[T, K any](slice []T, f func(t T) K) []K {
	converted := make([]K, 0, len(slice))
	for _, el := range slice {
		converted = append(converted, f(el))
	}
	return converted
}

func Filter[T any](slice []T, f func(t T) bool) ([]T, error) {
	if f == nil {
		return nil, errors.NewNilErr("f", errors.NewInternalError())
	}
	var sorted []T
	for _, item := range slice {
		if f(item) {
			sorted = append(sorted, item)
		}
	}
	return sorted, nil
}
