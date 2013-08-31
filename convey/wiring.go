package convey

import (
	"github.com/smartystreets/goconvey/convey/execution"
	"runtime"
)

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

func So(actual interface{}, match expectation, expected ...interface{}) {
	// TODO: what if they have extracted the So() call into a helper method?
	//       (runtime.Caller(1) will not yield the correct stack entry!)
	failure := match(actual, expected)
	_, file, line, _ := runtime.Caller(1)
	report := execution.Report{file, line, failure, nil}
	execution.SpecReporter.Success(report)
}
