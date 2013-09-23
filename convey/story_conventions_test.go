package convey

import (
	"github.com/smartystreets/goconvey/execution"
	"testing"
)

func TestMissingTopLevelGoTestReferenceCausesPanic(t *testing.T) {
	runner = execution.NewRunner()

	output := map[string]bool{}

	defer expectEqual(t, false, output["good"])
	defer requireGoTestReference(t)

	Convey("Hi", func() {
		output["bad"] = true // this shouldn't happen
	})
}

func requireGoTestReference(t *testing.T) {
	err := recover()
	if err == nil {
		t.Error("We should have recovered a panic here (because of a missing *testing.T reference)!")
	} else {
		expectEqual(t, execution.MissingGoTest, err)
	}
}

func TestMissingTopLevelGoTestReferenceAfterGoodExample(t *testing.T) {
	runner = execution.NewRunner()

	output := map[string]bool{}

	defer func() {
		expectEqual(t, true, output["good"])
		expectEqual(t, false, output["bad"])
	}()
	defer requireGoTestReference(t)

	Convey("Good example", t, func() {
		output["good"] = true
	})

	Convey("Bad example", func() {
		output["bad"] = true // shouldn't happen
	})
}
