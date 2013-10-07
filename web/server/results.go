package main

import "github.com/smartystreets/goconvey/reporting"

var ( // PackageResult.Outcome values:
	passed         = "passed"
	failed         = "failed"
	panicked       = "panicked"
	buildFailure   = "build failure"
	noTestFile     = "no test files"
	noTestFunction = "no test functions"
	noGo           = "no go code"
)

type PackageResult struct {
	PackageName string
	Elapsed     float64
	Outcome     string
	BuildOutput string
	TestResults []TestResult
}

func newPackageResult(packageName string) *PackageResult {
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

	rawLines []string
}

func newTestResult(testName string) *TestResult {
	self := &TestResult{}
	self.Stories = []reporting.ScopeResult{}
	self.rawLines = []string{}
	self.TestName = testName
	return self
}
