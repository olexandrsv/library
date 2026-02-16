package request

import (
	"library/errors"
	"library/slices"
	"mime/multipart"
	"net/http"
)

type Request struct {
	Form   *Form
	Cookie *Cookie
	r      *http.Request
	h      errors.ErrorContainer
}

func New(r *http.Request) *Request {
	h := errors.NewErrorContainer()
	return &Request{
		Form: &Form{
			r: r,
			h: h,
		},
		Cookie: NewCookie(r, h),
		r:      r,
		h:      h,
	}
}

func (r *Request) Err() error {
	if r.h.FatalErr() != nil {
		return r.h.FatalErr()
	}
	return r.h.Error()
}

type Form struct {
	r *http.Request
	h errors.ErrorContainer
}

func (f *Form) GetIntSlice(name string) *ValueRequest[[]int] {
	return getSlice(f, name, ToInt)
}

func (f *Form) GetStringSlice(name string) *ValueRequest[[]string] {
	return getSlice(f, name, ToString)
}

func (f *Form) GetFloat64(name string) *ValueRequest[float64] {
	return get(f, name, ToFloat64)
}

func (f *Form) GetString(name string) *ValueRequest[string] {
	return get(f, name, ToString)
}

func (f *Form) GetInt(name string) *ValueRequest[int] {
	return get(f, name, ToInt)
}


func (f *Form) parse() error {
	if f.r.MultipartForm == nil {
		if err := f.r.ParseMultipartForm(1 << 20); err != nil {
			return errors.NewDataParseErr("form")
		}
	}
	return nil
}

func get[T any](f *Form, name string, convert func(string) (T, error)) *ValueRequest[T] {
	if f.h.FatalErr() != nil {
		return NewRequestBuilder[T](f.h).Create()
	}

	values, fatalErr, err := getValues(f, name)
	if fatalErr != nil || err != nil {
		return NewRequestBuilder[T](f.h).WithFatalError(fatalErr).WithError(err).Create()
	}

	if len(values) != 1 {
		return NewRequestBuilder[T](f.h).
			WithError(errors.NewWrongValueSizeError(name, len(values), "1")).
			Create()
	}

	rawValue := values[0]
	value, err := convert(rawValue)
	if e, ok := err.(errors.WrongTypeErr); ok {
		return NewRequestBuilder[T](f.h).
			WithError(errors.NewWrongFieldTypeErr(name, e.Type)).Create()
	}
	if err != nil {
		return NewRequestBuilder[T](f.h).WithError(errors.NewUnknownErr(err)).Create()
	}

	return NewRequestBuilder[T](f.h).WithValue(value).Create()
}

func getFile(f *Form, name string) *ValueRequest[*multipart.FileHeader] {
	if f.h.FatalErr() != nil {
		return NewRequestBuilder[*multipart.FileHeader](f.h).Create()
	}

	files, fatalErr, err := getFiles(f, name)
	if fatalErr != nil || err != nil {
		return NewRequestBuilder[*multipart.FileHeader](f.h).
			WithFatalError(fatalErr).WithError(err).Create()
	}

	if len(files) != 1 {
		return NewRequestBuilder[*multipart.FileHeader](f.h).
			WithError(errors.NewWrongValueSizeError(name, len(files), "1")).
			Create()
	}

	return NewRequestBuilder[*multipart.FileHeader](f.h).WithValue(files[0]).Create()
}

func getValues(f *Form, name string) ([]string, error, error) {
	if err := f.parse(); err != nil {
		return nil, err, nil
	}

	values, ok := f.r.MultipartForm.Value[name]
	if !ok {
		return nil, nil, errors.NewFieldNotExistsErr(name)
	}

	return values, nil, nil
}

func getFiles(f *Form, name string) ([]*multipart.FileHeader, error, error) {
	if err := f.parse(); err != nil {
		return nil, err, nil
	}

	files, ok := f.r.MultipartForm.File[name]
	if !ok {
		return nil, nil, errors.NewFieldNotExistsErr(name)
	}

	return files, nil, nil
}

func getSlice[T any](f *Form, name string, convert func(string) (T, error)) *ValueRequest[[]T] {
	if f.h.FatalErr() != nil {
		return NewRequestBuilder[[]T](f.h).Create()
	}

	rawValues, fatalErr, err := getValues(f, name)
	if fatalErr != nil || err != nil {
		return NewRequestBuilder[[]T](f.h).WithFatalError(fatalErr).WithError(err).Create()
	}

	values, err := slices.Convert(rawValues, convert)
	if e, ok := err.(errors.WrongTypeErr); ok {
		return NewRequestBuilder[[]T](f.h).
			WithError(errors.NewWrongFieldTypeErr(name, e.Type)).Create()
	}
	if err != nil {
		return NewRequestBuilder[[]T](f.h).WithError(errors.NewUnknownErr(err)).Create()
	}

	return NewRequestBuilder[[]T](f.h).WithValue(values).Create()
}
