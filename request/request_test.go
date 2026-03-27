package request

import (
	"library/errors"
	"library/test"
	"mime/multipart"
	"strconv"
	"testing"
)

func TestGetValues(t *testing.T) {
	fatalErr := errors.NewUnknownErr(nil)
	fieldName := "field1"

	type Input struct {
		name string
	}

	type Dependecies struct {
		parseError       error
		values           []string
		isValuesExpected bool
	}

	type Output struct {
		expectedValues   []string
		expectedFatalErr error
		expectedErr      error
	}

	testCases := []struct {
		Input
		Dependecies
		Output
	}{
		{
			Dependecies: Dependecies{
				parseError:       nil,
				values:           []string{"1"},
				isValuesExpected: true,
			},
			Input: Input{
				name: fieldName,
			},
			Output: Output{
				expectedValues:   []string{"1"},
				expectedFatalErr: nil,
				expectedErr:      nil,
			},
		},
		{
			Dependecies: Dependecies{
				parseError: fatalErr,
			},
			Output: Output{
				expectedValues:   nil,
				expectedFatalErr: fatalErr,
				expectedErr:      nil,
			},
		},
		{
			Dependecies: Dependecies{
				parseError:       nil,
				values:           []string{"1"},
				isValuesExpected: false,
			},
			Input: Input{
				name: fieldName,
			},
			Output: Output{
				expectedValues:   nil,
				expectedFatalErr: nil,
				expectedErr:      errors.NewFieldNotExistsErr(fieldName),
			},
		},
	}

	for _, testCase := range testCases {
		r := &MockHttpRequest{
			MockParseForm: func() error {
				return testCase.parseError
			},
			MockGetValue: func(name string) ([]string, bool) {
				if name != testCase.name {
					t.Errorf("unexpected name: got '%s', expected '%s'\n", name, testCase.name)
				}
				return testCase.values, testCase.isValuesExpected
			},
		}
		form := &form{
			r: r,
		}

		values, fatalErr, err := form.getValues(testCase.name)
		test.CompareErrors(t, "fatalErr", fatalErr, testCase.expectedFatalErr)
		test.CompareCustomErrors(t, "err", err, testCase.expectedErr)
		test.CompareSlices(t, "values", values, testCase.expectedValues)
	}
}

func TestGetFiles(t *testing.T) {
	fieldName := "field1"
	files := []*test.MockMultipartFile{
		{
			Name:    "file1",
			Content: "content1",
		},
	}
	fatalErr := errors.NewUnknownErr(nil)

	type Dependecies struct {
		parseError       error
		values           []*test.MockMultipartFile
		isValuesExpected bool
	}

	type Input struct {
		name string
	}

	type Output struct {
		expectedValues   []*test.MockMultipartFile
		expectedFatalErr error
		expectedErr      error
	}

	testCases := []struct {
		Dependecies
		Input
		Output
	}{
		{
			Dependecies: Dependecies{
				parseError:       nil,
				values:           files,
				isValuesExpected: true,
			},
			Input: Input{
				name: fieldName,
			},
			Output: Output{
				expectedValues:   files,
				expectedFatalErr: nil,
				expectedErr:      nil,
			},
		},
		{
			Dependecies: Dependecies{
				parseError: fatalErr,
			},
			Input: Input{
				name: fieldName,
			},
			Output: Output{
				expectedValues:   nil,
				expectedFatalErr: fatalErr,
				expectedErr:      nil,
			},
		},
		{
			Dependecies: Dependecies{
				parseError:       nil,
				isValuesExpected: false,
			},
			Input: Input{
				name: fieldName,
			},
			Output: Output{
				expectedValues:   nil,
				expectedFatalErr: nil,
				expectedErr:      errors.NewFieldNotExistsErr(fieldName),
			},
		},
	}

	for _, testCase := range testCases {
		r := &MockHttpRequest{
			MockParseForm: func() error {
				return testCase.parseError
			},
			MockGetFile: func(name string) ([]*multipart.FileHeader, bool) {
				test.Compare(t, "name", name, testCase.name)
				mockFiles := test.MockMultipartFiles(t, name, testCase.values)
				return mockFiles, testCase.isValuesExpected
			},
		}
		form := &form{
			r: r,
		}
		values, fatalErr, err := form.getFiles(testCase.name)
		test.CheckMultipartFiles(t, "values", values, testCase.expectedValues)
		test.CompareErrors(t, "fatalErr", fatalErr, testCase.expectedFatalErr)
		test.CompareCustomErrors(t, "err", err, testCase.expectedErr)
	}
}

