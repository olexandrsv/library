package slices

import "slices"

func ConvertToNumber[T, K int32 | int | int64 | float32 | float64](slice []T) []K {
	return Map(slice, func(el T) K {
		return K(el)
	})
}

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

func DeleteFunc[T any](slice []T, f func(t T) bool) []T {
	return slices.DeleteFunc(slice, f)
}

type Iterator[T any] interface {
	Next() (T, bool)
}