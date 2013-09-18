package execution

import (
	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/reporting"
)

type Runner interface {
	Begin(test gotest.T, situation string, action *Action)
	Register(situation string, action *Action)
	RegisterReset(action *Action)
	UpgradeReporter(out reporting.Reporter)
	Run()
}

func (self *runner) Begin(test gotest.T, situation string, action *Action) {
	self.out.BeginStory(test)
	self.Register(situation, action)
}

func (self *runner) Register(situation string, action *Action) {
	parentAction := self.link(action)
	parent := self.accessScope(parentAction)
	child := newScope(situation, action, self.out)
	parent.adopt(child)
}
func (self *runner) link(action *Action) string {
	parentAction := resolveExternalCaller()
	childAction := action.Name
	self.linkTo(topLevel, parentAction)
	self.linkTo(parentAction, childAction)
	return parentAction
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

func (self *runner) RegisterReset(action *Action) {
	parentAction := self.link(action)
	parent := self.accessScope(parentAction)
	parent.registerReset(action)
}

func (self *runner) Run() {
	for !self.top.visited() {
		self.top.visit()
	}
	self.out.EndStory()
}

type runner struct {
	top   *scope
	chain map[string]string
	out   reporting.Reporter
}

func NewRunner() *runner {
	self := runner{}
	self.out = NewNilReporter()
	self.top = newScope(topLevel, NewAction(func() {}), self.out)
	self.chain = make(map[string]string)
	return &self
}

func (self *runner) UpgradeReporter(out reporting.Reporter) {
	self.out = out
}

const topLevel = "TOP"
