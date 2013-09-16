package assertions

import (
	"fmt"
	"testing"
	"time"
)

func TestShouldHappenBefore(t *testing.T) {
	fail(t, so(0, ShouldHappenBefore), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so(0, ShouldHappenBefore, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	fail(t, so(0, ShouldHappenBefore, 1), shouldUseTimes)
	fail(t, so(0, ShouldHappenBefore, time.Now()), shouldUseTimes)
	fail(t, so(time.Now(), ShouldHappenBefore, 0), shouldUseTimes)

	fail(t, so(january3, ShouldHappenBefore, january1), fmt.Sprintf("Expected '%s' to happen before '%s' (it happened '48h0m0s' after)!", pretty(january3), pretty(january1)))
	fail(t, so(january3, ShouldHappenBefore, january3), fmt.Sprintf("Expected '%s' to happen before '%s' (it happened '0' after)!", pretty(january3), pretty(january3)))
	pass(t, so(january1, ShouldHappenBefore, january3))
}

func TestShouldHappenOnOrBefore(t *testing.T) {
	fail(t, so(0, ShouldHappenOnOrBefore), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so(0, ShouldHappenOnOrBefore, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	fail(t, so(0, ShouldHappenOnOrBefore, 1), shouldUseTimes)
	fail(t, so(0, ShouldHappenOnOrBefore, time.Now()), shouldUseTimes)
	fail(t, so(time.Now(), ShouldHappenOnOrBefore, 0), shouldUseTimes)

	// TODO: test actual logic
}

func TestShouldHappenAfter(t *testing.T) {
	fail(t, so(0, ShouldHappenAfter), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so(0, ShouldHappenAfter, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	fail(t, so(0, ShouldHappenAfter, 1), shouldUseTimes)
	fail(t, so(0, ShouldHappenAfter, time.Now()), shouldUseTimes)
	fail(t, so(time.Now(), ShouldHappenAfter, 0), shouldUseTimes)

	// TODO: test actual logic
}

func TestShouldHappenOnOrAfter(t *testing.T) {
	fail(t, so(0, ShouldHappenOnOrAfter), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so(0, ShouldHappenOnOrAfter, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	fail(t, so(0, ShouldHappenOnOrAfter, 1), shouldUseTimes)
	fail(t, so(0, ShouldHappenOnOrAfter, time.Now()), shouldUseTimes)
	fail(t, so(time.Now(), ShouldHappenOnOrAfter, 0), shouldUseTimes)

	// TODO: test actual logic
}

func TestShouldHappenBetween(t *testing.T) {
	fail(t, so(0, ShouldHappenBetween), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so(0, ShouldHappenBetween, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	fail(t, so(0, ShouldHappenBetween, 1), shouldUseTimes)
	fail(t, so(0, ShouldHappenBetween, time.Now()), shouldUseTimes)
	fail(t, so(time.Now(), ShouldHappenBetween, 0), shouldUseTimes)

	// TODO: test actual logic
}

func TestShouldHappenOnOrBetween(t *testing.T) {
	fail(t, so(0, ShouldHappenOnOrBetween), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so(0, ShouldHappenOnOrBetween, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	fail(t, so(0, ShouldHappenOnOrBetween, 1), shouldUseTimes)
	fail(t, so(0, ShouldHappenOnOrBetween, time.Now()), shouldUseTimes)
	fail(t, so(time.Now(), ShouldHappenOnOrBetween, 0), shouldUseTimes)

	// TODO: test actual logic
}

func TestShouldHappenWithin(t *testing.T) {
	fail(t, so(0, ShouldHappenWithin), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so(0, ShouldHappenWithin, 1, 2, 3), "This assertion requires exactly 1 comparison values (you provided 3).")

	fail(t, so(0, ShouldHappenWithin, 1), shouldUseTimes)
	fail(t, so(0, ShouldHappenWithin, time.Now()), shouldUseTimes)
	fail(t, so(time.Now(), ShouldHappenWithin, 0), shouldUseTimes)

	// TODO: test actual logic
}

const layout = "2006-01-02 15:04"

var january1, _ = time.Parse(layout, "2013-01-01 00:00")
var january2, _ = time.Parse(layout, "2013-01-02 00:00")
var january3, _ = time.Parse(layout, "2013-01-03 00:00")

var february1, _ = time.Parse(layout, "2013-02-01 00:00")
var february2, _ = time.Parse(layout, "2013-02-02 00:00")
var february3, _ = time.Parse(layout, "2013-02-03 00:00")

var march1, _ = time.Parse(layout, "2013-03-01 00:00")
var march2, _ = time.Parse(layout, "2013-03-02 00:00")
var march3, _ = time.Parse(layout, "2013-03-03 00:00")

func pretty(t time.Time) string {
	return fmt.Sprintf("%v", t)
}
