package trace

import (
	"bytes"
	"fmt"
	"library/errors"
	"library/log"

	"runtime"
)

var callers = runtime.Callers

func Save(args ...any) {
	frame := getPreviousFrame()
	var b bytes.Buffer
	for i, arg := range args {
		b.WriteString(fmt.Sprintf("%v", arg))
		if i < len(args)-1 {
			b.WriteString(", ")
		}
	}
	msg := fmt.Sprintf("%s(%s)  %s %d", frame.Function, b.String(), frame.File, frame.Line)
	fmt.Println(msg)
	log.Trace(msg)
}

func getPreviousFrame() runtime.Frame {
	frames, _ := getFrames(3)
	frame, _ := frames.Next()
	return frame
}

func getFrames(size int) (*runtime.Frames, error) {
	minSizeValue, maxSizeValue := 1, 1_000
	if size < 1 || size > 1_000 {
		return nil, errors.NewWrongValueRangeError(size, minSizeValue, maxSizeValue)
	}
	pcs := make([]uintptr, size)
	n := runtime.Callers(0, pcs)
	if size == n {
		return getFrames(size * 2)
	}
	return runtime.CallersFrames(pcs[:n]), nil
}
