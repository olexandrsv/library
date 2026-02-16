package errors

import "bytes"

type ErrorContainer interface {
	AddError(error)
	AddErrors([]error)
	AddFatalError(error)
	FatalErr() error
	Error() error
}

type errorContainer struct {
	errors     []error
	fatalError error
}

func NewErrorContainer() ErrorContainer {
	return &errorContainer{}
}

func (e *errorContainer) AddError(err error) {
	e.errors = append(e.errors, err)
}

func (e *errorContainer) AddErrors(errSlice []error) {
	e.errors = append(e.errors, errSlice...)
}

func (e *errorContainer) AddFatalError(err error) {
	e.fatalError = err
}

func (e *errorContainer) FatalErr() error {
	return e.fatalError
}

func (e *errorContainer) Error() error {
	if len(e.errors) == 0 {
		return nil
	}
	return NewMultipleErr(e.errors)
}

type MultipleErr struct {
	Errors []error
}

func NewMultipleErr(errors []error) MultipleErr {
	return MultipleErr{
		Errors: errors,
	}
}

func (e MultipleErr) Error() string {
	var b bytes.Buffer
	for _, err := range e.Errors{
		b.WriteString(err.Error())
		b.WriteString("\n")
	}
	return b.String()
}
