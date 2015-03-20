package reporting

type gotestReporter struct{ test T }

func (self *gotestReporter) Enter(scope *ScopeReport) {}

func (self *gotestReporter) Report(r *AssertionResult) {
	if !passed(r) {
		self.test.Fail()
	}
}

func (self *gotestReporter) Exit() {}

func (self *gotestReporter) Close() {}

func (self *gotestReporter) Write(content []byte) (written int, err error) {
	return len(content), nil // no-op
}

func NewGoTestReporter(t T) *gotestReporter {
	return &gotestReporter{t}
}

func passed(r *AssertionResult) bool {
	return r.Error == nil && r.Failure == ""
}
