package executor

import (
	"fmt"
	"github.com/smartystreets/goconvey/web/server/contract"
)

type ConcurrentTester struct {
	shell     contract.Shell
	batchSize int
}

func (self *ConcurrentTester) SetBatchSize(batchSize int) {
	self.batchSize = batchSize
}

func (self *ConcurrentTester) TestAll(folders []string) (output []string) {
	for _, folder := range folders {
		self.shell.Execute("go", "test", "-i", folder)
	}

	if self.batchSize == 1 {
		output = self.executeSynchronously(folders)
	} else {
		output = newCuncurrentCoordinator(folders, self.batchSize, self.shell).ExecuteConcurrently()
	}
	return
}

func (self *ConcurrentTester) executeSynchronously(folders []string) []string {
	all := make([]string, len(folders))
	for i, folder := range folders {
		all[i], _ = self.shell.Execute("go", "test", "-v", "-timeout=-42s", folder) // TODO: err
	}
	return all
}

func NewConcurrentTester(shell contract.Shell) *ConcurrentTester {
	self := &ConcurrentTester{}
	self.shell = shell
	self.batchSize = defaultBatchSize
	return self
}

const defaultBatchSize = 4

var _ = fmt.Sprintf("Hi")
