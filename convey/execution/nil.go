package execution

type nilReporter struct{}

func (self *nilReporter) Success(scope string)                {}
func (self *nilReporter) Failure(scope string, problem error) {}
func (self *nilReporter) Error(scope string, problem error)   {}
func (self *nilReporter) End(scope string)                    {}
