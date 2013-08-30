package execution

import (
	"reflect"
	"runtime"
	"strings"
)

var SpecRunner runner

func init() {
	SpecRunner = NewScopeRunner()
}

type runner interface {
	Begin(test GoTest, situation string, action func())
	Register(situation string, action func())
	RegisterReset(action func())
	Run()
}

type GoTest interface {
	Fail()
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
		if strings.HasSuffix(file, "test.go") || strings.HasSuffix(file, "tests.go") {
			return runtime.FuncForPC(caller_id).Name()
		}
	}
	return "<unknown caller!>" // panic?
}

func last(group []string) string {
	return group[len(group)-1]
}

const maxStackDepth = 10

var callStack []uintptr = make([]uintptr, maxStackDepth, maxStackDepth)
