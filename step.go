package goconvey

import (
	"testing"
)

type step struct {
	children []*step
	name string
	action func(C)
	test *testing.T
	complete bool
}

func NewStep(message string, action func(c C), t *testing.T) *step {
	return &step{name: message, action: action, test: t}
}

func (self *step) visit() {
	defer catch(self.test)()
	
	self.action(self)

	self.test.Errorf("Children: %d", len(self.children))

	index := 0
	for index = 0; index < len(self.children); index++ {
		child := self.children[index]

		if child.complete {
			continue
		}
		self.action(self)
		child.visit()
		break
	}

	if len(self.children) == 0 || index == len(self.children) - 1 {
		self.complete = true
	}
}

func (parent *step) Convey(situation string, action func(C)) {
	child := NewStep(situation, action, parent.test)
	parent.imprint(child)
}

func (parent *step) imprint(child *step) {
	for _, sibling := range parent.children {
		if sibling.name == child.name {
			return
		}
	}
	parent.test.Errorf("Imprinting child (%s) on parent (%s)", child.name, parent.name)
	parent.children = append(parent.children, child)
}

func catch(t *testing.T) func() {
	return func() {
		if r := recover(); r != nil {
			t.Errorf("%v", r)
		}
	}
}