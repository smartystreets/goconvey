package convey

import (
	"fmt"
	"runtime"
	"strings"
	"sync"

	"github.com/smartystreets/goconvey/execution"
	"github.com/smartystreets/goconvey/reporting"
)

type SuiteContext struct {
	runners   map[string]execution.Runner
	reporters map[string]reporting.Reporter
	lock      sync.Mutex
}

func (self *SuiteContext) Assign() execution.Runner {
	key := resolveExternalCallerWithTestName()
	reporter := buildReporter()
	runner := execution.NewRunner()
	runner.UpgradeReporter(reporter)

	self.lock.Lock()
	self.runners[key] = runner
	self.reporters[key] = reporter
	self.lock.Unlock()

	return runner
}

func (self *SuiteContext) CurrentRunner() execution.Runner {
	key := resolveExternalCallerWithTestName()
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.runners[key]
}

func (self *SuiteContext) CurrentReporter() reporting.Reporter {
	key := resolveExternalCallerWithTestName()
	self.lock.Lock()
	defer self.lock.Unlock()
	return self.reporters[key]
}

func NewSuiteContext() *SuiteContext {
	self := new(SuiteContext)
	self.runners = make(map[string]execution.Runner)
	self.reporters = make(map[string]reporting.Reporter)
	return self
}

func resolveExternalCallerWithTestName() string {
	// TODO: It turns out the more robust solution is to manually parse the debug.Stack()
	//       because we can then filter out non-test methods that start with "Test".

	var (
		caller_id uintptr
		testName  string
		file      string
	)
	callers := runtime.Callers(0, callStack)

	var x int
	for ; x < callers; x++ {
		caller_id, file, _, _ = runtime.Caller(x)
		if strings.HasSuffix(file, "test.go") {
			break
		}
	}

	for ; x < callers; x++ {
		caller_id, _, _, _ = runtime.Caller(x)
		packageAndTestName := runtime.FuncForPC(caller_id).Name()
		parts := strings.Split(packageAndTestName, ".")
		testName = parts[len(parts)-1]
		if strings.HasPrefix(testName, "Test") {
			break
		}
	}

	if testName == "" {
		testName = "<unkown test method name>" // panic?
	}
	return fmt.Sprintf("%s---%s", testName, file)
}

const maxStackDepth = 100 // This had better be enough...

var callStack []uintptr = make([]uintptr, maxStackDepth, maxStackDepth)
