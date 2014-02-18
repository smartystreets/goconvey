package convey

import "github.com/smartystreets/goconvey/gotest"

type registration struct {
	Situation string
	action    *action
	Test      t
	File      string
	Line      int
	Focus     bool
}

func (self *registration) IsTopLevel() bool {
	return self.Test != nil
}

func newRegistration(situation string, action *action, test t) *registration {
	file, line, _ := gotest.ResolveExternalCaller()
	self := new(registration)
	self.Situation = situation
	self.action = action
	self.Test = test
	self.File = file
	self.Line = line
	return self
}
