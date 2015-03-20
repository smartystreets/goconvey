package convey

import (
	"flag"
	"math/rand"
	"os"
	"sync/atomic"
	"time"

	"github.com/jtolds/gls"
	"github.com/smartystreets/assertions"
	"github.com/smartystreets/goconvey/convey/reporting"
)

func init() {
	assertions.GoConveyMode(true)

	declareFlags()

	ctxMgr = gls.NewContextManager()
}

func declareFlags() {
	flag.BoolVar(&json, "json", false, "When true, emits results in JSON blocks. Default: 'false'")
	flag.BoolVar(&silent, "silent", false, "When true, all output from GoConvey is suppressed.")
	flag.BoolVar(&story, "story", false, "When true, emits story output, otherwise emits dot output. When not provided, this flag mirros the value of the '-test.v' flag")
	flag.BoolVar(&randomizeTests, "randomize-tests", false, "When true, randomizes the order of tests run.")
	flag.Int64Var(&randomSeed, "random-seed", 0, "The randomization seed if -randomize-tests is specified. 0 uses the time.")

	if noStoryFlagProvided() {
		story = verboseEnabled
	}

	// FYI: flag.Parse() is called from the testing package.
}

func noStoryFlagProvided() bool {
	return !story && !storyDisabled
}

func buildReporter(test t, seed int64) reporting.Reporter {
	switch {
	case testReporter != nil:
		return testReporter
	case json:
		return reporting.BuildJsonReporter(test, seed)
	case silent:
		return reporting.BuildSilentReporter(test)
	case story:
		return reporting.BuildStoryReporter(test, seed)
	default:
		return reporting.BuildDotReporter(test, seed)
	}
}

var (
	ctxMgr *gls.ContextManager

	// only set by internal tests
	testReporter reporting.Reporter
)

var (
	json           bool
	silent         bool
	story          bool
	randomizeTests bool
	randomSeed     int64

	verboseEnabled = flagFound("-test.v=true")
	storyDisabled  = flagFound("-story=false")
)

type picker interface {
	Enter(remaining []string)
	New() picker
	Pick(string) bool
}

type fakePicker struct{}

func (fakePicker) New() picker      { return fakePicker{} }
func (fakePicker) Enter([]string)   {}
func (fakePicker) Pick(string) bool { return true }

type realPicker struct {
	r      *rand.Rand
	choice string
}

func (p *realPicker) New() picker             { return &realPicker{r: p.r} }
func (p *realPicker) Pick(choice string) bool { return choice == p.choice }
func (p *realPicker) Enter(choices []string)  { p.choice = choices[p.r.Intn(len(choices))] }

func newPicker() picker {
	if randomizeTests {
		return &realPicker{r: rand.New(rand.NewSource(getSeed()))}
	} else {
		return fakePicker{}
	}
}

func getSeed() int64 {
	if randomizeTests {
		if randomSeed == 0 {
			atomic.CompareAndSwapInt64(&randomSeed, 0, time.Now().UnixNano())
		}
		return randomSeed
	}
	return 0
}

// flagFound parses the command line args manually for flags defined in other
// packages. Like the '-v' flag from the "testing" package, for instance.
func flagFound(flagValue string) bool {
	for _, arg := range os.Args {
		if arg == flagValue {
			return true
		}
	}
	return false
}
