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
		name          string
		fatalErr      func() error
		cookie        func(string) (*http.Cookie, error)
		addFatalErr   func(error)
		addError      func(error)
		convert       func(string) (int, error)
		fieldName     string
		expectedValue int
	}{
		{
			name: "check with fatalErr",
			fatalErr: func() error {
				return fatalErr
			},
			expectedValue: 0,
		},
		{
			name: "check with FieldNotExistsErr error from Cookie()",
			fatalErr: func() error {
				return nil
			},
			cookie: func(s string) (*http.Cookie, error) {
				return nil, http.ErrNoCookie
			},
			addFatalErr: func(err error) {
				test.CompareCustomErrors(t, "fatalErr", err, errors.NewFieldNotExistsErr("name"))
			},
			fieldName: "name",
		},
		{
			name: "check with UnknownErr error from Cookie()",
			fatalErr: func() error {
				return nil
			},
			cookie: func(s string) (*http.Cookie, error) {
				return nil, cookieErr
			},
			addFatalErr: func(err error) {
				test.CompareCustomErrors(t, "fatalErr", err, errors.NewUnknownErr(cookieErr))
			},
			fieldName: "name",
		},
		{
			name: "check with WrongTypeErr error from convert()",
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
			fieldName: "name",
		},
		{
			name: "check with simple error from convert()",
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
			name: "check normal case",
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
		t.Log(testCase.name)
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

		request := getCookie(cookie, testCase.fieldName, testCase.convert)
		test.Compare(t, "value", request.value, testCase.expectedValue)
	}
}
