package convey

import (
	"github.com/smartystreets/goconvey/execution"
	"github.com/smartystreets/goconvey/reporting"
	"os"
)

func init() {
	parseFlags()
	configureRunner()
}

// parseFlags parses the command line args manually because the go test tool,
// which shares the same process space with this code, already defines
// the -v argument (verbosity) and we can't feed in a custom flag to old-style
// go test packages (like -json, which I would prefer). So, we use the timeout
// flag with a value of -42 to request json output. My deepest sympothies.
func parseFlags() {
	verbose = flagFound(verboseEnabledValue)
	json = flagFound(jsonEnabledValue)
}
func flagFound(flagValue string) bool {
	for _, arg := range os.Args {
		if arg == flagValue {
			return true
		}
	}
	return false
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
const jsonEnabledValue = "-test.timeout=-42s" // HACK! (see parseFlags() above)