func TestGet(t *testing.T) {
	fatalErr := errors.NewUnknownErr(nil)

	testCases := []struct {
		fatalErr    func() error
		err         func() error
		getValues   func(string) ([]string, error, error)
		convert     func(string) (int, error)
		addFatalErr func(error)
		addError    func(error)
		name        string
		value       int
	}{
		{
			fatalErr: func() error {
				return fatalErr
			},
		},
		{
			fatalErr: func() error {
				return nil
			},
			getValues: func(name string) ([]string, error, error) {
				test.Compare(t, "name", name, "field2")
				return nil, fatalErr, nil
			},
			addFatalErr: func(err error) {
				test.CompareErrors(t, "getValueFatalErr", err, fatalErr)
			},
			addError: func(err error) {
				test.CompareErrors(t, "err", err, nil)
			},
			name: "field2",
		},
		{
			fatalErr: func() error {
				return nil
			},
			getValues: func(name string) ([]string, error, error) {
				test.Compare(t, "name", name, "")
				return nil, nil, errors.NewUnknownErr(nil)
			},
			addFatalErr: func(err error) {
				test.CompareErrors(t, "getValueFatalErr", err, nil)
			},
			addError: func(err error) {
				test.CompareCustomErrors(t, "err", err, errors.NewUnknownErr(nil))
			},
			name: "",
		},
		{
			fatalErr: func() error {
				return nil
			},
			getValues: func(name string) ([]string, error, error) {
				test.Compare(t, "name", name, "name")
				return []string{"1", "2"}, nil, nil
			},
			addError: func(err error) {
				test.CompareCustomErrors(t, "err", err, errors.NewWrongValueSizeError("name", 2, "1"))
			},
			name: "name",
		},
		{
			fatalErr: func() error {
				return nil
			},
			getValues: func(name string) ([]string, error, error) {
				test.Compare(t, "name", name, "name")
				return []string{"1"}, nil, nil
			},
			addError: func(err error) {
				test.CompareCustomErrors(t, "err", err, errors.NewWrongFieldTypeErr("name", "string"))
			},
			convert: func(s string) (int, error) {
				test.Compare(t, "s", s, "1")
				return 0, errors.NewWrongTypeErr("string")
			},
			name: "name",
		},
		{
			fatalErr: func() error {
				return nil
			},
			getValues: func(name string) ([]string, error, error) {
				test.Compare(t, "name", name, "name")
				return []string{"1"}, nil, nil
			},
			addError: func(err error) {
				test.CompareCustomErrors(t, "err", err, errors.NewUnknownErr(errors.New("unknow error")))
			},
			convert: func(s string) (int, error) {
				test.Compare(t, "s", s, "1")
				return 0, errors.New("unknow error")
			},
			name: "name",
		},
		{
			fatalErr: func() error {
				return nil
			},
			getValues: func(name string) ([]string, error, error) {
				test.Compare(t, "name", name, "name")
				return []string{"1"}, nil, nil
			},
			convert: func(s string) (int, error) {
				test.Compare(t, "s", s, "1")
				return 1, nil
			},
			name:  "name",
			value: 1,
		},
	}

	for _, testCase := range testCases {
		form := &MockForm{
			MockErrorContainer: errors.MockErrorContainer{
				MockFatalErr:      testCase.fatalErr,
				MockAddFatalError: testCase.addFatalErr,
				MockAddError:      testCase.addError,
			},
			mockGetValues: testCase.getValues,
		}

		request := get(form, testCase.name, testCase.convert)
		test.Compare(t, "value", request.value, testCase.value)
	}
}

