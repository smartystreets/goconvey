package something

import (
	. "github.com/mdwhatcott/goconvey/convey"
	"testing"
)

func Test(t *testing.T) {
	Run(t, func() {
		Convey("Subject: Integer incrementation and decrementation", func() {
			x := 0

			Convey("Given a starting integer value", func() {
				x = 42

				Convey("When incremented", func() {
					x++

					Convey("The value should be greater by one", func() {
						So(x, ShouldEqual, 43)
					})
					Convey("The value should NOT be what it used to be", func() {
						// So(x, ShouldNotEqual, 42)
					})
				})
				Convey("When decremented", func() {
					x--

					Convey("The value should be lesser by one", func() {
						So(x, ShouldEqual, 41)
					})
					Convey("The value should NOT be what it used to be", func() {
						// So(x, ShouldNotEqual, 42)
					})
				})
				Reset(func() {
					x = 0
				})
			})
		})
	})
}

/*
Output:

	- Subject: Integer incrementation and decrementation
		- Given a starting integer value
			- When incremented
				- The value should be greater by one

	Subject: Integer incrementation and decrementation
		Given a starting integer value
			When incremented
				The value should NOT be what it used to be

	Subject: Integer incrementation and decrementation
		Given a starting integer value
			When decremented
				The value should be lesser by one

	Subject: Integer incrementation and decrementation
		Given a starting integer value
			When decremented
X				The value should NOT be what it used to be
X Discrepancy:
X   Stack trace (line 13): 42 should equal 43 but did not!

Running Total:

1 Stor[y|ies] with 4 Assertion[s] (1 Failed)[, 0 Errors]

*/
