package executor

import (
	"errors"
	"fmt"
	"log"
	"strings"

	"github.com/smartystreets/goconvey/web/server/contract"
)

type ConcurrentTester struct {
	shell     contract.Shell
	batchSize int
}

func (c *ConcurrentTester) SetBatchSize(batchSize int) {
	c.batchSize = batchSize
	log.Printf("Now configured to test %d packages concurrently.\n", c.batchSize)
}

func (c *ConcurrentTester) TestAll(folders []*contract.Package) {
	if c.batchSize == 1 {
		c.executeSynchronously(folders)
	} else {
		newConcurrentCoordinator(folders, c.batchSize, c.shell).ExecuteConcurrently()
	}
	return
}

func (c *ConcurrentTester) executeSynchronously(folders []*contract.Package) {
	for _, folder := range folders {
		packageName := strings.Replace(folder.Name, "\\", "/", -1)
		if !folder.Active() {
			log.Printf("Skipping execution: %s\n", packageName)
			continue
		}
		if folder.HasImportCycle {
			message := fmt.Sprintf("can't load package: import cycle not allowed\npackage %s\n\timports %s", packageName, packageName)
			log.Println(message)
			folder.Output, folder.Error = message, errors.New(message)
		} else {
			log.Printf("Executing tests: %s\n", packageName)
			folder.Output, folder.Error = c.shell.GoTest(folder.Path, packageName, folder.BuildTags, folder.TestArguments)
		}
	}
}

func NewConcurrentTester(shell contract.Shell) *ConcurrentTester {
	self := new(ConcurrentTester)
	self.shell = shell
	self.batchSize = defaultBatchSize
	return self
}

const defaultBatchSize = 10
