package executor

import (
	"github.com/smartystreets/goconvey/web/server/contract"
	"sync"
)

type concurrentCoordinator struct {
	folders   []string
	results   []string
	errors    []error
	batchSize int
	shell     contract.Shell
	waiter    sync.WaitGroup
	queue     chan string
}

func (self *concurrentCoordinator) ExecuteConcurrently() []string {
	self.enlistWorkers()
	self.scheduleTasks()
	self.awaitCompletion()
	self.checkForErrors()
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
		output, err := self.shell.Execute("go", "test", "-v", "-timeout=-42s", folder)
		self.errors[id] = err
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

func (self *concurrentCoordinator) checkForErrors() {
	for _, err := range self.errors {
		if err != nil {
			panic(err)
		}
	}
}

func newCuncurrentCoordinator(folders []string, batchSize int, shell contract.Shell) *concurrentCoordinator {
	self := &concurrentCoordinator{}
	self.results = make([]string, len(folders))
	self.errors = make([]error, len(folders))
	self.queue = make(chan string)
	self.folders = folders
	self.batchSize = batchSize
	self.shell = shell
	return self
}
