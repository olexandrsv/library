package trace

import (
	"library/slices"
	"library/test"
	"runtime"
	"testing"
)

func mock3(fn func() *runtime.Frames) *runtime.Frames {
	return mock2(fn)
}

func mock2(fn func() *runtime.Frames) *runtime.Frames {
	return mock1(fn)
}

func mock1(fn func() *runtime.Frames) *runtime.Frames {
	return fn()
}

func TestGetNFrames(t *testing.T) {
	getNFrames := &runtime.Frame{
		File:     "/home/george/IT/Projects/internal/library2/trace/trace.go",
		Line:     72,
		Function: "library/trace.(*trace).getNFrames",
	}
	fnFrame := &runtime.Frame{
		File:     "/home/george/IT/Projects/internal/library2/trace/trace_test.go",
		Function: "library/trace.TestGetNFrames.func1",
	}
	mock1Frame := &runtime.Frame{
		File:     "/home/george/IT/Projects/internal/library2/trace/trace_test.go",
		Line:     30,
		Function: "library/trace.mock1",
	}
	mock2Frame := &runtime.Frame{
		File:     "/home/george/IT/Projects/internal/library2/trace/trace_test.go",
		Line:     26,
		Function: "library/trace.mock2",
	}
	mock3Frame := &runtime.Frame{
		File:     "/home/george/IT/Projects/internal/library2/trace/trace_test.go",
		Line:     22,
		Function: "library/trace.mock3",
	}
	TestGetNFrames := &runtime.Frame{
		File:     "/home/george/IT/Projects/internal/library2/trace/trace_test.go",
		Function: "library/trace.TestGetNFrames",
	}
	runnerFrame := &runtime.Frame{
		File:     "/usr/local/go/src/testing/testing.go",
		Line:     2036,
		Function: "testing.tRunner",
	}
	exitFrame := &runtime.Frame{
		File:     "/usr/local/go/src/runtime/asm_amd64.s",
		Line:     1771,
		Function: "runtime.goexit",
	}
	_, _, _ = mock1Frame, mock2Frame, mock3Frame
	testCases := []struct {
		name         string
		framesNeeded int
		stackSize    int
		frames       []*runtime.Frame
	}{
		{
			name:         "check stackSize 0",
			framesNeeded: 1000,
			stackSize:    0,
			frames: []*runtime.Frame{
				getNFrames, fnFrame, TestGetNFrames, runnerFrame, exitFrame,
			},
		},
		{
			name:         "check stackSize 1",
			framesNeeded: 1000,
			stackSize:    1,
			frames: []*runtime.Frame{
				getNFrames, fnFrame, mock1Frame, TestGetNFrames, runnerFrame, exitFrame,
			},
		},
		{
			name:         "check stackSize 2",
			framesNeeded: 1000,
			stackSize:    2,
			frames: []*runtime.Frame{
				getNFrames, fnFrame, mock1Frame, mock2Frame, TestGetNFrames, runnerFrame, exitFrame,
			},
		},
		{
			name:         "check stackSize 3",
			framesNeeded: 1000,
			stackSize:    3,
			frames: []*runtime.Frame{
				getNFrames, fnFrame, mock1Frame, mock2Frame, mock3Frame, TestGetNFrames, runnerFrame, exitFrame,
			},
		},
		{
			name:         "check stackSize 3 and framesNeeded 1",
			framesNeeded: 1,
			stackSize:    3,
			frames: []*runtime.Frame{
				getNFrames,
			},
		},
		{
			name:         "check stackSize 3 and framesNeeded 2",
			framesNeeded: 2,
			stackSize:    3,
			frames: []*runtime.Frame{
				getNFrames, fnFrame,
			},
		},
		{
			name:         "check stackSize 3 and framesNeeded 3",
			framesNeeded: 3,
			stackSize:    3,
			frames: []*runtime.Frame{
				getNFrames, fnFrame, mock1Frame,
			},
		},
		{
			name:         "check stackSize 3 and framesNeeded 4",
			framesNeeded: 4,
			stackSize:    3,
			frames: []*runtime.Frame{
				getNFrames, fnFrame, mock1Frame, mock2Frame,
			},
		},
		{
			name:         "check stackSize 2 and framesNeeded 4",
			framesNeeded: 4,
			stackSize:    2,
			frames: []*runtime.Frame{
				getNFrames, fnFrame, mock1Frame, mock2Frame,
			},
		},
		{
			name:         "check stackSize 1 and framesNeeded 4",
			framesNeeded: 4,
			stackSize:    1,
			frames: []*runtime.Frame{
				getNFrames, fnFrame, mock1Frame, TestGetNFrames,
			},
		},
		{
			name:         "check stackSize 0 and framesNeeded 4",
			framesNeeded: 4,
			stackSize:    0,
			frames: []*runtime.Frame{
				getNFrames, fnFrame, TestGetNFrames, runnerFrame,
			},
		},
		{
			name:         "check framesNeeded 0",
			framesNeeded: 0,
			stackSize:    2,
			frames:       nil,
		},
		{
			name:         "check framesNeeded -1",
			framesNeeded: -1,
			stackSize:    3,
			frames:       nil,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.name)
		trace := &trace{}
		if testCase.stackSize < 0 || testCase.stackSize > 3 {
			t.Fatalf("stackSize expected in range [0, 3], got '%d'", testCase.stackSize)
		}
		fnFrame.Line = 189
		fn := func() *runtime.Frames {
			return trace.getNFrames(testCase.framesNeeded)
		}
		var frames *runtime.Frames
		TestGetFrameLine := 196
		switch testCase.stackSize {
		case 0:
			TestGetNFrames.Line = TestGetFrameLine
			frames = fn()
		case 1:
			TestGetNFrames.Line = TestGetFrameLine + 3
			frames = mock1(fn)
		case 2:
			TestGetNFrames.Line = TestGetFrameLine + 6
			frames = mock2(fn)
		case 3:
			TestGetNFrames.Line = TestGetFrameLine + 9
			frames = mock3(fn)
		}

		if frames == nil && testCase.frames == nil {
			continue
		}

		s := slices.IteratorToSlice(frames)
		for i, f := range s {
			exFrame := testCase.frames[i]
			if f.File != exFrame.File || f.Line != exFrame.Line || f.Function != exFrame.Function {
				t.Log("expected: ", f.File, f.Line, f.Function)
				t.Log("received: ", exFrame.File, exFrame.Line, exFrame.Function)
				t.Log()
				return
			}
		}
	}
}

