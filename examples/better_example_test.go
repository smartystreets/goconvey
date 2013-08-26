package main

import (
	. "github.com/mdwhatcott/goconvey"
)

func init() {
	// convey("one equals one", so(1, ShouldEqual, 1))  // single 'so()', wrapping func() optional. (requires that this be wired to the runner...)

	convey("Subject: Integer incrementation and decrementation", func() {
		x := 0

		convey("Given a starting integer value", func() {
			x = 42

			convey("When incremented", func() {
				x++

				convey("The value should be greater by one", func() {
					so(x, ShouldEqual, 43) 
				})
				convey("The value should NOT be what it used to be", func() {
					so(x, ShouldNotEqual, 42)
				})
			})
			convey("When decremented", func() {
				x--

				convey("The value should be lesser by one", func() {
					so(x, ShouldEqual, 41)
				})
				convey("The value should NOT be what it used to be", func() {
					so(x, ShouldNotEqual, 42)
				})
			})
		})
	})
}

/*
Output:

????
*/