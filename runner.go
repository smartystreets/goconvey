package goconvey

import (
	"runtime"
	"reflect"
	"fmt"
)

func (self *runner) Convey(situation string, action func()) {
	fmt.Sprintf("")
	// fmt.Println("--------")
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

	// imprint
	parent := self.accessConveyedScope(parentScope)
	// fmt.Println("parent:", parent.subject)
	parent.Convey(NewScope(situation, action))
}

func (self *runner) accessConveyedScope(current string) *scope {
	if self.chain[current] == "TOP" {
		return self.top
	}

	breadCrumbs := []string{current, self.chain[current]}
	for {
		next := self.chain[breadCrumbs[len(breadCrumbs) - 1]]
		if next == "TOP" {
			break
		} else {
			breadCrumbs = append(breadCrumbs, next)
		}
	}
	accessed := self.top.convey[breadCrumbs[len(breadCrumbs) - 2]] // why do I have to start at "- 2"?

	for x := len(breadCrumbs) - 3; x >=0; x-- {
		accessed = accessed.convey[breadCrumbs[x]]
	}
	return accessed
}

func (self *runner) Reset(what string, action func()) {
	childScope := functionName(action)
	parentScope := resolveExternalCaller()

	if self.chain[childScope] == "" {
		self.chain[childScope] = parentScope
	}

	// TODO: self.accessConveyedScope doesn't know how to deal with resets yet (or so I think)... Need to drive this out with a test. (reset at a nested level)
	parent := self.accessConveyedScope(parentScope)
	parent.reset[what] = NewScope(what, action)
}

func (self *runner) Run() {
	for !self.top.Run() {
	}
}

type runner struct {
	top *scope
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