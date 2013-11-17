package system

import (
	"fmt"
	"strings"
)

type FakeShell struct {
	outputByCommand map[string]string // name + args: output
	errorsByOutput  map[string]error  // output: err
	environment     map[string]string
}

func (self *FakeShell) Register(fullCommand string, output string, err error) {
	self.outputByCommand[fullCommand] = output
	self.errorsByOutput[output] = err
}

func (self *FakeShell) ChangeDirectory(directory string) error {
	if self.Getenv("__deleted__") == directory {
		return fmt.Errorf("Directory does not exist: %s", directory)
	}
	return self.Setenv("cwd", directory)
}

func (self *FakeShell) RemoveDirectory(directory string) {
	self.Setenv("__deleted__", directory)
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

func (self *FakeShell) Getenv(key string) string {
	return self.environment[key]
}

func (self *FakeShell) Setenv(key, value string) error {
	self.environment[key] = value
	return nil
}

func NewFakeShell() *FakeShell {
	self := &FakeShell{}
	self.outputByCommand = map[string]string{}
	self.errorsByOutput = map[string]error{}
	self.environment = map[string]string{}
	return self
}
