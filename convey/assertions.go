package convey

import (
	"github.com/smartystreets/goconvey/assertions"
)

var (
	ShouldEqual       = assertions.ShouldEqual
	ShouldNotEqual    = assertions.ShouldNotEqual
	ShouldResemble    = assertions.ShouldResemble
	ShouldNotResemble = assertions.ShouldNotResemble
	ShouldBeNil       = assertions.ShouldBeNil
	ShouldNotBeNil    = assertions.ShouldNotBeNil
	ShouldBeTrue      = assertions.ShouldBeTrue
	ShouldBeFalse     = assertions.ShouldBeFalse

	ShouldBeGreaterThan          = assertions.ShouldBeGreaterThan
	ShouldBeGreaterThanOrEqualTo = assertions.ShouldBeGreaterThanOrEqualTo
	ShouldBeLessThan             = assertions.ShouldBeLessThan
	ShouldBeLessThanOrEqualTo    = assertions.ShouldBeLessThanOrEqualTo
	ShouldBeBetween              = assertions.ShouldBeBetween
)
