package goconvey

import "fmt"

func (parent *scope) Visit() {
	parent.action()
	parent.visitNextChild()
}
func (parent *scope) visitNextChild() {
	if parent.child >= len(parent.birthOrder) {
		return
	}
	child := parent.children[parent.birthOrder[parent.child]]
	child.Visit()
	if child.Visited() {
		parent.child++
	}
}

func (parent *scope) Adopt(child *scope) {
	if !parent.hasChild(child) {
		parent.adopt(child)
	}
}
func (parent *scope) hasChild(child *scope) bool {
	for _, name := range parent.birthOrder {
		if name == functionName(child.action) {
			return true
		}
	}
	return false
}
func (parent *scope) adopt(child *scope) {
	name := functionName(child.action)
	parent.birthOrder = append(parent.birthOrder, name)
	parent.children[name] = child
}

func (self *scope) Visited() bool {
	return self.child >= len(self.birthOrder)
}

func NewScope(name string, action func()) *scope {
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
