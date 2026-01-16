package commonerrors

import (
	"context"
	"errors"
	"fmt"
	"reflect"
	"runtime"
	"strconv"
	"strings"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
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

// func getInput(input ...string) string {
// 	if len(input) > 0 {
// 		return input[0]
// 	}
// 	return ""
// }

type namedValue struct {
	name  string
	value any
}

func Named(name string, value any) namedValue {
	return namedValue{name: name, value: value}
}

func getInput(args ...any) string {
	var b strings.Builder

	for _, arg := range args {
		switch v := arg.(type) {
		case namedValue:
			b.WriteString(fmt.Sprintf("%s = %s\n", v.name, formatValue(v.value)))
		default:
			b.WriteString(formatValue(arg))
			b.WriteString("\n")
		}
	}

	return strings.TrimRight(b.String(), "\n")
}

func formatValue(v any) string {
	if v == nil {
		return "nil"
	}

	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	switch val.Kind() {

	case reflect.Struct:
		var b strings.Builder
		b.WriteString(fmt.Sprintf("%s {\n", typ.Name()))
		for i := 0; i < val.NumField(); i++ {
			field := typ.Field(i)
			b.WriteString(fmt.Sprintf(
				"  %s: %v\n",
				field.Name,
				val.Field(i).Interface(),
			))
		}
		b.WriteString("}")
		return b.String()

	case reflect.Map:
		var b strings.Builder
		b.WriteString("map {\n")
		for _, key := range val.MapKeys() {
			b.WriteString(fmt.Sprintf(
				"  %v: %v\n",
				key.Interface(),
				val.MapIndex(key).Interface(),
			))
		}
		b.WriteString("}")
		return b.String()

	case reflect.Slice, reflect.Array:
		var b strings.Builder
		b.WriteString("[ ")
		for i := 0; i < val.Len(); i++ {
			b.WriteString(fmt.Sprintf("%v, ", val.Index(i).Interface()))
		}
		b.WriteString("]")
		return b.String()

	default:
		return fmt.Sprintf("%v", v)
	}
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
	var count int
	for {
		count++
		frame, more := frames.Next()
		name := frame.Function
		start := strings.LastIndex(name, "/")
		builder.WriteString("level ")
		builder.WriteString(strconv.Itoa(count))
		builder.WriteString(" -> ")
		builder.WriteString(name[start+1:])
		builder.WriteString(" at l. ")
		builder.WriteString(strconv.Itoa(frame.Line))
		if !more {
			break
		}
		builder.WriteString("\n          ")
	}

	return builder.String()
}

// Helper mapper from error to grpc code.
func GetCode(err error) codes.Code {
	if err == nil {
		return codes.OK
	}

	// Propagate gRPC status errors
	if st, ok := status.FromError(err); ok {
		return st.Code()
	}

	// Handle context errors
	if errors.Is(err, context.DeadlineExceeded) {
		return codes.DeadlineExceeded
	}
	if errors.Is(err, context.Canceled) {
		return codes.Canceled
	}

	// Handle domain error
	var e *Error
	if errors.As(err, &e) {
		if code, ok := classToGRPC[e.class]; ok {
			return code
		}
	}

	// 4. Fallback
	return codes.Unknown
}
