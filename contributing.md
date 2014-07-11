### "A pull request, with test, is best."

Do you have any ideas for this project? Want to tackle any of the
[open issues](https://github.com/smartystreets/goconvey/issues?direction=desc&sort=created&state=open)?
Great! As a starting point please take the time to peruse the tests in the part of the project you want 
to work on. That will help you get to know the existing code and will give you an idea of the kinds of
test cases you should include with your pull request. We expect to maintain a high level of confidence
in the code for GoConvey. Your [disciplined efforts](http://butunclebob.com/ArticleS.UncleBob.TheThreeRulesOfTdd)
will make the project much more maintainable in the long-term.

Or, if you're feeling brave, how about tackling one of our [deferred/known issues](Known Issues)? Not for the faint of heart. You have been warned...

Want to contribute but not sure how? There are a few modules that still aren't up to snuff as far as test coverage is concerned (yes, I know test coverage isn't a silver bullet, but where it can be improved why not do so?). Maybe a good way to get your feet wet would be to implement tests that prove the functionality of the packages in question. Here's a current breakdown of test coverage (2013-11-16):

```
PASS
coverage: 33.8% of statements
ok    github.com/smartystreets/goconvey/reporting 0.016s

PASS
coverage: 50.0% of statements
ok    github.com/smartystreets/goconvey/web/server/contract 0.025s

PASS
coverage: 71.6% of statements
ok    github.com/smartystreets/goconvey/web/server/system 0.049s

PASS
coverage: 87.5% of statements
ok    github.com/smartystreets/goconvey/printing  0.011s

PASS
coverage: 92.3% of statements
ok    github.com/smartystreets/goconvey/convey  0.044s

PASS
coverage: 97.4% of statements
ok    github.com/smartystreets/goconvey/web/server/parser 0.025s

PASS
coverage: 98.6% of statements
ok    github.com/smartystreets/goconvey/web/server/executor 1.075s

PASS
coverage: 99.8% of statements
ok    github.com/smartystreets/goconvey/assertions  0.019s

PASS
coverage: 100.0% of statements
ok    github.com/smartystreets/goconvey/examples  0.063s

PASS
coverage: 100.0% of statements
ok    github.com/smartystreets/goconvey/web/server/api  0.266s

PASS
coverage: 100.0% of statements
ok    github.com/smartystreets/goconvey/web/server/watcher  0.083s

```

Of course, you'll want to get more details from the the built-in (as of Go 1.2) cover tool so you know where to start within each module. I'll leave that as your first exercise...

