package request

import (
	"library/errors"
	"testing"
)

func TestDo(t *testing.T) {
	testCases := []struct {
		value         int
		defaultValue  int
		fatalError    error
		err           error
		expectedValue int
	}{
		{
			value:      1,
			fatalError: nil,
			err:        nil,
			expectedValue: 1,
		},
		{
			value: 1,
			fatalError: errors.NewUnknownErr(nil),
			defaultValue: 10,
			expectedValue: 0,
		},
		{
			value: 1,
			err: errors.NewUnknownErr(nil),
			defaultValue: 10,
			expectedValue: 10,
		},
		{
			value: 1,
			err: errors.NewUnknownErr(nil),
			expectedValue: 0,
		},
	}

	for _, testCase := range testCases {
		v := ValueRequest[int]{
			value:        testCase.value,
			defaultValue: &testCase.defaultValue,
			errorContainer: &errors.MockErrorContainer{
				MockFatalErr: func() error {
					return testCase.fatalError
				},
				MockError: func() error {
					return testCase.err
				},
			},
		}

		value := v.Do()
		if value != testCase.expectedValue {
			t.Errorf("expected value = '%d' but got '%d'", testCase.expectedValue, value)
		}
	}
}
