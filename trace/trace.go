package trace

import (
	"bytes"
	"fmt"
	"library/log"

	"runtime"
)

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
	pcs := make([]uintptr, 2)
	n := runtime.Callers(3, pcs)
	frames := runtime.CallersFrames(pcs[:n])

	frame, _ := frames.Next()
	return frame
}
