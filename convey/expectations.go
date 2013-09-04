package convey

import (
	"fmt"
	_ "reflect"
)

type expectation func(actual interface{}, expected ...interface{}) string

func ShouldEqual(actual interface{}, expected ...interface{}) string {
	if fail := onlyOne(expected); fail != "" {
		return fail
	} else if actual != expected[0] {
		return fmt.Sprintf(shouldHaveBeenEqual, actual, expected[0])
	}
	return success
}

func ShouldNotEqual(actual interface{}, expected ...interface{}) string {
	if fail := onlyOne(expected); fail != "" {
		return fail
	} else if actual == expected[0] {
		return fmt.Sprintf(shouldNotHaveBeenEqual, actual, expected[0])
	}
	return success
}

func ShouldBeNil(actual interface{}, expected ...interface{}) string {
	if fail := none(expected); fail != "" {
		return fail
	} else if actual != nil {
		return fmt.Sprintf(shouldHaveBeenNil, actual)
	}
	return success
}
