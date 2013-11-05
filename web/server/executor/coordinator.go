package executor

import (
	"github.com/smartystreets/goconvey/web/server/contract"
	"sync"
)

type concurrentCoordinator struct {
	batchSize int
	queue     chan *contract.Package
	folders   []*contract.Package
	shell     contract.Shell
	waiter    sync.WaitGroup
}

func (self *concurrentCoordinator) ExecuteConcurrently() {
	self.enlistWorkers()
	self.scheduleTasks()
	self.awaitCompletion()
	self.checkForErrors()
}

func (self *concurrentCoordinator) enlistWorkers() {
	for i := 0; i < self.batchSize; i++ {
		self.waiter.Add(1)
		go self.worker(i)
	}
}
func (self *concurrentCoordinator) worker(id int) {
	for folder := range self.queue {
		if !folder.Active {
			continue
		}
		output, err := self.shell.Execute("go", "test", "-v", "-timeout=-42s", folder.Name)
		folder.Output = output
		folder.Error = err
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
	for _, folder := range self.folders {
		if folder.Error != nil {
			panic(folder.Error)
		}
	}
}

func newCuncurrentCoordinator(folders []*contract.Package, batchSize int, shell contract.Shell) *concurrentCoordinator {
	self := &concurrentCoordinator{}
	self.queue = make(chan *contract.Package)
	self.folders = folders
	self.batchSize = batchSize
	self.shell = shell
	return self
}
