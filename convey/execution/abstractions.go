package execution

var SpecRunner runner
var SpecReporter Reporter

func init() {
	SpecRunner = NewScopeRunner()
	SpecReporter = NewStatisticsReporter() // TODO: package with dot or story reporter
}

type runner interface {
	Begin(test GoTest, situation string, action func())
	Register(situation string, action func())
	RegisterReset(action func())
	Run()
	UpgradeReporter(out Reporter)
}

type Reporter interface {
	Success(scope string)
	Failure(scope string, problem error)
	Error(scope string, problem error)
	End(scope string)
}

type GoTest interface {
	Fail()
}
