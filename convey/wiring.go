package convey

import (
	"flag"
	"github.com/smartystreets/goconvey/execution"
	"github.com/smartystreets/goconvey/reporting"
	"os"
)

func init() {
	parseFlags()
	configureRunner()
}

func parseFlags() {
	flag.BoolVar(&json, "json", false, "internal: output json for goconvey-server")
	flag.Parse()
	parseVerbosity()
}

// parseVerbosity parses the command line args manually because the go test tool,
// which shares the same process space with this code, already defines the -v argument.
func parseVerbosity() {
	for _, arg := range os.Args {
		if verbose = arg == verboseEnabledValue; verbose {
			return
		}
	}
}

func configureRunner() {
	reporter = buildReporter()
	runner = execution.NewRunner()
	runner.UpgradeReporter(reporter)
}

func buildReporter() reporting.Reporter {
	if json {
		return reporting.BuildJsonReporter()
	} else if verbose {
		return reporting.BuildStoryReporter()
	} else {
		return reporting.BuildDotReporter()
	}
}

var (
	runner   execution.Runner
	reporter reporting.Reporter
)

var (
	json    bool
	verbose bool
)

const verboseEnabledValue = "-test.v=true"
