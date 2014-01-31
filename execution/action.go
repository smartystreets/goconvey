package execution

import "github.com/smartystreets/goconvey/gotest"

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

	// The choice to use the filename and line number as the action name
	// reflects the need for something unique but also that corresponds
	// in a determinist way to the action itself.
	self.name = gotest.FormatExternalFileAndLine()
	self.action = action
	return self
}
