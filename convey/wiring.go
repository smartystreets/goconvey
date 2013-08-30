package convey

import (
	"fmt"
	"github.com/smartystreets/goconvey/convey/execution"
)

func init() {
	// TODO: hook the notifier up to the spec runner...
}

func Convey(items ...interface{}) {
	name, action, test := parseRegistration(items)

	if test != nil {
		execution.SpecRunner.Begin(test, name, action)
		execution.SpecRunner.Run()
	} else {
		execution.SpecRunner.Register(name, action)
	}
}

func Reset(action func()) {
	execution.SpecRunner.RegisterReset(action)
}

// TODO: hook into runner (or reporter?)
func So(actual interface{}, match constraint, expected ...interface{}) func() {
	assertion := func() {
		err := match(actual, expected)
		fmt.Println(err)
	}
	return assertion
}
