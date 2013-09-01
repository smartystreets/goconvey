package execution

type nilReporter struct{}

func (self *nilReporter) Enter(scope string) {}
func (self *nilReporter) Report(r *Report)   {}
func (self *nilReporter) Exit()              {}
