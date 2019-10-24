package reporting

import (
	"fmt"
	"sync"
)

func (s *statistics) BeginStory(story *StoryReport) {}

func (s *statistics) Enter(scope *ScopeReport) {}

func (s *statistics) Report(report *AssertionResult) {
	s.Lock()
	defer s.Unlock()

	if !s.failing && report.Failure != "" {
		s.failing = true
	}
	if !s.erroring && report.Error != nil {
		s.erroring = true
	}
	if report.Skipped {
		s.skipped += 1
	} else {
		s.total++
	}
}

func (s *statistics) Exit() {}

func (s *statistics) EndStory() {
	s.Lock()
	defer s.Unlock()

	if !s.suppressed {
		s.printSummaryLocked()
	}
}

func (s *statistics) Suppress() {
	s.Lock()
	defer s.Unlock()
	s.suppressed = true
}

func (s *statistics) PrintSummary() {
	s.Lock()
	defer s.Unlock()
	s.printSummaryLocked()
}

func (s *statistics) printSummaryLocked() {
	s.reportAssertionsLocked()
	s.reportSkippedSectionsLocked()
	s.completeReportLocked()
}
func (s *statistics) reportAssertionsLocked() {
	s.decideColorLocked()
	s.out.Print("\n%d total %s", s.total, plural("assertion", s.total))
}
func (s *statistics) decideColorLocked() {
	if s.failing && !s.erroring {
		fmt.Print(yellowColor)
	} else if s.erroring {
		fmt.Print(redColor)
	} else {
		fmt.Print(greenColor)
	}
}
func (s *statistics) reportSkippedSectionsLocked() {
	if s.skipped > 0 {
		fmt.Print(yellowColor)
		s.out.Print(" (one or more sections skipped)")
	}
}
func (s *statistics) completeReportLocked() {
	fmt.Print(resetColor)
	s.out.Print("\n")
	s.out.Print("\n")
}

func (s *statistics) Write(content []byte) (written int, err error) {
	return len(content), nil // no-op
}

func NewStatisticsReporter(out *Printer) *statistics {
	self := statistics{}
	self.out = out
	return &self
}

type statistics struct {
	sync.Mutex

	out        *Printer
	total      int
	failing    bool
	erroring   bool
	skipped    int
	suppressed bool
}

func plural(word string, count int) string {
	if count == 1 {
		return word
	}
	return word + "s"
}
