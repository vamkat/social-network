package commonerrors

import (
	"runtime"
	"strconv"
	"strings"
)

// Returns c error if c is not nil and is a defined error
// in commonerrors else returns ErrUnknown
func parseCode(c error) error {
	_, ok := classToGRPC[c]
	if c == nil || !ok {
		c = ErrUnknown
	}
	return c
}

func getInput(input ...string) string {
	if len(input) > 0 {
		return input[0]
	}
	return ""
}

func getStack(depth int, skip int) string {
	var builder strings.Builder
	builder.Grow(150)
	pc := make([]uintptr, depth)
	n := runtime.Callers(skip, pc)
	if n == 0 {
		return "(no caller data)"
	}
	pc = pc[:n] // pass only valid pcs to runtime.CallersFrames
	frames := runtime.CallersFrames(pc)
	for {
		frame, more := frames.Next()
		name := frame.Function
		start := strings.LastIndex(name, "/")
		builder.WriteString("by ")
		builder.WriteString(name[start+1:])
		builder.WriteString(" at ")
		builder.WriteString(strconv.Itoa(frame.Line))
		if !more {
			break
		}
	}

	return builder.String()
}
