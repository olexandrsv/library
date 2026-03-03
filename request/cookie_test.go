package request

import (
	"library/errors"
	"library/test"
	"net/http"
	"testing"
)

func TestGetCookie(t *testing.T) {
	fatalErr := errors.New("fatal error")
	cookieErr := errors.New("cookie error")
	convertErr := errors.New("convert error")

	testCases := []struct {
		fatalErr      func() error
		cookie        func(string) (*http.Cookie, error)
		addFatalErr   func(error)
		addError      func(error)
		convert       func(string) (int, error)
		name          string
		expectedValue int
	}{
		{
			fatalErr: func() error {
				return fatalErr
			},
			expectedValue: 0,
		},
		{
			fatalErr: func() error {
				return nil
			},
			cookie: func(s string) (*http.Cookie, error) {
				return nil, http.ErrNoCookie
			},
			addFatalErr: func(err error) {
				test.CompareCustomErrors(t, "fatalErr", err, errors.NewFieldNotExistsErr("name"))
			},
			name: "name",
		},
		{
			fatalErr: func() error {
				return nil
			},
			cookie: func(s string) (*http.Cookie, error) {
				return nil, cookieErr
			},
			addFatalErr: func(err error) {
				test.CompareCustomErrors(t, "fatalErr", err, errors.NewUnknownErr(cookieErr))
			},
			name: "name",
		},
		{
			fatalErr: func() error {
				return nil
			},
			cookie: func(s string) (*http.Cookie, error) {
				return &http.Cookie{
					Value: "1",
				}, nil
			},
			convert: func(s string) (int, error) {
				test.Compare(t, "s", s, "1")
				return 0, errors.NewWrongTypeErr("string")
			},
			addError: func(err error) {
				test.CompareCustomErrors(t, "err", err, errors.NewWrongFieldTypeErr("name", "string"))
			},
			name: "name",
		},
		{
			fatalErr: func() error {
				return nil
			},
			cookie: func(s string) (*http.Cookie, error) {
				return &http.Cookie{
					Value: "1",
				}, nil
			},
			convert: func(s string) (int, error) {
				test.Compare(t, "s", s, "1")
				return 0, convertErr
			},
			addError: func(err error) {
				test.CompareCustomErrors(t, "err", err, errors.NewUnknownErr(convertErr))
			},
		},
		{
			fatalErr: func() error {
				return nil
			},
			cookie: func(s string) (*http.Cookie, error) {
				return &http.Cookie{
					Value: "1",
				}, nil
			},
			convert: func(s string) (int, error) {
				test.Compare(t, "s", s, "1")
				return 1, nil
			},
			expectedValue: 1,
		},
	}

	for _, testCase := range testCases {
		cookie := &Cookie{
			r: &MockHttpRequest{
				MockCookie: testCase.cookie,
			},
			h: &errors.MockErrorContainer{
				MockFatalErr:      testCase.fatalErr,
				MockAddFatalError: testCase.addFatalErr,
				MockAddError:      testCase.addError,
			},
		}

		request := getCookie(cookie, testCase.name, testCase.convert)
		test.Compare(t, "value", request.value, testCase.expectedValue)
	}
}
