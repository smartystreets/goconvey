package execution

import (
	"fmt"

	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/reporting"
)

type Runner interface {
	Begin(entry *Registration)
	Register(entry *Registration)
	RegisterReset(action *Action)
	UpgradeReporter(out reporting.Reporter)
	Run()
}

func (self *runner) Begin(entry *Registration) {
	self.ensureStoryCanBegin()
	self.out.BeginStory(reporting.NewStoryReport(entry.Test))
	self.Register(entry)
}
func (self *runner) ensureStoryCanBegin() {
	if self.awaitingNewStory {
		self.awaitingNewStory = false
	} else {
		panic(fmt.Sprintf("%s (See %s)", ExtraGoTest, gotest.FormatExternalFileAndLine()))
	}
}

func (self *runner) Register(entry *Registration) {
	self.ensureStoryAlreadyStarted()
	parentAction := self.link(entry.Action)
	parent := self.accessScope(parentAction)
	child := newScope(entry, self.out)
	parent.adopt(child)
}
func (self *runner) ensureStoryAlreadyStarted() {
	if self.awaitingNewStory {
		panic(MissingGoTest)
	}
}
func (self *runner) link(action *Action) string {
	_, _, parentAction := gotest.ResolveExternalCaller()
	childAction := action.name
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
	self.awaitingNewStory = true
}

type runner struct {
	top   *scope
	chain map[string]string
	out   reporting.Reporter

	awaitingNewStory bool
}

func NewRunner() *runner {
	self := runner{}
	self.out = NewNilReporter()
	self.top = newScope(NewRegistration(topLevel, NewAction(func() {}), nil), self.out)
	self.chain = make(map[string]string)
	self.awaitingNewStory = true
	return &self
}

func (self *runner) UpgradeReporter(out reporting.Reporter) {
	self.out = out
}

const topLevel = "TOP"
const MissingGoTest = `Top-level calls to Convey(...) need a reference to the *testing.T. 
    Hint: Convey("description here", t, func() { /* notice that the second argument was the *testing.T (t)! */ }) `
const ExtraGoTest = `Only the top-level call to Convey(...) needs a reference to the *testing.T.`
