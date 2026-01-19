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
			b.WriteString(fmt.Sprintf("%s = %s\n", v.name, FormatValue(v.value)))
		default:
			b.WriteString(FormatValue(arg))
			b.WriteString("\n")
		}
	}

	return strings.TrimRight(b.String(), "\n")
}

func FormatValue(v any) string {
	return formatValueIndented(v, 0, make(map[uintptr]bool))
}

// TODO: Needs more testing with nested values
func formatValueIndented(v any, depth int, seen map[uintptr]bool) (out string) {
	defer func() {
		if r := recover(); r != nil {
			out = "<unprintable>"
		}
	}()

	if v == nil {
		return "nil"
	}

	val := reflect.ValueOf(v)
	typ := reflect.TypeOf(v)

	// Unwrap interfaces
	for val.Kind() == reflect.Interface {
		if val.IsNil() {
			return "nil"
		}
		val = val.Elem()
		typ = val.Type()
	}

	// Handle pointers (with cycle detection)
	if val.Kind() == reflect.Pointer {
		if val.IsNil() {
			return "nil"
		}
		ptr := val.Pointer()
		if seen[ptr] {
			return "<cycle>"
		}
		seen[ptr] = true
		return formatValueIndented(val.Elem().Interface(), depth, seen)
	}

	indent := strings.Repeat("  ", depth)
	nextIndent := strings.Repeat("  ", depth+1)
	stringerType := reflect.TypeOf((*fmt.Stringer)(nil)).Elem()

	// Value implements Stringer
	if typ.Implements(stringerType) {
		return val.Interface().(fmt.Stringer).String()
	}

	// Pointer implements Stringer
	if val.CanAddr() {
		ptrVal := val.Addr()
		if ptrVal.Type().Implements(stringerType) {
			return ptrVal.Interface().(fmt.Stringer).String()
		}
	}

	switch val.Kind() {

	case reflect.Struct:
		var b strings.Builder
		name := typ.Name()
		if name == "" {
			name = "struct"
		}

		b.WriteString(indent + name + " {\n")

		for i := 0; i < val.NumField(); i++ {
			fieldType := typ.Field(i)
			fieldVal := val.Field(i)

			b.WriteString(nextIndent + fieldType.Name + ": ")

			if fieldVal.CanInterface() {
				b.WriteString(formatValueIndented(
					fieldVal.Interface(),
					depth+1,
					seen,
				))
			} else {
				b.WriteString("<unexported>")
			}
			b.WriteString("\n")
		}

		b.WriteString(indent + "}")
		return b.String()

	case reflect.Map:
		var b strings.Builder
		b.WriteString("map {\n")

		for _, key := range val.MapKeys() {
			b.WriteString(nextIndent)
			b.WriteString(fmt.Sprintf(
				"%v: %s\n",
				key.Interface(),
				formatValueIndented(val.MapIndex(key).Interface(), depth+1, seen),
			))
		}

		b.WriteString(indent + "}")
		return b.String()

	case reflect.Slice, reflect.Array:
		var b strings.Builder
		b.WriteString("[\n")

		for i := 0; i < val.Len(); i++ {
			b.WriteString(nextIndent)
			b.WriteString(formatValueIndented(
				val.Index(i).Interface(),
				depth+1,
				seen,
			))
			b.WriteString("\n")
		}

		b.WriteString(indent + "]")
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
