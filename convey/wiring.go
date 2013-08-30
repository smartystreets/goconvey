package convey

import (
	"fmt"
	"github.com/smartystreets/goconvey/convey/execution"
)

func init() {
	// TODO: hook the notifier up to the spec runner...
}

func Run(t execution.GoTest, action func()) {
	execution.SpecRunner.Begin(t)
	action()
	execution.SpecRunner.Run()
}

func Convey(name string, action func()) {
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
