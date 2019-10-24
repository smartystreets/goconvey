package executor

import (
	"log"
	"time"

	"github.com/smartystreets/goconvey/web/server/contract"
)

const (
	Idle      = "idle"
	Executing = "executing"
)

type Executor struct {
	tester     Tester
	parser     Parser
	status     string
	statusChan chan chan string
	statusFlag bool
}

func (e *Executor) Status() string {
	return e.status
}

func (e *Executor) ClearStatusFlag() bool {
	hasNewStatus := e.statusFlag
	e.statusFlag = false
	return hasNewStatus
}

func (e *Executor) ExecuteTests(folders []*contract.Package) *contract.CompleteOutput {
	defer func() { e.setStatus(Idle) }()
	e.execute(folders)
	result := e.parse(folders)
	return result
}

func (e *Executor) execute(folders []*contract.Package) {
	e.setStatus(Executing)
	e.tester.TestAll(folders)
}

func (e *Executor) parse(folders []*contract.Package) *contract.CompleteOutput {
	result := &contract.CompleteOutput{Revision: now().String()}
	e.parser.Parse(folders)
	for _, folder := range folders {
		result.Packages = append(result.Packages, folder.Result)
	}
	return result
}

func (e *Executor) setStatus(status string) {
	e.status = status
	e.statusFlag = true

Loop:
	for {
		select {
		case c := <-e.statusChan:
			e.statusFlag = false
			c <- status
		default:
			break Loop
		}
	}

	log.Printf("Executor status: '%s'\n", e.status)
}

func NewExecutor(tester Tester, parser Parser, ch chan chan string) *Executor {
	return &Executor{
		tester:     tester,
		parser:     parser,
		status:     Idle,
		statusChan: ch,
		statusFlag: false,
	}
}

var now = func() time.Time {
	return time.Now()
}
