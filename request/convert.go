package request

import (
	"library/errors"
	"strconv"
)

func ToString(s string) (string, error) {
	return s, nil
}

func ToInt(s string) (int, error) {
	i, err := strconv.Atoi(s)
	if err != nil {
		return 0, errors.NewWrongTypeErr("int")
	}
	return i, nil
}

func ToFloat64(s string) (float64, error) {
	f, err := strconv.ParseFloat(s, 64)
	if err != nil {
		return 0, errors.NewWrongTypeErr("float64")
	}
	return f, nil
}