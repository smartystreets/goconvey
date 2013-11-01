package executor

import (
	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/web/server/parser"
	"strings"
	"testing"
	"time"
)

func TestExecutor(t *testing.T) {
	var fixture *ExecutorFixture

	Convey("Subject: Execution of test packages and aggregation of parsed results", t, func() {
		fixture = newExecutorFixture()

		Convey("When tests packages are executed", func() {
			fixture.ExecuteTests()

			Convey("The result should include parsed results for each test package.",
				fixture.ResultShouldBePopulated)
		})

		Convey("When the executor is idle", func() {
			Convey("The status of the executor should be 'idle'", func() {
				So(fixture.executor.Status(), ShouldEqual, Idle)
			})
		})

		Convey("During test execution", func() {
			status := fixture.CaptureStatusDuringExecutionPhase()

			Convey("The status of the executor should be 'executing'", func() {
				So(status, ShouldEqual, Executing)
			})
		})

		Convey("During test output parsing", func() {
			status := fixture.CaptureStatusDuringParsingPhase()

			Convey("The status of the executor should be 'parsing'", func() {
				So(status, ShouldEqual, Parsing)
			})
		})
	})
}

type ExecutorFixture struct {
	executor *Executor
	tester   *FakeTester
	parser   *FakeParser
	folders  []string
	result   *parser.CompleteOutput
	expected *parser.CompleteOutput
	stamp    time.Time
}

func (self *ExecutorFixture) ExecuteTests() {
	self.result = self.executor.ExecuteTests(self.folders)
}

func (self *ExecutorFixture) CaptureStatusDuringExecutionPhase() string {
	nap, _ := time.ParseDuration("25ms")
	self.tester.addDelay(nap)
	return self.delayedExecution(nap)
}

func (self *ExecutorFixture) CaptureStatusDuringParsingPhase() string {
	nap, _ := time.ParseDuration("25ms")
	self.parser.addDelay(nap)
	return self.delayedExecution(nap)
}

func (self *ExecutorFixture) delayedExecution(nap time.Duration) string {
	go self.ExecuteTests()
	time.Sleep(nap)
	return self.executor.Status()
}

func (self *ExecutorFixture) ResultShouldBePopulated() {
	So(self.result, ShouldResemble, self.expected)
}

var (
	prefix   = "/Users/blah/gopath/src/"
	packageA = "github.com/smartystreets/goconvey/a"
	packageB = "github.com/smartystreets/goconvey/b"
	resultA  = &parser.PackageResult{PackageName: packageA}
	resultB  = &parser.PackageResult{PackageName: packageB}
)

func newExecutorFixture() *ExecutorFixture {
	self := &ExecutorFixture{}
	self.tester = newFakeTester()
	self.parser = newFakeParser()
	self.executor = NewExecutor(self.tester, self.parser)
	self.folders = []string{
		prefix + packageA,
		prefix + packageB,
	}
	self.stamp = time.Now()
	now = func() time.Time { return self.stamp }

	self.expected = &parser.CompleteOutput{
		Packages: []*parser.PackageResult{
			resultA,
			resultB,
		},
		Revision: self.stamp.String(),
	}
	return self
}

type FakeTester struct {
	nap time.Duration
}

func (self *FakeTester) SetBatchSize(batchSize int) { panic("NOT SUPPORTED") }
func (self *FakeTester) TestAll(folders []string) (output []string) {
	time.Sleep(self.nap)
	return folders
}
func (self *FakeTester) addDelay(nap time.Duration) {
	self.nap = nap
}

func newFakeTester() *FakeTester {
	self := &FakeTester{}
	zero, _ := time.ParseDuration("0")
	self.nap = zero
	return self
}

type FakeParser struct {
	nap time.Duration
}

func (self *FakeParser) Parse(packageName, output string) *parser.PackageResult {
	time.Sleep(self.nap)
	if packageName == packageA && strings.HasSuffix(output, packageA) {
		return resultA
	}
	if packageName == packageB && strings.HasSuffix(output, packageB) {
		return resultB
	}
	return nil
}

func (self *FakeParser) addDelay(nap time.Duration) {
	self.nap = nap
}

func newFakeParser() *FakeParser {
	self := &FakeParser{}
	zero, _ := time.ParseDuration("0")
	self.nap = zero
	return self
}
