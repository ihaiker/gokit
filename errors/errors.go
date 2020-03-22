package errors

import (
	"errors"
	"fmt"
	"runtime"
	"strings"
)

var ErrAssert = errors.New("assert")

func New(format string, obj ...interface{}) error {
	return fmt.Errorf(format, obj...)
}

func Convert(rev interface{}) error {
	if rev == nil {
		return nil
	}
	switch tt := rev.(type) {
	case error:
		return tt
	default:
		return fmt.Errorf("%v", rev)
	}
}

func Stack() string {
	stackBuf := make([]uintptr, 50)
	length := runtime.Callers(3, stackBuf[:])
	stack := stackBuf[:length]
	trace := ""
	frames := runtime.CallersFrames(stack)
	for {
		frame, more := frames.Next()
		if strings.Contains(frame.File, "/gokit/errors/") ||
			strings.HasSuffix(frame.File, "/src/runtime/panic.go") ||
			strings.HasSuffix(frame.File, "/testing/testing.go") ||
			frame.Function == "runtime.goexit" ||
			frame.Function == "" {
		} else {
			trace = trace + fmt.Sprintf("  %s:%d , Function: %s,\n", frame.File, frame.Line, frame.Function)
		}
		if !more {
			break
		}
	}
	return trace
}
