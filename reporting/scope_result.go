package reporting

type ScopeResult struct {
	Title      string
	File       string
	Line       int
	Depth      int
	Assertions []*AssertionResult
}

func newScopeResult(title string, depth int, file string, line int) *ScopeResult {
	self := &ScopeResult{}
	self.Title = title
	self.Depth = depth
	self.File = file
	self.Line = line
	self.Assertions = []*AssertionResult{}
	return self
}
