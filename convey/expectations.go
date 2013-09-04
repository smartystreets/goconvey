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

/*

1.0		// Equality
*	X	So(thing, ShouldEqual, thing2)
*	X	So(thing, ShouldNotEqual, thing2)
*		So(thing, ShouldBeLike, thing2)
*		So(thing, ShouldNotBeLike, thing2)
		So(thing, ShouldPointTo, thing2)
		So(thing, ShouldNotPointTo, thing2)
*	X	So(thing, ShouldBeNil, thing2)
*		So(thing, ShouldNotBeNil, thing2)
*		So(thing, ShouldBeTrue)
*		SO(thing, ShouldBeFalse)

		// Interfaces
		So(1, ShouldImplement, Interface)
		So(1, ShouldNotImplement, OtherInterface)

		// Type checking
		So(1, ShouldBeAn, int)
		So(1, ShouldNotBeAn, int)
		So("1", ShouldBeA, string)
		So("1", ShouldNotBeA, string)

		// Quantity comparison (Time and Duration might be special cases?)
*		So(1, ShouldBeGreaterThan, 0)
*		So(1, ShouldBeGreaterThanOrEqualTo, 0)
*		So(1, ShouldBeLessThan, 2)
*		So(1, ShouldBeLessThanOrEqualTo, 2)

		// Tolerences
*		So(1.1, ShouldBeWithin, .1, 1)
*		So(1.1, ShouldNotBeWithin, .1, 2)

		// Collections
*		So([]int{}, ShouldBeEmpty)
*		So([]int{1}, ShouldNotBeEmpty)
*		So([]int{1, 2, 3}, ShouldContain, 1, 2) // This could receive several final arguments as proposed members
*		So([]int{1, 2, 3}, ShouldNotContain, 4, 5) // This could receive several final arguments as proposed members
		So(1, ShouldBeIn, []int{1, 2, 3})
		So(4, ShouldNotBeIn, []int{1, 2, 3})

		// Strings
*		So("asdf", ShouldStartWith, "as")
*		So("asdf", ShouldNotStartWith, "df")
*		So("asdf", ShouldEndWith, "df")
*		So("asdf", ShouldNotEndWith, "df")
		So("(asdf)", ShouldBeSurroundedWith, "()")
		So("(asdf)", ShouldNotBeSurroundedWith, "[]")

*/
