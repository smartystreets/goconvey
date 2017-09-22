package convey

import "github.com/smartystreets/assertions"

var (
	ShouldEqual          assertion = assertions.ShouldEqual
	ShouldNotEqual       assertion = assertions.ShouldNotEqual
	ShouldAlmostEqual    assertion = assertions.ShouldAlmostEqual
	ShouldNotAlmostEqual assertion = assertions.ShouldNotAlmostEqual
	ShouldResemble       assertion = assertions.ShouldResemble
	ShouldNotResemble    assertion = assertions.ShouldNotResemble
	ShouldPointTo        assertion = assertions.ShouldPointTo
	ShouldNotPointTo     assertion = assertions.ShouldNotPointTo
	ShouldBeNil          assertion = assertions.ShouldBeNil
	ShouldNotBeNil       assertion = assertions.ShouldNotBeNil
	ShouldBeTrue         assertion = assertions.ShouldBeTrue
	ShouldBeFalse        assertion = assertions.ShouldBeFalse
	ShouldBeZeroValue    assertion = assertions.ShouldBeZeroValue

	ShouldBeGreaterThan          assertion = assertions.ShouldBeGreaterThan
	ShouldBeGreaterThanOrEqualTo assertion = assertions.ShouldBeGreaterThanOrEqualTo
	ShouldBeLessThan             assertion = assertions.ShouldBeLessThan
	ShouldBeLessThanOrEqualTo    assertion = assertions.ShouldBeLessThanOrEqualTo
	ShouldBeBetween              assertion = assertions.ShouldBeBetween
	ShouldNotBeBetween           assertion = assertions.ShouldNotBeBetween
	ShouldBeBetweenOrEqual       assertion = assertions.ShouldBeBetweenOrEqual
	ShouldNotBeBetweenOrEqual    assertion = assertions.ShouldNotBeBetweenOrEqual

	ShouldContain       assertion = assertions.ShouldContain
	ShouldNotContain    assertion = assertions.ShouldNotContain
	ShouldContainKey    assertion = assertions.ShouldContainKey
	ShouldNotContainKey assertion = assertions.ShouldNotContainKey
	ShouldBeIn          assertion = assertions.ShouldBeIn
	ShouldNotBeIn       assertion = assertions.ShouldNotBeIn
	ShouldBeEmpty       assertion = assertions.ShouldBeEmpty
	ShouldNotBeEmpty    assertion = assertions.ShouldNotBeEmpty
	ShouldHaveLength    assertion = assertions.ShouldHaveLength

	ShouldStartWith           assertion = assertions.ShouldStartWith
	ShouldNotStartWith        assertion = assertions.ShouldNotStartWith
	ShouldEndWith             assertion = assertions.ShouldEndWith
	ShouldNotEndWith          assertion = assertions.ShouldNotEndWith
	ShouldBeBlank             assertion = assertions.ShouldBeBlank
	ShouldNotBeBlank          assertion = assertions.ShouldNotBeBlank
	ShouldContainSubstring    assertion = assertions.ShouldContainSubstring
	ShouldNotContainSubstring assertion = assertions.ShouldNotContainSubstring

	ShouldPanic        assertion = assertions.ShouldPanic
	ShouldNotPanic     assertion = assertions.ShouldNotPanic
	ShouldPanicWith    assertion = assertions.ShouldPanicWith
	ShouldNotPanicWith assertion = assertions.ShouldNotPanicWith

	ShouldHaveSameTypeAs    assertion = assertions.ShouldHaveSameTypeAs
	ShouldNotHaveSameTypeAs assertion = assertions.ShouldNotHaveSameTypeAs
	ShouldImplement         assertion = assertions.ShouldImplement
	ShouldNotImplement      assertion = assertions.ShouldNotImplement

	ShouldHappenBefore         assertion = assertions.ShouldHappenBefore
	ShouldHappenOnOrBefore     assertion = assertions.ShouldHappenOnOrBefore
	ShouldHappenAfter          assertion = assertions.ShouldHappenAfter
	ShouldHappenOnOrAfter      assertion = assertions.ShouldHappenOnOrAfter
	ShouldHappenBetween        assertion = assertions.ShouldHappenBetween
	ShouldHappenOnOrBetween    assertion = assertions.ShouldHappenOnOrBetween
	ShouldNotHappenOnOrBetween assertion = assertions.ShouldNotHappenOnOrBetween
	ShouldHappenWithin         assertion = assertions.ShouldHappenWithin
	ShouldNotHappenWithin      assertion = assertions.ShouldNotHappenWithin
	ShouldBeChronological      assertion = assertions.ShouldBeChronological

	ShouldBeError assertion = assertions.ShouldBeError
)