func TestGetAllFrames(t *testing.T) {
	getAllFrames := &runtime.Frame{
		File:     "/home/george/IT/Projects/internal/library2/trace/trace.go",
		Line:     80,
		Function: "library/trace.(*trace).getAllFrames",
	}
	fnFrame := &runtime.Frame{
		File:     "/home/george/IT/Projects/internal/library2/trace/trace_test.go",
		Function: "library/trace.TestGetAllFrames.func1",
	}
	mock1Frame := &runtime.Frame{
		File:     "/home/george/IT/Projects/internal/library2/trace/trace_test.go",
		Line:     30,
		Function: "library/trace.mock1",
	}
	mock2Frame := &runtime.Frame{
		File:     "/home/george/IT/Projects/internal/library2/trace/trace_test.go",
		Line:     26,
		Function: "library/trace.mock2",
	}
	mock3Frame := &runtime.Frame{
		File:     "/home/george/IT/Projects/internal/library2/trace/trace_test.go",
		Line:     22,
		Function: "library/trace.mock3",
	}
	TestGetNFrames := &runtime.Frame{
		File:     "/home/george/IT/Projects/internal/library2/trace/trace_test.go",
		Function: "library/trace.TestGetAllFrames",
	}
	runnerFrame := &runtime.Frame{
		File:     "/usr/local/go/src/testing/testing.go",
		Line:     2036,
		Function: "testing.tRunner",
	}
	exitFrame := &runtime.Frame{
		File:     "/usr/local/go/src/runtime/asm_amd64.s",
		Line:     1771,
		Function: "runtime.goexit",
	}
	testCases := []struct {
		name      string
		stackSize int
		frames    []*runtime.Frame
	}{
		{
			name:      "check with stack 3",
			stackSize: 3,
			frames: []*runtime.Frame{
				getAllFrames, fnFrame, mock1Frame, mock2Frame, mock3Frame, TestGetNFrames, runnerFrame, exitFrame,
			},
		},
		{
			name:      "check with stack 2",
			stackSize: 2,
			frames: []*runtime.Frame{
				getAllFrames, fnFrame, mock1Frame, mock2Frame, TestGetNFrames, runnerFrame, exitFrame,
			},
		},
		{
			name:      "check with stack 1",
			stackSize: 1,
			frames: []*runtime.Frame{
				getAllFrames, fnFrame, mock1Frame, TestGetNFrames, runnerFrame, exitFrame,
			},
		},
		{
			name:      "check with stack 0",
			stackSize: 0,
			frames: []*runtime.Frame{
				getAllFrames, fnFrame, TestGetNFrames, runnerFrame, exitFrame,
			},
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.name)
		trace := &trace{}
		if testCase.stackSize < 0 || testCase.stackSize > 3 {
			t.Fatalf("stackSize expected in range [0, 3], got '%d'", testCase.stackSize)
		}
		fnFrame.Line = 308
		fn := func() *runtime.Frames {
			return trace.getAllFrames()
		}
		var frames *runtime.Frames
		TestGetFrameLine := 315
		switch testCase.stackSize {
		case 0:
			TestGetNFrames.Line = TestGetFrameLine
			frames = fn()
		case 1:
			TestGetNFrames.Line = TestGetFrameLine + 3
			frames = mock1(fn)
		case 2:
			TestGetNFrames.Line = TestGetFrameLine + 6
			frames = mock2(fn)
		case 3:
			TestGetNFrames.Line = TestGetFrameLine + 9
			frames = mock3(fn)
		}

		if frames == nil && testCase.frames == nil {
			continue
		}

		s := slices.IteratorToSlice(frames)
		for i, f := range s {
			exFrame := testCase.frames[i]
			if f.File != exFrame.File || f.Line != exFrame.Line || f.Function != exFrame.Function {
				t.Log("expected: ", f.File, f.Line, f.Function)
				t.Log("received: ", exFrame.File, exFrame.Line, exFrame.Function)
				t.Log()
				return
			}
		}
	}
}

