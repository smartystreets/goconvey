"""
This integration test script will execute various HTTP actions against the 
goconvey-server to ensure that tests are being run and that the watcher 
can be updated.
"""


import urllib2
import json


def main():
    # GET: /watcher: returns working directory
    # POST: /watcher?new-dir: returns 200
    # GET: /watcher: returns new-dir
    
    # create new package
    # GET: /latest: returns latest output
    # add new test file
    # GET: /latest: returns latest output (including new test)
    # remove test file
    # GET: /latest: returns first output
    # remove package
    # GET: /latest: returns blank stuff?
    pass


if __name__ == '__main__':
    main()


TEST_CODE = """package testing

import (
    "testing"
    . "github.com/smartystreets/goconvey/convey"
)

func TestSomething(t *testing.T) {
    Convey("Something", t, func() {
        Convey("should happen", func() {
            So(true, ShouldBeTrue)
        })
    })
}

"""

TESTS = {
    'initial': ('initial_test.go', TEST_CODE),
    'additional': ('additional_test.go', TEST_CODE.replace("Something", "Something2")),
}
