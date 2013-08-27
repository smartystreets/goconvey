package goconvey

import (
	// "reflect"
	"fmt"
)

func Convey(items ...interface{}) {
	name, action, test := parseRegistration(items)

	// if test != nil {
	specRunner.begin(test)
	// }

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

func parseRegistration(items []interface{}) (name string, action func(), test Test) {
	if len(items) < 2 {
		panic("You must provide a name (string), then a *testing.T (if in outermost scope), and then an action (func()).")
	}
	name, _ = items[0].(string)
	test, _ = items[1].(Test)
	if test == nil {
		action, _ = items[1].(func())
	} else {
		action, _ = items[2].(func())
	}
	return
}

type Test interface {
	Fail()
}

type runner interface {
	begin(test Test)
	register(situation string, action func())
	run()
}

var specRunner runner = newSpecRunner()
