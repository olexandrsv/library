package request

import (
	"library/errors"
	"net/http"
)

type Cookie struct {
	r *http.Request
	h errors.ErrorContainer
}

func NewCookie(r *http.Request, h errors.ErrorContainer) *Cookie {
	return &Cookie{
		r: r,
		h: h,
	}
}

func (c *Cookie) GetFloat64(name string) *ValueRequest[float64] {
	return getCookie(c, name, ToFloat64)
}

func (c *Cookie) GetInt(name string) *ValueRequest[int] {
	return getCookie(c, name, ToInt)
}

func (c *Cookie) GetString(name string) *ValueRequest[string] {
	return getCookie(c, name, ToString)
}

func getCookie[T any](c *Cookie, name string, convert func(string) (T, error)) *ValueRequest[T] {
	if c.h.FatalErr() != nil {
		return NewRequestBuilder[T](c.h).Create()
	}

	cookie, err := c.r.Cookie(name)
	if err == http.ErrNoCookie {
		return NewRequestBuilder[T](c.h).
			WithFatalError(errors.NewFieldNotExistsErr(name)).Create()
	}

	v, err := convert(cookie.Value)
	if e, ok := err.(errors.WrongTypeErr); ok {
		return NewRequestBuilder[T](c.h).
			WithError(errors.NewWrongFieldTypeErr(name, e.Type)).Create()
	}
	if err != nil {
		return NewRequestBuilder[T](c.h).WithError(errors.NewUnknownErr(err)).Create()
	}

	return NewRequestBuilder[T](c.h).WithValue(v).Create()
}
