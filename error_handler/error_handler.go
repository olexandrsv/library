package error_handler

import (
	"encoding/json"
	"library/errors"
	"library/log"
	"net/http"
	"strings"
	"testing"
)

type ErrorList struct {
	Errors []ProgrammerError `json:"errors"`
}

type ProgrammerError struct {
	OriginalErr string   `json:"error"`
	Message     string   `json:"message"`
	Time        string   `json:"time"`
	Trace       []string `json:"trace"`
}

type ErrorHandler struct {
	userMessages     []string
	programmerErrors ErrorList
	mockErrorHandler
}

func newErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		programmerErrors: ErrorList{},
	}
}

type mockErrorHandler struct {
	mockWriteError         func(w http.ResponseWriter) error
	mockParseError         func(err error) error
	mockParseMultipleError func(err *errors.MultipleErr) error
	mockParseCustomError   func(err errors.CustomError) error
}

func HandleError(programError error, w http.ResponseWriter) {
	handler := newErrorHandler()
	err := handler.parseError(programError)
	if err != nil {
		log.Error(err)
		return
	}
	if err := handler.writeError(w); err != nil {
		log.Error(err)
		return
	}
}

func concatenate(slice []string, delimeter string) string {
	var msg strings.Builder
	for _, item := range slice {
		msg.WriteString(item)
		msg.WriteString(delimeter)
	}
	return msg.String()
}

func (h *ErrorHandler) writeError(w http.ResponseWriter) error {
	if testing.Testing() && h.mockWriteError != nil {
		return h.mockWriteError(w)
	}
	userMsg := concatenate(h.userMessages, "\n")
	_, err := w.Write([]byte(userMsg))
	if err != nil {
		return err
	}
	bytes, err := json.MarshalIndent(h.programmerErrors, "", " ")
	if err != nil {
		return err
	}
	programmerMsg := string(bytes)
	log.Error(errors.New(programmerMsg))
	return nil
}

func (h *ErrorHandler) parseError(err error) error {
	if testing.Testing() && h.mockParseError != nil {
		return h.mockParseError(err)
	}
	switch v := err.(type) {
	case *errors.MultipleErr:
		return h.parseMultipleError(v)
	case errors.CustomError:
		return h.parseCustomError(v)
	}
	return nil
}

func (h *ErrorHandler) parseMultipleError(err *errors.MultipleErr) error {
	if testing.Testing() && h.mockParseMultipleError != nil {
		return h.mockParseMultipleError(err)
	}
	for _, e := range err.Errors {
		err := h.parseError(e)
		if err != nil {
			return err
		}
	}
	return nil
}

func (h *ErrorHandler) parseCustomError(err errors.CustomError) error {
	if testing.Testing() && h.mockParseCustomError != nil {
		return h.mockParseCustomError(err)
	}
	if err == nil {
		return errors.NewNilErr("err", errors.NewInternalError())
	}
	originalErr := "nil"
	if err.OriginalErr() != nil {
		originalErr = err.OriginalErr().Error()
	}
	h.programmerErrors.Errors = append(h.programmerErrors.Errors, ProgrammerError{
		OriginalErr: originalErr,
		Time:        err.Time().Format("2006-01-02 15:04:05"),
		Message:     err.Message(),
		Trace:       err.Trace(),
	})
	h.userMessages = append(h.userMessages, err.UserMessage())
	return nil
}
