package goconvey

import (
	"fmt"
)

func Convey(items ...interface{}) {
	name, action, test := parseRegistration(items)

	if test != nil {
		specRunner.begin(test)
	}

	specRunner.register(name, action)
}

func Reset(action func()) {
	// TODO: hook into runner
}

// TODO: hook into runner
func So(actual interface{}, match constraint, expected ...interface{}) func() {
	assertion := func() {
		err := match(actual, expected)
		fmt.Println(err)
	}
	return assertion
}

type goTest interface {
	Fail()
}

type runner interface {
	begin(test goTest)
	register(situation string, action func())
	run()
}

var specRunner runner = newSpecRunner()
