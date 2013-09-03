package gotest

type T interface {
	Fail()
	Fatalf(format string, args ...interface{})
}
