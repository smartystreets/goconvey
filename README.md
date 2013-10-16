GoConvey makes Go testing awesome
=================================

Welcome to GoConvey, a yummy Go testing tool for gophers. Works with `go test`. Use it in the terminal or browser according your viewing pleasure.

**Features:**

- Directly integrates with `go test`
- Huge suite of regression tests
- Readable, colorized console output (understandable by any manager, IT or not)
- Fully-automatic web UI (works with native Go tests, too)
- Test code generator
- Idler script automatically runs tests in the terminal
- Immediately open problem lines in [Sublime Text](http://www.sublimetext.com) ([some assembly required](https://github.com/asuth/subl-handler))

**Jump to:**

- [Installation](#installation)
- [Quick start](#quick-start)
- [Web UI](#web-ui)
- [Terminal output](#terminal-output-with-go-test--v)
- [Writing tests](#writing-tests)
- [Documentation](#documentation)
- [Execution](#execution)
- [Assertions](#assertions)
- [Writing your own assertions](#writing-your-own-assertions)
- [Skipping Convey() registrations](#skipping-convey-registrations)
- [Unimplemented Convey() registrations](#unimplemented-convey-registrations)
- [Skipping So() assertions](#skipping-so-assertions)
- [Contributors](#contributors-thanks)


Installation
------------

	$ go get github.com/smartystreets/goconvey
	$ go install github.com/smartystreets/goconvey/web/goconvey-server


Quick start
-----------

Make a test, for example:

```go
func TestSpec(t *testing.T) {
	Convey("Given some integer with a starting value", t, func() {
		x := 1

		Convey("When the integer is incremented", func() {
			x++

			Convey("The value should be greater by one", func() {
				So(x, ShouldEqual, 2)
			})
		})
	})
}
```

### In the browser

Start up the GoConvey web server at your project's path:

    $ go install github.com/smartystreets/goconvey/web/goconvey-server
    ...
    $ $GOPATH/bin/goconvey-server

Then open your browser to:

	http://localhost:8080

There you have it. As long as GoConvey is running, test results will automatically update in your browser window. The design is responsive, so you can squish the browser real tight if you need to put it beside your code.

The browser UI supports traditional Go tests, so feel free to use it even if you're not using the GoConvey style of testing.


### In the terminal

Just do what you do best:

    $ go test

Or if you want the output to include the story:

    $ go test -v


Web UI
-----------
![GoConvey rocks](http://i.imgur.com/O7uVvoq.png)


Terminal output (with `go test -v`):
---------------

**Tests pass:**

![Pass](http://i.imgur.com/c2qAQcR.png)

**Test fail:**

![Fail](http://i.imgur.com/sRcyZBr.png)

**Test panic:**

![Panic](http://i.imgur.com/iG2EZ5C.png)



Writing tests
-------------

You can write GoConvey tests manually or with a nice code generator.


### Code Generator

From the web UI served by GoConvey, click "Code Gen" in the top-right. Then describe your program's behavior in a natural, flowing way, for example (make sure you convert indents to tabs, as GitHub transformed them to spaces):


	TestSpec
		Subject: Integer incrementation and decrementation
			Given a starting integer value
				When incremented
					The value should be greater by one
					The value should NOT be what it used to be
				When decremented
					The value should be lesser by one
					The value should NOT be what it used to be


The skeleton of your test file will be stubbed out automatically as you type. There
are a few things to notice about this:

- Lines starting with "Test" (case-sensitive), without indentation, are treated as the name of the function in which all nested tests will be included
- Indentation defines scope
- Assertions are not made here; you'll do that later after pasting the generated code into your `_test.go` file.


### Manually

See the [examples folder](https://github.com/smartystreets/goconvey/tree/master/examples).
We recommend reviewing [isolated_execution_test.go](https://github.com/smartystreets/goconvey/blob/master/convey/isolated_execution_test.go) for a more thorough understanding of how tests are
composed and how they actually work.



Documentation:
--------------

Public functions are documented using GoDoc conventions. See the [godoc.org page for this project]
(http://godoc.org/github.com/smartystreets/goconvey) for the details.



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
	        The score should be zero ✔
	      When all throws knock down only one pin 
	        The score should be 20 ✔
	      When a spare is thrown 
	        The score should include a spare bonus. ✔
	      When a strike is thrown 
	        The score should include a strike bonus. ✔
	      When all strikes are thrown 
	        The score should be 300. ✔

	5 assertions and counting

	--- PASS: TestScoring (0.00 seconds)
	=== RUN TestSpec

	  Subject: Integer incrementation and decrementation 
	    Given a starting integer value 
	      When incremented 
	        The value should be greater by one ✔
	        The value should NOT be what it used to be ✔
	      When decremented 
	        The value should be lesser by one ✔
	        The value should NOT be what it used to be ✔

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

### General Equality

	So(thing, ShouldEqual, thing2)
	So(thing, ShouldNotEqual, thing2)
	So(thing, ShouldResemble, thing2)
	So(thing, ShouldNotResemble, thing2)
	So(thing, ShouldPointTo, thing2)
	So(thing, ShouldNotPointTo, thing2)
	So(thing, ShouldBeNil, thing2)
	So(thing, ShouldNotBeNil, thing2)
	So(thing, ShouldBeTrue)
	So(thing, ShouldBeFalse)

### Numeric Quantity comparison

	So(1, ShouldBeGreaterThan, 0)
	So(1, ShouldBeGreaterThanOrEqualTo, 0)
	So(1, ShouldBeLessThan, 2)
	So(1, ShouldBeLessThanOrEqualTo, 2)
	So(1.1, ShouldBeBetween, .8, 1.2)
	So(1.1, ShouldNotBeBetween, 2, 3)
	So(1.1, ShouldBeBetweenOrEqual, .9, 1.1)
	So(1.1, ShouldNotBeBetweenOrEqual, 1000, 2000)

### Collections

	So([]int{2, 4, 6}, ShouldContain, 4)
	So([]int{2, 4, 6}, ShouldNotContain, 5)
	So(4, ShouldBeIn, ...[]int{2, 4, 6})
	So(4, ShouldNotBeIn, ...[]int{1, 3, 5})

### Strings

	So("asdf", ShouldStartWith, "as")
	So("asdf", ShouldNotStartWith, "df")
	So("asdf", ShouldEndWith, "df")
	So("asdf", ShouldNotEndWith, "df")
	So("asdf", ShouldContain, "sd")  // optional 'expected occurences' arguments?
	So("asdf", ShouldNotContain, "er")
	So("adsf", ShouldBeBlank)
	So("asdf", ShouldNotBeBlank)

### panic

	So(func(), ShouldPanic)
	So(func(), ShouldNotPanic)
	So(func(), ShouldPanicWith, "") // or errors.New("something")
	So(func(), ShouldNotPanicWith, "") // or errors.New("something")

### Type checking

	So(1, ShouldHaveSameTypeAs, 0)
	So(1, ShouldNotHaveSameTypeAs, "asdf")

### time.Time (and time.Duration)

	So(time.Now(), ShouldHappenBefore, time.Now())
	So(time.Now(), ShouldHappenOnOrBefore, time.Now())
	So(time.Now(), ShouldHappenAfter, time.Now())
	So(time.Now(), ShouldHappenOnOrAfter, time.Now())
	So(time.Now(), ShouldHappenBetween, time.Now(), time.Now())
	So(time.Now(), ShouldHappenOnOrBetween, time.Now(), time.Now())
	So(time.Now(), ShouldNotHappenOnOrBetween, time.Now(), time.Now())
	So(time.New(), ShouldHappenWithin, duration, time.Now())
	So(time.New(), ShouldNotHappenWithin, duration, time.Now())

Thanks to [github.com/jacobsa](https://github.com/jacobsa/oglematchers) for his excellent 
[oglematchers](https://github.com/smartystreets/oglmatchers) library, which
is what many of these methods make use of to do their jobs.


Writing your own assertions:
----------------------------

Sometimes a test suite will might need an assertion that is too
specific to be included in this tool. Not to worry, simply implement
a function with the following signature (fill in the bracketed parts
and string values):

```go
func should<do-something>(actual interface, expected ...interface{}) string {
    if <some-important-condition-is-met(actual, expected)> {
        return "" // empty string means the assertion passed
    } else {
        return "some descriptive message detailing why the assertion failed..."
    }
}
```

Suppose I implemented the following assertion:

```go
func shouldScareGophersMoreThan(actual interface, expected ...interface{}) string {
    if actual == "BOO!" && expected[0] == "boo" {
        return ""
    } else {
        return "Ha! You'll have to get a lot friendlier with the capslock if you want to scare a gopher!"
    }
}
```

I can then make use of the assertion function when calling the `So(...)` method in the tests:

```go
Convey("All caps always makes text more meaningful", func() {
    So("BOO!", shouldScareGophersMoreThan, "boo")
})
```

Skipping `Convey` Registrations:
--------------------------------

Changing a `Convey` to `SkipConvey` prevents the `func()` associated with
that call from running. This also has the consequence of preventing any nested 
`Convey` registrations from running. The reporter will indicate that the 
registration was skipped.

```go
SkipConvey("Important stuff", func() { // This func() will not be executed!

    Convey("More important stuff", func() {
        So("asdf", ShouldEqual, "asdf")
    })

})
```

Using `SkipConvey` has nearly the same effect as commenting out the test
entirely. However, this is preferred over commenting out tests to avoid the
usual "declared/imported but not used" errors. Usage of `SkipConvey` is
intended for temporary code alterations.


Unimplemented `Convey` Registrations:
-------------------------------------

When composing `Convey` registrations sometimes it's convenient to use `nil`
instead of an actual `func()`. This allows you to do that and it also provides
an indication in the report that the registration is not complete.

```go
Convey("Some stuff", func() {

    // This will show up as 'skipped' in the report
    Convey("Should go boink", nil)

}
```

Skipping `So` Assertions:
-------------------------

Similar to `SkipConvey`, changing a `So` to `SkipSo` prevents the execution of
that assertion. The report will show that the assertion was skipped.

```go
Convey("1 Should Equal 2", func() {
    
    // This assertion will not be executed and will show up as 'skipped' in the report
    SkipSo(1, ShouldEqual, 2)

})
```

And like `SkipConvey`, this function is only intended for use during
temporary code alterations.


Contributors (Thanks!):
-----------------------

 - [Michael Whatcott](https://github.com/mdwhatcott)
 - [Matt Holt](https://github.com/mholt)
