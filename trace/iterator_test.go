package trace

import (
	"library/test"
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

func TestIteratorToSlice_NilIterator(t *testing.T) {
	testCases := []struct {
		iterator      TraceIterator[int]
		expectedItems []int
	}{
		{
			iterator:      nil,
			expectedItems: nil,
		},
	}

	for _, testCase := range testCases {
		receivedItems := IteratorToSlice(testCase.iterator)
		test.CompareSlices(t, "receivedItems", receivedItems, testCase.expectedItems)
	}
}

func TestIterator(t *testing.T) {
	testCases := []struct {
		name          string
		items         []int
		expectedItems []int
	}{
		{
			name:          "test common slice",
			items:         []int{1, 2, 3},
			expectedItems: []int{1, 2, 3},
		},
		{
			name:          "test nil items",
			items:         nil,
			expectedItems: []int{0},
		},
		{
			name:          "test empty items",
			items:         []int{},
			expectedItems: []int{0},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.name)
		iterator := NewIterator(testCase.items)
		receivedItems := IteratorToSlice(iterator)
		test.CompareSlices(t, "items", receivedItems, testCase.expectedItems)
	}
}
