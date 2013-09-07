package convey

import (
	"github.com/smartystreets/goconvey/execution"
	"github.com/smartystreets/goconvey/printing"
	"github.com/smartystreets/goconvey/reporting"
	"os"
)

// Convey is the method intended for use when declaring the scopes
// of a specification. Each scope has a description and a func()
// which may contain other calls to Convey(), Reset() or Should-style
// assertions. Convey calls can be nested as far as you see fit.
//
// IMPORTANT NOTE: The top-level Convey() within a Test method
// must conform to the following signature:
//
//     Convey(description string, t *testing.T, action func())
//
// All other calls should like like this (no need to pass in *testing.T):
//
//     Convey(description string, action func())
//
// See the examples package for, well, examples.
func Convey(items ...interface{}) {
	name, action, test := parseRegistration(items)

	if test != nil {
		runner.Begin(test, name, action)
		runner.Run()
	} else {
		runner.Register(name, action)
	}
}

// Reset registers a cleanup function to be run after each Convey()
// in the same scope. See the examples package for a simple use case.
func Reset(action func()) {
	runner.RegisterReset(action)
}

// So is the means by which assertions are made against the system under test.
// The majority of exported names in this package begin with the word 'Should'
// and describe how the first argument (actual) should compare with any of the
// final (expected) arguments. How many final arguments are accepted depends on
// the particular assertion that is passed in as the assert argument.
// See the examples package for use cases.
func So(actual interface{}, assert assertion, expected ...interface{}) {
	if result := assert(actual, expected...); result == success {
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
		reporting.NewProblemReporter(printer),
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
