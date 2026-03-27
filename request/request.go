package request

import (
	"library/errors"
	"library/slices"
	"mime/multipart"
	"net/http"
)

type Request struct {
	Form   *form
	Cookie *Cookie
	r      *http.Request
	c      errors.ErrorContainer
}

func New(r *http.Request) *Request {
	h := errors.NewErrorContainer()
	return &Request{
		Form: &form{
			r:              newHttpRequest(r),
			ErrorContainer: h,
		},
		Cookie: NewCookie(r, h),
		r:      r,
		c:      h,
	}
}

func (r *Request) Err() error {
	if r.c.FatalErr() != nil {
		return r.c.FatalErr()
	}
	return r.c.Error()
}

type Form interface {
	errors.ErrorContainer
	getValues(string) ([]string, error, error)
	getFiles(name string) ([]*multipart.FileHeader, error, error)
}

type form struct {
	r HttpRequest
	errors.ErrorContainer
}

func (f *form) GetIntSlice(name string) *ValueRequest[[]int] {
	return getSlice(f, name, ToInt)
}

func (f *form) GetStringSlice(name string) *ValueRequest[[]string] {
	return getSlice(f, name, ToString)
}

func (f *form) GetFloat64(name string) *ValueRequest[float64] {
	return get(f, name, ToFloat64)
}

func (f *form) GetString(name string) *ValueRequest[string] {
	return get(f, name, ToString)
}

func (f *form) GetInt(name string) *ValueRequest[int] {
	return get(f, name, ToInt)
}

func (f *form) GetFile(name string) *ValueRequest[*multipart.FileHeader]{
	return getFile(f, name)
}

func get[T any](f Form, name string, convert func(string) (T, error)) *ValueRequest[T] {
	if f.FatalErr() != nil {
		return NewRequestBuilder[T](f).Create()
	}

	values, fatalErr, err := f.getValues(name)
	if fatalErr != nil || err != nil {
		f.AddFatalError(fatalErr)
		f.AddError(err)
		return NewRequestBuilder[T](f).Create()
	}

	if len(values) != 1 {
		f.AddError(errors.NewWrongValueSizeError(name, len(values), "1"))
		return NewRequestBuilder[T](f).Create()
	}

	rawValue := values[0]
	value, err := convert(rawValue)
	if e, ok := err.(*errors.WrongTypeErr); ok {
		f.AddError(errors.NewWrongFieldTypeErr(name, e.Type))
		return NewRequestBuilder[T](f).Create()
	}
	if err != nil {
		f.AddError(errors.NewUnknownErr(err))
		return NewRequestBuilder[T](f).Create()
	}

	return NewRequestBuilder[T](f).WithValue(value).Create()
}

func getFile(f Form, fieldName string) *ValueRequest[*multipart.FileHeader] {
	if f.FatalErr() != nil {
		return NewRequestBuilder[*multipart.FileHeader](f).Create()
	}

	files, fatalErr, err := f.getFiles(fieldName)
	if fatalErr != nil || err != nil {
		f.AddFatalError(fatalErr)
		f.AddError(err)
		return NewRequestBuilder[*multipart.FileHeader](f).Create()
	}

	if len(files) != 1 {
		f.AddError(errors.NewWrongValueSizeError(fieldName, len(files), "1"))
		return NewRequestBuilder[*multipart.FileHeader](f).Create()
	}

	return NewRequestBuilder[*multipart.FileHeader](f).WithValue(files[0]).Create()
}

func (f *form) getValues(name string) ([]string, error, error) {
	if err := f.r.ParseForm(); err != nil {
		return nil, err, nil
	}

	values, ok := f.r.GetValue(name)
	if !ok {
		return nil, nil, errors.NewFieldNotExistsErr(name)
	}

	return values, nil, nil
}

func (f *form) getFiles(name string) ([]*multipart.FileHeader, error, error) {
	if err := f.r.ParseForm(); err != nil {
		return nil, err, nil
	}

	files, ok := f.r.GetFile(name)
	if !ok {
		return nil, nil, errors.NewFieldNotExistsErr(name)
	}

	return files, nil, nil
}

func getSlice[T any](f Form, name string, convert func(string) (T, error)) *ValueRequest[[]T] {
	if f.FatalErr() != nil {
		return NewRequestBuilder[[]T](f).Create()
	}

	rawValues, fatalErr, err := f.getValues(name)
	if fatalErr != nil || err != nil {
		f.AddFatalError(fatalErr)
		f.AddError(err)
		return NewRequestBuilder[[]T](f).Create()
	}

	values, err := slices.Convert(rawValues, convert)
	if e, ok := err.(*errors.WrongTypeErr); ok {
		f.AddError(errors.NewWrongFieldTypeErr(name, e.Type))
		return NewRequestBuilder[[]T](f).Create()
	}
	if err != nil {
		f.AddError(errors.NewUnknownErr(err))
		return NewRequestBuilder[[]T](f).Create()
	}

	return NewRequestBuilder[[]T](f).WithValue(values).Create()
}

type MockForm struct {
	errors.MockErrorContainer
	mockGetValues func(string) ([]string, error, error)
	mockGetFiles func(string) ([]*multipart.FileHeader, error, error)
}

func (f *MockForm) getValues(name string) ([]string, error, error) {
	return f.mockGetValues(name)
}

func (f *MockForm) getFiles(name string) ([]*multipart.FileHeader, error, error){
	return f.mockGetFiles(name)
}