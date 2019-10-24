package reporting

import "fmt"

type dot struct{ out *Printer }

func (d *dot) BeginStory(story *StoryReport) {}

func (d *dot) Enter(scope *ScopeReport) {}

func (d *dot) Report(report *AssertionResult) {
	if report.Error != nil {
		fmt.Print(redColor)
		d.out.Insert(dotError)
	} else if report.Failure != "" {
		fmt.Print(yellowColor)
		d.out.Insert(dotFailure)
	} else if report.Skipped {
		fmt.Print(yellowColor)
		d.out.Insert(dotSkip)
	} else {
		fmt.Print(greenColor)
		d.out.Insert(dotSuccess)
	}
	fmt.Print(resetColor)
}

func (d *dot) Exit() {}

func (d *dot) EndStory() {}

func (d *dot) Write(content []byte) (written int, err error) {
	return len(content), nil // no-op
}

func NewDotReporter(out *Printer) *dot {
	self := new(dot)
	self.out = out
	return self
}
