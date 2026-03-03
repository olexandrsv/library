package request

import (
	"library/errors"
	"mime/multipart"
	"net/http"
)

type HttpRequest interface {
	ParseForm() error
	GetValue(string) ([]string, bool)
	GetFile(string) ([]*multipart.FileHeader, bool)
	Cookie(string) (*http.Cookie, error)
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

func (r *httpRequest) Cookie(name string) (*http.Cookie, error) {
	return r.r.Cookie(name)
}

type MockHttpRequest struct {
	MockParseForm func() error
	MockGetValue  func(name string) ([]string, bool)
	MockGetFile   func(name string) ([]*multipart.FileHeader, bool)
	MockCookie    func(string) (*http.Cookie, error)
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

func (c *MockHttpRequest) Cookie(name string) (*http.Cookie, error) {
	return c.MockCookie(name)
}
