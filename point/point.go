package point

import (
	"bytes"
	"fmt"
	"library/errors"
	"library/log"
	"library/trace"
)

func Save(args ...any) {
	frames := trace.New(1).Frames()
	if len(frames) > 1 {
		log.Error(errors.Newf("got frames with len '%d' expcted len '1'", len(frames)))
		return
	}
	msg := format(frames[0], args...)
	log.Trace(msg)
}

func format(frame trace.Frame, args ...any) string {
	var b bytes.Buffer
	for i, arg := range args {
		fmt.Fprintf(&b, "%v", arg)
		if i < len(args)-1 {
			b.WriteString(", ")
		}
	}
	return fmt.Sprintf("%s %d  %s(%s)", frame.File(), frame.Line(), frame.Function(), b.String())
}
