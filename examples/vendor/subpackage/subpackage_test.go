package subpackage

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestEcho(t *testing.T) {
	s := "Hello"
	Convey("Test Echo", t, func() {
		So(Echo(s), ShouldEqual, s)
	})
}
