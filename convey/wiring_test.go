package convey

import (
	"github.com/smartystreets/goconvey/execution"
	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/reporting"
	"testing"
)

func TestParseTopLevelRegistration(t *testing.T) {
	myRunner := newFakeRunner()
	situation := "Hello, World!"
	runner = myRunner
	var test gotest.T = &fakeGoTest{}
	executed := false
	action := func() {
		executed = true
	}

	Convey(situation, test, action)

	if myRunner.test != test {
		t.Errorf("Should have received a reference to the test object (was '%v').", myRunner.test)
	}

	if myRunner.situation != situation {
		t.Errorf("Situation should have been '%s', was '%s'.", situation, myRunner.situation)
	}

	if !executed {
		t.Error("Action should have been captured but was not.")
	}

	if !myRunner.runnerStarted {
		t.Error("Runner should have been .Run().")
	}
}

func TestParseRegistrationWithoutIncludingGoTestObject(t *testing.T) {
	myRunner := newFakeRunner()
	situation := "Hello, World!"
	runner = myRunner
	executed := false
	action := func() { executed = true }

	Convey(situation, action)

	if myRunner.test != nil {
		t.Errorf("goTest object should have been nil (was '%v').", myRunner.test)
	}

	if myRunner.situation != situation {
		t.Errorf("Situation should have been '%s', was '%s'.", situation, myRunner.situation)
	}

	if !executed {
		t.Error("Action should have been captured but was not.")
	}

	if myRunner.runnerStarted {
		t.Error("Runner should NOT have been .Run(), but it was.")
	}
}

func TestParseRegistrationMissingRequiredElements(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != "You must provide a name (string), then a *testing.T (if in outermost scope), and then an action (func())." {
				t.Errorf("Incorrect panic message.")
			}
		}
	}()

	Convey()

	t.Errorf("goTest should have panicked in Convey(...) and then recovered in the defer func().")
}

func TestParseRegistration_MissingNameString(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			if r != parseError {
				t.Errorf("Incorrect panic message.")
			}
		}
	}()

	myRunner := newFakeRunner()
	runner = myRunner
	action := func() {}

	Convey(action)

	t.Errorf("goTest should have panicked in Convey(...) and then recovered in the defer func().")
}

func TestParseRegistration_MissingActionFunc(t *testing.T) {
	myRunner := newFakeRunner()

	defer func() {
		if r := recover(); r != nil {
			if r != parseError {
				t.Errorf("Incorrect panic message: '%s'", r)
			}
		}
	}()

	runner = myRunner

	Convey("Hi there", 12345)

	t.Errorf("goTest should have panicked in Convey(...) and then recovered in the defer func().")
}

func TestParseFirstRegistrationAndNextRegistration_PreservesGoTest(t *testing.T) {
	myRunner := newFakeRunner()
	situation := "Hello, World!"
	nextSituation := "Goodbye, World!"
	runner = myRunner
	var test gotest.T = &fakeGoTest{}
	executed := 0
	action := func() { executed++ }

	Convey(situation, test, action)
	Convey(nextSituation, action)

	if myRunner.test == nil {
		t.Errorf("goTest object should NOT havebeen nil, but was.")
	}

	if myRunner.situation != nextSituation {
		t.Errorf("Situation should have been '%s', was '%s'.", nextSituation, myRunner.situation)
	}

	if executed != 2 {
		t.Error("Action should have been captured but was not.")
	}
}

type fakeRunner struct {
	test          gotest.T
	situation     string
	action        *execution.Action
	runnerStarted bool
}

func newFakeRunner() *fakeRunner {
	f := fakeRunner{}
	f.action = execution.NewAction(func() {})
	return &f
}

func (self *fakeRunner) Begin(entry *execution.Registration) {
	self.test = entry.Test
	self.Register(entry)
}
func (self *fakeRunner) Register(entry *execution.Registration) {
	self.situation = entry.Situation
	self.action = entry.Action
	if self.action != nil {
		self.action.Invoke()
	}
}
func (self *fakeRunner) RegisterReset(action *execution.Action) {}
func (self *fakeRunner) Run() {
	self.runnerStarted = true
}
func (self *fakeRunner) UpgradeReporter(out reporting.Reporter) {}

type fakeGoTest struct{}

func (self *fakeGoTest) Fail()                                     {}
func (self *fakeGoTest) Fatalf(format string, args ...interface{}) {}

var test gotest.T = &fakeGoTest{}
