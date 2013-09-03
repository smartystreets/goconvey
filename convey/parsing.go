package convey

import "github.com/smartystreets/goconvey/gotest"

func parseRegistration(items []interface{}) (name string, action func(), test gotest.T) {
	ensureEnough(items)

	name = parseName(items)
	test = parseGoTest(items)
	action = parseAction(items, test)

	return name, action, test
}
func ensureEnough(items []interface{}) {
	if len(items) < 2 {
		panic(parseError)
	}
}
func parseName(items []interface{}) string {
	if name, parsed := items[0].(string); parsed {
		return name
	}
	panic(parseError)
}
func parseGoTest(items []interface{}) gotest.T {
	if test, parsed := items[1].(gotest.T); parsed {
		return test
	}
	return nil
}
func parseAction(items []interface{}, test gotest.T) func() {
	var index = 1
	if test != nil {
		index = 2
	}

	if action, parsed := items[index].(func()); parsed {
		return action
	}
	panic(parseError)
}

const parseError = "You must provide a name (string), then a *testing.T (if in outermost scope), and then an action (func())."
