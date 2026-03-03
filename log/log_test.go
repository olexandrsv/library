package log

import (
	"library/errors"
	"library/test"
	"net/http"
	"testing"
)

func TestEndpoint(t *testing.T) {
	content := "body content"
	testCases := []struct {
		read         func() (string, error)
		writeHeader  func(int)
		write        func([]byte) (int, error)
		expectedCode int
		expectedContent string
	}{
		{
			read: func() (string, error) {
				return "", errors.New("error")
			},
			writeHeader: func(code int) {
				test.Compare(t, "code", code, http.StatusOK)
			},
			expectedCode: http.StatusOK,
		},
		{
			read: func() (string, error) {
				return content, nil
			},
			writeHeader: func(code int) {
				test.Compare(t, "code", code, http.StatusOK)
			},
			expectedCode: http.StatusOK,
			write: func(b []byte) (int, error) {
				test.Compare(t, "content", string(b), content)
				return 0, nil
			},
			expectedContent: content,
		},
	}

	for _, testCase := range testCases {
		rw := &test.MockHttpResposeWriter{
			MockHeader: func() http.Header {
				m := make(map[string][]string)
				return http.Header(m)
			},
			MockWriteHeader: testCase.writeHeader,
			MockWrite:       testCase.write,
		}
		l = &logger{
			logFile: &MockFile{
				mockRead: testCase.read,
			},
		}
		Endpoint(rw, nil)
	}
}
