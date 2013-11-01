package executor

import (
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
	"time"
)

func TestConcurrentTester(t *testing.T) {
	var fixture *TesterFixture

	Convey("Subject: Controlled (and concurrent) execution of test packages", t, func() {
		fixture = NewTesterFixture()

		Convey("When tests for each package are executed", func() {
			fixture.RunTests()

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
			fixture.RunTests()

			Convey("Each package should be run synchronously and in the given order",
				fixture.CheckContiguousExecution)
		})

		Convey("When packages are tested in batches", func() {
			Convey("packages should be tested in batches while maintaining the given order", nil)
		})
	})
}

type TesterFixture struct {
	tester       *ConcurrentTester
	shell        *TimedShell
	results      []string
	compilations []*ShellCommand
	executions   []*ShellCommand
	packages     []string
}

func NewTesterFixture() *TesterFixture {
	self := &TesterFixture{}
	self.shell = NewTimedShell()
	self.tester = NewConcurrentTester(self.shell)
	self.packages = []string{"a", "b", "c", "d"}
	return self
}

func (self *TesterFixture) RunTests() {
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

func (self *TesterFixture) CheckContiguousExecution() {
	for i := 0; i < len(self.executions)-1; i++ {
		current := self.executions[i]
		next := self.executions[i+1]
		So(current.Started, ShouldHappenBefore, next.Started)
		So(current.Ended, ShouldHappenOnOrBefore, next.Started)
	}
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
}

func (self *TimedShell) Compilations() []*ShellCommand {
	return self.compilations
}

func (self *TimedShell) Executions() []*ShellCommand {
	return self.executions
}

func (self *TimedShell) Execute(name string, args ...string) (output string, err error) {
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
