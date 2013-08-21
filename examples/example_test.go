package examples

import (
	"testing"
	. "github.com/mdwhatcott/goconvey"
)

func TestSimple(t *testing.T) {
	Convey("a single step", t, func(c C) {
		So(1, ShouldEqual, 1)
		So(1, ShouldNotEqual, 2)
	})
}

func TestIsolatedNesting(t *testing.T) {
	Convey("given an integer value", t, func(c C) {
		x := 42

		c.Convey("when incremented", func(c C) {
			x++

			c.Convey("then the value should be greater by one", func(c C) {
				So(x, ShouldEqual, 43)
			})

			c.Convey("and the value should NOT be what it used to be", func(c C) {
				So(x, ShouldNotEqual, 42)
			})
		})

		c.Convey("when decremented", func(c C) {
			x--

			c.Convey("then the value should be lesser by one", func(c C) {
				So(x, ShouldEqual, 41)
			})

			c.Convey("then the value should NOT be what it used to be", func(c C) {
				So(x, ShouldNotEqual, 42)
			})
		})

		c.Convey("cleanup after", func(c C) {
			x = 0
			So(x, ShouldEqual, 0)
		})
	})
}
