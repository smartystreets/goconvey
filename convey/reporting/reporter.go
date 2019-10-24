package reporting

import "io"

type Reporter interface {
	BeginStory(story *StoryReport)
	Enter(scope *ScopeReport)
	Report(r *AssertionResult)
	Exit()
	EndStory()
	io.Writer
}

type reporters struct{ collection []Reporter }

func (r *reporters) BeginStory(s *StoryReport) { r.foreach(func(r Reporter) { r.BeginStory(s) }) }
func (r *reporters) Enter(s *ScopeReport)      { r.foreach(func(r Reporter) { r.Enter(s) }) }
func (r *reporters) Report(a *AssertionResult) { r.foreach(func(r Reporter) { r.Report(a) }) }
func (r *reporters) Exit()                     { r.foreach(func(r Reporter) { r.Exit() }) }
func (r *reporters) EndStory()                 { r.foreach(func(r Reporter) { r.EndStory() }) }

func (r *reporters) Write(contents []byte) (written int, err error) {
	r.foreach(func(r Reporter) {
		written, err = r.Write(contents)
	})
	return written, err
}

func (r *reporters) foreach(action func(Reporter)) {
	for _, r := range r.collection {
		action(r)
	}
}

func NewReporters(collection ...Reporter) *reporters {
	self := new(reporters)
	self.collection = collection
	return self
}
