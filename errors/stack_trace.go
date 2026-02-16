package errors

import (
	"fmt"
	"runtime"
)

type Tracer interface {
	Trace() []string
}

type trace struct {
	msgs []string
}

func (t *trace) Trace() []string {
	return t.msgs
}

func newTrace() Tracer {
	return &trace{
		msgs: getStack(),
	}
}

func getStack() []string {
	pcs := make([]uintptr, 10)
	n := runtime.Callers(2, pcs)
	frames := runtime.CallersFrames(pcs[:n])
	var records []string

	for {
		frame, more := frames.Next()
		record := fmt.Sprintf("%s  %s %d", frame.Function, frame.File, frame.Line)
		records = append(records, record)
		if !more {
			break
		}
	}

	return records
}

// type stackTracer interface {
// 	StackTrace() errors.StackTrace
// }

// type tracesSkiper struct {
// 	err  error
// 	skip int
// }

// func newTracesSkiper(err error, skip int) tracesSkiper {
// 	return tracesSkiper{
// 		err:  err,
// 		skip: skip,
// 	}
// }

// func (ts tracesSkiper) Error() string {
// 	tracer, ok := ts.err.(stackTracer)
// 	if !ok {
// 		return fmt.Sprintf("%+v\n", ts.err)
// 	}
// 	trace := tracer.StackTrace()
// 	if len(trace) < ts.skip {
// 		return fmt.Sprintf("%+v\n", ts.err)
// 	}
// 	trace = trace[ts.skip:]
// 	return fmt.Sprintf("%+v\n", trace)
// }
