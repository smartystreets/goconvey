package gotest

import (
	"fmt"
	"runtime"
	"strings"
)

func FormatExternalFileAndLine() string {
	file, line, _ := ResolveExternalCaller()
	if line == -1 {
		return "<unknown caller!>" // panic?
	}
	return fmt.Sprintf("%s:%d", file, line)
}

func ResolveExternalCaller() (file string, line int, name string) {
	var caller_id uintptr
	callers := runtime.Callers(0, callStack)

	for x := 0; x < callers; x++ {
		caller_id, file, line, _ = runtime.Caller(x)
		if strings.HasSuffix(file, "test.go") {
			name = runtime.FuncForPC(caller_id).Name()
			return
		}
	}
	file, line, name = "<unkown file>", -1, "<unknown name>"
	return // panic?
}

// Much like ResolveExternalCaller, but goes a bit deeper to get the test method name.
func ResolveExternalCallerWithTestName() (file string, line int, testName string) {
	// TODO: It turns out the more robust solution is to manually parse the debug.Stack()

	var caller_id uintptr
	callers := runtime.Callers(0, callStack)

	var x int
	for ; x < callers; x++ {
		caller_id, file, line, _ = runtime.Caller(x)
		if strings.HasSuffix(file, "test.go") {
			break
		}
	}

	for ; x < callers; x++ {
		caller_id, _, _, _ = runtime.Caller(x)
		packageAndTestName := runtime.FuncForPC(caller_id).Name()
		parts := strings.Split(packageAndTestName, ".")
		testName = parts[len(parts)-1]
		if strings.HasPrefix(testName, "Test") {
			return
		}
	}

	testName = "<unkown test method name>"
	return // panic?
}

const maxStackDepth = 100 // This had better be enough...

var callStack []uintptr = make([]uintptr, maxStackDepth, maxStackDepth)
