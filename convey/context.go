// Oh the stack trace scanning!
// The density of comments in this file is evidence that
// the code doesn't exactly explain itself. Tread with care...
package convey

import (
	"fmt"
	"github.com/jtolds/gls"
	"reflect"
	"runtime"
	"strings"

	"github.com/smartystreets/goconvey/convey/reporting"
)

const (
	missingGoTest = `Top-level calls to Convey(...) need a reference to the *testing.T.
		Hint: Convey("description here", t, func() { /* notice that the second argument was the *testing.T (t)! */ }) `
	extraGoTest = `Only the top-level call to Convey(...) needs a reference to the *testing.T.`

	failureHalt = "___FAILURE_HALT___"

	nodeKey = "node"
)

type actionSpecifier uint8

const (
	noSpecifier actionSpecifier = iota
	skipConvey
	focusConvey
)

// suiteContext magically handles all coordination of reporter, runners as they handle calls
// to Convey, So, and the like. It does this via runtime call stack inspection, making sure
// that each test function has its own runner, and routes all live registrations
// to the appropriate runner.
type suiteContextNode struct {
	name string

	reporter reporting.Reporter
	test     t

	curIdx   int
	children []*suiteContextNode

	resets []func()

	executedOnce   bool
	expectChildRun *bool
	result         VisitResult

	focus       bool
	failureMode FailureMode
}

type VisitResult uint8

const (
	VisitedIncomplete VisitResult = iota
	VisitedOK
	VisitedPanic
)

func (c *suiteContextNode) shouldVisit() bool {
	return c.result == VisitedIncomplete && *c.expectChildRun
}

func getCurrentContext() *suiteContextNode {
	ctx, ok := ctxMgr.GetValue(nodeKey)
	if ok {
		return ctx.(*suiteContextNode)
	}
	return nil
}

func mustGetCurrentContext() *suiteContextNode {
	ctx := getCurrentContext()
	if ctx == nil {
		panic("cannot perform operation outside of a convey context")
	}
	return ctx
}

func (ctx *suiteContextNode) conveyInner(situation string, f func(C)) {
	defer func() {
		ctx.executedOnce = true
		ctx.curIdx = 0
		*ctx.expectChildRun = false
	}()

	fname := runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
	ctx.reporter.Enter(reporting.NewScopeReport(situation, fname))
	defer ctx.reporter.Exit()

	ctx.resets = []func(){}
	defer func() {
		for _, r := range ctx.resets {
			// TODO(riannucci): handle panics?
			r()
		}
	}()

	defer func() {
		if problem := recover(); problem != nil {
			if strings.HasPrefix(fmt.Sprintf("%v", problem), extraGoTest) {
				panic(problem)
			}
			if problem != failureHalt {
				ctx.reporter.Report(reporting.NewErrorReport(problem))
			}
			ctx.result = VisitedPanic
		} else {
			ctx.result = VisitedOK
			for _, child := range ctx.children {
				if child.result == VisitedIncomplete {
					ctx.result = VisitedIncomplete
					return
				}
			}
		}
	}()

	if f == nil {
		// skipped
		ctx.reporter.Report(reporting.NewSkipReport())
	} else {
		f(ctx)
	}
}

func computeNewFailureMode(parent, cur FailureMode) FailureMode {
	if cur == FailureInherits {
		if parent == "" {
			return defaultFailureMode
		}
		return parent
	}
	return cur
}

func RootConvey(items ...interface{}) {
	entry := discover(items)

	if entry.Test == nil {
		panic(missingGoTest)
	}
	expectChildRun := true
	ctx := &suiteContextNode{
		name: entry.Situation,

		test:     entry.Test,
		reporter: buildReporter(),

		expectChildRun: &expectChildRun,

		focus:       entry.Focus,
		failureMode: computeNewFailureMode("", entry.FailMode),
	}
	ctxMgr.SetValues(gls.Values{nodeKey: ctx}, func() {
		ctx.reporter.BeginStory(reporting.NewStoryReport(ctx.test))
		defer ctx.reporter.EndStory()

		for ctx.shouldVisit() {
			ctx.conveyInner(entry.Situation, entry.Func)
			expectChildRun = true
		}
	})
}

func (ctx *suiteContextNode) SkipConvey(items ...interface{}) {
	ctx.Convey(items, skipConvey)
}

func (ctx *suiteContextNode) FocusConvey(items ...interface{}) {
	ctx.Convey(items, focusConvey)
}

func (ctx *suiteContextNode) Convey(items ...interface{}) {
	entry := discover(items)

	// we're a branch, or leaf (on the wind)
	if entry.Test != nil {
		panic(extraGoTest)
	}
	if ctx.focus && !entry.Focus {
		return
	}

	var inner_ctx *suiteContextNode
	if ctx.executedOnce {
		if ctx.curIdx >= len(ctx.children) {
			panic("different set of Convey statements on subsequent pass!")
		}
		inner_ctx = ctx.children[ctx.curIdx]
		if inner_ctx.name != entry.Situation {
			panic("different set of Convey statements on subsequent pass!")
		}
		ctx.curIdx++
	} else {
		inner_ctx = &suiteContextNode{
			name:     entry.Situation,
			test:     ctx.test,
			reporter: ctx.reporter,

			expectChildRun: ctx.expectChildRun,

			focus:       entry.Focus,
			failureMode: computeNewFailureMode(ctx.failureMode, entry.FailMode),
		}
		ctx.children = append(ctx.children, inner_ctx)
	}

	if inner_ctx.shouldVisit() {
		ctxMgr.SetValues(gls.Values{nodeKey: inner_ctx}, func() {
			inner_ctx.conveyInner(entry.Situation, entry.Func)
		})
	}
}

func (ctx *suiteContextNode) SkipSo(stuff ...interface{}) {
	ctx.assertionReport(reporting.NewSkipReport())
}

func (ctx *suiteContextNode) So(actual interface{}, assert assertion, expected ...interface{}) {
	if result := assert(actual, expected...); result == assertionSuccess {
		ctx.assertionReport(reporting.NewSuccessReport())
	} else {
		ctx.assertionReport(reporting.NewFailureReport(result))
	}
}

func (ctx *suiteContextNode) assertionReport(r *reporting.AssertionResult) {
	ctx.reporter.Report(r)
	if r.Failure != "" && ctx.failureMode == FailureHalts {
		panic(failureHalt)
	}
}

func (ctx *suiteContextNode) Reset(action func()) {
	/* TODO: Failure mode configuration */
	ctx.resets = append(ctx.resets, action)
}

func (ctx *suiteContextNode) Print(items ...interface{}) (int, error) {
	fmt.Fprint(ctx.reporter, items...)
	return fmt.Print(items...)
}

func (ctx *suiteContextNode) Println(items ...interface{}) (int, error) {
	fmt.Fprintln(ctx.reporter, items...)
	return fmt.Println(items...)
}

func (ctx *suiteContextNode) Printf(format string, items ...interface{}) (int, error) {
	fmt.Fprintf(ctx.reporter, format, items...)
	return fmt.Printf(format, items...)
}
