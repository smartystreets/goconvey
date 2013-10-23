package system

import (
	"fmt"
	"strings"
)

type FakeShell struct {
	outputByCommand map[string]string // name + args: output
	errorsByOutput  map[string]error  // output: err
}

func (self *FakeShell) Register(fullCommand string, output string, err error) {
	self.outputByCommand[fullCommand] = output
	self.errorsByOutput[output] = err
}

func (self *FakeShell) Execute(name string, args ...string) (output string, err error) {
	fullCommand := name + " " + strings.Join(args, " ")
	var exists bool = false
	if output, exists = self.outputByCommand[fullCommand]; !exists {
		panic(fmt.Sprintf("Missing command output for %s", fullCommand))
	}
	err = self.errorsByOutput[output]
	return
}

func NewFakeShell() *FakeShell {
	self := &FakeShell{}
	self.outputByCommand = map[string]string{}
	self.errorsByOutput = map[string]error{}
	return self
}
