package execution

import (
	"fmt"

	"github.com/smartystreets/goconvey/gotest"
)

type Registration struct {
	Situation string
	Action    *Action
	Test      gotest.T
	File      string
	Line      int
	TestName  string
}

func (self *Registration) IsTopLevel() bool {
	return self.Test != nil
}

func (self *Registration) KeyName() string {
	return fmt.Sprintf("%s:%s", self.File, self.TestName)
}

func NewRegistration(situation string, action *Action, test gotest.T) *Registration {
	file, line, testName := gotest.ResolveExternalCallerWithTestName()
	self := &Registration{}
	self.Situation = situation
	self.Action = action
	self.Test = test
	self.File = file
	self.Line = line
	self.TestName = testName
	return self
}
