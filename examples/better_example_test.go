package examples

import (
	"testing"
	. "github.com/mdwhatcott/goconvey"
)

func Test(t *testing.T) {
	Subject("blah blah blah", t, func() {
		Convey("a single step", t, func() {
			So(1, ShouldEqual, 1)
			So(1, ShouldNotEqual, 2)
		})	
	})

	Subject("nested example", t, func() {
		Convey("given an integer value", func() {
			x := 42

			Convey("when incremented", func() {
				x++

				Convey("then the value should be greater by one", func() {
					So(x, ShouldEqual, 43) 
				})

				Convey("and the value should NOT be what it used to be", func() {
					So(x, ShouldNotEqual, 42)
				})
			})

			Convey("when decremented", func() {
				x--

				Convey("then the value should be lesser by one", func() {
					So(x, ShouldEqual, 41)
				})

				Convey("then the value should NOT be what it used to be", func() {
					So(x, ShouldNotEqual, 42)
				})
			})

			Convey("cleanup after", func() {
				x = 0
				So(x, ShouldEqual, 0)
			})
		})
	})
}
