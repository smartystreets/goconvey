package execution

type nilReporter struct{}

func (self *nilReporter) BeginStory(test GoTest) {}
func (self *nilReporter) Enter(title, id string) {}
func (self *nilReporter) Report(r *Report)       {}
func (self *nilReporter) Exit()                  {}
func (self *nilReporter) EndStory()              {}
