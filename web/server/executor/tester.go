package executor

import (
	"github.com/smartystreets/goconvey/web/server/contract"
)

type ConcurrentTester struct {
	shell contract.Shell
}

func (self *ConcurrentTester) TestAll(folders []string) (output []string) {
	for _, folder := range folders {
		self.shell.Execute("go", "test", "-i", folder)
	}

	for _, folder := range folders {
		o, _ := self.shell.Execute("go", "test", "-v", "-timeout=-42s", folder) // TODO: err
		output = append(output, o)
	}
	return
}

func NewConcurrentTester(shell contract.Shell) *ConcurrentTester {
	self := &ConcurrentTester{}
	self.shell = shell
	return self
}
