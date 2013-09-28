package convey

import (
	"flag"
	"github.com/smartystreets/goconvey/execution"
	"github.com/smartystreets/goconvey/printing"
	"github.com/smartystreets/goconvey/reporting"
	"os"
)

func init() {
	flag.BoolVar(&json, "json", false, "internal: output json for goconvey-server")
	flag.Parse()
	reporter = buildReporter()
	runner = execution.NewRunner()
	runner.UpgradeReporter(reporter)
}

func buildReporter() reporting.Reporter {
	var consoleReporter reporting.Reporter
	console := printing.NewConsole()
	printer := printing.NewPrinter(console)

	if json {
		return reporting.NewJsonReporter(printer)
	} else if verbose() {
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

// verbose parses the command line args manually because the go test tool,
// which shares the same process space with this code, already defines
// the -v argument.
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

var json bool

const verboseEnabled = "-test.v=true"
