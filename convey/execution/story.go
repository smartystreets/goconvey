package execution

func (self *story) BeginStory(test GoTest) {
	self.test = test
}

func (self *story) Enter(title, id string) {
	self.out.indent()
	self.currentId = id

	if _, found := self.titlesById[id]; !found {
		self.out.println("")
		self.out.print("- " + title)
		self.out.insert(" ")
		self.titlesById[id] = title
	}
}

func (self *story) Report(r *Report) {
	if r.Error != nil {
		self.out.insert(error_)
		self.errors++
		self.out.println("")
		self.out.indent()
		self.out.indent()
		self.out.print("* %s \n* Line: %d: %v \n* %s", r.File, r.Line, r.Error, r.stackTrace)
		self.out.dedent()
		self.out.dedent()
		self.test.Fail()
	} else if r.Failure != "" {
		self.out.insert(failure)
		self.failures++
		self.out.println("")
		self.out.indent()
		self.out.indent()
		self.out.print("* %s \n* Line %d: %s", r.File, r.Line, r.Failure)
		self.out.dedent()
		self.out.dedent()
		self.test.Fail()
	} else {
		self.out.insert(success)
		self.successes++
	}
}

func (self *story) Exit() {
	self.out.dedent()
}

func (self *story) EndStory() {
	self.currentId = ""
	self.titlesById = make(map[string]string)
	self.out.println("\n")
	self.out.println("Passed: %d | Failed: %d | Errors: %d\n", self.successes, self.failures, self.errors)
}

func NewStoryReporter(out *printer) *story {
	self := story{}
	self.out = out
	self.titlesById = make(map[string]string)
	return &self
}

type story struct {
	successes  int
	failures   int
	errors     int
	inner      Reporter
	out        *printer
	titlesById map[string]string
	currentId  string
	test       GoTest
}

const success = "✓"
const failure = "✗"
const error_ = "🔥"
