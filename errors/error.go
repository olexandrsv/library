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
	return customError{
		originalErr: err,
		message:     message,
		t:           time.Now(),
		trace: trace.New(-1).Format(func(f trace.Frame) string {
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
