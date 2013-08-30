package execution

import (
	"fmt"
)

func (self *scopeRunner) Begin(test GoTest, situation string, action func()) {
	self.currentTest = test
	/*child := */ self.registerChildScope(situation, action)
	// child.aggregate = true
}

func (parent *scopeRunner) Register(situation string, action func()) {
	parent.registerChildScope(situation, action)
}
func (self *scopeRunner) registerChildScope(situation string, action func()) *scope {
	parentAction := self.link(action)
	parent := self.accessScope(parentAction)
	child := newScope(situation, action)
	parent.adopt(child)
	return child
}
func (self *scopeRunner) link(action func()) string {
	parentAction := resolveExternalCaller()
	childAction := functionName(action)
	self.linkTo(topLevel, parentAction)
	self.linkTo(parentAction, childAction)
	return parentAction
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
	out         reporter
}

type reporter interface {
	Success(scope string)
	Failure(scope string, problem error)
	Error(scope string, problem error)
	End(scope string)
}

func NewScopeRunner() *scopeRunner {
	fmt.Sprintf("")

	self := scopeRunner{}
	self.top = newScope(topLevel, func() {})
	self.chain = make(map[string]string)
	return &self
}

const topLevel = "TOP"
