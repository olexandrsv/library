package errors

type Response interface {
	UserMessage() string
	UserCode() int
}

type userError struct {
	msg  string
	code int
}

func NewResponse(msg string, code int) *userError {
	return &userError{
		msg:  msg,
		code: code,
	}
}

func (err *userError) UserMessage() string {
	return err.msg
}

func (err *userError) UserCode() int {
	return err.code
}

func NewForbiddenError() Response {
	return NewResponse("Forbidden", 403)
}

func NewInternalError() Response {
	return NewResponse("Internal error", 500)
}