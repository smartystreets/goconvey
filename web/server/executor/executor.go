package executor

import (
	"log"
	"time"

	"github.com/smartystreets/goconvey/web/server/contract"
)

const (
	Idle      = "idle"
	Executing = "executing"
	Parsing   = "parsing"
)

type Executor struct {
	tester       Tester
	parser       Parser
	status       string
	statusUpdate chan bool
}

func (self *Executor) Status() string {
	return self.status
}

func (self *Executor) ExecuteTests(folders []*contract.Package) *contract.CompleteOutput {
	defer func() { self.setStatus(Idle) }()
	self.execute(folders)
	result := self.parse(folders)
	return result
}

func (self *Executor) execute(folders []*contract.Package) {
	self.setStatus(Executing)
	self.tester.TestAll(folders)
}

func (self *Executor) parse(folders []*contract.Package) *contract.CompleteOutput {
	self.setStatus(Parsing)
	result := &contract.CompleteOutput{Revision: now().String()}
	self.parser.Parse(folders)
	for _, folder := range folders {
		result.Packages = append(result.Packages, folder.Result)
	}
	return result
}

func (self *Executor) setStatus(status string) {
	self.status = status

	select {
	case self.statusUpdate <- true:
	default:
	}

	log.Printf("Executor status: '%s'\n", self.status)
}

func NewExecutor(tester Tester, parser Parser, statusUpdate chan bool) *Executor {
	self := new(Executor)
	self.tester = tester
	self.parser = parser
	self.status = Idle
	self.statusUpdate = statusUpdate
	return self
}

var now = func() time.Time {
	return time.Now()
}
