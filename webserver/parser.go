package main

import (
	"encoding/json"
	"fmt"
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
	return (strings.HasPrefix(line, "FAIL") ||
		strings.HasPrefix(line, "exit status") ||
		strings.HasPrefix(line, "PASS") ||
		strings.HasPrefix(line, "ok  \t"))
}

func (self *outputParser) registerTestFunction() {
	self.test = &TestResult{}
	self.test.Stories = []reporting.ScopeResult{}
	self.test.rawLines = []string{}
	self.test.TestName = self.line[len("=== RUN "):]
	self.tests = append(self.tests, self.test)
}
func (self *outputParser) recordTestMetadata() {
	self.test.Passed = strings.HasPrefix(self.line, "--- PASS: ")
	self.test.Elapsed = parseTestFunctionDuration(self.line)
}
func (self *outputParser) saveLineForParsingLater() {
	self.line = strings.TrimSpace(self.line)
	self.test.rawLines = append(self.test.rawLines, self.line)
}
func (self *outputParser) recordPackageMetadata() {
	if strings.HasPrefix(self.line, "FAIL\t") {
		self.parseLastLine()
		self.result.Passed = false
	} else if strings.HasPrefix(self.line, "ok  \t") {
		self.parseLastLine()
		self.result.Passed = true
	}
}
func (self *outputParser) parseLastLine() {
	fields := strings.Split(self.line, "\t")
	self.result.PackageName = strings.TrimSpace(fields[1])
	self.result.Elapsed = parseDurationInSeconds(fields[2], 3)
}

func (self *outputParser) parseTestFunctions() {
	for _, self.test = range self.tests {
		if len(self.test.rawLines) == 0 {
			continue
		} else if isJson(self.test.rawLines[0]) {
			self.deserializeScopes()
		} else {
			self.parseGoTestMessage()
		}
	}
}
func isJson(line string) bool {
	return strings.HasPrefix(line, "{")
}
func (self *outputParser) deserializeScopes() {
	// TODO: clean up
	rawJson := strings.Join(self.test.rawLines, "")
	var scopes []reporting.ScopeResult
	if strings.HasSuffix(rawJson, ",") { // Shouldn't need this...
		rawJson = rawJson[:len(rawJson)-1]
	}
	rawJson = "[" + rawJson + "]"
	err := json.Unmarshal([]byte(rawJson), &scopes)
	if err != nil {
		fmt.Println(err) // panic?
	}
	self.test.Stories = scopes
}
func (self *outputParser) parseGoTestMessage() {
	lineFields := self.test.rawLines[0]
	fields := strings.Split(lineFields, ":")
	self.test.File = strings.TrimSpace(fields[0])
	self.test.Line, _ = strconv.Atoi(fields[1])
	self.test.Message = strings.TrimSpace(fields[2])
	if len(self.test.rawLines) > 1 {
		additionalLines := strings.Join(self.test.rawLines[1:], "\n")
		self.test.Message = self.test.Message + "\n" + additionalLines
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
