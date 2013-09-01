package execution

func (self *story) Enter(scope string) {

}

func (self *story) Report(r *Report) {

}

func (self *story) Exit() {

}

func NewStoryReporter() *story {
	self := story{}
	return &self
}

type story struct {
	inner Reporter
}
