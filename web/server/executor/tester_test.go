package executor

import (
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
	"time"
)

func TestConcurrentTester(t *testing.T) {
	var fixture *TesterFixture

	Convey("Subject: Controlled execution of test packages", t, func() {
		fixture = NewTesterFixture()

		Convey("Whenever tests for each package are executed", func() {
			fixture.InBatchesOf(1).RunTests()

			Convey("The tester should build all dependencies of input packages",
				fixture.ShouldHaveRecordOfCompilationCommands)

			Convey("The tester should execute the tests in each package with the correct arguments",
				fixture.ShouldHaveRecordOfExecutionCommands)

			Convey("There should be a test output result for each package",
				fixture.ShouldHaveOneOutputPerInput)

			Convey("The output should be as expected",
				fixture.OutputShouldBeAsExpected)
		})

		Convey("When the tests for each package are executed synchronously", func() {
			fixture.InBatchesOf(1).RunTests()

			Convey("Each package should be run synchronously and in the given order",
				fixture.TestsShouldHaveRunContiguously)
		})

		Convey("When packages are tested concurrently", func() {
			fixture.InBatchesOf(concurrentBatchSize).RunTests()

			Convey("Packages should be arranged and tested in batches of the appropriate size",
				fixture.TestsShouldHaveRunInBatchesOfTwo)
		})

		Convey("When running a test package produces no output and exits with an error", func() {
			fixture.InBatchesOf(1).SetupAbnormalError("This really shouldn't happen...").RunTests()

			Convey("Panic should ensue", func() {
				So(fixture.recovered.Error(), ShouldEqual, "This really shouldn't happen...")
			})
		})

		Convey("When running test packages concurrently and a test package produces no output and exits with an error", func() {
			fixture.InBatchesOf(2).SetupAbnormalError("This really shouldn't happen...").RunTests()

			Convey("Panic should ensue", func() {
				So(fixture.recovered.Error(), ShouldEqual, "This really shouldn't happen...")
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
	packages     []string
	recovered    error
}

func NewTesterFixture() *TesterFixture {
	self := &TesterFixture{}
	self.shell = NewTimedShell()
	self.tester = NewConcurrentTester(self.shell)
	self.packages = []string{"a", "b", "c", "d", "e", "f"}
	return self
}

func (self *TesterFixture) InBatchesOf(batchSize int) *TesterFixture {
	self.tester.SetBatchSize(batchSize)
	return self
}

func (self *TesterFixture) SetupAbnormalError(message string) *TesterFixture {
	self.shell.setTripWire(message)
	return self
}

func (self *TesterFixture) RunTests() {
	defer func() {
		if r := recover(); r != nil {
			self.recovered = r.(error)
		}
	}()

	self.results = self.tester.TestAll(self.packages)
	self.compilations = self.shell.Compilations()
	self.executions = self.shell.Executions()
}

func (self *TesterFixture) ShouldHaveRecordOfCompilationCommands() {
	for i, pkg := range self.packages {
		command := self.compilations[i].Command
		So(command, ShouldEqual, "go test -i "+pkg)
	}
}

func (self *TesterFixture) ShouldHaveRecordOfExecutionCommands() {
	for i, pkg := range self.packages {
		So(self.executions[i].Command, ShouldEqual, "go test -v -timeout=-42s "+pkg)
	}
}

func (self *TesterFixture) ShouldHaveOneOutputPerInput() {
	So(len(self.results), ShouldEqual, len(self.packages))
}

func (self *TesterFixture) OutputShouldBeAsExpected() {
	for i, _ := range self.packages {
		So(self.results[i], ShouldEqual, self.executions[i].Command)
	}
}

func (self *TesterFixture) TestsShouldHaveRunContiguously() {
	So(self.shell.MaxConcurrentCommands(), ShouldEqual, 1)

	for i := 0; i < len(self.executions)-1; i++ {
		current := self.executions[i]
		next := self.executions[i+1]
		So(current.Started, ShouldHappenBefore, next.Started)
		So(current.Ended, ShouldHappenOnOrBefore, next.Started)
	}
}

func (self *TesterFixture) TestsShouldHaveRunInBatchesOfTwo() {
	So(self.shell.MaxConcurrentCommands(), ShouldEqual, concurrentBatchSize)
}

/**** Fakes ****/

type ShellCommand struct {
	Command string
	Started time.Time
	Ended   time.Time
}

type TimedShell struct {
	executions   []*ShellCommand
	compilations []*ShellCommand
	panicMessage string
}

func (self *TimedShell) Compilations() []*ShellCommand {
	return self.compilations
}

func (self *TimedShell) Executions() []*ShellCommand {
	return self.executions
}

func (self *TimedShell) MaxConcurrentCommands() int {
	var concurrent int

	for x, current := range self.executions {
		concurrentWith_x := 1
		for y, comparison := range self.executions {
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
	return comparison.Started.After(current.Started) && comparison.Started.Before(current.Ended)
}

func (self *TimedShell) setTripWire(message string) {
	self.panicMessage = message
}

func (self *TimedShell) Execute(name string, args ...string) (output string, err error) {
	if self.panicMessage != "" && args[1] == "-v" {
		return "", errors.New(self.panicMessage)
	}

	command := self.composeCommand(name + " " + strings.Join(args, " "))
	output = command.Command

	if strings.Contains(command.Command, " -i ") {
		self.compilations = append(self.compilations, command)
	} else {
		self.executions = append(self.executions, command)
	}
	return
}
func (self *TimedShell) composeCommand(commandText string) *ShellCommand {
	start := time.Now()
	time.Sleep(nap)
	end := time.Now()
	return &ShellCommand{commandText, start, end}
}
func (self *TimedShell) Getenv(key string) string {
	panic("NOT SUPPORTED")
}
func (self *TimedShell) Setenv(key, value string) error {
	panic("NOT SUPPORTED")
}

func NewTimedShell() *TimedShell {
	self := &TimedShell{}
	self.executions = []*ShellCommand{}
	self.compilations = []*ShellCommand{}
	return self
}

var nap, _ = time.ParseDuration("10ms")
var _ = fmt.Sprintf("fmt")
