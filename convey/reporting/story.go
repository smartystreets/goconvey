package reporting

type StoryReporter struct {
	nestedReporter
	out  *Printer
	seed int64
}

func (s *StoryReporter) Close() {
	o := s.out

	if s.seed != 0 {
		o.Insert(whiteColor)
		o.Suite("Random Seed")
		o.Statement(s.seed)
		o.Exit()
		o.Insert(resetColor)
	}

	s.Walk(func(obj interface{}) {
		switch obj := obj.(type) {
		case *NestedScopeResult:
			o.Suite(obj.Title)

		case ScopeExit:
			o.Exit()

		case string:
			o.Insert(whiteColor)
			o.Statement(obj)
			o.Insert(resetColor)

		case *AssertionResult:
			color, tok := greenColor, success
			if obj.Error != nil {
				color, tok = redColor, error_
			} else if obj.Failure != "" {
				color, tok = yellowColor, failure
			} else if obj.Skipped {
				color, tok = yellowColor, skip
			}
			o.Insert(color)
			o.Expression(tok)
			o.Insert(resetColor)
		}
	})
	o.Insert("\n\n")
}

func NewStoryReporter(out *Printer, seed int64) Reporter {
	return &StoryReporter{out: out, seed: seed}
}
