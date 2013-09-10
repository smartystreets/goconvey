package assertions

import (
	"fmt"
	"reflect"
)

import (
	"github.com/smartystreets/oglematchers"
)

// ShouldContain receives exactly two parameters. The first is a slice and the
// second is a proposed member. Membership is determined using ShouldEqual.
func ShouldContain(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}
	if matchError := oglematchers.Contains(expected[0]).Matches(actual); matchError != nil {
		typeName := reflect.TypeOf(actual)

		if fmt.Sprintf("%v", matchError) == "which is not a slice or array" {
			return fmt.Sprintf(shouldHaveBeenAValidCollection, typeName)
		}
		return fmt.Sprintf(shouldHaveContained, typeName, expected[0])
	}
	return success
}

// ShouldNotContain receives exactly two parameters. The first is a slice and the
// second is a proposed member. Membership is determinied using ShouldEqual.
func ShouldNotContain(actual interface{}, expected ...interface{}) string {
	if fail := need(1, expected); fail != success {
		return fail
	}
	typeName := reflect.TypeOf(actual)

	if matchError := oglematchers.Contains(expected[0]).Matches(actual); matchError != nil {
		if fmt.Sprintf("%v", matchError) == "which is not a slice or array" {
			return fmt.Sprintf(shouldHaveBeenAValidCollection, typeName)
		}
		return success
	}
	return fmt.Sprintf(shouldNotHaveContained, typeName, expected[0])
}
