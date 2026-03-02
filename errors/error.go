package errors

import "time"

type CustomError interface {
	OriginalErr() error
	Message() string
	Time() time.Time
	Tracer
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

func (e *err) Error() string {
	return e.msg
}

type customError struct {
	originalErr error
	message     string
	t           time.Time
	Tracer
	Response
	error
}

func NewCustomError(err error, message string, userMessage Response) CustomError {
	return customError{
		originalErr: err,
		message:     message,
		t:           time.Now(),
		Tracer:      newTrace(),
		Response:    userMessage,
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
