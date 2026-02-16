package error_handler

import (
	"encoding/json"
	"fmt"
	"library/errors"
	"library/log"
	"net/http"
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

func HandleError(err error, w http.ResponseWriter) {
	l := &ErrorList{}
	userMsg, _ := handleError(err, "", l)
	w.Write([]byte(userMsg))
	programmerMsg, _ := json.MarshalIndent(*l, "", " ")
	fmt.Println(string(programmerMsg))
	log.Error(string(programmerMsg))
}

func handleError(err error, userError string, programmerError *ErrorList) (string, *ErrorList) {
	switch v := err.(type) {
	case errors.MultipleErr:
		for _, e := range v.Errors {
			ue, _ := handleError(e, userError, programmerError)
			userError += ue + "\n"
		}
	case errors.CustomError:
		originalErr := "nil"
		if v.OriginalErr() != nil {
			originalErr = v.OriginalErr().Error()
		}
		programmerError.Errors = append(programmerError.Errors, ProgrammerError{
			OriginalErr: originalErr,
			Time:        v.Time().Format("2006-01-02 15:04:05"),
			Message:     v.Message(),
			Trace:       v.Trace(),
		})
		return v.UserMessage(), programmerError
	}
	return userError, programmerError
}
