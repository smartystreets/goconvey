package convey

import (
	"fmt"
	"github.com/smartystreets/goconvey/execution"
	"strings"
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

func TestExtraReferencePanics(t *testing.T) {
	runner = execution.NewRunner()
	output := map[string]bool{}

	defer func() {
		err := recover()
		if err == nil {
			t.Error("We should have recovered a panic here (because of an extra *testing.T reference)!")
		} else if !strings.HasPrefix(fmt.Sprintf("%v", err), execution.ExtraGoTest) {
			t.Error("Should have panicked with the 'extra go test' error!")
		}
		if output["bad"] {
			t.Error("We should NOT have run the bad example!")
		}
	}()

	Convey("Good example", t, func() {
		Convey("Bad example - passing in *testing.T a second time!", t, func() {
			output["bad"] = true // shouldn't happen
		})
	})
}