func TestGetFrames(t *testing.T) {
	frames := &runtime.Frames{}
	testCases := []struct {
		name         string
		framesNeeded int
		getNFrames   func(framesNeeded int) *runtime.Frames
		getAllFrames func() *runtime.Frames
		frames       *runtime.Frames
	}{
		{
			name:         "check framesNeeded -2",
			framesNeeded: -2,
			frames:       nil,
		},
		{
			name:         "check getAllFrames",
			framesNeeded: -1,
			getAllFrames: func() *runtime.Frames {
				return frames
			},
			frames: frames,
		},
		{
			name:         "check getNFrames with framesNeeded 0",
			framesNeeded: 0,
			getNFrames: func(framesNeeded int) *runtime.Frames {
				test.Compare(t, "framesNeeded", framesNeeded, 0)
				return nil
			},
			frames: nil,
		},
		{
			name:         "check getNFrames",
			framesNeeded: 1,
			getNFrames: func(framesNeeded int) *runtime.Frames {
				test.Compare(t, "framesNeeded", framesNeeded, 1)
				return frames
			},
			frames: frames,
		},
	}

	for _, testCase := range testCases {
		t.Log(testCase.name)
		trace := &trace{
			mockTrace: mockTrace{
				mockGetNFrames:   testCase.getNFrames,
				mockGetAllFrames: testCase.getAllFrames,
			},
		}
		receivedFrames := trace.getFrames(testCase.framesNeeded)
		if receivedFrames == nil && testCase.frames == nil {
			return
		}
		if receivedFrames.(*runtime.Frames) != testCase.frames {
			t.Errorf("expected frames with value '%v', got '%v'", testCase.frames, receivedFrames)
		}
	}
}

func TestLoadFrames(t *testing.T) {
	frame1 := runtime.Frame{
		File:     "file1",
		Line:     1,
		Function: "func1",
	}
	frame2 := runtime.Frame{
		File:     "file2",
		Line:     2,
		Function: "func2",
	}
	testCases := []struct {
		name         string
		framesNeeded int
		getFrames    func(framesNeeded int) slices.Iterator[runtime.Frame]
		frames       []runtime.Frame
	}{
		{
			framesNeeded: 1,
			frames: []runtime.Frame{
				frame1,
			},
			getFrames: func(framesNeeded int) slices.Iterator[runtime.Frame] {
				test.Compare(t, "framesNeeded", framesNeeded, 1)
				return slices.NewIterator([]runtime.Frame{frame1})
			},
		},
		{
			framesNeeded: 0,
			frames:       nil,
		},
		{
			framesNeeded: 1,
			frames:       []runtime.Frame{},
			getFrames: func(framesNeeded int) slices.Iterator[runtime.Frame] {
				test.Compare(t, "framesNeeded", framesNeeded, 1)
				return nil
			},
		},
		{
			framesNeeded: 2,
			frames: []runtime.Frame{
				frame1, frame2,
			},
			getFrames: func(framesNeeded int) slices.Iterator[runtime.Frame] {
				test.Compare(t, "framesNeeded", framesNeeded, 2)
				return slices.NewIterator([]runtime.Frame{frame1, frame2})
			},
		},
	}

	for _, testCase := range testCases {
		trace := &trace{
			mockTrace: mockTrace{
				mockGetFrames: testCase.getFrames,
			},
		}
		trace.loadFrames(testCase.framesNeeded)
		test.CompareSlicesWithFunc(t, "frames", trace.frames, testCase.frames, func(received Frame, expected runtime.Frame) bool {
			return received.File() == expected.File && received.Line() == expected.Line && received.Function() == expected.Function
		})
	}
}

func TestString(t *testing.T) {
	testCases := []struct {
		name           string
		text           []string
		expectedString string
	}{
		{
			text:           []string{"1", "2", "3"},
			expectedString: "1\n2\n3\n",
		},
		{
			text: nil,
			expectedString: "",
		},
		{
			text: []string{},
			expectedString: "",
		},
		{
			text: []string{"1"},
			expectedString: "1\n",
		},
	}

	for _, testCase := range testCases {
		var frames []Frame
		for _, s := range testCase.text {
			frames = append(frames, &MockFrame{
				MockString: func() string {
					return s
				},
			})
		}
		trace := &trace{
			frames: frames,
		}
		s := trace.String()
		test.Compare(t, "s", s, testCase.expectedString)
	}
}
