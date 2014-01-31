package execution

import "github.com/smartystreets/goconvey/gotest"

type Registration struct {
	Situation string
	Action    *Action
	Test      gotest.T
	File      string
	Line      int
}

func (self *Registration) IsTopLevel() bool {
	return self.Test != nil
}

func NewRegistration(situation string, action *Action, test gotest.T) *Registration {
	file, line, _ := gotest.ResolveExternalCaller()
	self := &Registration{}
	self.Situation = situation
	self.Action = action
	self.Test = test
	self.File = file
	self.Line = line
	return self
}
