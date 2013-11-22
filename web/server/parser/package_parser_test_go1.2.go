// +build go1.2

package parser

import (
	"github.com/smartystreets/goconvey/web/server/contract"
	"testing"
)

func TestParsePackage_OldSchoolWithSuccessAndBogusCoverage_ReturnsCompletePackageResult(t *testing.T) {
	actual := &contract.PackageResult{PackageName: expectedOldSchool_PassesButCoverageIsBogus.PackageName}
	ParsePackageResults(actual, inputOldSchool_PassesButCoverageIsBogus)
	assertEqual(t, expectedOldSchool_PassesButCoverageIsBogus, *actual)
}
