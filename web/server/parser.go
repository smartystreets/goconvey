package main

import (
	"encoding/json"
	"fmt"
	"github.com/smartystreets/goconvey/reporting"
	"strconv"
	"strings"
)

func parsePackageResults(packageName, raw string) *PackageResult {
	parser := newOutputParser(packageName, raw)
	return parser.Parse()
}

func newOutputParser(packageName, raw string) *outputParser {
	self := &outputParser{}
	self.raw = strings.TrimSpace(raw)
	self.lines = strings.Split(self.raw, "\n")
	self.result = &PackageResult{}
	self.result.PackageName = packageName
	self.tests = []*TestResult{}
	self.result.TestResults = []TestResult{}
	return self
}

func (self *outputParser) Parse() *PackageResult {
	self.separateTestFunctionsAndMetadata()
	self.parseEachTestFunction()
	self.attachParsedTestFunctionsToResult()
	return self.result
}

func (self *outputParser) separateTestFunctionsAndMetadata() {
	for _, self.line = range self.lines {
		if self.processNonTestOutput() {
			break
		}
		self.processTestOutput()
	}
}
func (self *outputParser) processNonTestOutput() bool {
	if noGoFiles(self.line) {
		self.recordOutcome(noGo)

	} else if buildFailed(self.line) {
		self.recordOutcome(buildFailure)

	} else if noTestFiles(self.line) {
		self.recordOutcome(noTestFile)

	} else if noTestFunctions(self.line) {
		self.recordOutcome(noTestFunction)

	} else {
		return false
	}
	return true
}

func (self *outputParser) recordOutcome(outcome string) {
	self.result.Outcome = outcome
	self.result.BuildOutput = strings.Join(self.lines, "\n")
}

func (self *outputParser) processTestOutput() {
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
func (self *outputParser) recordPackageMetadata() {
	if strings.HasPrefix(self.line, "FAIL\t") {
		self.parseLastLine()
		self.result.Outcome = failed
	} else if strings.HasPrefix(self.line, "ok  \t") {
		self.parseLastLine()
		self.result.Outcome = passed
	}
}
func (self *outputParser) parseLastLine() {
	fields := strings.Split(self.line, "\t")
	self.result.PackageName = strings.TrimSpace(fields[1])
	self.result.Elapsed = parseDurationInSeconds(fields[2], 3)
}
func (self *outputParser) saveLineForParsingLater() {
	self.line = strings.TrimSpace(self.line)
	if self.test == nil {
		fmt.Println("LINE:", self.line)
		return
	}
	self.test.rawLines = append(self.test.rawLines, self.line)
}

func (self *outputParser) parseEachTestFunction() {
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
	// TODO: clean up!
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
	// TODO: clean up!
	if strings.HasPrefix(self.test.rawLines[0], "panic: ") {
		self.result.Outcome = panicked
		for i, line := range self.test.rawLines {
			if strings.HasPrefix(line, "goroutine") && strings.Contains(line, "[running]") {
				metaLine := self.test.rawLines[i+4]
				fields := strings.Split(metaLine, " ")
				fileAndLine := strings.Split(fields[0], ":")
				self.test.File = fileAndLine[0]
				self.test.Line, _ = strconv.Atoi(fileAndLine[1])
			}
			if strings.Contains(line, "+") || (i > 0 && strings.Contains(line, "panic: ")) {
				self.test.rawLines[i] = "\t" + line
			}
		}
		self.test.Error = strings.Join(self.test.rawLines, "\n")
	} else {
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
}

func (self *outputParser) attachParsedTestFunctionsToResult() {
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
	Outcome     string
	BuildOutput string
	TestResults []TestResult
}

var (
	passed         = "passed"
	failed         = "failed"
	panicked       = "panicked"
	buildFailure   = "build failure"
	noTestFile     = "no test files"
	noTestFunction = "no test functions"
	noGo           = "no go code"
)

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

func noGoFiles(line string) bool {
	return strings.HasPrefix(line, "can't load package: ") &&
		strings.Contains(line, ": no Go source files in ")
}
func buildFailed(line string) bool {
	return strings.HasPrefix(line, "can't load package: ") &&
		!strings.Contains(line, ": no Go source files in ")
}
func noTestFunctions(line string) bool {
	return line == "testing: warning: no tests to run"
}
func noTestFiles(line string) bool {
	return strings.HasPrefix(line, "?") && strings.Contains(line, "[no test files]")
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
