package executor

import (
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
	"time"
)

func TestTester(t *testing.T) {
	var (
		tester       *ConcurrentTester
		shell        *TimedShell
		results      []string
		compilations []*ShellCommand
		executions   []*ShellCommand
	)

	Convey("Subject: Test controlled execution of tests", t, func() {
		shell = NewTimedShell()
		tester = NewConcurrentTester(shell)
		folders := []string{"a", "b", "c"}

		Convey("When packages are executed synchronously", func() {
			results = tester.TestAll(folders)
			compilations = shell.Compilations()
			executions = shell.Executions()

			Convey("The tester should build all dependencies of input folders", func() {
				for i, input := range folders {
					So(compilations[i].Command, ShouldEqual, "go test -i "+input)
				}
			})

			Convey("The tester should execute the tests in each folder with the correct arguments", func() {
				for i, input := range folders {
					So(executions[i].Command, ShouldEqual, "go test -v -timeout=-42s "+input)
				}
			})

			Convey("Each package should be run after the other in the given order", func() {
				for i := 0; i < len(executions)-1; i++ {
					current := executions[i]
					next := executions[i+1]
					So(current.Started, ShouldHappenBefore, next.Started)
					So(current.Ended, ShouldHappenOnOrBefore, next.Started)
				}
			})

			Convey("There should be a test output result for each input folder", func() {
				So(len(results), ShouldEqual, len(folders))

			})

			Convey("The output should be as expected", func() {
				for i, _ := range folders {
					So(results[i], ShouldEqual, executions[i].Command)
				}
			})
		})

		Convey("When packages are tested in batches", func() {
			Convey("packages should be tested in batches while maintaining the given order", nil)
		})
	})
}

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
	output = name + " " + strings.Join(args, " ")
	start := time.Now()
	nap, err := time.ParseDuration("10ms")
	if err != nil {
		panic(err)
	}
	time.Sleep(nap)
	end := time.Now()
	command := &ShellCommand{output, start, end}
	if strings.Contains(output, " -i ") {
		self.compilations = append(self.compilations, command)
	} else {
		self.executions = append(self.executions, command)
	}
	return
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
