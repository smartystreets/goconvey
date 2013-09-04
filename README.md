GoConvey - BDD in Go - by SmartyStreets, LLC
============================================

Welcome to GoConvey, a yummy BDD tool for gophers. You'll soon be enjoying the benefits of
this robust, descriptive and fun-to-use tool. Among those benefits are the following:

- An ever-growing suite of regression tests
- Tests which are formatted to the console as a readable specification, accessible to any manager (IT or not).
- Integration with the already excellent `go test` tool
- Constant updates on the working state of your application via the bundled `idle.py` script (more below).


Installation:
-------------

	(assuming you have set your $GOPATH environment variable) (*link*)
	
	$ go get github.com/smartystreets/goconvey



Composition:
------------

See the [examples folder](https://github.com/smartystreets/goconvey/tree/master/examples).


Execution:
----------

Concise mode (default):

	$ cd $GOPATH/src/github.com/smartystreets/goconvey/examples
	$ go test
	.....

	5 assertions and counting

	....

	9 assertions and counting

	PASS
	ok  	github.com/smartystreets/goconvey/examples	0.022s



Verbose mode:

	$ cd $GOPATH/src/github.com/smartystreets/goconvey/examples
	$ go test -v
	=== RUN TestScoring

	  Subject: Bowling Game Scoring 
	    Given a fresh score card 
	      When all gutter balls are thrown 
	        The score should be zero ✓
	      When all throws knock down only one pin 
	        The score should be 20 ✓
	      When a spare is thrown 
	        The score should include a spare bonus. ✓
	      When a strike is thrown 
	        The score should include a strike bonus. ✓
	      When all strikes are thrown 
	        The score should be 300. ✓

	5 assertions and counting

	--- PASS: TestScoring (0.00 seconds)
	=== RUN TestSpec

	  Subject: Integer incrementation and decrementation 
	    Given a starting integer value 
	      When incremented 
	        The value should be greater by one ✓
	        The value should NOT be what it used to be ✓
	      When decremented 
	        The value should be lesser by one ✓
	        The value should NOT be what it used to be ✓

	9 assertions and counting

	--- PASS: TestSpec (0.00 seconds)
	PASS
	ok  	github.com/smartystreets/goconvey/examples	0.023s


Happy? Well, it gets even better with the auto-reload script.  This script detects changes below
the directory where it is run and recurses that location, running `go test` wherever it finds
`*_test.go` files. This means that once it's running you don't ever have to leave your editor
to run tests.  It also accepts the `-v` argument for verbose mode. Check it out:

	$ cd $GOPATH/src/github.com/smartystreets/goconvey/examples
	$ $GOPATH/src/github.com/smartystreets/goconvey/scripts/idle.py

	-------------------------------------- 1 --------------------------------------

	.....

	5 assertions and counting

	....

	9 assertions and counting

	PASS
	ok  	github.com/smartystreets/goconvey/examples	0.022s

	(now waiting for you to save changes to the application under test...)


Assertions:
-----------

Here's the listing of assertions that this project aims to implement 
(see the examples folder for actual usage):


 1.0  | completed |usage
:----:|:---------:|-----
      |           |__Equality__
*     |X          |So(thing, ShouldEqual, thing2)
*     |X          |So(thing, ShouldNotEqual, thing2)
*     |           |So(thing, ShouldBeLike, thing2)
*     |           |So(thing, ShouldNotBeLike, thing2)
      |           |So(thing, ShouldPointTo, thing2)
      |           |So(thing, ShouldNotPointTo, thing2)
*     |X          |So(thing, ShouldBeNil, thing2)
*     |           |So(thing, ShouldNotBeNil, thing2)
*     |           |So(thing, ShouldBeTrue)
*     |           |SO(thing, ShouldBeFalse)
      |           |__Interfaces__
      |           |So(1, ShouldImplement, Interface)
      |           |So(1, ShouldNotImplement, OtherInterface)
      |           |__Type checking__
      |           |So(1, ShouldBeAn, int)
      |           |So(1, ShouldNotBeAn, int)
      |           |So("1", ShouldBeA, string)
      |           |So("1", ShouldNotBeA, string)
      |           |__Quantity comparison__
*     |           |So(1, ShouldBeGreaterThan, 0)
*     |           |So(1, Shou|ldBeGreaterThanOrEqualTo, 0)
*     |           |So(1, ShouldBeLessThan, 2)
*     |           |So(1, ShouldBeLessThanOrEqualTo, 2)
      |           |__Tolerences__
*     |           |So(1.1, ShouldBeWithin, .1, 1)
*     |           |So(1.1, ShouldNotBeWithin, .1, 2)
      |           |__Collections__
*     |           |So([]int{}, ShouldBeEmpty)
*     |           |So([]int{1}, ShouldNotBeEmpty)
*     |           |So([]int{1, 2, 3}, ShouldContain, 1, 2)
*     |           |So([]int{1, 2, 3}, ShouldNotContain, 4, 5)
      |           |So(1, ShouldBeIn, []int{1, 2, 3})
      |           |So(4, ShouldNotBeIn, []int{1, 2, 3})
      |           |__Strings__
*     |           |So(|"asdf", ShouldStartWith, "as")
*     |           |So("asdf", ShouldNotStartWith, "df")
*     |           |So("asdf", ShouldEndWith, "df")
*     |           |So("asdf", ShouldNotEndWith, "df")
      |           |So("(asdf)", ShouldBeSurroundedWith, "()")
      |           |So("(asdf)", ShouldNotBeSurroundedWith, "[]")



RoadMap:
--------

Still in development:

	- Randomized execution of stories (including resets)
	- Full suite of 'Should' assertions


Would be awesome:

	- Auto-reloading local http endpoint:
		- Output Story presentation to HTML file (https://github.com/russross/blackfriday)
		- Create http endpoint that serves the html output
		- make http endpoint poll for updates reload report (collapse all but failed and erred stuff)
		- clicking on filename/line-number in web report shows that file as a web page w/ problem line highlighted