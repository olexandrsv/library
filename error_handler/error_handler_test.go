package error_handler

import (
	"encoding/json"
	"library/errors"
	"library/log"
	"library/test"
	"testing"
	"time"
)

func TestWriteError(t *testing.T) {
	testCases := []struct {
		name             string
		userMessages     []string
		Write            func([]byte) (int, error)
		programmerErrors ErrorList
		expectedError    error
	}{
		{
			name:         "check concatenate and write with error",
			userMessages: []string{"msg1", "msg2", "msg3"},
			Write: func(b []byte) (int, error) {
				test.Compare(t, "b", string(b), "msg1\nmsg2\nmsg3\n")
				return 0, errors.New("write error")
			},
			expectedError: errors.New("write error"),
		},
		{
			name:         "check concatenate with messages len 1 and write with error",
			userMessages: []string{"msg1"},
			Write: func(b []byte) (int, error) {
				test.Compare(t, "b", string(b), "msg1\n")
				return 0, errors.New("write error")
			},
			expectedError: errors.New("write error"),
		},
		{
			name:         "check concatenate with messages nil and write with error",
			userMessages: nil,
			Write: func(b []byte) (int, error) {
				test.Compare(t, "b", string(b), "")
				return 0, errors.NewDataParseErr("int")
			},
			expectedError: errors.NewDataParseErr("int"),
		},
		{
			name: "check json",
			Write: func(b []byte) (int, error) {
				return 0, nil
			},
			programmerErrors: ErrorList{
				Errors: []ProgrammerError{
					{
						OriginalErr: "original error",
						Message:     "msg",
						Time:        "2006-01-01 11:32:44",
						Trace: []string{
							"record1", "recrod2",
						},
					},
				},
			},
		},
		{
			name: "check json with nil ErrorList",
			Write: func(b []byte) (int, error) {
				return 0, nil
			},
			programmerErrors: ErrorList{
				Errors: nil,
			},
		},
	}

	for _, testCase := range testCases {
		log.InitMock(log.MockLog{
			MockError: func(errToLog error) {
				bytes, err := json.MarshalIndent(testCase.programmerErrors, "", " ")
				if err != nil {
					t.Fatal(err)
				}
				test.Compare(t, "msg", errToLog.Error(), string(bytes))
			},
		})
		w := &test.MockHttpResposeWriter{
			MockWrite: testCase.Write,
		}
		h := ErrorHandler{
			userMessages:     testCase.userMessages,
			programmerErrors: testCase.programmerErrors,
			mockErrorHandler: mockErrorHandler{},
		}
		err := h.writeError(w)
		test.CompareCustomErrors(t, "err", err, testCase.expectedError)
	}
}

func TestConcatenate(t *testing.T) {
	testCases := []struct {
		slice       []string
		delimeter   string
		expectedMsg string
	}{
		{
			slice:       []string{"1", "2", "3"},
			delimeter:   "\n",
			expectedMsg: "1\n2\n3\n",
		},
		{
			slice:       []string{"1", "2", "hello"},
			delimeter:   "",
			expectedMsg: "12hello",
		},
		{
			slice:       nil,
			delimeter:   "\n",
			expectedMsg: "",
		},
		{
			slice:       []string{},
			delimeter:   "",
			expectedMsg: "",
		},
	}

	for _, testCase := range testCases {
		msg := concatenate(testCase.slice, testCase.delimeter)
		test.Compare(t, "msg", msg, testCase.expectedMsg)
	}
}

func TestParseError(t *testing.T) {
	testCases := []struct {
		err                error
		parseMultipleError func(*errors.MultipleErr) error
		parseCustomError   func(errors.CustomError) error
		expectedError      error
	}{
		{
			err: errors.NewDataParseErr("string"),
			parseCustomError: func(ce errors.CustomError) error {
				test.CompareCustomErrors(t, "ce", ce, errors.NewDataParseErr("string"))
				return nil
			},
			expectedError: nil,
		},
		{
			err: errors.NewMultipleErr(nil),
			parseMultipleError: func(me *errors.MultipleErr) error {
				test.CompareCustomErrors(t, "me", me, errors.NewMultipleErr(nil))
				return nil
			},
			expectedError: nil,
		},
		{
			err:           errors.New("not custom or multiple error"),
			expectedError: nil,
		},
		{
			err:           nil,
			expectedError: nil,
		},
	}

	for _, testCase := range testCases {
		h := ErrorHandler{
			mockErrorHandler: mockErrorHandler{
				mockParseMultipleError: testCase.parseMultipleError,
				mockParseCustomError:   testCase.parseCustomError,
			},
		}
		err := h.parseError(testCase.err)
		test.CompareCustomErrors(t, "err", err, testCase.expectedError)
	}
}

