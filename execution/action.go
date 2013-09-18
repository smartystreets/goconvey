package execution

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
	self.Name = resolveExternalFileAndLine()
	self.action = action
	return self
}
