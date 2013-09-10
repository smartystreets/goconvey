package assertions

import (
	"runtime"
	"testing"
)

func pass(t *testing.T, result string) {
	if result != success {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("Expectation should have passed but failed (see line %d): '%s'", line, result)
	}
}

func fail(t *testing.T, actual string, expected string) {
	if actual != expected {
		if actual == "" {
			actual = "(empty)"
		}
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("Expectation should have failed but passed (see line %d). \nExpected: %s\nActual:   %s\n",
			line, expected, actual)
	}
}

func so(actual interface{}, assert func(interface{}, ...interface{}) string, expected ...interface{}) string {
	return assert(actual, expected...)
}

type Thing1 struct {
	a string
}
type Thing2 struct {
	a string
}

type Thinger interface {
	Hi()
}

type Thing struct{}

func (self *Thing) Hi() {}
