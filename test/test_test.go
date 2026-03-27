package test

import (
	"errors"
	"testing"
)

type mockTB struct{
	errored bool
	fataled bool
	testing.TB
}

func (m *mockTB) Errorf(format string, args ...interface{}) {
	m.errored = true
}

func (m *mockTB) Fatal(args ...interface{}) {
	m.fataled = true
	panic("fatal called")
}

func TestCompare(t *testing.T) {
	testCases := []struct {
		name     string
		received any
		expected any
		wantErr  bool
	}{
		{"equal ints", 1, 1, false},
		{"different ints", 1, 2, true},
		{"equal strings", "foo", "foo", false},
		{"different strings", "foo", "bar", true},
		{"equal bools", true, true, false},
		{"different bools", true, false, true},
		{"zero values", 0, 0, false},
		{"empty strings", "", "", false},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockTB{}
			switch v := tc.received.(type) {
			case int:
				Compare(m, "int", v, tc.expected.(int))
			case string:
				Compare(m, "str", v, tc.expected.(string))
			case bool:
				Compare(m, "bool", v, tc.expected.(bool))
			}
			if m.errored != tc.wantErr {
				t.Errorf("wantErr=%v got=%v", tc.wantErr, m.errored)
			}
		})
	}
}

func TestComparePointers(t *testing.T) {
	a, b, c := 1, 1, 2
	testCases := []struct {
		name     string
		received *int
		expected *int
		wantErr  bool
	}{
		{"both nil", nil, nil, false},
		{"equal values", &a, &b, false},
		{"different values", &a, &c, true},
		{"one nil", &a, nil, true},
		{"other nil", nil, &a, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockTB{}
			ComparePointers(m, "ptr", tc.received, tc.expected)
			if m.errored != tc.wantErr {
				t.Errorf("wantErr=%v got=%v", tc.wantErr, m.errored)
			}
		})
	}
}

func TestCompareSlices(t *testing.T) {
	testCases := []struct {
		name     string
		received []int
		expected []int
		wantErr  bool
	}{
		{"both nil", nil, nil, false},
		{"both empty", []int{}, []int{}, false},
		{"equal slices", []int{1, 2}, []int{1, 2}, false},
		{"different order", []int{1, 2}, []int{2, 1}, true},
		{"different lengths", []int{1, 2}, []int{1, 2, 3}, true},
		{"one nil", []int{1, 2}, nil, true},
		{"other nil", nil, []int{1, 2}, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockTB{}
			CompareSlices(m, "slice", tc.received, tc.expected)
			if m.errored != tc.wantErr {
				t.Errorf("wantErr=%v got=%v", tc.wantErr, m.errored)
			}
		})
	}
}

func TestCompareSlicesWithFunc(t *testing.T) {
	       testCases := []struct {
		       name      string
		       received  []int
		       expected  []int
		       isEqual   func(int, int) bool
		       wantErr   bool
		       wantFatal bool
	       }{
		       {"custom equal", []int{1, 2}, []int{2, 3}, func(a, b int) bool { return a+1 == b }, false, false},
		       {"custom unequal", []int{1, 2}, []int{2, 4}, func(a, b int) bool { return a+1 == b }, true, false},
		       {"empty slices", []int{}, []int{}, func(a, b int) bool { return a == b }, false, false},
		       {"nil slices", nil, nil, func(a, b int) bool { return a == b }, false, false},
		       {"different lengths", []int{1}, []int{1, 2}, func(a, b int) bool { return a == b }, true, false},
		       {"isEqual nil", []int{1}, []int{1}, nil, false, true},
	       }
	       for _, tc := range testCases {
		       t.Run(tc.name, func(t *testing.T) {
			       m := &mockTB{}
			       func() {
				       defer func() {
					       if r := recover(); r != nil {
						       
					       }
				       }()
				       CompareSlicesWithFunc(m, "slice", tc.received, tc.expected, tc.isEqual)
			       }()
			       if m.fataled != tc.wantFatal {
				       t.Errorf("wantFatal=%v got=%v", tc.wantFatal, m.fataled)
			       }
			       if !m.fataled && m.errored != tc.wantErr {
				       t.Errorf("wantErr=%v got=%v", tc.wantErr, m.errored)
			       }
		       })
	       }
}

func TestCompareErrors(t *testing.T) {
	err1 := errors.New("fail1")
	err2 := errors.New("fail2")
	testCases := []struct {
		name     string
		received error
		expected error
		wantErr  bool
	}{
		{"both nil", nil, nil, false},
		{"same error", err1, err1, false},
		{"different errors", err1, err2, true},
		{"one nil", err1, nil, true},
		{"other nil", nil, err2, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockTB{}
			CompareErrors(m, "err", tc.received, tc.expected)
			if m.errored != tc.wantErr {
				t.Errorf("wantErr=%v got=%v", tc.wantErr, m.errored)
			}
		})
	}
}

func TestCompareCustomErrors(t *testing.T) {
	err1 := errors.New("fail1")
	err2 := errors.New("fail2")
	err1copy := errors.New("fail1")
	testCases := []struct {
		name     string
		received error
		expected error
		wantErr  bool
	}{
		{"both nil", nil, nil, false},
		{"same error instance", err1, err1, false},
		{"same error message", err1, err1copy, false},
		{"different errors", err1, err2, true},
		{"one nil", err1, nil, true},
		{"other nil", nil, err2, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockTB{}
			CompareCustomErrors(m, "err", tc.received, tc.expected)
			if m.errored != tc.wantErr {
				t.Errorf("wantErr=%v got=%v", tc.wantErr, m.errored)
			}
		})
	}
}

type testStruct struct {
	A int
	B string
}

type testIface interface {
	Value() int
}

type testImpl struct {
	v int
}

func (t testImpl) Value() int { return t.v }

func TestCompare_Structs(t *testing.T) {
	testCases := []struct {
		name     string
		received testStruct
		expected testStruct
		wantErr  bool
	}{
		{"equal structs", testStruct{1, "a"}, testStruct{1, "a"}, false},
		{"different int field", testStruct{1, "a"}, testStruct{2, "a"}, true},
		{"different string field", testStruct{1, "a"}, testStruct{1, "b"}, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockTB{}
			Compare(m, "struct", tc.received, tc.expected)
			if m.errored != tc.wantErr {
				t.Errorf("wantErr=%v got=%v", tc.wantErr, m.errored)
			}
		})
	}
}

func TestCompare_Interfaces(t *testing.T) {
	testCases := []struct {
		name     string
		received testIface
		expected testIface
		wantErr  bool
	}{
		{"equal interface values", testImpl{1}, testImpl{1}, false},
		{"different interface values", testImpl{1}, testImpl{2}, true},
		{"both nil", nil, nil, false},
		{"one nil", testImpl{1}, nil, true},
		{"other nil", nil, testImpl{1}, true},
	}
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			m := &mockTB{}
			Compare(m, "iface", tc.received, tc.expected)
			if m.errored != tc.wantErr {
				t.Errorf("wantErr=%v got=%v", tc.wantErr, m.errored)
			}
		})
	}
}
