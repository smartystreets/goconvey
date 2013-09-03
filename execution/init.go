package execution

import (
	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/printing"
	"github.com/smartystreets/goconvey/reporting"
)

var SpecRunner runner
var SpecReporter Reporter

func init() {
	console := printing.NewConsole()
	printer := printing.NewPrinter(console)
	SpecReporter = reporting.NewStoryReporter(printer)
	SpecRunner = NewScopeRunner()
	SpecRunner.UpgradeReporter(SpecReporter)
}

type runner interface {
	Begin(test gotest.T, situation string, action func())
	Register(situation string, action func())
	RegisterReset(action func())
	Run()
	UpgradeReporter(out Reporter)
}
