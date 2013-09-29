package execution

import (
	"reflect"
	"runtime"
)

func functionName(action func()) string {
	return runtime.FuncForPC(functionId(action)).Name()
}

func functionId(action func()) uintptr {
	return reflect.ValueOf(action).Pointer()
}

func last(group []string) string {
	return group[len(group)-1]
}
