package reporting

import "fmt"

import (
	"github.com/smartystreets/goconvey/gotest"
	"github.com/smartystreets/goconvey/printing"
)

func (self *story) BeginStory(test gotest.T) {}

func (self *story) Enter(title, id string) {
	self.out.Indent()
	self.currentId = id

	if _, found := self.titlesById[id]; !found {
		self.out.Println("")
		self.out.Print(title)
		self.out.Insert(" ")
		self.titlesById[id] = title
	}
}

func (self *story) Report(r *Report) {
	if r.Error != nil {
		fmt.Print(redColor)
		self.out.Insert(error_)
	} else if r.Failure != "" {
		fmt.Print(yellowColor)
		self.out.Insert(failure)
	} else {
		fmt.Print(greenColor)
		self.out.Insert(success)
	}
	fmt.Print(resetColor)
}

func (self *story) Exit() {
	self.out.Dedent()
}

func (self *story) EndStory() {
	self.currentId = ""
	self.titlesById = make(map[string]string)
}

func NewStoryReporter(out *printing.Printer) *story {
	self := story{}
	self.out = out
	self.titlesById = make(map[string]string)
	return &self
}

type story struct {
	out        *printing.Printer
	titlesById map[string]string
	currentId  string
}
