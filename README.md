GoConvey is awesome Go testing
==============================

Welcome to GoConvey, a yummy Go testing tool for gophers. Works with `go test`. Use it in the terminal or browser according your viewing pleasure. **[View full feature tour!](http://smartystreets.github.io/goconvey)**

**Features:**

- Directly integrates with `go test`
- Huge suite of regression tests
- Readable, colorized console output (understandable by any manager, IT or not)
- Fully-automatic web UI (works with native Go tests, too)
- Test code generator
- Desktop notifications (optional)
- Auto-test script automatically runs tests in the terminal
- Immediately open problem lines in [Sublime Text](http://www.sublimetext.com) ([some assembly required](https://github.com/asuth/subl-handler))

**Menu:**

- [Installation](#installation)
- [Quick start](#quick-start)
- [Documentation](#documentation)
- [Screenshots](#screenshots)
- [Contributors](#contributors-thanks)




Installation
------------

	$ go get github.com/smartystreets/goconvey





[Quick start](https://github.com/smartystreets/goconvey/wiki#get-going-in-25-seconds)
-----------

Make a test, for example:

```go
func TestSpec(t *testing.T) {
	var x int
	
	// Only pass t into top-level Convey calls
	Convey("Given some integer with a starting value", t, func() {
		x = 1

		Convey("When the integer is incremented", func() {
			x++

			Convey("The value should be greater by one", func() {
				So(x, ShouldEqual, 2)
			})
		})
	})
}
```


#### [In the browser](https://github.com/smartystreets/goconvey/wiki/Web-UI)

Start up the GoConvey web server at your project's path:

	$ $GOPATH/bin/goconvey

Then open your browser to:

	http://localhost:8080

There you have it. As long as GoConvey is running, test results will automatically update in your browser window. The design is responsive, so you can squish the browser real tight if you need to put it beside your code.

The [web UI](https://github.com/smartystreets/goconvey/wiki/Web-UI) supports traditional Go tests, so use it even if you're not using GoConvey tests.



#### [In the terminal](https://github.com/smartystreets/goconvey/wiki/Execution)

Just do what you do best:

    $ go test

Or if you want the output to include the story:

    $ go test -v





[Documentation](https://github.com/smartystreets/goconvey/wiki)
-----------

Check out the 

- [GoConvey wiki](https://github.com/smartystreets/goconvey/wiki),
- [godoc.org](http://godoc.org/github.com/smartystreets/goconvey),
- and the *_test.go files scattered throughout this project.





[Screenshots](http://smartystreets.github.io/goconvey)
-----------

For web UI and terminal screenshots, check out [the full feature tour](http://smartystreets.github.io/goconvey).






Contributors (Thanks!)
----------------------

We appreciate everyone's contributions to the project! Please see the [contributor graphs](https://github.com/smartystreets/goconvey/graphs/contributors) provided by GitHub for all the credits.

GoConvey is brought to you by [SmartyStreets](https://github.com/smartystreets); in particular:

 - [Michael Whatcott](https://github.com/mdwhatcott)
 - [Matt Holt](https://github.com/mholt)
