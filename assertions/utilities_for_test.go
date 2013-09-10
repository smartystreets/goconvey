package assertions

import (
	"path"
	"runtime"
	"testing"
)

func pass(t *testing.T, result string) {
	if result != success {
		_, file, line, _ := runtime.Caller(1)
		base := path.Base(file)
		t.Errorf("Expectation should have passed but failed (see %s: line %d): '%s'", base, line, result)
	}
}

func fail(t *testing.T, actual string, expected string) {
	if actual != expected {
		if actual == "" {
			actual = "(empty)"
		}
		_, file, line, _ := runtime.Caller(1)
		base := path.Base(file)
		t.Errorf("Expectation should have failed but passed (see %s: line %d). \nExpected: %s\nActual:   %s\n",
			base, line, expected, actual)
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
