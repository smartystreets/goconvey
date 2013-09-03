package convey

import "github.com/smartystreets/goconvey/convey/execution"

func Convey(items ...interface{}) {
	name, action, test := parseRegistration(items)

	if test != nil {
		// TODO: we need to be able to know and expect this branch of the if and panic if they forget to pass in *testing.T
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
	if result := match(actual, expected...); result == success {
		execution.SpecReporter.Report(execution.NewSuccessReport())
	} else {
		execution.SpecReporter.Report(execution.NewFailureReport(result))
	}
}
