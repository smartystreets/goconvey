package convey

import "github.com/smartystreets/goconvey/gotest"

func (self *action) Invoke() {
	self.wrapped()
}

type action struct {
	wrapped func()
	name    string
}

func newAction(wrapped func()) *action {
	self := new(action)
	self.name = functionName(wrapped)
	self.wrapped = wrapped
	return self
}

func newSkippedAction(wrapped func()) *action {
	self := new(action)

	// The choice to use the filename and line number as the action name
	// reflects the need for something unique but also that corresponds
	// in a determinist way to the action itself.
	self.name = gotest.FormatExternalFileAndLine()
	self.wrapped = wrapped
	return self
}
