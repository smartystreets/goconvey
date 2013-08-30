package execution

func (self *story) Enter(scope string) {

}

func (self *story) Success(r Report) {

}

func (self *story) Failure(r Report) {

}

func (self *story) Error(r Report) {

}

func (self *story) Exit() {

}

func NewStoryReporter() *story {
	self := story{}
	return &self
}

type story struct {
}
