package convey

import (
	"github.com/smartystreets/goconvey/execution"
	"github.com/smartystreets/goconvey/printing"
	"github.com/smartystreets/goconvey/reporting"
	"os"
)

func Convey(items ...interface{}) {
	name, action, test := parseRegistration(items)

	if test != nil {
		runner.Begin(test, name, action)
		runner.Run()
	} else {
		runner.Register(name, action)
	}
}

func Reset(action func()) {
	runner.RegisterReset(action)
}

func So(actual interface{}, match expectation, expected ...interface{}) {
	if result := match(actual, expected...); result == success {
		reporter.Report(reporting.NewSuccessReport())
	} else {
		reporter.Report(reporting.NewFailureReport(result))
	}
}

func init() {
	console := printing.NewConsole()
	printer := printing.NewPrinter(console)
	reporter = reporting.NewReporters(
		reporting.NewGoTestReporter(),
		consoleReporter(printer),
		reporting.NewStatisticsReporter(printer))
	runner = execution.NewRunner()
	runner.UpgradeReporter(reporter)
}
func consoleReporter(printer *printing.Printer) reporting.Reporter {
	for _, arg := range os.Args {
		if arg == "-test.v=true" {
			return reporting.NewStoryReporter(printer)
		}
	}
	return reporting.NewDotReporter(printer)
}

var runner execution.Runner
var reporter reporting.Reporter
