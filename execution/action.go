package execution

import (
	"github.com/smartystreets/goconvey/gotest"
)

func (self *Action) Invoke() {
	self.action()
}

type Action struct {
	action func()
	name   string
}

func NewAction(action func()) *Action {
	return &Action{action: action, name: functionName(action)}
}

func NewSkippedAction(action func()) *Action {
	self := &Action{}
	self.name = gotest.ResolveExternalFileAndLine()
	self.action = action
	return self
}
