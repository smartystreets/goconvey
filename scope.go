package goconvey

import "fmt"

func (self *scope) Run() bool {
	fmt.Sprintf("")
	// fmt.Println("Running:", self.subject)
	self.action()
	childFinished := true

	// fmt.Println(self.conveyIndex, self.conveyOrder, self.convey)

	if self.conveyIndex < len(self.conveyOrder) {
		child := self.convey[self.conveyOrder[self.conveyIndex]]
		childFinished = child.Run()
		if childFinished {
			self.conveyIndex++
		}
	}
	
	for _, resetName := range self.resetOrder {
		reset := self.reset[resetName]
		reset.action()
	}

	finished := childFinished && self.conveyIndex >= len(self.conveyOrder)
	// fmt.Printf("Finished (%s): %t\n", self.subject, finished)
	return finished
}

func (self *scope) Convey(child *scope) {
	for _, name := range self.conveyOrder {
		if name == functionName(child.action) {
			return
		}
	}
	//fmt.Println("Attaching child:", child.subject)
	self.convey, self.conveyOrder = self.add(child, self.convey, self.conveyOrder)
}
func (self *scope) Reset(child *scope) {
	for _, name := range self.resetOrder {
		if name == functionName(child.action) {
			return
		}
	}
	self.reset, self.resetOrder = self.add(child, self.reset, self.resetOrder)
}
func(self *scope) add(child *scope, group map[string]*scope, order []string) (map[string]*scope, []string) {
	name := functionName(child.action)
	order = append(order, name)
	group[name] = child
	return group, order
}

func NewScope(subject string, action func()) *scope {
	self := scope{subject: subject, action: action}
	self.convey = make(map[string]*scope)
	self.reset = make(map[string]*scope)
	self.conveyOrder = []string{}
	self.resetOrder = []string{}
	return &self
}

type scope struct {
	subject string
	action func()
	convey map[string]*scope
	reset map[string]*scope

	conveyOrder []string
	resetOrder []string

	conveyIndex int
}