package convey

type suite struct {
	Situation string
	Test      t
	Focus     bool
	Func      func() // nil means skipped
	FailMode  FailureMode
}

func (self *suite) ShouldBeTopLevel() bool {
	return self.Test != nil
}

func newSuite(situation string, failureMode FailureMode, f func(), test t) *suite {
	return &suite{
		Situation: situation,
		Test:      test,
		Func:      f,
		FailMode:  failureMode,
	}
}

func discover(items []interface{}) *suite {
	name, items := parseName(items)
	test, items := parseGoTest(items)
	failure, items := parseFailureMode(items)
	action, items := parseAction(items)

	if len(items) != 0 {
		panic(parseError)
	}

	return newSuite(name, failure, action, test)
}
func item(items []interface{}) interface{} {
	if len(items) == 0 {
		panic(parseError)
	}
	return items[0]
}
func parseName(items []interface{}) (string, []interface{}) {
	if name, parsed := item(items).(string); parsed {
		return name, items[1:]
	}
	panic(parseError)
}
func parseGoTest(items []interface{}) (t, []interface{}) {
	if test, parsed := item(items).(t); parsed {
		return test, items[1:]
	}
	return nil, items
}
func parseFailureMode(items []interface{}) (FailureMode, []interface{}) {
	if mode, parsed := item(items).(FailureMode); parsed {
		return mode, items[1:]
	}
	return FailureInherits, items
}
func parseAction(items []interface{}) (func(), []interface{}) {
	switch x := item(items).(type) {
	case nil:
		return nil, items[1:]
	case func():
		return x, items[1:]
	}
	panic(parseError)
}

// This interface allows us to pass the *testing.T struct
// throughout the internals of this tool without ever
// having to import the "testing" package.
type t interface {
	Fail()
}

const parseError = "You must provide a name (string), then a *testing.T (if in outermost scope), an optional FailureMode, and then an action (func())."
