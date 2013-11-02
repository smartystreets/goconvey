package executor

import (
	"github.com/smartystreets/goconvey/web/server/contract"
	"github.com/smartystreets/goconvey/web/server/parser"
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

func (self *Executor) ExecuteTests(folders []string) *parser.CompleteOutput {
	defer func() { self.status = Idle }()
	output := self.execute(folders)
	result := self.parse(output, folders)
	return result
}

func (self *Executor) execute(folders []string) []string {
	self.status = Executing
	return self.tester.TestAll(folders)
}

func (self *Executor) parse(outputs, folders []string) *parser.CompleteOutput {
	self.status = Parsing
	result := &parser.CompleteOutput{Revision: now().String()}
	for i, output := range outputs {
		packageName := contract.ResolvePackageName(folders[i])
		parsed := self.parser.Parse(packageName, output)
		result.Packages = append(result.Packages, parsed)
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
