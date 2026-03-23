package trace

import (
	"bytes"
	"fmt"
	"library/slices"
	"runtime"
)

type Trace interface {
	Frames() []Frame
	fmt.Stringer
}

type MockTrace struct {
	MockFrames func() []Frame
	MockString func() string
}

func (m *MockTrace) Frames() []Frame {
	return m.MockFrames()
}

func (m *MockTrace) String() string {
	return m.MockString()
}

type trace struct {
	frames []Frame
	mockTrace
}

type mockTrace struct {
	mockFrames       func() []Frame
	mockString       func() string
	mockLoadFrames   func(framesNeeded int)
	mockGetFrames    func(framesNeeded int) slices.Iterator[runtime.Frame]
	mockGetNFrames   func(framesNeeded int) *runtime.Frames
	mockGetAllFrames func() *runtime.Frames
}

func New(framesNeeded int) Trace {
	t := &trace{}
	t.loadFrames(framesNeeded)
	return t
}

func (t *trace) Frames() []Frame {
	if t.mockFrames != nil {
		return t.mockFrames()
	}
	return t.frames
}

func (t *trace) String() string {
	if t.mockString != nil {
		return t.mockString()
	}
	var b bytes.Buffer
	for _, frame := range t.frames {
		b.WriteString(frame.String())
		b.WriteString("\n")
	}
	return b.String()
}

func (t *trace) loadFrames(framesNeeded int) {
	if t.mockLoadFrames != nil {
		t.mockLoadFrames(framesNeeded)
		return
	}
	if framesNeeded < 1 {
		return
	}
	runtimeFrames := t.getFrames(framesNeeded)
	frames := slices.IteratorToSlice(runtimeFrames)
	t.frames = slices.Map(frames, func(frame runtime.Frame) Frame {
		return NewFrame(frame.File, frame.Line, frame.Function)
	})
}

func (t *trace) getFrames(framesNeeded int) slices.Iterator[runtime.Frame] {
	if t.mockGetFrames != nil {
		return t.mockGetFrames(framesNeeded)
	}
	if framesNeeded < -1 {
		return nil
	}
	if framesNeeded == -1 {
		return t.getAllFrames()
	}
	return t.getNFrames(framesNeeded)
}

func (t *trace) getNFrames(framesNeeded int) *runtime.Frames {
	if t.mockGetNFrames != nil {
		return t.mockGetNFrames(framesNeeded)
	}
	if framesNeeded < 1 {
		return nil
	}
	pointers := make([]uintptr, framesNeeded)
	n := runtime.Callers(1, pointers)
	return runtime.CallersFrames(pointers[:n])
}

func (t *trace) getAllFrames() *runtime.Frames {
	if t.mockGetAllFrames != nil {
		return t.mockGetAllFrames()
	}
	defaultFramesAmount := 10
	for i := defaultFramesAmount; i < 1_000; i += 2 {
		pointers := make([]uintptr, i)
		n := runtime.Callers(1, pointers)
		if i != n {
			return runtime.CallersFrames(pointers[:n])
		}
	}
	return nil
}
