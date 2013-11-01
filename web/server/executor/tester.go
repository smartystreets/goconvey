package executor

import (
	"fmt"
	"github.com/smartystreets/goconvey/web/server/contract"
	"sync"
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
	self.batchSize = 1
	return self
}

///////////

type concurrentCoordinator struct {
	folders   []string
	results   []string
	batchSize int
	shell     contract.Shell
	waiter    sync.WaitGroup
	queue     chan string
}

func (self *concurrentCoordinator) ExecuteConcurrently() []string {
	self.enlistWorkers()
	self.scheduleTasks()
	self.awaitCompletion()
	return self.results
}

func (self *concurrentCoordinator) enlistWorkers() {
	for i := 0; i < self.batchSize; i++ {
		self.waiter.Add(1)
		go self.worker(i)
	}
}
func (self *concurrentCoordinator) worker(id int) {
	for folder := range self.queue {
		output, _ := self.shell.Execute("go", "test", "-v", "-timeout=-42s", folder) // TODO: err
		self.results[id] = output
	}
	self.waiter.Done()
}

func (self *concurrentCoordinator) scheduleTasks() {
	for _, folder := range self.folders {
		self.queue <- folder
	}
}

func (self *concurrentCoordinator) awaitCompletion() {
	close(self.queue)
	self.waiter.Wait()
}

func newCuncurrentCoordinator(folders []string, batchSize int, shell contract.Shell) *concurrentCoordinator {
	self := &concurrentCoordinator{}
	self.results = make([]string, len(folders))
	self.queue = make(chan string)
	self.folders = folders
	self.batchSize = batchSize
	self.shell = shell
	return self
}

var _ = fmt.Sprintf("Hi")
