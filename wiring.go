package goconvey

import (
	"fmt"
	"github.com/mdwhatcott/goconvey/execution"
)

func Convey(items ...interface{}) {
	name, action, test := parseRegistration(items)

	if test != nil {
		specRunner.Begin(test)
	}

	specRunner.Register(name, action)
}

func Reset(action func()) {
	//specRunner.RegisterReset(action)
}

// TODO: hook into runner
func So(actual interface{}, match constraint, expected ...interface{}) func() {
	assertion := func() {
		err := match(actual, expected)
		fmt.Println(err)
	}
	return assertion
}

type runner interface {
	Begin(test execution.GoTest)
	Register(situation string, action func())
	Run()
}

var specRunner runner = execution.NewSpecRunner()
