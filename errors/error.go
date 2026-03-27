package errors

import (
	"fmt"
	"library/trace"
	"time"
)

type CustomError interface {
	OriginalErr() error
	Message() string
	Time() time.Time
	Trace() []string
	Response
	error
}

type err struct {
	msg string
}

func New(msg string) error {
	return &err{
		msg: msg,
	}
}

func Newf(pattern string, args ...any) error {
	msg := fmt.Sprintf(pattern, args...)
	return New(msg)
}

func (e *err) Error() string {
	return e.msg
}

type customError struct {
	originalErr error
	message     string
	t           time.Time
	trace       []string
	Response
	error
}

func NewCustomError(err error, message string, userMessage Response) CustomError {
	t := trace.New(-1)
	fmt.Printf("custom error frames: %+v\n", t.Frames())
	return customError{
		originalErr: err,
		message:     message,
		t:           time.Now(),
		trace: t.Format(func(f trace.Frame) string {
			return fmt.Sprintf("%s %d %s", f.File(), f.Line(), f.Function())
		}),
		Response: userMessage,
	}
}

func (err customError) Error() string {
	return err.Response.UserMessage()
}

func (err customError) Message() string {
	return err.message
}

func (err customError) OriginalErr() error {
	return err.originalErr
}

func (err customError) Time() time.Time {
	return err.t
}

func (err customError) Trace() []string {
	return err.trace
}

type MockCustomError struct {
	MockError       func() string
	MockMessage     func() string
	MockOriginalErr func() error
	MockTime        func() time.Time
	MockTrace       func() []string
	MockResponse
}

func (err MockCustomError) Error() string {
	return err.MockError()
}

func (err MockCustomError) Message() string {
	return err.MockMessage()
}

func (err MockCustomError) OriginalErr() error {
	return err.MockOriginalErr()
}

func (err MockCustomError) Time() time.Time {
	return err.MockTime()
}

func (err MockCustomError) Trace() []string {
	return err.MockTrace()
}
