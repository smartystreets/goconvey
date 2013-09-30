package reporting

import (
	"github.com/smartystreets/goconvey/gotest"
)

type Reporter interface {
	BeginStory(test gotest.T)
	Enter(title, id string)
	Report(r *AssertionReport)
	Exit()
	EndStory()
}

func (self *reporters) BeginStory(test gotest.T) {
	for _, r := range self.collection {
		r.BeginStory(test)
	}
}
func (self *reporters) Enter(title, id string) {
	for _, r := range self.collection {
		r.Enter(title, id)
	}
}
func (self *reporters) Report(r *AssertionReport) {
	for _, x := range self.collection {
		x.Report(r)
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
