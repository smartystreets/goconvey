package convey

import (
	"fmt"
	"regexp"
	"strings"
	"testing"

	"github.com/smartystreets/goconvey/convey/reporting"
)

var goroutineRE = regexp.MustCompile(`goroutine \d+ \[`)

// countGoroutines takes a test output file and counts the number of goroutines
// that were mentioned inside it. This does this by hunting for lines such as
// "goroutine 42 [running]", while excluding secondary mentions of already-counted
// goroutines.
func countGoroutines(testOutput string) int {
	return len(goroutineRE.FindAllStringSubmatch(testOutput, -1))
}

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
	stackCount := countGoroutines(file.String())
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

	stackCount := countGoroutines(file.String())
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

	stackCount := countGoroutines(file.String())
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

func (mf *memoryFile) Write(p []byte) (n int, err error) {
	mf.buffer += string(p)
	return len(p), nil
}

func (mf *memoryFile) String() string {
	return mf.buffer
}

func newMemoryFile() *memoryFile {
	return new(memoryFile)
}

// func monochrome() {
// 	greenColor, yellowColor, redColor, resetColor = "", "", "", ""
// }
