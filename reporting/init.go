package reporting

import (
	"os"
	"fmt"
	"strings"
)

const newline = "\n"

const (
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
	yellowColor = "\033[33m"
	redColor    = "\033[31m"
	resetColor  = "\033[0m"
)

func init() {
	if !xterm() {
		greenColor, yellowColor, redColor, resetColor = "", "", "", ""
	}
}

func xterm() bool {
	return strings.Contains(fmt.Sprintf("%v", os.Environ()), " TERM=xterm")
}
