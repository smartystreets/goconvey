package system

import (
	"errors"
	"fmt"
	. "github.com/smartystreets/goconvey/convey"
	"strings"
	"testing"
)

type FakeShell struct {
	outputByCommand map[string]string // name + args: output
	errorsByOutput  map[string]error  // output: err
}

func (self *FakeShell) Register(fullCommand string, output string, err error) {
	self.outputByCommand[fullCommand] = output
	self.errorsByOutput[output] = err
}

func (self *FakeShell) Execute(name string, args ...string) (output string, err error) {
	fullCommand := name + " " + strings.Join(args, " ")
	var exists bool = false
	if output, exists = self.outputByCommand[fullCommand]; !exists {
		panic(fmt.Sprintf("Missing command output for %s", fullCommand))
	}
	err = self.errorsByOutput[output]
	return
}

func NewFakeShell() *FakeShell {
	self := &FakeShell{}
	self.outputByCommand = map[string]string{}
	self.errorsByOutput = map[string]error{}
	return self
}

func TestFakeShell(t *testing.T) {
	var shell *FakeShell
	var output string
	var err error

	Convey("Subject: FakeShell", t, func() {
		shell = NewFakeShell()

		Convey("When executing an unrecognized command and arguments", func() {
			execute := func() { shell.Execute("Hello,", "World!") }

			Convey("panic ensues", func() {
				So(execute, ShouldPanic)
			})
		})

		Convey("When executing a known command with no error", func() {
			shell.Register("Hello, World!", "OUTPUT", nil)
			output, err = shell.Execute("Hello,", "World!")

			Convey("The expected output should be returned", func() {
				So(output, ShouldEqual, "OUTPUT")
			})

			Convey("No error should be returned", func() {
				So(err, ShouldBeNil)
			})
		})

		Convey("When executing a known command with a corresponding error", func() {
			shell.Register("Hello, World!", "OUTPUT", errors.New("Hi"))
			output, err = shell.Execute("Hello,", "World!")

			Convey("The expected output should be returned", func() {
				So(output, ShouldEqual, "OUTPUT")
			})

			Convey("The error should be returned", func() {
				So(err.Error(), ShouldEqual, "Hi")
			})
		})
	})
}
