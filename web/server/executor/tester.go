package executor

import (
	"log"
	"strings"

	"github.com/smartystreets/goconvey/web/server/contract"
)

type ConcurrentTester struct {
	shell     contract.Shell
	batchSize int
}

func (self *ConcurrentTester) SetBatchSize(batchSize int) {
	self.batchSize = batchSize
	log.Printf("Now configured to test %d packages concurrently.\n", self.batchSize)
}

func (self *ConcurrentTester) TestAll(folders []*contract.Package) {
	if self.batchSize == 1 {
		self.executeSynchronously(folders)
	} else {
		newConcurrentCoordinator(folders, self.batchSize, self.shell).ExecuteConcurrently()
	}
	return
}

func (self *ConcurrentTester) KillRunningTests() {
	log.Print("Killing executing tests")
	self.shell.AbortGoTest()
}

func (self *ConcurrentTester) executeSynchronously(folders []*contract.Package) {
	for _, folder := range folders {
		packageName := strings.Replace(folder.Name, "\\", "/", -1)
		if !folder.Active() {
			log.Printf("Skipping execution: %s\n", packageName)
			continue
		}
		log.Printf("Executing tests: %s\n", packageName)
		folder.Output, folder.Error = self.shell.GoTest(folder.Path, packageName, folder.TestArguments)
	}
}

func NewConcurrentTester(shell contract.Shell) *ConcurrentTester {
	self := new(ConcurrentTester)
	self.shell = shell
	self.batchSize = defaultBatchSize
	return self
}

const defaultBatchSize = 10
