package execution

import (
	"reflect"
	"runtime"
	"strings"
)

func caller() (file string, line int, stack string) {
	// TODO: what if they have extracted the So() call into a helper method?
	//       (runtime.Caller(3) will not yield the correct stack entry!)
	_, file, line, _ = runtime.Caller(3)
	stack = stackTrace()
	return
}
func stackTrace() string {
	// TODO: what if the stack trace is larger than the buffer? What should the max size of buffer be?
	buffer := make([]byte, 1024*10)
	runtime.Stack(buffer, false)
	return strings.Trim(string(buffer), string([]byte{0}))
}

func functionName(action func()) string {
	return runtime.FuncForPC(functionId(action)).Name()
}

func functionId(action func()) uintptr {
	return reflect.ValueOf(action).Pointer()
}

func resolveExternalCaller() string {
	callers := runtime.Callers(0, callStack)

	for x := 0; x < callers; x++ {
		caller_id, file, _, _ := runtime.Caller(x)
		if strings.HasSuffix(file, "test.go") {
			return runtime.FuncForPC(caller_id).Name()
		}
	}
	return "<unknown caller!>" // panic?
}

func last(group []string) string {
	return group[len(group)-1]
}

const maxStackDepth = 100 // This had better be enough...

var callStack []uintptr = make([]uintptr, maxStackDepth, maxStackDepth)
