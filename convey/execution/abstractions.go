package execution

type runner interface {
	Begin(test GoTest)
	Register(situation string, action func())
	RegisterReset(action func())
	Run()
}

var SpecRunner runner = NewScopeRunner()

type GoTest interface {
	Fail()
}
