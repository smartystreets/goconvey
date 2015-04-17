package reporting

import (
	"fmt"
	"os"
	"runtime"
	"strings"
)

func init() {
	if !isXterm() {
		monochrome()
	}

	if runtime.GOOS == "windows" {
		success, failure, error_ = dotSuccess, dotFailure, dotError
	}
}

func BuildJsonReporter(t T, seed int64) Reporter {
	out := NewPrinter(NewConsole())
	return NewReporters(
		NewGoTestReporter(t),
		NewJsonReporter(out, seed))
}
func BuildDotReporter(t T, seed int64) Reporter {
	out := NewPrinter(NewConsole())
	return NewReporters(
		NewGoTestReporter(t),
		NewDotReporter(out, seed),
		NewProblemReporter(out),
		consoleStatistics)
}
func BuildStoryReporter(t T, seed int64) Reporter {
	out := NewPrinter(NewConsole())
	return NewReporters(
		NewGoTestReporter(t),
		NewStoryReporter(out, seed),
		NewProblemReporter(out),
		consoleStatistics)
}
func BuildSilentReporter(t T) Reporter {
	out := NewPrinter(NewConsole())
	return NewReporters(
		NewGoTestReporter(t),
		NewProblemReporter(out))
}

var (
	newline         = "\n"
	success         = "✔"
	failure         = "✘"
	error_          = "🔥"
	skip            = "⚠"
	dotSuccess      = "."
	dotFailure      = "x"
	dotError        = "E"
	dotSkip         = "S"
	errorTemplate   = "* %s \nLine %d: - %v \n%s\n"
	failureTemplate = "* %s \nLine %d:\n%s\n"
)

var (
	greenColor  = "\033[32m"
	redColor    = "\033[31m"
	resetColor  = "\033[0m"
	whiteColor  = "\033[37;1m"
	yellowColor = "\033[33m"
)

var consoleStatistics = NewStatisticsReporter(NewPrinter(NewConsole()))

// QuiteMode disables all console output symbols. This is only meant to be used
// for tests that are internal to goconvey where the output is distracting or
// otherwise not needed in the test output.
func QuietMode() {
	success, failure, error_, skip, dotSuccess, dotFailure, dotError, dotSkip = "", "", "", "", "", "", "", ""
}

func monochrome() {
	greenColor, yellowColor, redColor, resetColor = "", "", "", ""
}

func isXterm() bool {
	env := fmt.Sprintf("%v", os.Environ())
	return strings.Contains(env, " TERM=isXterm") ||
		strings.Contains(env, " TERM=xterm")
}

// This interface allows us to pass the *testing.T struct
// throughout the internals of this tool without ever
// having to import the "testing" package.
type T interface {
	Fail()
}