func TestParseMultipleError_parseError(t *testing.T) {
	testCases := []struct {
		err              *errors.MultipleErr
		parseErrorReturn error
	}{
		{
			err: &errors.MultipleErr{
				Errors: []error{
					errors.NewNilErr("nil", errors.NewInternalError()),
					errors.NewDataParseErr("float64"),
				},
			},
			parseErrorReturn: nil,
		},
		{
			err:              &errors.MultipleErr{},
			parseErrorReturn: nil,
		},
		{
			err: &errors.MultipleErr{
				Errors: []error{},
			},
			parseErrorReturn: nil,
		},
		{
			err: &errors.MultipleErr{
				Errors: []error{
					errors.NewNilErr("nil", errors.NewInternalError()),
					errors.NewDataParseErr("float64"),
					errors.New("not custom error"),
				},
			},
			parseErrorReturn: nil,
		},
		{
			err: &errors.MultipleErr{
				Errors: []error{
					errors.New("not custom error"),
				},
			},
			parseErrorReturn: nil,
		},
		{
			err: &errors.MultipleErr{
				Errors: []error{
					errors.New("not custom error"),
				},
			},
			parseErrorReturn: errors.New("error"),
		},
	}

	for _, testCase := range testCases {
		var i int
		h := ErrorHandler{
			mockErrorHandler: mockErrorHandler{
				mockParseError: func(err error) error {
					test.CompareCustomErrors(t, "err", err, testCase.err.Errors[i])
					i++
					return testCase.parseErrorReturn
				},
			},
		}
		err := h.parseMultipleError(testCase.err)
		test.CompareErrors(t, "err", err, testCase.parseErrorReturn)
	}
}

func TestParseCustomError(t *testing.T) {
	testCases := []struct {
		name             string
		originalErrorMsg string
		time             time.Time
		message          string
		trace            []string
		userMessage      string
		expectedError    error
	}{
		{
			name:             "check usual case",
			originalErrorMsg: "original error msg",
			time:             time.Now(),
			message:          "programmer internal error",
			trace: []string{
				"record1", "record2",
			},
			userMessage: "user internal error",
		},
		{
			name:             "check fields empty and nil",
			originalErrorMsg: "",
			time:             time.Time{},
			message:          "",
			trace:            nil,
			userMessage:      "",
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.name)
		h := newErrorHandler()
		e := errors.MockCustomError{
			MockOriginalErr: func() error {
				return errors.New(testCase.originalErrorMsg)
			},
			MockTime: func() time.Time {
				return testCase.time
			},
			MockTrace: func() []string {
				return testCase.trace
			},
			MockMessage: func() string {
				return testCase.message
			},
			MockResponse: errors.MockResponse{
				MockUserMessage: func() string {
					return testCase.userMessage
				},
			},
		}
		err := h.parseCustomError(e)
		programmerErrors := h.programmerErrors.Errors
		addedError := programmerErrors[len(programmerErrors)-1]
		addedMessage := h.userMessages[len(h.userMessages)-1]
		test.Compare(t, "originalErrorMsg", addedError.OriginalErr, testCase.originalErrorMsg)
		test.Compare(t, "time", addedError.Time, testCase.time.Format("2006-01-02 15:04:05"))
		test.CompareSlices(t, "trace", addedError.Trace, testCase.trace)
		test.Compare(t, "messages", addedError.Message, testCase.message)
		test.Compare(t, "userMessage", addedMessage, testCase.userMessage)
		test.CompareCustomErrors(t, "err", err, testCase.expectedError)
	}
}

func TestParseCustomError_NilFunctions(t *testing.T) {
	now := time.Now()
	originalErrorMsg := "original error"
	testCases := []struct {
		name                     string
		customError              error
		OriginalErr              func() error
		Time                     func() time.Time
		Trace                    func() []string
		Message                  func() string
		UserMessage              func() string
		expectedError            error
		expectedProgrammerErrors []ProgrammerError
		expectedUserMessages     []string
	}{
		{
			name:          "check with error nil",
			customError:   nil,
			expectedError: errors.NewNilErr("err", errors.NewInternalError()),
		},
		{
			name:        "check with originalError() not nil",
			customError: errors.New("not nil"),
			OriginalErr: func() error {
				return errors.New(originalErrorMsg)
			},
			Time: func() time.Time {
				return now
			},
			Trace: func() []string {
				return nil
			},
			Message: func() string {
				return ""
			},
			UserMessage: func() string {
				return "user message"
			},
			expectedUserMessages: []string{"user message"},
			expectedProgrammerErrors: []ProgrammerError{
				{
					OriginalErr: "original error",
					Message:     "",
					Time:        now.Format("2006-01-02 15:04:05"),
					Trace:       nil,
				},
			},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.name)
		var e errors.CustomError = errors.MockCustomError{
			MockOriginalErr: testCase.OriginalErr,
			MockTime:        testCase.Time,
			MockTrace:       testCase.Trace,
			MockMessage:     testCase.Message,
			MockResponse: errors.MockResponse{
				MockUserMessage: testCase.UserMessage,
			},
		}
		if testCase.customError == nil {
			e = nil
		}
		h := newErrorHandler()
		err := h.parseCustomError(e)

		test.CompareCustomErrors(t, "err", err, testCase.expectedError)
		test.CompareSlices(t, "user messages", h.userMessages, testCase.expectedUserMessages)
		test.CompareSlicesWithFunc(t, "programmerErrors", h.programmerErrors.Errors,
			testCase.expectedProgrammerErrors, func(e1, e2 ProgrammerError) bool {
				test.CompareSlices(t, "trace", e1.Trace, e2.Trace)
				return e1.OriginalErr == e2.OriginalErr && e1.Message == e2.Message &&
					e1.Time == e2.Time
			})
	}
}
