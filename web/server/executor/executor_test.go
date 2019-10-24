package executor

import (
	"strings"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/web/server/contract"
)

func TestExecutor(t *testing.T) {
	t.Skip("BROKEN!")

	Convey("Subject: Execution of test packages and aggregation of parsed results", t, func() {
		fixture := newExecutorFixture()

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

		Convey("When the status is updated", func() {
			fixture.executor.setStatus(Executing)

			Convey("The status flag should be set to true", func() {
				So(fixture.executor.statusFlag, ShouldBeTrue)
			})
		})

		Convey("During test execution", func() {
			status := fixture.CaptureStatusDuringExecutionPhase()

			Convey("The status of the executor should be 'executing'", func() {
				So(status, ShouldEqual, Executing)
			})
		})
	})
}

type ExecutorFixture struct {
	executor *Executor
	tester   *FakeTester
	parser   *FakeParser
	folders  []*contract.Package
	result   *contract.CompleteOutput
	expected *contract.CompleteOutput
	stamp    time.Time
}

func (e *ExecutorFixture) ExecuteTests() {
	e.result = e.executor.ExecuteTests(e.folders)
}

func (e *ExecutorFixture) CaptureStatusDuringExecutionPhase() string {
	nap, _ := time.ParseDuration("25ms")
	e.tester.addDelay(nap)
	return e.delayedExecution(nap)
}

func (e *ExecutorFixture) delayedExecution(nap time.Duration) string {
	go e.ExecuteTests()
	time.Sleep(nap)
	return e.executor.Status()
}

func (e *ExecutorFixture) ResultShouldBePopulated() {
	So(e.result, ShouldResemble, e.expected)
}

var (
	prefix   = "/Users/blah/gopath/src/"
	packageA = "github.com/smartystreets/goconvey/a"
	packageB = "github.com/smartystreets/goconvey/b"
	resultA  = &contract.PackageResult{PackageName: packageA}
	resultB  = &contract.PackageResult{PackageName: packageB}
)

func newExecutorFixture() *ExecutorFixture {
	self := new(ExecutorFixture)
	self.tester = newFakeTester()
	self.parser = newFakeParser()
	self.executor = NewExecutor(self.tester, self.parser, make(chan chan string))
	self.folders = []*contract.Package{
		&contract.Package{Path: prefix + packageA, Name: packageA},
		&contract.Package{Path: prefix + packageB, Name: packageB},
	}
	self.stamp = time.Now()
	now = func() time.Time { return self.stamp }

	self.expected = &contract.CompleteOutput{
		Packages: []*contract.PackageResult{
			resultA,
			resultB,
		},
		Revision: self.stamp.String(),
	}
	return self
}

/******** FakeTester ********/

type FakeTester struct {
	nap time.Duration
}

func (f *FakeTester) SetBatchSize(batchSize int) { panic("NOT SUPPORTED") }
func (f *FakeTester) TestAll(folders []*contract.Package) {
	for _, p := range folders {
		p.Output = p.Path
	}
	time.Sleep(f.nap)
}
func (f *FakeTester) addDelay(nap time.Duration) {
	f.nap = nap
}

func newFakeTester() *FakeTester {
	self := new(FakeTester)
	zero, _ := time.ParseDuration("0")
	self.nap = zero
	return self
}

/******** FakeParser ********/

type FakeParser struct {
	nap time.Duration
}

func (f *FakeParser) Parse(packages []*contract.Package) {
	time.Sleep(f.nap)
	for _, package_ := range packages {
		if package_.Name == packageA && strings.HasSuffix(package_.Output, packageA) {
			package_.Result = resultA
		}
		if package_.Name == packageB && strings.HasSuffix(package_.Output, packageB) {
			package_.Result = resultB
		}
	}
}

func (f *FakeParser) addDelay(nap time.Duration) {
	f.nap = nap
}

func newFakeParser() *FakeParser {
	self := new(FakeParser)
	zero, _ := time.ParseDuration("0")
	self.nap = zero
	return self
}
