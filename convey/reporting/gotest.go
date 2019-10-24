package reporting

type gotestReporter struct{ test T }

func (g *gotestReporter) BeginStory(story *StoryReport) {
	g.test = story.Test
}

func (g *gotestReporter) Enter(scope *ScopeReport) {}

func (g *gotestReporter) Report(r *AssertionResult) {
	if !passed(r) {
		g.test.Fail()
	}
}

func (g *gotestReporter) Exit() {}

func (g *gotestReporter) EndStory() {
	g.test = nil
}

func (g *gotestReporter) Write(content []byte) (written int, err error) {
	return len(content), nil // no-op
}

func NewGoTestReporter() *gotestReporter {
	return new(gotestReporter)
}

func passed(r *AssertionResult) bool {
	return r.Error == nil && r.Failure == ""
}
