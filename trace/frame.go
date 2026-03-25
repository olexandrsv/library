package trace

import "fmt"

type Frame interface {
	File() string
	Line() int
	Function() string
	Format(func(Frame) string) string
	fmt.Stringer
}

type MockFrame struct {
	MockFile     func() string
	MockLine     func() int
	MockFunction func() string
	MockFormat   func(func(Frame) string) string
	MockString   func() string
}

func (f *MockFrame) File() string {
	return f.File()
}

func (f *MockFrame) Line() int {
	return f.Line()
}

func (f *MockFrame) Function() string {
	return f.Function()
}

func (f *MockFrame) Format(formatter func(Frame) string) string {
	return f.MockFormat(formatter)
}

func (f *MockFrame) String() string {
	return f.MockString()
}

type frame struct {
	file     string
	line     int
	function string
}

func NewFrame(file string, line int, function string) Frame {
	return &frame{
		file:     file,
		line:     line,
		function: function,
	}
}

func (f *frame) File() string {
	return f.file
}

func (f *frame) Line() int {
	return f.line
}

func (f *frame) Function() string {
	return f.function
}

func (f *frame) Format(formatter func(Frame) string) string {
	return formatter(f)
}

func (f *frame) String() string {
	return fmt.Sprintf("%s %d %s", f.file, f.line, f.function)
}
