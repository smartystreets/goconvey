package assertions

import "testing"

type TestObjectA interface {
	Foo()
}

type testObjectA struct {

}

func NewTestObjectA() *testObjectA {
	return new(testObjectA);
}

func (t *testObjectA) Foo() {

}

type TestObjectB interface {
	Bar()
}

type testObjectB struct {

}

func NewTestObjectB() *testObjectB {
	return new(testObjectB);
}

func (t *testObjectB) Bar() {

}

func TestShouldImplement(t *testing.T) {
	var i *TestObjectA = nil
	fail(t, so(NewTestObjectA(), ShouldImplement), "This assertion requires exactly 1 comparison values (you provided 0).")
	fail(t, so(NewTestObjectA(), ShouldImplement, i, i), "This assertion requires exactly 1 comparison values (you provided 2).")
	fail(t, so(NewTestObjectA(), ShouldImplement, i, i, i), "This assertion requires exactly 1 comparison values (you provided 3).")

	fail(t, so(NewTestObjectA(), ShouldImplement, "foo"), "This assertion requires a pointer with the interface type")
	fail(t, so(NewTestObjectA(), ShouldImplement, 1), "This assertion requires a pointer with the interface type")
	fail(t, so(NewTestObjectA(), ShouldImplement, nil), "This assertion requires a pointer with the interface type")

	fail(t, so(nil, ShouldImplement, i), "Expected '<nil>' implement 'assertions.TestObjectA' (but was: '<nil>')!")
	fail(t, so(1, ShouldImplement, i), "Expected '1' implement 'assertions.TestObjectA' (but was: 'int')!")
	fail(t, so(1, ShouldImplement, i), "Expected '1' implement 'assertions.TestObjectA' (but was: 'int')!")

	fail(t, so(NewTestObjectB(), ShouldImplement, i), "Expected '*assertions.testObjectB' implement 'assertions.TestObjectA' (but was: '*assertions.testObjectB')!")
	pass(t, so(NewTestObjectA(), ShouldImplement, i))
}
