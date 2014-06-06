package convey

func discover(items []interface{}) *registration {
	ensureEnough(items)

	name, items := parseName(items)
	test, items := parseGoTest(items)
	action := parseAction(items)

	return newRegistration(name, action, test)
}
func ensureEnough(items []interface{}) {
	if len(items) < 2 {
		panic(parseError)
	}
}
func parseName(items []interface{}) (string, []interface{}) {
	if name, parsed := items[0].(string); parsed {
		return name, items[1:]
	}
	panic(parseError)
}
func parseGoTest(items []interface{}) (t, []interface{}) {
	if test, parsed := items[0].(t); parsed {
		return test, items[1:]
	}
	return nil, items
}
func parseFailureMode(items []interface{}) (FailureMode, []interface{}) {
	if mode, parsed := items[0].(FailureMode); parsed {
		return mode, items[1:]
	}
	return FailureInherits, items
}
func parseAction(items []interface{}) *action {
	failure, items := parseFailureMode(items)

	if action, parsed := items[0].(func()); parsed {
		return newAction(action, failure)
	}
	if items[0] == nil {
		return newSkippedAction(skipReport, failure)
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
