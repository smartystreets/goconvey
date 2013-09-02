package execution

func (self *story) Enter(scope string) {
	// self.out.println(scope)
	// self.out.indent()
}

func (self *story) Report(r *Report) {

}

func (self *story) Exit() {
	// self.out.dedent()
}

func NewStoryReporter(out *printer) *story {
	self := story{}
	self.out = out
	return &self
}

type story struct {
	inner Reporter
	out   *printer
}
