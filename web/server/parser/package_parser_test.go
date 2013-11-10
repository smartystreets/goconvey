package parser

import (
	"encoding/json"
	"fmt"
	"github.com/smartystreets/goconvey/web/server/contract"
	"strings"
	"testing"
)

func TestParsePackage_NoGoFiles_ReturnsPackageResult(t *testing.T) {
	actual := &contract.PackageResult{PackageName: expected_NoGoFiles.PackageName}
	ParsePackageResults(actual, input_NoGoFiles)
	assertEqual(t, expected_NoGoFiles, *actual)
}

func TestParsePackage_NoTestFiles_ReturnsPackageResult(t *testing.T) {
	actual := &contract.PackageResult{PackageName: expected_NoTestFiles.PackageName}
	ParsePackageResults(actual, input_NoTestFiles)
	assertEqual(t, expected_NoTestFiles, *actual)
}

func TestParsePacakge_NoTestFunctions_ReturnsPackageResult(t *testing.T) {
	actual := &contract.PackageResult{PackageName: expected_NoTestFunctions.PackageName}
	ParsePackageResults(actual, input_NoTestFunctions)
	assertEqual(t, expected_NoTestFunctions, *actual)
}

func TestParsePackage_BuildFailed_ReturnsPackageResult(t *testing.T) {
	actual := &contract.PackageResult{PackageName: expected_BuildFailed_InvalidPackageDeclaration.PackageName}
	ParsePackageResults(actual, input_BuildFailed_InvalidPackageDeclaration)
	assertEqual(t, expected_BuildFailed_InvalidPackageDeclaration, *actual)

	actual = &contract.PackageResult{PackageName: expected_BuildFailed_OtherErrors.PackageName}
	ParsePackageResults(actual, input_BuildFailed_OtherErrors)
	assertEqual(t, expected_BuildFailed_OtherErrors, *actual)

	actual = &contract.PackageResult{PackageName: expected_BuildFailed_CantFindPackage.PackageName}
	ParsePackageResults(actual, input_BuildFailed_CantFindPackage)
	assertEqual(t, expected_BuildFailed_CantFindPackage, *actual)
}

func TestParsePackage_OldSchoolWithFailureOutput_ReturnsCompletePackageResult(t *testing.T) {
	actual := &contract.PackageResult{PackageName: expectedOldSchool_Fails.PackageName}
	ParsePackageResults(actual, inputOldSchool_Fails)
	assertEqual(t, expectedOldSchool_Fails, *actual)
}

func TestParsePackage_OldSchoolWithSuccessOutput_ReturnsCompletePackageResult(t *testing.T) {
	actual := &contract.PackageResult{PackageName: expectedOldSchool_Passes.PackageName}
	ParsePackageResults(actual, inputOldSchool_Passes)
	assertEqual(t, expectedOldSchool_Passes, *actual)
}

func TestParsePackage_OldSchoolWithPanicOutput_ReturnsCompletePackageResult(t *testing.T) {
	actual := &contract.PackageResult{PackageName: expectedOldSchool_Panics.PackageName}
	ParsePackageResults(actual, inputOldSchool_Panics)
	assertEqual(t, expectedOldSchool_Panics, *actual)
}

func TestParsePackage_GoConveyOutput_ReturnsCompletePackageResult(t *testing.T) {
	actual := &contract.PackageResult{PackageName: expectedGoConvey.PackageName}
	ParsePackageResults(actual, inputGoConvey)
	assertEqual(t, expectedGoConvey, *actual)
}

func TestParsePackage_ActualPackageNameDifferentThanDirectoryName_ReturnsActualPackageName(t *testing.T) {
	actual := &contract.PackageResult{PackageName: strings.Replace(expectedGoConvey.PackageName, "examples", "stuff", -1)}
	ParsePackageResults(actual, inputGoConvey)
	assertEqual(t, expectedGoConvey, *actual)
}

func TestParsePackage_GoConveyOutputMalformed_CausesPanic(t *testing.T) {
	defer func() {
		if r := recover(); r != nil {
			message := fmt.Sprintf("%v", r)
			if !strings.Contains(message, "bug report") {
				t.Errorf("Should have panicked with a request to file a bug report but we received this error instead: %s", message)
			}
		} else {
			t.Errorf("Should have panicked with a request to file a bug report but we received no error.")
		}
	}()

	actual := &contract.PackageResult{PackageName: expectedGoConvey.PackageName}
	ParsePackageResults(actual, inputGoConvey_Malformed)
}

func TestParsePackage_GoConveyWithRandomOutput_ReturnsPackageResult(t *testing.T) {
	actual := &contract.PackageResult{PackageName: expectedGoConvey_WithRandomOutput.PackageName}
	ParsePackageResults(actual, inputGoConvey_WithRandomOutput)
	assertEqual(t, expectedGoConvey_WithRandomOutput, *actual)
}

func assertEqual(t *testing.T, expected, actual interface{}) {
	a, _ := json.Marshal(expected)
	b, _ := json.Marshal(actual)
	if string(a) != string(b) {
		t.Errorf(failureTemplate, string(a), string(b))
	}
}

const failureTemplate = "Comparison failed:\n  Expected: %v\n    Actual: %v\n"

/*
Test output for these tests was generated from the following test code:

Old School style tests:

	package examples

	import "testing"

	func TestOldSchool_Passes(t *testing.T) {
		// passes implicitly
	}

	func TestOldSchool_PassesWithMessage(t *testing.T) {
		t.Log("I am a passing test.\nWith a newline.")
	}

	func TestOldSchool_Failure(t *testing.T) {
		t.Fail() // no message
	}

	func TestOldSchool_FailureWithReason(t *testing.T) {
		t.Error("I am a failing test.")
	}

GoConvey style tests:

	package examples

	import (
		. "github.com/smartystreets/goconvey/convey"
		"testing"
	)

	func TestPassingStory(t *testing.T) {
		Convey("A passing story", t, func() {
			So("This test passes", ShouldContainSubstring, "pass")
		})
	}

GoConvey style tests with random output:

	package examples

	import (
		"fmt"
		. "github.com/smartystreets/goconvey/convey"
		"testing"
	)

	func TestPassingStory(t *testing.T) {
		fmt.Println("*** Hello, World! (1) ***")

		Convey("A passing story", t, func() {
			fmt.Println("*** Hello, World! (2) ***")
			So("This test passes", ShouldContainSubstring, "pass")
			fmt.Println("*** Hello, World! (3) ***")
		})

		Convey("A passing story", t, func() {
			fmt.Println("*** Hello, World! (4)***")
			So("This test passes", ShouldContainSubstring, "pass")
			fmt.Println("*** Hello, World! (5) ***")
		})

		fmt.Println("*** Hello, World! (6) ***")
	}


*/
