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
	missingGoTest string = `Top-level calls to Convey(...) need a reference to the *testing.T.
		Hint: Convey("description here", t, func() { /* notice that the second argument was the *testing.T (t)! */ }) `
	extraGoTest string = `Only the top-level call to Convey(...) needs a reference to the *testing.T.`

	failureHalt = "___FAILURE_HALT___"

	nodeKey string = "node"
)

// suiteContext magically handles all coordination of reporter, runners as they handle calls
// to Convey, So, and the like. It does this via runtime call stack inspection, making sure
// that each test function has its own runner, and routes all live registrations
// to the appropriate runner.
type suiteContextNode struct {
	name string

	reporter reporting.Reporter
	test     t

	parent   *suiteContextNode
	curIdx   int
	children []*suiteContextNode

	resets []func()

	executedOnce   bool
	expectChildRun bool
	result         VisitResult

	focus       bool
	failureMode FailureMode
}

func (c *suiteContextNode) Write(p []byte) (int, error) {
	return c.reporter.Write(p)
}

type VisitResult uint8

const (
	VisitedIncomplete VisitResult = iota
	VisitedOK
	VisitedPanic
)

func (c *suiteContextNode) shouldVisit() bool {
	return c.result == VisitedIncomplete && (c.parent == nil || c.parent.expectChildRun)
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

func conveyanceInner(situation string, f func(), ctx *suiteContextNode) {
	ctx.expectChildRun = true
	defer func() {
		ctx.executedOnce = true
		ctx.curIdx = 0
		tmp := ctx
		for tmp != nil {
			tmp.expectChildRun = false
			tmp = tmp.parent
		}
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
			fmt.Println(situation, "marked as panic")
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
		f()
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

func Conveyance(entry *suite) {
	ctx := getCurrentContext()
	if ctx == nil {
		// we're the root
		if entry.Test == nil {
			panic(missingGoTest)
		}
		ctx = &suiteContextNode{
			name: entry.Situation,

			test:     entry.Test,
			reporter: buildReporter(),

			focus:       entry.Focus,
			failureMode: computeNewFailureMode("", entry.FailMode),
		}
		ctxMgr.SetValues(gls.Values{nodeKey: ctx}, func() {
			ctx.reporter.BeginStory(reporting.NewStoryReport(ctx.test))
			defer ctx.reporter.EndStory()

			rootRun := func() {
				fmt.Println("root", entry.Situation)
				defer func() {
					el := recover()
					fmt.Println("done", entry.Situation, el)
					ctx.expectChildRun = true
					if el != nil {
						panic(el)
					}
				}()
				conveyanceInner(entry.Situation, entry.Func, ctx)
			}

		keepRunning:
			rootRun()
			if ctx.shouldVisit() {
				for _, c := range ctx.children {
					if c.shouldVisit() {
						goto keepRunning
					}
				}
			}
		})
	} else {
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

				parent: ctx,

				focus:       entry.Focus,
				failureMode: computeNewFailureMode(ctx.failureMode, entry.FailMode),
			}
			ctx.children = append(ctx.children, inner_ctx)
		}

		if inner_ctx.shouldVisit() {
			fmt.Println("visiting", entry.Situation)
			defer func() {
				el := recover()
				fmt.Println("done_visit", entry.Situation, el)
				if el != nil {
					panic(el)
				}
			}()
			ctxMgr.SetValues(gls.Values{nodeKey: inner_ctx}, func() {
				conveyanceInner(entry.Situation, entry.Func, inner_ctx)
			})
		} else {
			fmt.Println("deferring", entry.Situation)
		}
	}
}

func assertionReport(r *reporting.AssertionResult) {
	ctx := mustGetCurrentContext()
	ctx.reporter.Report(r)
	if r.Failure != "" && ctx.failureMode == FailureHalts {
		panic(failureHalt)
	}
}

func registerReset(action func()) {
	ctx := mustGetCurrentContext()
	ctx.resets = append(ctx.resets, action)
}
