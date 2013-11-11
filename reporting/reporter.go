package reporting

type Reporter interface {
	BeginStory(story *StoryReport)
	Enter(scope *ScopeReport)
	Report(r *AssertionResult)
	Exit()
	EndStory()
}

func (self *reporters) BeginStory(story *StoryReport) {
	for _, r := range self.collection {
		r.BeginStory(story)
	}
}
func (self *reporters) Enter(scope *ScopeReport) {
	for _, r := range self.collection {
		r.Enter(scope)
	}
}
func (self *reporters) Report(report *AssertionResult) {
	for _, x := range self.collection {
		x.Report(report)
	}
}
func (self *reporters) Exit() {
	for _, r := range self.collection {
		r.Exit()
	}
}
func (self *reporters) EndStory() {
	for _, r := range self.collection {
		r.EndStory()
	}
}

type reporters struct {
	collection []Reporter
}

func NewReporters(collection ...Reporter) *reporters {
	self := reporters{collection}
	return &self
}
