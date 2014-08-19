// Copyright 2011 Aaron Jacobs. All Rights Reserved.
// Author: aaronjjacobs@gmail.com (Aaron Jacobs)
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package ogletest

import (
	"errors"
	. "github.com/jacobsa/oglematchers"
	"testing"
)

////////////////////////////////////////////////////////////////////////
// Helpers
////////////////////////////////////////////////////////////////////////

// Set up a new test state with empty fields.
func setUpCurrentTest() {
	currentlyRunningTest = newTestInfo()
}

type fakeExpectThatMatcher struct {
	desc string
	err  error
}

func (m *fakeExpectThatMatcher) Matches(c interface{}) error {
	return m.err
}

func (m *fakeExpectThatMatcher) Description() string {
	return m.desc
}

func assertEqInt(t *testing.T, e, c int) {
	if e != c {
		t.Fatalf("Expected %d, got %d", e, c)
	}
}

func expectEqInt(t *testing.T, e, c int) {
	if e != c {
		t.Errorf("Expected %v, got %v", e, c)
	}
}

func expectEqStr(t *testing.T, e, c string) {
	if e != c {
		t.Errorf("Expected %s, got %s", e, c)
	}
}

////////////////////////////////////////////////////////////////////////
// Tests
////////////////////////////////////////////////////////////////////////

func TestNoCurrentTest(t *testing.T) {
	panicked := false

	defer func() {
		if !panicked {
			t.Errorf("Expected panic; got none.")
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()

	currentlyRunningTest = nil
	ExpectThat(17, Equals(17))
}

func TestNoFailure(t *testing.T) {
	setUpCurrentTest()
	matcher := &fakeExpectThatMatcher{"", nil}
	ExpectThat(17, matcher)

	assertEqInt(t, 0, len(currentlyRunningTest.failureRecords))
}

func TestInvalidFormatString(t *testing.T) {
	panicked := false

	defer func() {
		if !panicked {
			t.Errorf("Expected panic; got none.")
		}
	}()

	defer func() {
		if r := recover(); r != nil {
			panicked = true
		}
	}()

	setUpCurrentTest()
	matcher := &fakeExpectThatMatcher{"", errors.New("")}
	ExpectThat(17, matcher, 19, "blah")
}

func TestNoMatchWithoutErrorText(t *testing.T) {
	setUpCurrentTest()
	matcher := &fakeExpectThatMatcher{"taco", errors.New("")}
	ExpectThat(17, matcher)

	assertEqInt(t, 1, len(currentlyRunningTest.failureRecords))

	record := currentlyRunningTest.failureRecords[0]
	expectEqStr(t, "expect_that_test.go", record.FileName)
	expectEqInt(t, 118, record.LineNumber)
	expectEqStr(t, "Expected: taco\nActual:   17", record.GeneratedError)
	expectEqStr(t, "", record.UserError)
}

func TestNoMatchWithErrorTExt(t *testing.T) {
	setUpCurrentTest()
	matcher := &fakeExpectThatMatcher{"taco", errors.New("which is foo")}
	ExpectThat(17, matcher)

	assertEqInt(t, 1, len(currentlyRunningTest.failureRecords))
	record := currentlyRunningTest.failureRecords[0]

	expectEqStr(t, "Expected: taco\nActual:   17, which is foo", record.GeneratedError)
}

func TestFailureWithUserMessage(t *testing.T) {
	setUpCurrentTest()
	matcher := &fakeExpectThatMatcher{"taco", errors.New("")}
	ExpectThat(17, matcher, "Asd: %d %s", 19, "taco")

	assertEqInt(t, 1, len(currentlyRunningTest.failureRecords))
	record := currentlyRunningTest.failureRecords[0]

	expectEqStr(t, "Asd: 19 taco", record.UserError)
}

func TestAdditionalFailure(t *testing.T) {
	setUpCurrentTest()
	matcher := &fakeExpectThatMatcher{"", errors.New("")}

	// Fail twice.
	ExpectThat(17, matcher, "taco")
	ExpectThat(19, matcher, "burrito")

	assertEqInt(t, 2, len(currentlyRunningTest.failureRecords))
	record1 := currentlyRunningTest.failureRecords[0]
	record2 := currentlyRunningTest.failureRecords[1]

	expectEqStr(t, "taco", record1.UserError)
	expectEqStr(t, "burrito", record2.UserError)
}
