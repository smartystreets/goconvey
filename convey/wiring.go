package convey

import (
	"fmt"
	"github.com/mdwhatcott/goconvey/convey/execution"
)

func Convey(items ...interface{}) {
	name, action, test := parseRegistration(items)

	if test != nil {
		execution.SpecRunner.Begin(test)
	}

	execution.SpecRunner.Register(name, action)
}

func Reset(action func()) {
	execution.SpecRunner.RegisterReset(action)
}

// TODO: hook into runner
func So(actual interface{}, match constraint, expected ...interface{}) func() {
	assertion := func() {
		err := match(actual, expected)
		fmt.Println(err)
	}
	return assertion
}
