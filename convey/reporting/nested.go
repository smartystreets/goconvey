package reporting

// nestedReporter buffers scopes from Enter, Exit, Report, and Write so you only
// need to implement the pertinent details of Close().
//
// In particular, it's smart about the stack of scopes so that it will render:
//
//   A
//     B
//       C
//   A
//     B
//       D
//
// As:
//
//   A
//     B
//       C
//       D
type nestedReporter struct {
	root *NestedScopeResult
	cur  *NestedScopeResult
}

func (self *nestedReporter) Walk(cb func(interface{})) {
	self.root.Walk(cb)
}

func (self *nestedReporter) Enter(scope *ScopeReport) {
	if self.root == nil {
		self.root = &NestedScopeResult{}
		self.cur = self.root
	}

	// If the last thing at the current level of the tree is a scope whose Title
	// matches ours, then use that!
	if len(self.cur.Items) > 0 {
		lastItm := self.cur.Items[len(self.cur.Items)-1]
		if s, ok := lastItm.(*NestedScopeResult); ok {
			if s.Title == scope.Title {
				self.cur = s
				return
			}
		}
	}

	// Otherwise make a new one
	newScope := NewNestedScopeResult(self.cur, scope)
	self.cur.Items = append(self.cur.Items, newScope)
	self.cur = newScope
}

func (self *nestedReporter) Report(report *AssertionResult) {
	self.cur.Items = append(self.cur.Items, report)
}

func (self *nestedReporter) Exit() {
	self.cur = self.cur.Parent
	if self.cur == nil { // in case the caller overshoots
		self.cur = self.root
	}
}

func (self *nestedReporter) Write(content []byte) (written int, err error) {
	self.cur.Items = append(self.cur.Items, string(content))
	return len(content), nil
}
