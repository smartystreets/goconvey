package convey

import "fmt"

type constraint func(actual interface{}, expected []interface{}) string

func ShouldEqual(actual interface{}, expected []interface{}) string {
	if actual != expected[0] {
		message := fmt.Sprintf("'%v' should equal '%v' (but it doesn't)!", actual, expected[0])
		return message
	}
	return ""
}

func ShouldBeNil(actual interface{}, expected []interface{}) string {
	if actual != nil {
		return fmt.Sprintf("'%v' should have been nil!", actual)
	}
	return ""
}
