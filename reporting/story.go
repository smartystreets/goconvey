package reporting

import "github.com/smartystreets/goconvey/gotest"
import "github.com/smartystreets/goconvey/printing"

func (self *story) BeginStory(test gotest.T) {
	self.test = test
}

func (self *story) Enter(title, id string) {
	self.out.Indent()
	self.currentId = id

	if _, found := self.titlesById[id]; !found {
		self.out.Println("")
		self.out.Print("- " + title)
		self.out.Insert(" ")
		self.titlesById[id] = title
	}
}

func (self *story) Report(r *Report) {
	if r.Error != nil {
		self.out.Insert(error_)
		self.errors++
		self.out.Println("")
		self.out.Indent()
		self.out.Indent()
		self.out.Print("* %s \n* Line: %d - %v \n%s", r.File, r.Line, r.Error, r.stackTrace)
		self.out.Dedent()
		self.out.Dedent()
		self.test.Fail()
	} else if r.Failure != "" {
		self.out.Insert(failure)
		self.failures++
		self.out.Println("")
		self.out.Indent()
		self.out.Indent()
		self.out.Print("* %s \n* Line %d: %s", r.File, r.Line, r.Failure)
		self.out.Dedent()
		self.out.Dedent()
		self.test.Fail()
	} else {
		self.out.Insert(success)
		self.successes++
	}
}

func (self *story) Exit() {
	self.out.Dedent()
}

func (self *story) EndStory() {
	self.currentId = ""
	self.titlesById = make(map[string]string)
	self.out.Println("\n")
	self.out.Println("Passed: %d | Failed: %d | Errors: %d\n", self.successes, self.failures, self.errors)
}

func NewStoryReporter(out *printing.Printer) *story {
	self := story{}
	self.out = out
	self.titlesById = make(map[string]string)
	return &self
}

type story struct {
	successes  int
	failures   int
	errors     int
	out        *printing.Printer
	titlesById map[string]string
	currentId  string
	test       gotest.T
}

const success = "✓"
const failure = "✗"
const error_ = "🔥"
