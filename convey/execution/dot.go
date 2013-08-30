package execution

func (self *dot) Enter(scope string) {

}

func (self *dot) Success(r Report) {

}

func (self *dot) Failure(r Report) {

}

func (self *dot) Error(r Report) {

}

func (self *dot) Exit() {

}

func NewDotReporter() *dot {
	self := dot{}
	return &self
}

type dot struct {
}
