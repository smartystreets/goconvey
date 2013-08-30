package execution

var SpecRunner runner
var SpecReporter Reporter

func init() {
	SpecRunner = NewScopeRunner()
	SpecReporter = NewStatisticsReporter()
	SpecRunner.UpgradeReporter(SpecReporter)
}

type runner interface {
	Begin(test GoTest, situation string, action func())
	Register(situation string, action func())
	RegisterReset(action func())
	Run()
	UpgradeReporter(out Reporter)
}

type Reporter interface {
	Enter(scope string)
	Success(r Report)
	Failure(r Report)
	Error(r Report)
	Exit()
}

type Report struct {
	File    string
	Line    int
	Failure string
	Error   error
}

type GoTest interface {
	Fail()
}