func Test_getFiles(t *testing.T) {
	fatalErr := errors.New("fatal error")
	files := []*test.MockMultipartFile{
		{
			Name:    "file1",
			Content: "content1",
		},
		{
			Name:    "file2",
			Content: "content2",
		},
	}

	testCases := []struct {
		parseForm func() error
		getFile   func(string) ([]*multipart.FileHeader, bool)
		name      string
		files     []*test.MockMultipartFile
		fatalErr  error
		err       error
	}{
		{
			parseForm: func() error {
				return fatalErr
			},
			fatalErr: fatalErr,
		},
		{
			parseForm: func() error {
				return nil
			},
			getFile: func(name string) ([]*multipart.FileHeader, bool) {
				test.Compare(t, "name", name, "n")
				return nil, false
			},
			name: "n",
			err:  errors.NewFieldNotExistsErr("n"),
		},
		{
			parseForm: func() error {
				return nil
			},
			getFile: func(name string) ([]*multipart.FileHeader, bool) {
				return test.MockMultipartFiles(t, "files", files), true
			},
			files: files,
		},
	}

	for _, testCase := range testCases {
		r := &MockHttpRequest{
			MockParseForm: testCase.parseForm,
			MockGetFile:   testCase.getFile,
		}
		form := &form{
			r: r,
		}
		files, fatalErr, err := form.getFiles(testCase.name)
		test.CheckMultipartFiles(t, "files", files, testCase.files)
		test.CompareErrors(t, "fatalErr", fatalErr, testCase.fatalErr)
		test.CompareCustomErrors(t, "err", err, testCase.err)
	}
}

func TestGetFile(t *testing.T) {
	fatalErr := errors.New("fatal error")
	err := errors.New("err")
	files := []*test.MockMultipartFile{
		{
			Name:    "1",
			Content: "",
		},
		{
			Name:    "1",
			Content: "new",
		},
		{
			Name:    "file",
			Content: "simple content",
		},
	}
	testCases := []struct {
		testCaseName string
		fieldName    string
		FatalErr     func() error
		getFiles     func(string) ([]*multipart.FileHeader, error, error)
		file         *test.MockMultipartFile
		fatalErr     error
		err          error
	}{
		{
			testCaseName: "check with fatalErr",
			FatalErr: func() error {
				return fatalErr
			},
			fatalErr: fatalErr,
		},
		{
			testCaseName: "check with getFiles fatalErr",
			FatalErr: func() error {
				return nil
			},
			getFiles: func(fieldName string) ([]*multipart.FileHeader, error, error) {
				test.Compare(t, "fieldName", fieldName, "")
				return nil, fatalErr, nil
			},
			fatalErr:  fatalErr,
			fieldName: "",
		},
		{
			testCaseName: "check with getFiles return simple err",
			FatalErr: func() error {
				return nil
			},
			getFiles: func(fieldName string) ([]*multipart.FileHeader, error, error) {
				test.Compare(t, "fieldName", fieldName, "field1")
				return nil, nil, err
			},
			err:       err,
			fieldName: "field1",
		},
		{
			testCaseName: "check with files len > 1",
			FatalErr: func() error {
				return nil
			},
			getFiles: func(fieldName string) ([]*multipart.FileHeader, error, error) {
				test.Compare(t, "fieldName", fieldName, "s")
				files := test.MockMultipartFiles(t, "files", files)
				t.Log("len files:", len(files))
				return files, nil, nil
			},
			err:       errors.NewWrongValueSizeError("s", len(files), "1"),
			fieldName: "s",
		},
		{
			testCaseName: "check sucessfull case",
			FatalErr: func() error {
				return nil
			},
			getFiles: func(fieldName string) ([]*multipart.FileHeader, error, error) {
				test.Compare(t, "fieldName", fieldName, "s")
				multipartFiles := test.MockMultipartFiles(t, "files", []*test.MockMultipartFile{files[0]})
				return multipartFiles, nil, nil
			},
			fieldName: "s",
			file:      files[0],
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.testCaseName)
		addFatalErr := func(fatalErr error) {
			test.CompareCustomErrors(t, "fatalErr", fatalErr, testCase.fatalErr)
		}
		addErr := func(err error) {
			test.CompareCustomErrors(t, "err", err, testCase.err)
		}
		form := &MockForm{
			MockErrorContainer: errors.MockErrorContainer{
				MockFatalErr:      testCase.FatalErr,
				MockAddFatalError: addFatalErr,
				MockAddError:      addErr,
			},
			mockGetFiles: testCase.getFiles,
		}

		request := getFile(form, testCase.fieldName)
		test.CheckMultipartFiles(t, "value", []*multipart.FileHeader{request.value},
			[]*test.MockMultipartFile{testCase.file})
	}
}

