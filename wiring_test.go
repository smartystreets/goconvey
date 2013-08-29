/*

Problems:
	- tester forgets to pass in *testing.T the first time. How do we know it's the first time?

*/

package goconvey

import (
	// "fmt"
	"testing"
)

func TestParseRegistration(t *testing.T) {
	myRunner := newFakeRunner()
	situation := "Hello, World!"
	specRunner = myRunner
	var test goTest = &fakeGoTest{}
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
}

func TestParseRegistrationWithoutIncludinggoTestObject(t *testing.T) {
	myRunner := newFakeRunner()
	situation := "Hello, World!"
	specRunner = myRunner
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
	specRunner = myRunner
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

	specRunner = myRunner

	Convey("Hi there", 12345)

	t.Errorf("goTest should have panicked in Convey(...) and then recovered in the defer func().")
}

func TestParseFirstRegistrationAndNextRegistration_PreservesgoTest(t *testing.T) {
	myRunner := newFakeRunner()
	situation := "Hello, World!"
	nextSituation := "Goodbye, World!"
	specRunner = myRunner
	var test goTest = &fakeGoTest{}
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
	test      goTest
	situation string
	action    func()
	executed  bool
}

func newFakeRunner() *fakeRunner {
	f := fakeRunner{}
	f.action = func() {}
	return &f
}

func (self *fakeRunner) begin(test goTest) {
	self.test = test
}
func (self *fakeRunner) register(situation string, action func()) {
	self.situation = situation
	self.action = action
	if self.action != nil {
		self.action()
	}
}
func (self *fakeRunner) run() {
	self.executed = true
}

type fakeGoTest struct{}

func (self *fakeGoTest) Fail() {}
