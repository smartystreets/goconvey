package reporting

import (
	"errors"
	"testing"

	"github.com/smartystreets/goconvey/printing"
)

func TestDotReporterAssertionPrinting(t *testing.T) {
	monochrome()
	file := newMemoryFile()
	printer := printing.NewPrinter(file)
	reporter := NewDotReporter(printer)

	reporter.Report(NewSuccessReport())
	reporter.Report(NewFailureReport("failed"))
	reporter.Report(NewErrorReport(errors.New("error")))
	reporter.Report(NewSkipReport())

	expected := dotSuccess + dotFailure + dotError + dotSkip

	if file.buffer != expected {
		t.Errorf("\nExpected: '%s'\nActual:  '%s'", expected, file.buffer)
	}
}

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
