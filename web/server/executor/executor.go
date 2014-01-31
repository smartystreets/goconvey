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
	tester      Tester
	parser      Parser
	status      string
	statusNotif chan bool
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
	case self.statusNotif <- true:
	default:
	}

	log.Printf("Executor status: '%s'\n", self.status)
}

func NewExecutor(tester Tester, parser Parser, ch chan bool) *Executor {
	return &Executor{
		tester,
		parser,
		Idle,
		ch,
	}
}

var now = func() time.Time {
	return time.Now()
}
