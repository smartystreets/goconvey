package execution

func (self *story) Success(scope string) {

}

func (self *story) Failure(scope string, problem error) {

}

func (self *story) Error(scope string, problem error) {

}

func (self *story) End(scope string) {

}

func NewStoryReporter() *story {
	self := story{}
	return &self
}

type story struct {
}
