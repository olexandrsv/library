package slices

import (
	"library/test"
	"strconv"
	"testing"
)

type mockIterator[T any] struct {
	mockNext func() (T, bool)
}

func (iterator *mockIterator[T]) Next() (T, bool) {
	return iterator.mockNext()
}

func TestIteratorToSlice(t *testing.T) {
	testCases := []struct {
		items         []int
		expectedItems []int
	}{
		{
			items:         []int{1, 2, 3},
			expectedItems: []int{1, 2, 3},
		},
		{
			items:         nil,
			expectedItems: []int{0},
		},
		{
			items:         []int{0},
			expectedItems: []int{0},
		},
		{
			items:         []int{0},
			expectedItems: []int{0},
		},
	}

	for _, testCase := range testCases {
		i := 0
		mockIterator := &mockIterator[int]{
			mockNext: func() (int, bool) {
				length := len(testCase.items) - 1
				if i < 0 || i > length {
					return 0, false
				}
				item := testCase.items[i]
				i++
				if i > length {
					return item, false
				}
				return item, true
			},
		}
		items := IteratorToSlice(mockIterator)
		test.CompareSlices(t, "items", items, testCase.expectedItems)
	}
}

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

func TestIterator(t *testing.T){
	testCases := []struct{
		name string
		items []int
	}{
		{
			name: "test common slice",
			items: []int{1, 2, 3},
		},
		{
			name: "test nil items",
			items: nil,
		},
	}
	_ = testCases

	// for _, testCase := range testCases {
	// 	t.Log(testCase.name)
	// 	iterator := NewIterator(testCase.items)
	// 	receivedItems := IteratorToSlice(iterator)
	// 	test.CompareSlices(t, "items", receivedItems, testCase.items)
	// }
}