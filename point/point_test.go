package point

import (
	"library/log"
	"library/test"
	"library/trace"
	"strings"
	"testing"
)

func TestSave(t *testing.T) {
	testCases := []struct {
		args        []any
		expectedMsg string
	}{
		{
			args:        []any{1, 2, "hello", nil},
			expectedMsg: "1, 2, hello, <nil>",
		},
		{
			args:        nil,
			expectedMsg: "",
		},
	}
	for _, testCase := range testCases {
		log.InitMock(log.MockLog{
			MockTrace: func(msg string) {
				expected := testCase.expectedMsg + ")"
				if !strings.HasSuffix(msg, expected) {
					t.Errorf("got: '%s', expected string with ending: '%s'", msg, expected)
				}
			},
		})
		Save(testCase.args...)
	}
}

func TestFormat(t *testing.T) {
	testCases := []struct {
		frame       trace.Frame
		args        []any
		expectedMsg string
	}{
		{
			frame:       trace.NewFrame("file.go", 23, "test"),
			args:        []any{1, 2, "hello", nil},
			expectedMsg: "file.go 23  test(1, 2, hello, <nil>)",
		},
		{
			frame:       trace.NewFrame("", 0, ""),
			args:        []any{1, 2, "hello", nil},
			expectedMsg: " 0  (1, 2, hello, <nil>)",
		},
		{
			frame:       trace.NewFrame("", 0, ""),
			args:        nil,
			expectedMsg: " 0  ()",
		},
	}
	for _, testCase := range testCases {
		msg := format(testCase.frame, testCase.args...)
		test.Compare(t, "msg", msg, testCase.expectedMsg)
	}
}
