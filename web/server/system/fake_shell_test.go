package system

import (
	"errors"
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

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
