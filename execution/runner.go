package execution

import (
	"fmt"
	"reflect"
	"runtime"
)

func (self *runner) Begin(test GoTest) {
	self.currentTest = test
}

func (self *runner) Register(situation string, action func()) {
	parentAction := self.link(action)
	parent := self.accessScope(parentAction)
	child := newScope(situation, action)
	parent.adopt(child)
}
func (self *runner) link(action func()) (parentAction string) {
	parentAction, childAction := resolveParentChild(action)
	self.linkTo(topLevel, parentAction)
	self.linkTo(parentAction, childAction)
	return
}
func (self *runner) linkTo(value, name string) {
	if self.chain[name] == "" {
		self.chain[name] = value
	}
}
func (self *runner) accessScope(current string) *scope {
	if self.chain[current] == topLevel {
		return self.top
	}
	breadCrumbs := self.trail(current)
	return self.follow(breadCrumbs)
}
func (self *runner) trail(start string) []string {
	breadCrumbs := []string{start, self.chain[start]}
	for {
		next := self.chain[last(breadCrumbs)]
		if next == topLevel {
			break
		} else {
			breadCrumbs = append(breadCrumbs, next)
		}
	}
	return breadCrumbs[:len(breadCrumbs)-1]
}
func (self *runner) follow(trail []string) *scope {
	var accessed = self.top

	for x := len(trail) - 1; x >= 0; x-- {
		accessed = accessed.children[trail[x]]
	}
	return accessed
}

func (self *runner) RegisterReset(func()) {

}

func (self *runner) Run() {
	for !self.top.visited() {
		self.top.visit()
	}
}

type runner struct {
	top         *scope
	chain       map[string]string
	currentTest GoTest
}

func NewSpecRunner() *runner {
	fmt.Sprintf("")

	self := runner{}
	self.top = newScope(topLevel, func() {})
	self.chain = make(map[string]string)
	return &self
}

func resolveParentChild(function func()) (parent, child string) {
	parent = resolveExternalCaller()
	child = functionName(function)
	return
}

func functionName(function func()) string {
	return runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
}
func resolveExternalCaller() string {
	caller_id, _, _, _ := runtime.Caller(4) // TODO: how to better encapsulate this magic number (move it closer to user code)
	return runtime.FuncForPC(caller_id).Name()
}

func last(group []string) string {
	return group[len(group)-1]
}

var topLevel = "TOP"
