package main

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestCheckExcludedDirs(t *testing.T) {
	listOfTestExcludedDirs := map[string]string{
		"/go/src/package/vendor":        "vendor",
		"/go/src/package/vendor/":       "vendor",
		"/go/src/package/node_modules":  "node_modules",
		"/go/src/package/node_modules/": "node_modules",
	}

	Convey("Exclude excluded directory for fake pathes", t, func() {
		for working, excl := range listOfTestExcludedDirs {
			excludedDirItems := checkExcludedDirs(excludedDirs, working)
			for _, item := range excludedDirItems {
				So(item, ShouldNotEqual, excl)
			}
		}
	})
}
