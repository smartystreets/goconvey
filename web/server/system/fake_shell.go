package system

type FakeShell struct {
	environment map[string]string
	executions  []string
}

func (self *FakeShell) GoTest(directory string) (output string, err error) {
	self.executions = append(self.executions, directory)
	output = directory
	return
}

func (self *FakeShell) Executions() []string {
	return self.executions
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
	self.environment = map[string]string{}
	self.executions = []string{}
	return self
}
