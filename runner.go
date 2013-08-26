package goconvey

import (
	"fmt"
	"reflect"
	"runtime"
)

func (self *runner) Register(situation string, action func()) {
	fmt.Sprintf("")
	childScope := functionName(action)
	parentScope := resolveExternalCaller()

	// take care of this scope
	if self.chain[parentScope] == "" {
		self.chain[parentScope] = "TOP"
	}

	// prepare for the future...
	if self.chain[childScope] == "" {
		self.chain[childScope] = parentScope
	}

	parent := self.accessScope(parentScope)
	child := NewScope(situation, action)
	parent.Adopt(child)
}

func (self *runner) accessScope(current string) *scope {
	if self.chain[current] == "TOP" {
		return self.top
	}

	breadCrumbs := []string{current, self.chain[current]}
	for {
		next := self.chain[breadCrumbs[len(breadCrumbs)-1]]
		if next == "TOP" {
			break
		} else {
			breadCrumbs = append(breadCrumbs, next)
		}
	}
	accessed := self.top.children[breadCrumbs[len(breadCrumbs)-2]] // why do I have to start at "- 2"?

	for x := len(breadCrumbs) - 3; x >= 0; x-- {
		accessed = accessed.children[breadCrumbs[x]]
	}
	return accessed
}

func (self *runner) Run() {
	for !self.top.Visited() {
		self.top.Visit()
	}
}

type runner struct {
	top   *scope
	chain map[string]string
}

func NewRunner() *runner {
	self := runner{}
	self.top = NewScope("TOP", func() {})
	self.chain = make(map[string]string)
	return &self
}

func functionName(function func()) string {
	return runtime.FuncForPC(reflect.ValueOf(function).Pointer()).Name()
}
func resolveExternalCaller() string {
	caller_id, _, _, _ := runtime.Caller(2)
	return runtime.FuncForPC(caller_id).Name()
}
