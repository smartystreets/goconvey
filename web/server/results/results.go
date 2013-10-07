package results

import "github.com/smartystreets/goconvey/reporting"

type CompleteOutput struct {
	Packages []*PackageResult
	Revision string
}

var ( // PackageResult.Outcome values:
	Passed          = "passed"
	Failed          = "failed"
	Panicked        = "panicked"
	BuildFailure    = "build failure"
	NoTestFiles     = "no test files"
	NoTestFunctions = "no test functions"
	NoGoFiles       = "no go code"
)

type PackageResult struct {
	PackageName string
	Elapsed     float64
	Outcome     string
	BuildOutput string
	TestResults []TestResult
}

func NewPackageResult(packageName string) *PackageResult {
	self := &PackageResult{}
	self.PackageName = packageName
	self.TestResults = []TestResult{}
	return self
}

type TestResult struct {
	TestName string
	Elapsed  float64
	Passed   bool
	File     string
	Line     int
	Message  string
	Error    string
	Stories  []reporting.ScopeResult

	RawLines []string `json:",omitempty"`
}

func NewTestResult(testName string) *TestResult {
	self := &TestResult{}
	self.Stories = []reporting.ScopeResult{}
	self.RawLines = []string{}
	self.TestName = testName
	return self
}
