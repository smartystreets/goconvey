package execution

func (self *dot) BeginStory(test GoTest) {

}

func (self *dot) Enter(title, id string) {

}

func (self *dot) Report(r *Report) {

}

func (self *dot) Exit() {

}

func (self *dot) EndStory() {

}

func NewDotReporter() *dot {
	self := dot{}
	return &self
}

type dot struct {
	inner Reporter
}
