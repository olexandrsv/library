package errors

import (
	"fmt"
)

type WrongValueRangeError struct {
	CustomError
}

func NewWrongValueRangeError[T any](value, minValue, maxValue T) *WrongFieldTypeErr {
	msg := fmt.Sprintf("variable is expected in range '%v' and '%v', but got '%v'", minValue, maxValue, value)
	return &WrongFieldTypeErr{
		CustomError: NewCustomError(nil, msg, NewInternalError()),
	}
}

type ParsingError struct {
	CustomError
}

func NewParsingError(name string, kind string, resp Response) *ParsingError {
	msg := fmt.Sprintf("error while parsing field '%s' with type '%s'", name, kind)
	return &ParsingError{
		CustomError: NewCustomError(nil, msg, resp),
	}
}

type WrongValueSizeErr struct {
	CustomError
}

func NewWrongValueSizeError(name string, size int, expectedSize string) *WrongValueSizeErr {
	msg := fmt.Sprintf("value '%s' has size '%d', while expected size '%s'",
		name, size, expectedSize)
	return &WrongValueSizeErr{
		CustomError: NewCustomError(nil, msg, NewResponse(msg, 400)),
	}
}

type FieldNotExistsErr struct {
	CustomError
}

func NewFieldNotExistsErr(name string) *FieldNotExistsErr {
	msg := fmt.Sprintf("field %s not exists", name)
	return &FieldNotExistsErr{
		CustomError: NewCustomError(nil, msg, NewResponse(msg, 0)),
	}
}

type WrongTypeErr struct {
	Type string
	CustomError
}

func NewWrongTypeErr(fieldType string) *WrongTypeErr {
	msg := "can't convert to type " + fieldType
	return &WrongTypeErr{
		Type:        fieldType,
		CustomError: NewCustomError(nil, msg, NewResponse(msg, 400)),
	}
}

type WrongFieldTypeErr struct {
	Name string
	Type string
	CustomError
}

func NewWrongFieldTypeErr(name, fieldType string) *WrongFieldTypeErr {
	msg := fmt.Sprintf("can't convert field %s to type %s", name, fieldType)
	return &WrongFieldTypeErr{
		Name:        name,
		Type:        fieldType,
		CustomError: NewCustomError(nil, msg, NewResponse(msg, 400)),
	}
}

type UnknownErr struct {
	CustomError
}

func NewUnknownErr(err error) *UnknownErr {
	return &UnknownErr{
		CustomError: NewCustomError(err, "unknown error", NewInternalError()),
	}
}

type DataParseErr struct {
	DataType string
	CustomError
}

func NewDataParseErr(dataType string) *DataParseErr {
	msg := fmt.Sprintf("can't parse data: %s", dataType)
	return &DataParseErr{
		DataType:    dataType,
		CustomError: NewCustomError(nil, msg, NewResponse(msg, 400)),
	}
}

type NewFileError struct {
	CustomError
}

func NewFileErr(err error, filePath string, resp Response) *NewFileError {
	msg := fmt.Sprintf("got error while working with '%s' file", filePath)
	return &NewFileError{
		CustomError: NewCustomError(err, msg, resp),
	}
}

type NilError struct {
	CustomError
}

func NewNilErr(variableName string, response Response) *NilError {
	msg := fmt.Sprintf("variable '%s' can't be nil", variableName)
	return &NilError{
		CustomError: NewCustomError(nil, msg, response),
	}
}
