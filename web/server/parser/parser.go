package parser

import (
	"fmt"
	"github.com/smartystreets/goconvey/web/server/contract"
	"strings"
)

func ParsePackageResults(result *contract.PackageResult, rawOutput string) {
	newOutputParser(result, rawOutput).parse()
}

func newOutputParser(result *contract.PackageResult, rawOutput string) *outputParser {
	self := &outputParser{}
	self.raw = strings.TrimSpace(rawOutput)
	self.lines = strings.Split(self.raw, "\n")
	self.result = result
	self.tests = []*contract.TestResult{}
	return self
}

func (self *outputParser) parse() {
	self.separateTestFunctionsAndMetadata()
	self.parseEachTestFunction()
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
		self.recordFinalOutcome(contract.NoGoFiles)

	} else if buildFailed(self.line) {
		self.recordFinalOutcome(contract.BuildFailure)

	} else if noTestFiles(self.line) {
		self.recordFinalOutcome(contract.NoTestFiles)

	} else if noTestFunctions(self.line) {
		self.recordFinalOutcome(contract.NoTestFunctions)

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
	self.test = contract.NewTestResult(self.line[len("=== RUN "):])
	self.tests = append(self.tests, self.test)
}
func (self *outputParser) recordTestMetadata() {
	self.test.Passed = strings.HasPrefix(self.line, "--- PASS: ")
	self.test.Elapsed = parseTestFunctionDuration(self.line)
}
func (self *outputParser) recordPackageMetadata() {
	if packageFailed(self.line) {
		self.recordTestingOutcome(contract.Failed)

	} else if packagePassed(self.line) {
		self.recordTestingOutcome(contract.Passed)
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
		fmt.Println("Potential parsing output of", self.result.PackageName, "; couldn't handle this stray line:", self.line)
		return
	}
	self.test.RawLines = append(self.test.RawLines, self.line)
}

func (self *outputParser) parseEachTestFunction() {
	for _, self.test = range self.tests {
		self.test = parseTestOutput(self.test)
		if self.test.Error != "" {
			self.result.Outcome = contract.Panicked
		}
		self.test.RawLines = []string{}
		self.result.TestResults = append(self.result.TestResults, *self.test)
	}
}

type outputParser struct {
	raw    string
	lines  []string
	result *contract.PackageResult
	tests  []*contract.TestResult

	// place holders for loops
	line string
	test *contract.TestResult
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
