package examples

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestSpec(t *testing.T) {
	Convey("Subject: Integer incrementation and decrementation", t, func() {
		var x int

		Convey("Given a starting integer value", func() {
			x = 42

			Convey("When incremented", func() {
				x++

				Convey("The value should be greater by one", func() {
					So(x, ShouldEqual, 43)
				})
				Convey("The value should NOT be what it used to be", func() {
					So(x, ShouldNotEqual, 42)
				})
			})
			Convey("When decremented", func() {
				x--

				Convey("The value should be lesser by one", func() {
					So(x, ShouldEqual, 41)
				})
				Convey("The value should NOT be what it used to be", func() {
					So(x, ShouldNotEqual, 42)
				})
			})
			Reset(func() {
				x = 0
			})
		})
	})
}
