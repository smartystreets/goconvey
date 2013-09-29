package execution

import (
	"github.com/smartystreets/goconvey/gotest"
)

func (self *Action) Invoke() {
	self.action()
}

type Action struct {
	action func()
	Name   string
}

func NewAction(action func()) *Action {
	return &Action{action: action, Name: functionName(action)}
}

func NewSkippedAction(action func()) *Action {
	self := &Action{}
	self.Name = gotest.ResolveExternalFileAndLine()
	self.action = action
	return self
}
