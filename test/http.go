package test

import "net/http"

type MockHttpResposeWriter struct {
	MockHeader      func() http.Header
	MockWrite       func([]byte) (int, error)
	MockWriteHeader func(statusCode int)
}

func (rw *MockHttpResposeWriter) Header() http.Header {
	return rw.MockHeader()
}

func (rw *MockHttpResposeWriter) Write(bytes []byte) (int, error) {
	return rw.MockWrite(bytes)
}

func (rw *MockHttpResposeWriter) WriteHeader(statusCode int) {
	rw.MockWriteHeader(statusCode)
}
