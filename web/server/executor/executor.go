package executor

import (
	"github.com/smartystreets/goconvey/web/server/contract"
	"time"
)

const (
	Idle      = "idle"
	Executing = "executing"
	Parsing   = "parsing"
)

type Executor struct {
	tester Tester
	parser Parser
	status string
}

func (self *Executor) Status() string {
	return self.status
}

func (self *Executor) ExecuteTests(folders []*contract.Package) *contract.CompleteOutput {
	defer func() { self.status = Idle }()
	self.execute(folders)
	result := self.parse(folders)
	return result
}

func (self *Executor) execute(folders []*contract.Package) {
	self.status = Executing
	self.tester.TestAll(folders)
}

func (self *Executor) parse(folders []*contract.Package) *contract.CompleteOutput {
	self.status = Parsing
	result := &contract.CompleteOutput{Revision: now().String()}
	self.parser.Parse(folders)
	for _, folder := range folders {
		result.Packages = append(result.Packages, folder.Result)
	}
	return result
}

func NewExecutor(tester Tester, parser Parser) *Executor {
	self := &Executor{}
	self.tester = tester
	self.parser = parser
	self.status = "idle"
	return self
}

var now = func() time.Time {
	return time.Now()
}
