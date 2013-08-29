package execution

import (
	"fmt"
	"reflect"
	"runtime"
)

func (self *scopeRunner) Begin(test GoTest) {
	self.currentTest = test
}

func (self *scopeRunner) Register(situation string, action func()) {
	parentAction := self.link(action)
	parent := self.accessScope(parentAction)
	child := newScope(situation, action)
	parent.adopt(child)
}
func (self *scopeRunner) link(action func()) (parentAction string) {
	parentAction, childAction := resolveParentChild(action)
	self.linkTo(topLevel, parentAction)
	self.linkTo(parentAction, childAction)
	return
}
func (self *scopeRunner) linkTo(value, name string) {
	if self.chain[name] == "" {
		self.chain[name] = value
	}
}
func (self *scopeRunner) accessScope(current string) *scope {
	if self.chain[current] == topLevel {
		return self.top
	}
	breadCrumbs := self.trail(current)
	return self.follow(breadCrumbs)
}
func (self *scopeRunner) trail(start string) []string {
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
func (self *scopeRunner) follow(trail []string) *scope {
	var accessed = self.top

	for x := len(trail) - 1; x >= 0; x-- {
		accessed = accessed.children[trail[x]]
	}
	return accessed
}

func (self *scopeRunner) RegisterReset(action func()) {
	parentAction := self.link(action)
	parent := self.accessScope(parentAction)
	parent.registerReset(action)
}

func (self *scopeRunner) Run() {
	for !self.top.visited() {
		self.top.visit()
	}
}

type scopeRunner struct {
	top         *scope
	chain       map[string]string
	currentTest GoTest
}

func NewScopeRunner() *scopeRunner {
	fmt.Sprintf("")

	self := scopeRunner{}
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
	caller_id, _, _, _ := runtime.Caller(5)
	return runtime.FuncForPC(caller_id).Name()
}

func last(group []string) string {
	return group[len(group)-1]
}

var topLevel = "TOP"