func TestGetSlice(t *testing.T) {
	fatalErr := errors.New("fatal error")
	values := []string{"1"}

	testCases := []struct {
		fatalErr      func() error
		getValues     func(string) ([]string, error, error)
		addFatalError func(error)
		addError      func(error)
		convertValues []int
		convertErr    error
		expectedValue int
		name          string
	}{
		{
			fatalErr: func() error {
				return fatalErr
			},
			addFatalError: func(err error) {
				test.CompareErrors(t, "fatalErr", err, fatalErr)
			},
		},
		{
			fatalErr: func() error {
				return nil
			},
			getValues: func(name string) ([]string, error, error) {
				test.Compare(t, "name", name, "surname")
				return nil, fatalErr, nil
			},
			addFatalError: func(err error) {
				test.CompareErrors(t, "fatalErr", err, fatalErr)
			},
			addError: func(err error) {
				test.CompareErrors(t, "err", err, nil)
			},
			name: "surname",
		},
		{
			fatalErr: func() error {
				return nil
			},
			getValues: func(name string) ([]string, error, error) {
				test.Compare(t, "name", name, "")
				return nil, nil, errors.NewFieldNotExistsErr(name)
			},
			addFatalError: func(err error) {
				test.CompareErrors(t, "fatalErr", err, nil)
			},
			addError: func(err error) {
				test.CompareCustomErrors(t, "err", err, errors.NewFieldNotExistsErr(""))
			},
			name: "",
		},
		{
			fatalErr: func() error {
				return nil
			},
			getValues: func(name string) ([]string, error, error) {
				test.Compare(t, "name", name, "b")
				return values, nil, nil
			},
			convertValues: nil,
			convertErr:    errors.NewWrongTypeErr("string"),
			addError: func(err error) {
				test.CompareCustomErrors(t, "err", err, errors.NewWrongFieldTypeErr("b", "string"))
			},
			name: "b",
		},
		{
			fatalErr: func() error {
				return nil
			},
			getValues: func(name string) ([]string, error, error) {
				test.Compare(t, "name", name, "b")
				return values, nil, nil
			},
			convertValues: nil,
			convertErr:    errors.New("convert error"),
			addError: func(err error) {
				test.CompareCustomErrors(t, "err", err, errors.NewUnknownErr(errors.New("convert error")))
			},
			name: "b",
		},
		{
			fatalErr: func() error {
				return nil
			},
			getValues: func(name string) ([]string, error, error) {
				test.Compare(t, "name", name, "b")
				return values, nil, nil
			},
			convertValues: []int{1},
			convertErr:    nil,
			name:          "b",
		},
		{
			fatalErr: func() error {
				return nil
			},
			getValues: func(name string) ([]string, error, error) {
				test.Compare(t, "name", name, "b")
				return []string{"1", "-2", "203"}, nil, nil
			},
			convertValues: []int{1, -2, 203},
			convertErr:    nil,
			name:          "b",
		},
	}

	for _, testCase := range testCases {
		form := &MockForm{
			MockErrorContainer: errors.MockErrorContainer{
				MockFatalErr:      testCase.fatalErr,
				MockAddFatalError: testCase.addFatalError,
				MockAddError:      testCase.addError,
			},
			mockGetValues: testCase.getValues,
		}
		var convert func(string) (int, error)
		if testCase.convertErr != nil {
			convert = func(s string) (int, error) {
				return 0, testCase.convertErr
			}
		} else {
			convert = strconv.Atoi
		}
		getSlice(form, testCase.name, convert)
	}
}
