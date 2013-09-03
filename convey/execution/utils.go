package execution

import (
	"reflect"
	"runtime"
	"strings"
)

func caller() (file string, line int) {
	// TODO: what if they have extracted the So() call into a helper method?
	//       (runtime.Caller(3) will not yield the correct stack entry!)
	_, file, line, _ = runtime.Caller(3)
	return
}
func stackTrace() string {
	buffer := make([]byte, 1024*64)
	runtime.Stack(buffer, false)
	formatted := strings.Trim(string(buffer), string([]byte{0}))
	return filterStack(formatted)
}
func fullStackTrace() string {
	buffer := make([]byte, 1024*64)
	runtime.Stack(buffer, true)
	formatted := strings.Trim(string(buffer), string([]byte{0}))
	return filterStack(formatted)
}
func filterStack(stack string) string {
	lines := strings.Split(stack, newline)
	filtered := []string{}
	for _, line := range lines {
		if !strings.Contains(line, "goconvey/convey") {
			filtered = append(filtered, line)
		}
	}
	return strings.Join(filtered, newline)
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
