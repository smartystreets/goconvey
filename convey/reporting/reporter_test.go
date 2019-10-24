package reporting

import (
	"runtime"
	"testing"
)

func TestEachNestedReporterReceivesTheCallFromTheContainingReporter(t *testing.T) {
	fake1 := newFakeReporter()
	fake2 := newFakeReporter()
	reporter := NewReporters(fake1, fake2)

	reporter.BeginStory(nil)
	assertTrue(t, fake1.begun)
	assertTrue(t, fake2.begun)

	reporter.Enter(NewScopeReport("scope"))
	assertTrue(t, fake1.entered)
	assertTrue(t, fake2.entered)

	reporter.Report(NewSuccessReport())
	assertTrue(t, fake1.reported)
	assertTrue(t, fake2.reported)

	reporter.Exit()
	assertTrue(t, fake1.exited)
	assertTrue(t, fake2.exited)

	reporter.EndStory()
	assertTrue(t, fake1.ended)
	assertTrue(t, fake2.ended)

	content := []byte("hi")
	written, err := reporter.Write(content)
	assertTrue(t, fake1.written)
	assertTrue(t, fake2.written)
	assertEqual(t, written, len(content))
	assertNil(t, err)

}

func assertTrue(t *testing.T, value bool) {
	if !value {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("Value should have been true (but was false). See line %d", line)
	}
}

func assertEqual(t *testing.T, expected, actual int) {
	if actual != expected {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("Value should have been %d (but was %d). See line %d", expected, actual, line)
	}
}

func assertNil(t *testing.T, err error) {
	if err != nil {
		_, _, line, _ := runtime.Caller(1)
		t.Errorf("Error should have been <nil> (but wasn't). See line %d. %v", err, line)
	}
}

type fakeReporter struct {
	begun    bool
	entered  bool
	reported bool
	exited   bool
	ended    bool
	written  bool
}

func newFakeReporter() *fakeReporter {
	return &fakeReporter{}
}

func (f *fakeReporter) BeginStory(story *StoryReport) {
	f.begun = true
}
func (f *fakeReporter) Enter(scope *ScopeReport) {
	f.entered = true
}
func (f *fakeReporter) Report(report *AssertionResult) {
	f.reported = true
}
func (f *fakeReporter) Exit() {
	f.exited = true
}
func (f *fakeReporter) EndStory() {
	f.ended = true
}
func (f *fakeReporter) Write(content []byte) (int, error) {
	f.written = true
	return len(content), nil
}
