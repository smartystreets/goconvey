package executor

import (
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/smartystreets/goconvey/web/server/contract"
)

func init() {
	log.SetOutput(ioutil.Discard)
}

func TestConcurrentTester(t *testing.T) {
	t.Skip("BROKEN!")

	Convey("Subject: Controlled execution of test packages", t, func() {
		fixture := NewTesterFixture()

		Convey("Whenever tests for each package are executed", func() {
			fixture.InBatchesOf(1).RunTests()

			Convey("The tester should execute the tests in each active package with the correct arguments",
				fixture.ShouldHaveRecordOfExecutionCommands)

			Convey("There should be a test output result for each active package",
				fixture.ShouldHaveOneOutputPerInput)

			Convey("The output should be as expected",
				fixture.OutputShouldBeAsExpected)
		})

		Convey("When the tests for each package are executed synchronously", func() {
			fixture.InBatchesOf(1).RunTests()

			Convey("Each active package should be run synchronously and in the given order",
				fixture.TestsShouldHaveRunContiguously)
		})

		Convey("When the tests for each package are executed synchronously with failures", func() {
			fixture.InBatchesOf(1).SetupFailedTestSuites().RunTests()

			Convey("The failed test packages should not result in any panics", func() {
				So(fixture.recovered, ShouldBeNil)
			})
		})

		Convey("When packages are tested concurrently", func() {
			fixture.InBatchesOf(concurrentBatchSize).RunTests()

			Convey("Active packages should be arranged and tested in batches of the appropriate size",
				fixture.TestsShouldHaveRunInBatchesOfTwo)
		})

		Convey("When packages are tested concurrently with failures", func() {
			fixture.InBatchesOf(concurrentBatchSize).SetupFailedTestSuites().RunTests()

			Convey("The failed test packages should not result in any panics", func() {
				So(fixture.recovered, ShouldBeNil)
			})
		})
	})
}

const concurrentBatchSize = 2

type TesterFixture struct {
	tester       *ConcurrentTester
	shell        *TimedShell
	results      []string
	compilations []*ShellCommand
	executions   []*ShellCommand
	packages     []*contract.Package
	recovered    error
}

func NewTesterFixture() *TesterFixture {
	self := new(TesterFixture)
	self.shell = NewTimedShell()
	self.tester = NewConcurrentTester(self.shell)
	self.packages = []*contract.Package{
		{Path: "a"},
		{Path: "b"},
		{Path: "c"},
		{Path: "d"},
		{Path: "e", Ignored: true},
		{Path: "f"},
		{Path: "g", HasImportCycle: true},
	}
	return self
}

func (t *TesterFixture) InBatchesOf(batchSize int) *TesterFixture {
	t.tester.SetBatchSize(batchSize)
	return t
}

func (t *TesterFixture) SetupAbnormalError(message string) *TesterFixture {
	t.shell.setTripWire(message)
	return t
}

func (t *TesterFixture) SetupFailedTestSuites() *TesterFixture {
	t.shell.setExitWithError()
	return t
}

func (t *TesterFixture) RunTests() {
	defer func() {
		if r := recover(); r != nil {
			t.recovered = r.(error)
		}
	}()

	t.tester.TestAll(t.packages)
	for _, p := range t.packages {
		t.results = append(t.results, p.Output)
	}
	t.executions = t.shell.Executions()
}

func (t *TesterFixture) ShouldHaveRecordOfExecutionCommands() {
	executed := []string{"a", "b", "c", "d", "f"}
	ignored := "e"
	importCycle := "g"
	actual := []string{}
	for _, pkg := range t.executions {
		actual = append(actual, pkg.Command)
	}
	So(actual, ShouldResemble, executed)
	So(actual, ShouldNotContain, ignored)
	So(actual, ShouldNotContain, importCycle)
}

func (t *TesterFixture) ShouldHaveOneOutputPerInput() {
	So(len(t.results), ShouldEqual, len(t.packages))
}

func (t *TesterFixture) OutputShouldBeAsExpected() {
	for _, p := range t.packages {
		if p.HasImportCycle {
			So(p.Output, ShouldContainSubstring, "can't load package: import cycle not allowed")
			So(p.Error.Error(), ShouldContainSubstring, "can't load package: import cycle not allowed")
		} else {
			if p.Active() {
				So(p.Output, ShouldEndWith, p.Path)
			} else {
				So(p.Output, ShouldBeBlank)
			}
			So(p.Error, ShouldBeNil)
		}
	}
}

func (t *TesterFixture) TestsShouldHaveRunContiguously() {
	t.OutputShouldBeAsExpected()

	So(t.shell.MaxConcurrentCommands(), ShouldEqual, 1)

	for i := 0; i < len(t.executions)-1; i++ {
		current := t.executions[i]
		next := t.executions[i+1]
		So(current.Started, ShouldHappenBefore, next.Started)
		So(current.Ended, ShouldHappenOnOrBefore, next.Started)
	}
}

func (t *TesterFixture) TestsShouldHaveRunInBatchesOfTwo() {
	t.OutputShouldBeAsExpected()

	So(t.shell.MaxConcurrentCommands(), ShouldEqual, concurrentBatchSize)
}

/**** Fakes ****/

type ShellCommand struct {
	Command string
	Started time.Time
	Ended   time.Time
}

type TimedShell struct {
	executions   []*ShellCommand
	panicMessage string
	err          error
}

func (t *TimedShell) Executions() []*ShellCommand {
	return t.executions
}

func (t *TimedShell) MaxConcurrentCommands() int {
	var concurrent int

	for x, current := range t.executions {
		concurrentWith_x := 1
		for y, comparison := range t.executions {
			if y == x {
				continue
			} else if concurrentWith(current, comparison) {
				concurrentWith_x++
			}
		}
		if concurrentWith_x > concurrent {
			concurrent = concurrentWith_x
		}
	}
	return concurrent
}

func concurrentWith(current, comparison *ShellCommand) bool {
	return ((comparison.Started == current.Started || comparison.Started.After(current.Started)) &&
		(comparison.Started.Before(current.Ended)))
}

func (t *TimedShell) setTripWire(message string) {
	t.panicMessage = message
}

func (t *TimedShell) setExitWithError() {
	t.err = errors.New("Simulate test failure")
}

func (t *TimedShell) GoTest(directory, packageName string, arguments, tags []string) (output string, err error) {
	if t.panicMessage != "" {
		return "", errors.New(t.panicMessage)
	}

	output = directory
	err = t.err
	t.executions = append(t.executions, t.composeCommand(directory))
	return
}

func (t *TimedShell) composeCommand(commandText string) *ShellCommand {
	start := time.Now()
	time.Sleep(nap)
	end := time.Now()
	return &ShellCommand{commandText, start, end}
}

func NewTimedShell() *TimedShell {
	self := new(TimedShell)
	self.executions = []*ShellCommand{}
	return self
}

var nap, _ = time.ParseDuration("10ms")
var _ = fmt.Sprintf("fmt")
