package parser

import (
	"encoding/json"
	"fmt"
	"github.com/smartystreets/goconvey/reporting"
	"github.com/smartystreets/goconvey/web/server/results"
	"strconv"
	"strings"
)

func ParsePackageResults(packageName, raw string) *results.PackageResult {
	parser := newOutputParser(packageName, raw)
	return parser.parse()
}

func newOutputParser(packageName, raw string) *outputParser {
	self := &outputParser{}
	self.raw = strings.TrimSpace(raw)
	self.lines = strings.Split(self.raw, "\n")
	self.result = results.NewPackageResult(packageName)
	self.tests = []*results.TestResult{}
	return self
}

func (self *outputParser) parse() *results.PackageResult {
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
		self.recordFinalOutcome(results.NoGoFiles)

	} else if buildFailed(self.line) {
		self.recordFinalOutcome(results.BuildFailure)

	} else if noTestFiles(self.line) {
		self.recordFinalOutcome(results.NoTestFiles)

	} else if noTestFunctions(self.line) {
		self.recordFinalOutcome(results.NoTestFunctions)

	} else {
		return false
	}
	return true
}

func (self *outputParser) recordFinalOutcome(outcome string) {
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
	self.test = results.NewTestResult(self.line[len("=== RUN "):])
	self.tests = append(self.tests, self.test)
}
func (self *outputParser) recordTestMetadata() {
	self.test.Passed = strings.HasPrefix(self.line, "--- PASS: ")
	self.test.Elapsed = parseTestFunctionDuration(self.line)
}
func (self *outputParser) recordPackageMetadata() {
	if packageFailed(self.line) {
		self.recordTestingOutcome(results.Failed)

	} else if packagePassed(self.line) {
		self.recordTestingOutcome(results.Passed)
	}
}
func (self *outputParser) recordTestingOutcome(outcome string) {
	self.result.Outcome = outcome
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
	self.test.RawLines = append(self.test.RawLines, self.line)
}

func (self *outputParser) parseEachTestFunction() {
	for _, self.test = range self.tests {
		if len(self.test.RawLines) == 0 {
			continue
		} else if isJson(self.test.RawLines[0]) {
			self.deserializeScopes()
		} else {
			self.parseAdditionalGoTestOutput()
		}
	}
}
func (self *outputParser) deserializeScopes() {
	formatted := createArrayForJsonItems(self.test.RawLines)
	var scopes []reporting.ScopeResult
	err := json.Unmarshal(formatted, &scopes)
	if err != nil {
		panic(fmt.Sprintf(bugReportRequest, err, formatted))
	}
	self.test.Stories = scopes
}
func (self *outputParser) parseAdditionalGoTestOutput() {
	if strings.HasPrefix(self.test.RawLines[0], "panic: ") {
		self.parsePanicOutput()
	} else {
		self.parseLoggedOutput()
		self.compileCompleteMessage()
	}
}
func (self *outputParser) parsePanicOutput() {
	self.result.Outcome = results.Panicked
	for index, line := range self.test.RawLines {
		self.parsePanicMetadata(index, line)
		self.preserveStackTraceIndentation(index, line)
	}
	self.test.Error = strings.Join(self.test.RawLines, "\n")
}
func (self *outputParser) parsePanicMetadata(index int, line string) {
	if !panicLineHasMetadata(line) {
		return
	}
	metaLine := self.test.RawLines[index+4]
	fields := strings.Split(metaLine, " ")
	fileAndLine := strings.Split(fields[0], ":")
	self.test.File = fileAndLine[0]
	self.test.Line, _ = strconv.Atoi(fileAndLine[1])
}
func (self *outputParser) preserveStackTraceIndentation(index int, line string) {
	if panicLineShouldBeIndented(index, line) {
		self.test.RawLines[index] = "\t" + line
	}
}
func (self *outputParser) parseLoggedOutput() {
	lineFields := self.test.RawLines[0]
	fields := strings.Split(lineFields, ":")
	self.test.File = strings.TrimSpace(fields[0])
	self.test.Line, _ = strconv.Atoi(fields[1])
	self.test.Message = strings.TrimSpace(fields[2])
}
func (self *outputParser) compileCompleteMessage() {
	if len(self.test.RawLines) > 1 {
		additionalLines := strings.Join(self.test.RawLines[1:], "\n")
		self.test.Message = self.test.Message + "\n" + additionalLines
	}
}

func (self *outputParser) attachParsedTestFunctionsToResult() {
	for _, test := range self.tests {
		test.RawLines = []string{}
		self.result.TestResults = append(self.result.TestResults, *test)
	}
}

type outputParser struct {
	raw    string
	lines  []string
	result *results.PackageResult
	tests  []*results.TestResult

	// place holders for loops
	line string
	test *results.TestResult
}

func noGoFiles(line string) bool {
	return strings.HasPrefix(line, "can't load package: ") &&
		strings.Contains(line, ": no Go source files in ")
}
func buildFailed(line string) bool {
	return strings.HasPrefix(line, "# ") ||
		(strings.HasPrefix(line, "can't load package: ") &&
			!strings.Contains(line, ": no Go source files in "))
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

func packageFailed(line string) bool {
	return strings.HasPrefix(line, "FAIL\t")
}
func packagePassed(line string) bool {
	return strings.HasPrefix(line, "ok  \t")
}

func isJson(line string) bool {
	return strings.HasPrefix(line, "{")
}
func createArrayForJsonItems(lines []string) []byte {
	jsonArrayItems := strings.Join(lines, "")
	jsonArrayItems = removeTrailingComma(jsonArrayItems)
	return []byte(fmt.Sprintf("[%s]\n", jsonArrayItems))
}
func removeTrailingComma(rawJson string) string {
	if trailingComma(rawJson) {
		return rawJson[:len(rawJson)-1]
	}
	return rawJson
}
func trailingComma(value string) bool {
	return strings.HasSuffix(value, ",")
}

func panicLineHasMetadata(line string) bool {
	return strings.HasPrefix(line, "goroutine") && strings.Contains(line, "[running]")
}
func panicLineShouldBeIndented(index int, line string) bool {
	return strings.Contains(line, "+") || (index > 0 && strings.Contains(line, "panic: "))
}

const bugReportRequest = `
Uh-oh! Looks like something went wrong. Please copy the following text and file a bug report at: 

https://github.com/smartystreets/goconvey/issues?state=open

======= BEGIN BUG REPORT =======

ERROR: %v

OUTPUT: %s

======= END BUG REPORT =======

`
