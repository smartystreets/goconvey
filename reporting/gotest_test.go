package reporting

import (
	"testing"
)

func TestReporterRequiresGoTestToBegin(t *testing.T) {
	defer catch(t)
	reporter := NewGoTestReporter()
	reporter.BeginStory(nil)
}

func TestReporterRequiresGoTestToEnterScope(t *testing.T) {
	defer catch(t)
	reporter := NewGoTestReporter()
	reporter.Enter("hello", "world")
}

func TestReporterRequiresGoTestToReceiveReport(t *testing.T) {
	defer catch(t)
	reporter := NewGoTestReporter()
	reporter.Report(nil)
}

func TestReporterReceivesSuccessfulReport(t *testing.T) {
	reporter := NewGoTestReporter()
	test := &fakeTest{}
	reporter.BeginStory(test)
	reporter.Report(NewSuccessReport())

	if test.failed {
		t.Errorf("Should have have marked test as failed--the report reflected success.")
	}
}

func TestReporterReceivesFailureReport(t *testing.T) {
	reporter := NewGoTestReporter()
	test := &fakeTest{}
	reporter.BeginStory(test)
	reporter.Report(NewFailureReport("This is a failure."))

	if !test.failed {
		t.Errorf("Test should have been marked as failed (but it wasn't).")
	}
}

func TestReporterReceivesErrorReport(t *testing.T) {
	reporter := NewGoTestReporter()
	test := &fakeTest{}
	reporter.BeginStory(test)
	reporter.Report(NewErrorReport("This is an error."))

	if !test.failed {
		t.Errorf("Test should have been marked as failed (but it wasn't).")
	}
}

func TestReporterIsResetAtTheEndOfTheStory(t *testing.T) {
	defer catch(t)
	reporter := NewGoTestReporter()
	test := &fakeTest{}
	reporter.BeginStory(test)
	reporter.EndStory()
	reporter.Report(NewSuccessReport())
}

func catch(t *testing.T) {
	if r := recover(); r == registrationError {
		t.Log("Getting to this point means we've passed (because we caught a registration error appropriately).")
	} else {
		t.Errorf("Should have recovered a panic with: \n'%s'\n...but was: \n'%v'.", registrationError, r)
	}
}

type fakeTest struct {
	failed bool
}

func (self *fakeTest) Fail() {
	self.failed = true
}
