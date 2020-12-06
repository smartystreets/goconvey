package convey

import (
	"fmt"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey/reporting"
)

func TestStackTrace(t *testing.T) {
	file, test := setupFileReporter()

	Convey("A", test, func() {
		So(1, ShouldEqual, 2)
	})

	if !strings.Contains(file.String(), "Failures:\n") {
		t.Errorf("Expected errors, found none.")
	}
	if strings.Contains(file.String(), "goroutine ") {
		t.Errorf("Found stack trace, expected none.")
	}

	Convey("A", test, StackFail, func() {
		So(1, ShouldEqual, 2)
	})

	if !strings.Contains(file.String(), "goroutine ") {
		t.Errorf("Expected stack trace, found none.")
	}
}

func TestSetDefaultStackMode(t *testing.T) {
	file, test := setupFileReporter()
	SetDefaultStackMode(StackFail) // the default is normally StackError
	defer SetDefaultStackMode(StackError)

	Convey("A", test, func() {
		So(1, ShouldEqual, 2)
	})

	if !strings.Contains(file.String(), "goroutine ") {
		t.Errorf("Expected stack trace, found none.")
	}
}

func TestStackModeMultipleInvocationInheritance(t *testing.T) {
	file, test := setupFileReporter()

	// initial convey should default to StaskError, so no stack trace
	Convey("A", test, FailureContinues, func() {
		So(1, ShouldEqual, 2)

		// nested convey has explicit StaskFail, so should emit stack trace
		Convey("B", StackFail, func() {
			So(1, ShouldEqual, 2)
		})
	})

	stackCount := strings.Count(file.String(), "goroutine ")
	if stackCount != 1 {
		t.Errorf("Expected 1 stack trace, found %d.", stackCount)
		fmt.Printf("RESULT: %s \n", file.String())
	}
}

func TestStackModeMultipleInvocationInheritance2(t *testing.T) {
	file, test := setupFileReporter()

	// Explicit StackFail, expect stack trace
	Convey("A", test, FailureContinues, StackFail, func() {
		So(1, ShouldEqual, 2)

		// Nested Convey inherits StackFail, expect stack trace
		Convey("B", func() {
			So(1, ShouldEqual, 2)
		})
	})

	stackCount := strings.Count(file.String(), "goroutine ")
	if stackCount != 2 {
		t.Errorf("Expected 2 stack traces, found %d.", stackCount)
	}
}

func TestStackModeMultipleInvocationInheritance3(t *testing.T) {
	file, test := setupFileReporter()

	// Explicit StackFail, expect stack trace
	Convey("A", test, FailureContinues, StackFail, func() {
		So(1, ShouldEqual, 2)

		// Nested Convey explicitly sets StackError, so no stack trace
		Convey("B", StackError, func() {
			So(1, ShouldEqual, 2)
		})
	})

	stackCount := strings.Count(file.String(), "goroutine ")
	if stackCount != 1 {
		t.Errorf("Expected 1 stack trace1, found %d.", stackCount)
	}
}

func setupFileReporter() (*memoryFile, *fakeGoTest) {
	//monochrome()
	file := newMemoryFile()
	printer := reporting.NewPrinter(file)
	reporter := reporting.NewProblemReporter(printer)
	testReporter = reporter

	return file, new(fakeGoTest)
}

////////////////// memoryFile ////////////////////

type memoryFile struct {
	buffer string
}

func (self *memoryFile) Write(p []byte) (n int, err error) {
	self.buffer += string(p)
	return len(p), nil
}

func (self *memoryFile) String() string {
	return self.buffer
}

func newMemoryFile() *memoryFile {
	return new(memoryFile)
}

// func monochrome() {
// 	greenColor, yellowColor, redColor, resetColor = "", "", "", ""
// }
