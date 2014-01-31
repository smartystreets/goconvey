package executor

import (
	"fmt"
	"log"
	"sync"

	"github.com/smartystreets/goconvey/web/server/contract"
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
			log.Printf("Skipping concurrent execution: %s\n", folder.Name)
			continue
		}
		log.Printf("Executing concurrent tests: %s\n", folder.Name)
		folder.Output, folder.Error = self.shell.GoTest(folder.Path)
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
		if folder.Error != nil && folder.Output == "" {
			fmt.Println(folder.Path, folder.Error)
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
