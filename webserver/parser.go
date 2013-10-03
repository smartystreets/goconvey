package main

import (
	_ "fmt"
	"github.com/smartystreets/goconvey/reporting"
	"strconv"
	"strings"
)

func parsePackageResults(raw string) *PackageResult {
	parser := newOutputParser(raw)
	return parser.Parse()
}

func newOutputParser(raw string) *outputParser {
	self := &outputParser{}
	self.raw = strings.TrimSpace(raw)
	self.lines = strings.Split(self.raw, "\n")
	self.result = &PackageResult{}
	self.tests = []*TestResult{}
	self.result.TestResults = []TestResult{}
	return self
}

func (self *outputParser) Parse() *PackageResult {
	self.gatherTestFunctionsAndMetadata()
	self.parseTestFunctions()
	self.attachTestFunctionsToResult()
	return self.result
}

func (self *outputParser) gatherTestFunctionsAndMetadata() {
	for _, self.line = range self.lines {
		self.processNextLine()
	}
}
func (self *outputParser) processNextLine() {
	if isNewTest(self.line) {
		self.registerTestFunction()

	} else if isTestResult(self.line) {
		self.recordTestMetadata()

	} else if isPackageReport(self.line) {
		self.recordPackageMetadata()

	} else {
		self.saveLineForParsingLater()
	}
}

func isNewTest(line string) bool {
	return strings.HasPrefix(line, "=== ")
}
func isTestResult(line string) bool {
	return strings.HasPrefix(line, "--- ")
}
func isPackageReport(line string) bool {
	// TODO: passing test: strings.HasPrefix(line, "PASS\nok  \t")
	return strings.HasPrefix(line, "FAIL") || strings.HasPrefix(line, "exit status")
}

func (self *outputParser) registerTestFunction() {
	self.test = &TestResult{}
	self.test.Stories = []reporting.ScopeResult{}
	self.test.rawLines = []string{}
	self.test.TestName = self.line[len("=== RUN "):]
	self.tests = append(self.tests, self.test)
}
func (self *outputParser) recordTestMetadata() {
	if strings.Contains(self.line, "--- PASS: ") {
		self.test.Passed = true
		// TODO: parse duration
	}
}
func (self *outputParser) saveLineForParsingLater() {
	self.line = strings.TrimSpace(self.line)
	self.test.rawLines = append(self.test.rawLines, self.line)
}
func (self *outputParser) recordPackageMetadata() {
	if strings.HasPrefix(self.line, "FAIL\t") {
		fields := strings.Split(self.line, "\t")
		self.result.PackageName = strings.TrimSpace(fields[1])
		self.result.Elapsed = parseDuration(fields[2], 3)
	}
}

func (self *outputParser) parseTestFunctions() {
	for _, test := range self.tests {
		if len(test.rawLines) > 0 {
			lineFields := test.rawLines[0]
			fields := strings.Split(lineFields, ":")
			test.File = strings.TrimSpace(fields[0])
			test.Line, _ = strconv.Atoi(fields[1])
			test.Message = strings.TrimSpace(fields[2])
			if len(test.rawLines) > 1 {
				test.Message = test.Message + "\n" + strings.Join(test.rawLines[1:], "\n")
			}
		}
	}
}

func (self *outputParser) attachTestFunctionsToResult() {
	for _, test := range self.tests {
		test.rawLines = []string{}
		self.result.TestResults = append(self.result.TestResults, *test)
	}
}

type outputParser struct {
	raw    string
	lines  []string
	result *PackageResult
	tests  []*TestResult

	// place holders for loops
	line string
	test *TestResult
}

type PackageResult struct {
	PackageName string
	Elapsed     float64
	Passed      bool
	TestResults []TestResult
}

type TestResult struct {
	TestName string
	Elapsed  float64
	Passed   bool
	File     string
	Line     int
	Message  string
	Stories  []reporting.ScopeResult

	rawLines []string
}
