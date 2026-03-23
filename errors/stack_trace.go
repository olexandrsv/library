package errors

import (
	"fmt"
	"library/slices"
	"library/trace"
)

type Tracer interface {
	Trace() []string
}

type tracer struct {
	msgs []string
}

func (t *tracer) Trace() []string {
	return t.msgs
}

func newTrace() Tracer {
	return &tracer{
		msgs: getStack(),
	}
}

func getStack() []string {
	frames := trace.New(-1).Frames()
	records := slices.Map(frames, func(frame trace.Frame) string {
		return frameToRecord(frame)
	})

	return records
}

func frameToRecord(frame trace.Frame) string {
	return fmt.Sprintf("%s %d  %s", frame.File(), frame.Line(), frame.Function())
}
