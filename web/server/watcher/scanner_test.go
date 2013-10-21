package watcher

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestScanner(t *testing.T) {
	Convey("Subject: Scanner", t, func() {
		Convey("As a result of calling Scan()", func() {
			Convey("For the first time", func() {
				Convey("The scanner should report a change in state", nil)
			})

			Convey("On subsequent calls", func() {
				Convey("When the file system has not changed", func() {
					Convey("The scanner should NOT report any change in state", nil)
				})

				Convey("When an existing file on the file system has been modified", func() {
					Convey("The scanner should report the change via its return value", nil)
				})

				Convey("When the file system has received a new directory", func() {
					Convey("The scanner should report the addition to the watcher", nil)
				})

				Convey("When the file system loses a directory", func() {
					Convey("The scanner should report the deletion to the watcher", nil)
				})
			})
		})
	})
}
