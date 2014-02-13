package execution

import "github.com/smartystreets/goconvey/gotest"

type Registration struct {
	Situation string
	Action    *Action
	Test      T
	File      string
	Line      int
	Focus     bool
}

func (self *Registration) IsTopLevel() bool {
	return self.Test != nil
}

func NewRegistration(situation string, action *Action, test T) *Registration {
	file, line, _ := gotest.ResolveExternalCaller()
	self := new(Registration)
	self.Situation = situation
	self.Action = action
	self.Test = test
	self.File = file
	self.Line = line
	return self
}

// This interface allows us to pass the *testing.T struct
// throughout the internals of this tool without ever
// having to import the "testing" package.
type T interface {
	Fail()
}
