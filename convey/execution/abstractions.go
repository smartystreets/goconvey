package execution

var SpecRunner runner

func init() {
	SpecRunner = NewScopeRunner()
}

type runner interface {
	Begin(test GoTest)
	Register(situation string, action func())
	RegisterReset(action func())
	Run()
}


type GoTest interface {
	Fail()
}
