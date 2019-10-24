// TODO: in order for this reporter to be completely honest
// we need to retrofit to be more like the json reporter such that:
// 1. it maintains ScopeResult collections, which count assertions
// 2. it reports only after EndStory(), so that all tick marks
//    are placed near the appropriate title.
// 3. Under unit test

package reporting

import (
	"fmt"
	"strings"
)

type story struct {
	out        *Printer
	titlesById map[string]string
	currentKey []string
}

func (s *story) BeginStory(story *StoryReport) {}

func (s *story) Enter(scope *ScopeReport) {
	s.out.Indent()

	s.currentKey = append(s.currentKey, scope.Title)
	ID := strings.Join(s.currentKey, "|")

	if _, found := s.titlesById[ID]; !found {
		s.out.Println("")
		s.out.Print(scope.Title)
		s.out.Insert(" ")
		s.titlesById[ID] = scope.Title
	}
}

func (s *story) Report(report *AssertionResult) {
	if report.Error != nil {
		fmt.Print(redColor)
		s.out.Insert(error_)
	} else if report.Failure != "" {
		fmt.Print(yellowColor)
		s.out.Insert(failure)
	} else if report.Skipped {
		fmt.Print(yellowColor)
		s.out.Insert(skip)
	} else {
		fmt.Print(greenColor)
		s.out.Insert(success)
	}
	fmt.Print(resetColor)
}

func (s *story) Exit() {
	s.out.Dedent()
	s.currentKey = s.currentKey[:len(s.currentKey)-1]
}

func (s *story) EndStory() {
	s.titlesById = make(map[string]string)
	s.out.Println("\n")
}

func (s *story) Write(content []byte) (written int, err error) {
	return len(content), nil // no-op
}

func NewStoryReporter(out *Printer) *story {
	self := new(story)
	self.out = out
	self.titlesById = make(map[string]string)
	return self
}
