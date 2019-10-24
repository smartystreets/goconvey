package executor

import (
	"errors"
	"fmt"
	"log"
	"strings"
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

func (c *concurrentCoordinator) ExecuteConcurrently() {
	c.enlistWorkers()
	c.scheduleTasks()
	c.awaitCompletion()
}

func (c *concurrentCoordinator) enlistWorkers() {
	for i := 0; i < c.batchSize; i++ {
		c.waiter.Add(1)
		go c.worker(i)
	}
}
func (c *concurrentCoordinator) worker(id int) {
	for folder := range c.queue {
		packageName := strings.Replace(folder.Name, "\\", "/", -1)
		if !folder.Active() {
			log.Printf("Skipping concurrent execution: %s\n", packageName)
			continue
		}

		if folder.HasImportCycle {
			message := fmt.Sprintf("can't load package: import cycle not allowed\npackage %s\n\timports %s", packageName, packageName)
			log.Println(message)
			folder.Output, folder.Error = message, errors.New(message)
		} else {
			log.Printf("Executing concurrent tests: %s\n", packageName)
			folder.Output, folder.Error = c.shell.GoTest(folder.Path, packageName, folder.BuildTags, folder.TestArguments)
		}
	}
	c.waiter.Done()
}

func (c *concurrentCoordinator) scheduleTasks() {
	for _, folder := range c.folders {
		c.queue <- folder
	}
}

func (c *concurrentCoordinator) awaitCompletion() {
	close(c.queue)
	c.waiter.Wait()
}

func newConcurrentCoordinator(folders []*contract.Package, batchSize int, shell contract.Shell) *concurrentCoordinator {
	self := new(concurrentCoordinator)
	self.queue = make(chan *contract.Package)
	self.folders = folders
	self.batchSize = batchSize
	self.shell = shell
	return self
}
