// Package gotest contains internal functionality. Although this package
// contains one or more exported names it is not intended for public
// consumption. See the examples package for how to use this project.
package gotest

import (
	"runtime"
	"strings"
)

func ResolveExternalCaller() (string, int) {
	callers := runtime.Callers(0, callStack)

	for x := 0; x < callers; x++ {
		_, file, line, _ := runtime.Caller(x)
		if strings.HasSuffix(file, "_test.go") {
			return file, line
		}
	}

	// panic?
	return "<unkown file>", -1
}

const maxStackDepth = 100 // This had better be enough...

var callStack []uintptr = make([]uintptr, maxStackDepth, maxStackDepth)
