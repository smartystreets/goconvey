package execution

var SpecRunner runner
var SpecReporter Reporter

func init() {
	console := newConsole()
	printer := newPrinter(console)
	SpecReporter = NewStoryReporter(printer)
	SpecRunner = NewScopeRunner()
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
	BeginStory(test GoTest)
	Enter(title, id string)
	Report(r *Report)
	Exit()
	EndStory()
}

type Report struct {
	File       string
	Line       int
	Failure    string
	Error      interface{}
	stackTrace string
}

func NewFailureReport(failure string) *Report {
	file, line := caller()
	stack := stackTrace()
	report := Report{file, line, failure, nil, stack}
	return &report
}
func NewErrorReport(err interface{}) *Report {
	file, line := caller()
	stack := fullStackTrace()
	report := Report{file, line, "", err, stack}
	return &report
}
func NewSuccessReport() *Report {
	file, line := caller()
	stack := stackTrace()
	report := Report{file, line, "", nil, stack}
	return &report
}

type GoTest interface {
	Fail()
}
