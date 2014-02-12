package convey

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"strings"
	"sync"

	"github.com/smartystreets/goconvey/execution"
	"github.com/smartystreets/goconvey/reporting"
)

// suiteContext magically handles all coordination of reporter, runners as they handle calls
// to Convey, So, and the like. It does this via runtime call stack inspection, making sure
// that each test function has its own runner and reporter, and routes all live registrations
// to the appropriate runner/reporter.
type suiteContext struct {
	locations map[string]string             // key: file:line; value: testName
	runners   map[string]execution.Runner   // key: testName;
	reporters map[string]reporting.Reporter // key: testName;
	lock      sync.Mutex
}

func (self *suiteContext) Run(entry *execution.Registration) {
	testName, location, _ := resolveAnchorConvey()

	if self.currentRunner() != nil {
		panic(execution.ExtraGoTest)
	}

	reporter := buildReporter()
	runner := execution.NewRunner()
	runner.UpgradeReporter(reporter)

	self.lock.Lock()
	self.locations[location] = testName
	self.runners[testName] = runner
	self.reporters[testName] = reporter
	self.lock.Unlock()

	runner.Begin(entry)
	runner.Run()

	self.lock.Lock()
	delete(self.locations, location)
	delete(self.runners, testName)
	delete(self.reporters, testName)
	self.lock.Unlock()
}

func (self *suiteContext) CurrentRunner() execution.Runner {
	runner := self.currentRunner()

	if runner == nil {
		panic(execution.MissingGoTest)
	}

	return runner
}
func (self *suiteContext) currentRunner() execution.Runner {
	self.lock.Lock()
	defer self.lock.Unlock()
	testName, _, _ := resolveAnchorConvey()
	return self.runners[testName]
}

func (self *suiteContext) CurrentReporter() reporting.Reporter {
	self.lock.Lock()
	defer self.lock.Unlock()
	testName, _, err := resolveAnchorConvey()

	if err != nil {
		file, line := resolveTestFileAndLine()
		closest := -1
		for location, registeredTestName := range self.locations {
			parts := strings.Split(location, ":")
			locationFile := parts[0]
			if locationFile != file {
				continue
			}

			locationLine, err := strconv.Atoi(parts[1])
			if err != nil || locationLine < line {
				continue
			}

			if closest == -1 || locationLine < closest {
				closest = locationLine
				testName = registeredTestName
			}
		}
	}

	return self.reporters[testName]
}

func newSuiteContext() *suiteContext {
	self := new(suiteContext)
	self.locations = make(map[string]string)
	self.runners = make(map[string]execution.Runner)
	self.reporters = make(map[string]reporting.Reporter)
	return self
}

// resolveAnchorConvey traverses the call stack in reverse, looking for
// the go testing harnass call ("testing.tRunner") and then grabs the very next entry,
// which represents the test function name, including the package name as a prefix.
// It also returns the file:line combo of the top-level Convey. Voila!
func resolveAnchorConvey() (testName, location string, err error) {
	callers := runtime.Callers(0, callStack)

	for y := callers; y > 0; y-- {
		callerId, file, conveyLine, found := runtime.Caller(y)
		if !found {
			continue
		}

		if name := runtime.FuncForPC(callerId).Name(); name != goTestHarness {
			continue
		}

		callerId, file, conveyLine, _ = runtime.Caller(y - 1)
		testName = runtime.FuncForPC(callerId).Name()
		location = fmt.Sprintf("%s:%d", file, conveyLine)
		return
	}
	return "", "", errors.New("Can't resolve test method name! Are you calling Convey() from a `*_test.go` file and a `Test*` method (because you should be)?")
}

// resolveTestFileAndLine is used as a last-ditch effort to correlate an
// assertion with the right executor and runner.
func resolveTestFileAndLine() (file string, line int) {
	callers := runtime.Callers(0, callStack)
	var found bool

	for y := callers; y > 0; y-- {
		_, file, line, found = runtime.Caller(y)
		if !found {
			continue
		}

		if strings.HasSuffix(file, "_test.go") {
			return
		}
	}
	return "", 0
}

const maxStackDepth = 100               // This had better be enough...
const goTestHarness = "testing.tRunner" // I hope this doesn't change...

var callStack []uintptr = make([]uintptr, maxStackDepth, maxStackDepth)
