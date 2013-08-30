package execution

func (self *dot) Success(scope string) {

}

func (self *dot) Failure(scope string, problem error) {

}

func (self *dot) Error(scope string, problem error) {

}

func (self *dot) End(scope string) {

}

func NewDotReporter() *dot {
	self := dot{}
	return &self
}

type dot struct {
}
