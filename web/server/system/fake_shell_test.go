package system

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestFakeShell(t *testing.T) {
	var shell *FakeShell
	var output string
	var err error

	Convey("Subject: FakeShell", t, func() {
		shell = NewFakeShell()

		Convey("When executing go test", func() {
			output, err = shell.GoTest("/hi")
			shell.GoTest("/bye")

			Convey("The output should be an echo of the input", func() {
				So(output, ShouldEqual, "/hi")
			})

			Convey("There should be no error", func() {
				So(err, ShouldBeNil)
			})

			Convey("The shell should remember the directory of execution", func() {
				So(shell.Executions(), ShouldResemble, []string{"/hi", "/bye"})
			})
		})

		Convey("When setting an environment variable", func() {
			err := shell.Setenv("variable", "42")

			Convey("The value should persist", func() {
				So(shell.Getenv("variable"), ShouldEqual, "42")
			})

			Convey("The error should be nil", func() {
				So(err, ShouldBeNil)
			})
		})
	})
}
