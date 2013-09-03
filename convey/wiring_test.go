package convey

import (
	"github.com/smartystreets/goconvey/convey/execution"
	"testing"
)

func TestParseTopLevelRegistration(t *testing.T) {
	myRunner := newFakeRunner()
	situation := "Hello, World!"
	execution.SpecRunner = myRunner
	var test execution.GoTest = &fakeGoTest{}
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
	execution.SpecRunner = myRunner
	executed := false
	action := func() {
		executed = true
	}

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
	execution.SpecRunner = myRunner
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

	execution.SpecRunner = myRunner

	Convey("Hi there", 12345)

	t.Errorf("goTest should have panicked in Convey(...) and then recovered in the defer func().")
}

func TestParseFirstRegistrationAndNextRegistration_PreservesGoTest(t *testing.T) {
	myRunner := newFakeRunner()
	situation := "Hello, World!"
	nextSituation := "Goodbye, World!"
	execution.SpecRunner = myRunner
	var test execution.GoTest = &fakeGoTest{}
	executed := 0
	action := func() {
		executed++
	}

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
	test          execution.GoTest
	situation     string
	action        func()
	runnerStarted bool
}

func newFakeRunner() *fakeRunner {
	f := fakeRunner{}
	f.action = func() {}
	return &f
}

func (self *fakeRunner) Begin(test execution.GoTest, situation string, action func()) {
	self.test = test
	self.Register(situation, action)
}
func (self *fakeRunner) Register(situation string, action func()) {
	self.situation = situation
	self.action = action
	if self.action != nil {
		self.action()
	}
}
func (self *fakeRunner) RegisterReset(action func()) {}
func (self *fakeRunner) Run() {
	self.runnerStarted = true
}
func (self *fakeRunner) UpgradeReporter(out execution.Reporter) {}

type fakeGoTest struct{}

func (self *fakeGoTest) Fail() {}

var test execution.GoTest = &fakeGoTest{}
