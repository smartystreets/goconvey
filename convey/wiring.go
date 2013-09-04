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
	reporter = buildReporter()
	runner = execution.NewRunner()
	runner.UpgradeReporter(reporter)
}
func buildReporter() reporting.Reporter {
	var consoleReporter reporting.Reporter
	console := printing.NewConsole()
	printer := printing.NewPrinter(console)

	if verbose() {
		consoleReporter = reporting.NewStoryReporter(printer)
	} else {
		consoleReporter = reporting.NewDotReporter(printer)
	}

	return reporting.NewReporters(
		reporting.NewGoTestReporter(),
		consoleReporter,
		reporting.NewStatisticsReporter(printer))
}
func verbose() bool {
	for _, arg := range os.Args {
		if arg == verboseEnabled {
			return true
		}
	}
	return false
}

var runner execution.Runner
var reporter reporting.Reporter

const verboseEnabled = "-test.v=true"
