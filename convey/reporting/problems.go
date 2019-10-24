package reporting

import "fmt"

type problem struct {
	silent   bool
	out      *Printer
	errors   []*AssertionResult
	failures []*AssertionResult
}

func (p *problem) BeginStory(story *StoryReport) {}

func (p *problem) Enter(scope *ScopeReport) {}

func (p *problem) Report(report *AssertionResult) {
	if report.Error != nil {
		p.errors = append(p.errors, report)
	} else if report.Failure != "" {
		p.failures = append(p.failures, report)
	}
}

func (p *problem) Exit() {}

func (p *problem) EndStory() {
	p.show(p.showErrors, redColor)
	p.show(p.showFailures, yellowColor)
	p.prepareForNextStory()
}
func (p *problem) show(display func(), color string) {
	if !p.silent {
		fmt.Print(color)
	}
	display()
	if !p.silent {
		fmt.Print(resetColor)
	}
	p.out.Dedent()
}
func (p *problem) showErrors() {
	for i, e := range p.errors {
		if i == 0 {
			p.out.Println("\nErrors:\n")
			p.out.Indent()
		}
		p.out.Println(errorTemplate, e.File, e.Line, e.Error, e.StackTrace)
	}
}
func (p *problem) showFailures() {
	for i, f := range p.failures {
		if i == 0 {
			p.out.Println("\nFailures:\n")
			p.out.Indent()
		}
		p.out.Println(failureTemplate, f.File, f.Line, f.Failure, f.StackTrace)
	}
}

func (p *problem) Write(content []byte) (written int, err error) {
	return len(content), nil // no-op
}

func NewProblemReporter(out *Printer) *problem {
	self := new(problem)
	self.out = out
	self.prepareForNextStory()
	return self
}

func NewSilentProblemReporter(out *Printer) *problem {
	self := NewProblemReporter(out)
	self.silent = true
	return self
}

func (p *problem) prepareForNextStory() {
	p.errors = []*AssertionResult{}
	p.failures = []*AssertionResult{}
}
