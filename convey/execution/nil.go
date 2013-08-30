package execution

type nilReporter struct{}

func (self *nilReporter) Enter(scope string) {}
func (self *nilReporter) Success(r Report)   {}
func (self *nilReporter) Failure(r Report)   {}
func (self *nilReporter) Error(r Report)     {}
func (self *nilReporter) Exit()              {}
