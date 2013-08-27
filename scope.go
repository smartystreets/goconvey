package goconvey

import "fmt"

func (parent *scope) visit() {
	parent.action()
	parent.visitNextChild()
}
func (parent *scope) visitNextChild() {
	if parent.child >= len(parent.birthOrder) {
		return
	}
	child := parent.children[parent.birthOrder[parent.child]]
	child.visit()
	if child.visited() {
		parent.child++
	}
}

func (parent *scope) adopt(child *scope) {
	if parent.hasChild(child) {
		return
	}
	name := functionName(child.action)
	parent.birthOrder = append(parent.birthOrder, name)
	parent.children[name] = child
}
func (parent *scope) hasChild(child *scope) bool {
	for _, name := range parent.birthOrder {
		if name == functionName(child.action) {
			return true
		}
	}
	return false
}

func (self *scope) visited() bool {
	return self.child >= len(self.birthOrder)
}

func newScope(name string, action func()) *scope {
	fmt.Sprintf("")

	self := scope{name: name, action: action}
	self.children = make(map[string]*scope)
	self.birthOrder = []string{}
	return &self
}

type scope struct {
	name       string
	action     func()
	children   map[string]*scope
	birthOrder []string
	child      int
}
