package reporting

import (
	"errors"
	"testing"

	"github.com/smartystreets/goconvey/printing"
)

func TestDotReporter(t *testing.T) {
	monochrome()

	file := newMemoryFile()
	printer := printing.NewPrinter(file)
	reporter := NewDotReporter(printer)

	reporter.Report(successReport)
	reporter.Report(failureReport)
	reporter.Report(erroredReport)
	reporter.Report(skippedReport)

	if file.buffer != dotSuccess+dotFailure+dotError+dotSkip {
		t.Errorf("\nExpected: '%s%s%s%s'\nActual:  '%s'", dotSuccess, dotFailure, dotError, dotSkip, file.buffer)
	}
}

var (
	successReport *AssertionResult = NewSuccessReport()
	failureReport *AssertionResult = NewFailureReport("failed")
	erroredReport *AssertionResult = NewErrorReport(errors.New("error"))
	skippedReport *AssertionResult = NewSkipReport()
)

type memoryFile struct {
	buffer string
}

func (self *memoryFile) Write(p []byte) (n int, err error) {
	self.buffer += string(p)
	return len(p), nil
}

func newMemoryFile() *memoryFile {
	return new(memoryFile)
}
