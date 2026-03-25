package slices

import (
	"library/errors"
	"library/test"
	"strconv"
	"testing"
)

type TestCase[T, K any] struct {
	items         []T
	expectedItems []K
	mapFunc       func(T) K
}

func TestMap_IntToString(t *testing.T) {
	testCases := []TestCase[int, string]{
		{
			items:         []int{1, 2, 3},
			expectedItems: []string{"1", "2", "3"},
			mapFunc: func(i int) string {
				return strconv.Itoa(i)
			},
		},
	}

	runMapTest(t, testCases)
}

func TestMap_NilAndEmpty(t *testing.T) {
	testCases := []TestCase[int, int]{
		{
			items:         nil,
			expectedItems: []int{},
		},
		{
			items:         []int{},
			expectedItems: []int{},
		},
	}

	runMapTest(t, testCases)
}

func TestMap_IntToInt(t *testing.T) {
	testCases := []TestCase[int, int]{
		{
			items:         []int{1, 2, 3},
			expectedItems: []int{1, 2, 3},
			mapFunc: func(i int) int {
				return i
			},
		},
	}

	runMapTest(t, testCases)
}

func runMapTest[T, K comparable](t *testing.T, testCases []TestCase[T, K]) {
	for _, testCase := range testCases {
		i := 0
		fun := func(n T) K {
			if i > len(testCase.items)-1 {
				t.Fatal("called more than needed")
			}
			test.Compare(t, "n", n, testCase.items[i])
			i++
			return testCase.mapFunc(n)
		}
		items := Map(testCase.items, fun)
		test.CompareSlices(t, "items", items, testCase.expectedItems)
	}
}

type ConvertTestCase[T, K comparable] struct {
	name          string
	items         []T
	expectedItems []K
	convert       func(T) (K, error)
	err           error
}

func TestConvert_StringToInt(t *testing.T) {
	testCases := []ConvertTestCase[string, int]{
		{
			name:          "check string to int",
			items:         []string{"1", "2", "3"},
			expectedItems: []int{1, 2, 3},
			convert:       strconv.Atoi,
			err:           nil,
		},
		{
			name:    "check string to int",
			items:   []string{"1s"},
			convert: strconv.Atoi,
			err:     errors.New(`strconv.Atoi: parsing "1s": invalid syntax`),
		},
	}

	runConvertTest(t, testCases)
}

type User struct {
	Name    string
	Surname string
}

type Author struct {
	Name    string
	Surname string
}

func TestConvert_StructToStruct(t *testing.T) {
	err := errors.New("user convert error")
	testCases := []ConvertTestCase[User, Author]{
		{
			name: "check struct to struct",
			items: []User{
				{
					Name:    "Bob",
					Surname: "Smith",
				},
				{
					Name:    "Ben",
					Surname: "Arnum",
				},
			},
			expectedItems: []Author{
				{
					Name:    "Bob",
					Surname: "Smith",
				},
				{
					Name:    "Ben",
					Surname: "Arnum",
				},
			},
			convert: func(u User) (Author, error) {
				return Author{
					Name:    u.Name,
					Surname: u.Surname,
				}, nil
			},
			err: nil,
		},
		{
			name: "check struct to struct",
			items: []User{
				{
					Name:    "Bob",
					Surname: "Smith",
				},
				{
					Name:    "Ben",
					Surname: "Arnum",
				},
			},
			expectedItems: nil,
			convert: func(u User) (Author, error) {
				return Author{}, err
			},
			err: err,
		},
	}

	runConvertTest(t, testCases)
}

func TestConvert_EdgeCases(t *testing.T) {
	testCases := []ConvertTestCase[string, int]{
		{
			name:          "check nil",
			items:         nil,
			expectedItems: []int{},
			convert:       strconv.Atoi,
			err:           nil,
		},
		{
			name:          "check empty",
			items:         []string{},
			expectedItems: []int{},
			convert:       strconv.Atoi,
			err:           nil,
		},
	}

	runConvertTest(t, testCases)
}

func runConvertTest[T, K comparable](t *testing.T, testCases []ConvertTestCase[T, K]) {
	for _, testCase := range testCases {
		t.Log(testCase.name)
		receivedItems, err := Convert(testCase.items, testCase.convert)
		t.Log(err)
		test.CompareSlices(t, "receivedItems", receivedItems, testCase.expectedItems)
		test.CompareCustomErrors(t, "err", err, testCase.err)
	}
}

func TestFilter(t *testing.T) {
	testCases := []struct {
		name          string
		items         []int
		expectedItems []int
		expectedError error
		fun           func(int) bool
	}{
		{
			name:          "check ususal case",
			items:         []int{1, 2, 3},
			expectedItems: []int{1, 3},
			fun: func(i int) bool {
				return i != 2
			},
		},
		{
			name:          "check when items is nil",
			items:         nil,
			expectedItems: nil,
			fun:           func(i int) bool {
				return true
			},
		},
		{
			name:          "check when fun is nil",
			items:         []int{1, 2, 3},
			expectedItems: nil,
			fun:           nil,
			expectedError: errors.NewNilErr("f", errors.NewInternalError()),
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.name)
		receivedItesm, err := Filter(testCase.items, testCase.fun)
		test.CompareSlices(t, "receivedItems", receivedItesm, testCase.expectedItems)
		test.CompareCustomErrors(t, "err", err, testCase.expectedError)
	}
}
