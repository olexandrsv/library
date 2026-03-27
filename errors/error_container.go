package errors

import (
	"bytes"
)

type ErrorContainer interface {
	AddError(error)
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

func (e *errorContainer) IsEmpty() bool {
	if e.fatalError != nil {
		return false
	}
	return len(e.errors) == 0
}

func (e *errorContainer) AddError(err error) {
	if err == nil {
		return
	}
	e.errors = append(e.errors, err)
}

func (e *errorContainer) AddErrors(errSlice []error) {
	for _, err := range errSlice {
		e.AddError(err)
	}
}

func (e *errorContainer) AddFatalError(err error) {
	e.fatalError = err
}

func (e *errorContainer) Add(fatalErr, err error) {
	e.AddFatalError(fatalErr)
	e.AddError(err)
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

func NewMultipleErr(errors []error) *MultipleErr {
	return &MultipleErr{
		Errors: errors,
	}
}

func (e *MultipleErr) Error() string {
	var b bytes.Buffer
	for _, err := range e.Errors {
		b.WriteString(err.Error())
		b.WriteString("\n")
	}
	return b.String()
}

type MockErrorContainer struct {
	MockAddError      func(error)
	MockAddErrors     func([]error)
	MockAddFatalError func(error)
	MockAdd           func(error, error)
	MockFatalErr      func() error
	MockError         func() error
	MockIsEmpty       func() bool
}

func (c *MockErrorContainer) AddError(err error) {
	c.MockAddError(err)
}

func (c *MockErrorContainer) AddErrors(errors []error) {
	c.MockAddErrors(errors)
}

func (c *MockErrorContainer) AddFatalError(err error) {
	c.MockAddFatalError(err)
}

func (c *MockErrorContainer) Add(fatalErr, err error) {
	c.MockAdd(fatalErr, err)
}

func (c *MockErrorContainer) FatalErr() error {
	return c.MockFatalErr()
}

func (c *MockErrorContainer) Error() error {
	return c.MockError()
}

func (c *MockErrorContainer) IsEmpty() bool {
	return c.MockIsEmpty()
}
