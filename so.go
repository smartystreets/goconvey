package goconvey

import "fmt"

func So(actual interface{}, match constraint, expected interface{}) {
	if !match(actual, expected) {
		panic(fmt.Sprintf("Doesn't match: '%v' vs '%v'", actual, expected))
	}
}

func ShouldEqual(actual interface{}, expected interface{}) bool {
	return actual == expected
}

func ShouldNotEqual(actual interface{}, expected interface{}) bool {
	return !ShouldEqual(actual, expected)
}

type constraint func(actual interface{}, expected interface{}) bool