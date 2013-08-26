package goconvey

import (
	"fmt"
	"errors"
)

func So(actual interface{}, match constraint, expected ...interface{}) func() {
	assertion := func() {
		err := match(actual, expected)
		fmt.Println(err)
	}
	return assertion
}

type constraint func(actual interface{}, expected []interface{}) error

func ShouldEqual(actual interface{}, expected []interface{}) error {
	if actual != expected[0] {
		message := fmt.Sprintf("'%v' should equal '%v' (but it doesn't)!", actual, expected[0])
		return errors.New(message)
	}
	return nil
}

func ShouldBeNil(actual interface{}, expected []interface{}) error {
	if actual != nil {
		return errors.New(fmt.Sprintf("'%v' should have been nil!", actual))
	}
	return nil
}