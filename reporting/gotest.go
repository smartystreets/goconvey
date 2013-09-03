package reporting

import "github.com/smartystreets/goconvey/gotest"

func (self *gotestReporter) BeginStory(test gotest.T) {
	self.test = test
	self.ensureReady()
}

func (self *gotestReporter) Enter(title, id string) {
	self.ensureReady()
}

func (self *gotestReporter) Report(r *Report) {
	self.ensureReady()

	if !passed(r) {
		self.test.Fail()
	}
}

func (self *gotestReporter) Exit() {}

func (self *gotestReporter) EndStory() {
	self.test = nil
}

func (self *gotestReporter) ensureReady() {
	if self.test == nil {
		panic(registrationError)
	}
}

func NewGoTestReporter() *gotestReporter {
	self := gotestReporter{}
	return &self
}

type gotestReporter struct {
	test gotest.T
}

const registrationError = "You must register the actual *testing.T reference for this test in the first call to Convey(...)."

func passed(r *Report) bool {
	return r.Error == nil && r.Failure == ""
}
