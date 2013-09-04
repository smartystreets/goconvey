package convey

import (
	"github.com/smartystreets/goconvey/execution"
	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/printing"
	"github.com/smartystreets/goconvey/reporting"
)

func Convey(items ...interface{}) {
	name, action, test := parseRegistration(items)

	if test != nil {
		SpecRunner.Begin(test, name, action)
		SpecRunner.Run()
	} else {
		SpecRunner.Register(name, action)
	}
}

func Reset(action func()) {
	SpecRunner.RegisterReset(action)
}

func So(actual interface{}, match expectation, expected ...interface{}) {
	if result := match(actual, expected...); result == success {
		SpecReporter.Report(reporting.NewSuccessReport())
	} else {
		SpecReporter.Report(reporting.NewFailureReport(result))
	}
}

// TODO: private...
var SpecRunner runner
var SpecReporter reporting.Reporter

func init() {
	console := printing.NewConsole()
	printer := printing.NewPrinter(console)
	SpecReporter = reporting.NewReporters(
		reporting.NewGoTestReporter(),
		reporting.NewStoryReporter(printer), // TODO: or a dot reporter (-v)
		reporting.NewStatisticsReporter(printer))
	SpecRunner = execution.NewScopeRunner()
	SpecRunner.UpgradeReporter(SpecReporter)
}

type runner interface {
	Begin(test gotest.T, situation string, action func())
	Register(situation string, action func())
	RegisterReset(action func())
	Run()
	UpgradeReporter(out reporting.Reporter)
}
