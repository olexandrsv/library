package request

import (
	"library/errors"
	"mime/multipart"
	"net/http"
)

type HttpRequest interface {
	ParseForm() error
	GetValue(name string) ([]string, bool)
	GetFile(name string) ([]*multipart.FileHeader, bool)
}

type httpRequest struct {
	r *http.Request
}

func newHttpRequest(r *http.Request) HttpRequest {
	return &httpRequest{
		r: r,
	}
}

func (r *httpRequest) ParseForm() error {
	if r.r.MultipartForm == nil {
		if err := r.r.ParseMultipartForm(1 << 20); err != nil {
			return errors.NewDataParseErr("form")
		}
	}
	return nil
}

func (r *httpRequest) GetValue(name string) ([]string, bool) {
	values, ok := r.r.MultipartForm.Value[name]
	return values, ok
}

func (r *httpRequest) GetFile(name string) ([]*multipart.FileHeader, bool) {
	values, ok := r.r.MultipartForm.File[name]
	return values, ok
}

type MockHttpRequest struct {
	MockParseForm func() error
	MockGetValue  func(name string) ([]string, bool)
	MockGetFile   func(name string) ([]*multipart.FileHeader, bool)
}

func (c *MockHttpRequest) ParseForm() error {
	return c.MockParseForm()
}

func (c *MockHttpRequest) GetValue(name string) ([]string, bool) {
	return c.MockGetValue(name)
}

func (c *MockHttpRequest) GetFile(name string) ([]*multipart.FileHeader, bool) {
	return c.MockGetFile(name)
}
