package execution

func (self *dot) Enter(scope string) {

}

func (self *dot) Report(r *Report) {

}

func (self *dot) Exit() {

}

func NewDotReporter() *dot {
	self := dot{}
	return &self
}

type dot struct {
	inner Reporter
}
