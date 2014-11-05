package system

import (
	"errors"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestShellCommandComposition(t *testing.T) {
	var (
		buildFailed      = Command{Error: errors.New("BUILD FAILURE!")}
		buildSucceeded   = Command{Output: "ok"}
		goConvey         = Command{Output: "[fmt github.com/smartystreets/goconvey/convey net/http net/http/httptest path runtime strconv strings testing time]"}
		noGoConvey       = Command{Output: "[fmt net/http net/http/httptest path runtime strconv strings testing time]"}
		noCoveragePassed = Command{Output: "PASS\nok  	github.com/smartystreets/goconvey/examples	0.012s"}
		coveragePassed   = Command{Output: "PASS\ncoverage: 100.0% of statements\nok  	github.com/smartystreets/goconvey/examples	0.012s"}
		coverageFailed   = Command{
			Error: errors.New("Tests bombed!"),
			Output: "--- FAIL: TestIntegerManipulation (0.00 seconds)\nFAIL\ncoverage: 100.0% of statements\nexit status 1\nFAIL	github.com/smartystreets/goconvey/examples	0.013s",
		}
		coverageFailedTimeout = Command{
			Error:  errors.New("Tests bombed!"),
			Output: "=== RUN SomeTest\n--- PASS: SomeTest (0.00 seconds)\n=== RUN TimeoutTest\npanic: test timed out after 5s\n\ngoroutine 27 [running]:\n",
		}
	)

	const (
		yesCoverage = true
		noCoverage  = false
	)

	Convey("When attempting to run tests with coverage flags", t, func() {
		Convey("And buildSucceeded failed", func() {
			result := runWithCoverage(buildFailed, goConvey, noCoverage, "", "", "", "", nil)

			Convey("Then no action should be taken", func() {
				So(result, ShouldResemble, buildFailed)
			})
		})

		Convey("And coverage is not wanted", func() {
			result := runWithCoverage(buildSucceeded, goConvey, noCoverage, "", "", "", "", nil)

			Convey("Then no action should be taken", func() {
				So(result, ShouldResemble, buildSucceeded)
			})
		})

		Convey("And the package being tested usees the GoConvey DSL (`convey` package)", func() {
			result := runWithCoverage(buildSucceeded, goConvey, yesCoverage, "reportsPath", "/directory", "go", "5s", []string{"-arg1", "-arg2"})

			Convey("The returned command should be well formed (and include the -json flag)", func() {
				So(result, ShouldResemble, Command{
					directory:  "/directory",
					executable: "go",
					arguments:  []string{"test", "-v", "-coverprofile=reportsPath", "-covermode=set", "-timeout=5s", "-json", "-arg1", "-arg2"},
				})
			})
		})

		Convey("And the package being tested does NOT use the GoConvey DSL", func() {
			result := runWithCoverage(buildSucceeded, noGoConvey, yesCoverage, "reportsPath", "/directory", "go", "1s", []string{"-arg1", "-arg2"})

			Convey("The returned command should be well formed (and NOT include the -json flag)", func() {
				So(result, ShouldResemble, Command{
					directory:  "/directory",
					executable: "go",
					arguments:  []string{"test", "-v", "-coverprofile=reportsPath", "-covermode=set", "-timeout=1s", "-arg1", "-arg2"},
				})
			})
		})

		Convey("And the package being tested specifies an alternate covermode", func() {
			result := runWithCoverage(buildSucceeded, noGoConvey, yesCoverage, "reportsPath", "/directory", "go", "1s", []string{"-covermode=atomic"})

			Convey("The returned command should allow the alternate value", func() {
				So(result, ShouldResemble, Command{
					directory:  "/directory",
					executable: "go",
					arguments:  []string{"test", "-v", "-coverprofile=reportsPath", "-timeout=1s", "-covermode=atomic"},
				})
			})
		})

		Convey("And the package being tested specifies an alternate timeout", func() {
			result := runWithCoverage(buildSucceeded, noGoConvey, yesCoverage, "reportsPath", "/directory", "go", "1s", []string{"-timeout=5s"})

			Convey("The returned command should allow the alternate value", func() {
				So(result, ShouldResemble, Command{
					directory:  "/directory",
					executable: "go",
					arguments:  []string{"test", "-v", "-coverprofile=reportsPath", "-covermode=set", "-timeout=5s"},
				})
			})
		})

	})

	Convey("When attempting to run tests without the coverage flags", t, func() {
		Convey("And tests already succeeded with coverage", func() {
			result := runWithoutCoverage(buildSucceeded, coveragePassed, goConvey, "/directory", "go", "1s", []string{"-arg1", "-arg2"})

			Convey("Then no action should be taken", func() {
				So(result, ShouldResemble, coveragePassed)
			})
		})

		Convey("And tests already failed (legitimately) with coverage", func() {
			result := runWithoutCoverage(buildSucceeded, coverageFailed, goConvey, "/directory", "go", "1s", []string{"-arg1", "-arg2"})

			Convey("Then no action should be taken", func() {
				So(result, ShouldResemble, coverageFailed)
			})
		})

		Convey("And tests already failed (timeout) with coverage", func() {
			result := runWithoutCoverage(buildSucceeded, coverageFailedTimeout, goConvey, "/directory", "go", "1s", []string{"-arg1", "-arg2"})

			Convey("Then no action should be taken", func() {
				So(result, ShouldResemble, coverageFailedTimeout)
			})
		})

		Convey("And the build failed earlier", func() {
			result := runWithoutCoverage(buildFailed, Command{}, goConvey, "/directory", "go", "1s", []string{"-arg1", "-arg2"})

			Convey("Then no action should be taken", func() {
				So(result, ShouldResemble, buildFailed)
			})
		})

		Convey("And the package being tested uses the GoConvey DSL (`convey` package)", func() {
			result := runWithoutCoverage(buildSucceeded, buildSucceeded, goConvey, "/directory", "go", "1s", []string{"-arg1", "-arg2"})

			Convey("Then the returned command should be well formed (and include the -json flag)", func() {
				So(result, ShouldResemble, Command{
					directory:  "/directory",
					executable: "go",
					arguments:  []string{"test", "-v", "-timeout=1s", "-json", "-arg1", "-arg2"},
				})
			})
		})

		Convey("And the package being tested does NOT use the GoConvey DSL", func() {
			result := runWithoutCoverage(buildSucceeded, noCoveragePassed, noGoConvey, "/directory", "go", "1s", []string{"-arg1", "-arg2"})

			Convey("Then the returned command should be well formed (and NOT include the -json flag)", func() {
				So(result, ShouldResemble, Command{
					directory:  "/directory",
					executable: "go",
					arguments:  []string{"test", "-v", "-timeout=1s", "-arg1", "-arg2"},
				})
			})
		})

		Convey("And the package being tested specifies an alternate timeout", func() {
			result := runWithoutCoverage(buildSucceeded, buildSucceeded, noGoConvey, "/directory", "go", "1s", []string{"-timeout=5s"})

			Convey("The returned command should allow the alternate value", func() {
				So(result, ShouldResemble, Command{
					directory:  "/directory",
					executable: "go",
					arguments:  []string{"test", "-v", "-timeout=5s"},
				})
			})
		})

	})

	Convey("When generating coverage reports", t, func() {
		Convey("And the previous command failed for any reason (compilation or failed tests)", func() {
			result := generateReports(buildFailed, yesCoverage, "/directory", "go", "reportData", "reportHTML")

			Convey("Then no action should be taken", func() {
				So(result, ShouldResemble, buildFailed)
			})
		})

		Convey("And coverage reports are unwanted", func() {
			result := generateReports(noCoveragePassed, noCoverage, "/directory", "go", "reportData", "reportHTML")

			Convey("Then no action should beg taken", func() {
				So(result, ShouldResemble, noCoveragePassed)
			})
		})

		Convey("And tests passed and coverage reports are wanted", func() {
			result := generateReports(coveragePassed, yesCoverage, "/directory", "go", "reportData", "reportHTML")

			Convey("Then the resulting command should be well-formed", func() {
				So(result, ShouldResemble, Command{
					directory:  "/directory",
					executable: "go",
					arguments:  []string{"tool", "cover", "-html=reportData", "-o", "reportHTML"},
				})
			})
		})
	})
}
