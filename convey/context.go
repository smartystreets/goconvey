package convey

import (
	"runtime"
	"sync"

	"github.com/smartystreets/goconvey/execution"
	"github.com/smartystreets/goconvey/reporting"
)

// SuiteContext magically handles all coordination of reporter, runners as they handle calls
// to Convey, So, and the like. It does this via runtime call stack inspection, making sure
// that each test function has its own runner and reporter, and routes all live registrations
// to the appropriate runner/reporter.
type SuiteContext struct {
	runners   map[string]execution.Runner
	reporters map[string]reporting.Reporter
	lock      sync.Mutex
}

func (self *SuiteContext) Run(entry *execution.Registration) {
	key := resolveTestPackageAndFunctionName()
	if self.currentRunner() != nil {
		panic(execution.ExtraGoTest)
	}
	reporter := buildReporter()
	runner := execution.NewRunner()
	runner.UpgradeReporter(reporter)

	self.lock.Lock()
	self.runners[key] = runner
	self.reporters[key] = reporter
	self.lock.Unlock()

	runner.Begin(entry)
	runner.Run()

	self.lock.Lock()
	delete(self.runners, key)
	delete(self.reporters, key)
	self.lock.Unlock()
}

func (self *SuiteContext) CurrentRunner() execution.Runner {
	runner := self.currentRunner()

	if runner == nil {
		panic(execution.MissingGoTest)
	}

	return runner
}
func (self *SuiteContext) currentRunner() execution.Runner {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.runners[resolveTestPackageAndFunctionName()]
}

func (self *SuiteContext) CurrentReporter() reporting.Reporter {
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.reporters[resolveTestPackageAndFunctionName()]
}

func NewSuiteContext() *SuiteContext {
	self := new(SuiteContext)
	self.runners = make(map[string]execution.Runner)
	self.reporters = make(map[string]reporting.Reporter)
	return self
}

// resolveTestPackageAndFunctionName traverses the call stack in reverse, looking for
// the go testing harnass call ("testing.tRunner") and then grabs the very next entry,
// which represents the package under test and the test function name. Voila!
func resolveTestPackageAndFunctionName() string {
	var callerId uintptr
	callers := runtime.Callers(0, callStack)

	for y := callers; y > 0; y-- {
		callerId, _, _, _ = runtime.Caller(y)
		packageAndTestFunctionName := runtime.FuncForPC(callerId).Name()
		if packageAndTestFunctionName == goTestHarness {
			callerId, _, _, _ = runtime.Caller(y - 1)
			name := runtime.FuncForPC(callerId).Name()
			return name
		}
	}
	panic("Can't resolve test method name! Are you calling Convey() from a `*_test.go` file and a `Test*` method (because you should be)?")
}

const maxStackDepth = 100               // This had better be enough...
const goTestHarness = "testing.tRunner" // I hope this doesn't change...

var callStack []uintptr = make([]uintptr, maxStackDepth, maxStackDepth)
