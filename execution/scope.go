package execution

import (
	"fmt"
	"github.com/smartystreets/goconvey/reporting"
	"strings"
)

func (parent *scope) adopt(child *scope) {
	if parent.hasChild(child) {
		return
	}
	parent.birthOrder = append(parent.birthOrder, child.name)
	parent.children[child.name] = child
}
func (parent *scope) hasChild(child *scope) bool {
	for _, name := range parent.birthOrder {
		if name == child.name {
			return true
		}
	}
	return false
}

func (self *scope) registerReset(action *Action) {
	self.resets[action.name] = action
}

func (self *scope) visited() bool {
	return self.panicked || self.child >= len(self.birthOrder)
}

func (parent *scope) visit() {
	defer parent.exit()
	parent.enter()
	parent.action.Invoke()
	parent.visitChildren()
}
func (parent *scope) enter() {
	parent.reporter.Enter(parent.report)
}
func (parent *scope) visitChildren() {
	if len(parent.children) == 0 {
		parent.cleanup()
	} else {
		parent.visitChild()
	}
}
func (parent *scope) visitChild() {
	child := parent.children[parent.birthOrder[parent.child]]
	child.visit()
	if child.visited() {
		parent.cleanup()
		parent.child++
	}
}
func (parent *scope) cleanup() {
	for _, reset := range parent.resets {
		reset.Invoke()
	}
}
func (parent *scope) exit() {
	if problem := recover(); problem != nil {
		if strings.HasPrefix(fmt.Sprintf("%v", problem), ExtraGoTest) {
			panic(problem)
		}
		parent.panicked = true
		parent.reporter.Report(reporting.NewErrorReport(problem))
	}
	parent.reporter.Exit()
}

// func newScope(title string, action *Action, reporter reporting.Reporter) *scope {
func newScope(entry *Registration, reporter reporting.Reporter) *scope {
	self := &scope{}
	self.reporter = reporter
	self.name = entry.Action.name
	self.title = entry.Situation
	self.action = entry.Action
	self.children = make(map[string]*scope)
	self.birthOrder = []string{}
	self.resets = make(map[string]*Action)
	self.report = reporting.NewScopeReport(self.title, self.name)
	return self
}

type scope struct {
	name       string
	title      string
	action     *Action
	children   map[string]*scope
	birthOrder []string
	child      int
	resets     map[string]*Action
	panicked   bool
	reporter   reporting.Reporter
	report     *reporting.ScopeReport
}
