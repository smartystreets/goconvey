package convey

/*

These tests will assert that the runner reports failures, errors and
successes to some reporting abstraction.  They will probably be similar
in feeling to the execution_tests.

The reporting abstraction will merely aggregate statistics, not present
output to the user (that will be a presenter abstraction)

So, the trick is to hook up the various scope instances to the runner
so the runner knows which scope is currently executing (enter, exit, etc..).

Then, the So method can be hooked up to the runner so that a success or
failure will make it to the reporter via the runner.

Errors will have to be passed to the reporter in the defer recovery method
(probably...).

* The reporter will depend on the accurate 'Convey'-ing of the *testing.T
  to the top-level scope registrations.

*/

import (
	"github.com/smartystreets/goconvey/convey/execution"
	"testing"
)

func TestSuccessesLogged(t *testing.T) {
	t.Skip()
}
