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

func parseRegistration(items []interface{}) (name string, action func(), test Test) {
	if len(items) < 2 {
		panic(parseError)
	}

	name = parseName(items)
	test = parseTest(items)
	action = parseAction(items, test)

	return name, action, test
}
func parseName(items []interface{}) string {
	if name, parsed := items[0].(string); parsed {
		return name
	}
	panic(parseError)
}
func parseTest(items []interface{}) Test {
	if test, parsed := items[1].(Test); parsed {
		return test
	}
	return nil
}
func parseAction(items []interface{}, test Test) func() {
	var index = 1
	if test != nil {
		index = 2
	}

	if action, parsed := items[index].(func()); parsed {
		return action
	}
	panic(parseError)
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

const parseError = "You must provide a name (string), then a *testing.T (if in outermost scope), and then an action (func())."
