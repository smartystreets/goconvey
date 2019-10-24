package parser

import (
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"

	"github.com/smartystreets/goconvey/web/server/contract"
)

var (
	testNamePattern = regexp.MustCompile("^=== RUN:? +(.+)$")
)

func ParsePackageResults(result *contract.PackageResult, rawOutput string) {
	newOutputParser(result, rawOutput).parse()
}

type outputParser struct {
	raw    string
	lines  []string
	result *contract.PackageResult
	tests  []*contract.TestResult

	// place holders for loops
	line    string
	test    *contract.TestResult
	testMap map[string]*contract.TestResult
}

func newOutputParser(result *contract.PackageResult, rawOutput string) *outputParser {
	self := new(outputParser)
	self.raw = strings.TrimSpace(rawOutput)
	self.lines = strings.Split(self.raw, "\n")
	self.result = result
	self.tests = []*contract.TestResult{}
	self.testMap = make(map[string]*contract.TestResult)
	return self
}

func (o *outputParser) parse() {
	o.separateTestFunctionsAndMetadata()
	o.parseEachTestFunction()
}

func (o *outputParser) separateTestFunctionsAndMetadata() {
	for _, o.line = range o.lines {
		if o.processNonTestOutput() {
			break
		}
		o.processTestOutput()
	}
}
func (o *outputParser) processNonTestOutput() bool {
	if noGoFiles(o.line) {
		o.recordFinalOutcome(contract.NoGoFiles)

	} else if buildFailed(o.line) {
		o.recordFinalOutcome(contract.BuildFailure)

	} else if noTestFiles(o.line) {
		o.recordFinalOutcome(contract.NoTestFiles)

	} else if noTestFunctions(o.line) {
		o.recordFinalOutcome(contract.NoTestFunctions)

	} else {
		return false
	}
	return true
}

func (o *outputParser) recordFinalOutcome(outcome string) {
	o.result.Outcome = outcome
	o.result.BuildOutput = strings.Join(o.lines, "\n")
}

func (o *outputParser) processTestOutput() {
	o.line = strings.TrimSpace(o.line)
	if isNewTest(o.line) {
		o.registerTestFunction()

	} else if isTestResult(o.line) {
		o.recordTestMetadata()

	} else if isPackageReport(o.line) {
		o.recordPackageMetadata()

	} else {
		o.saveLineForParsingLater()

	}
}

func (o *outputParser) registerTestFunction() {
	testNameReg := testNamePattern.FindStringSubmatch(o.line)
	if len(testNameReg) < 2 { // Test-related lines that aren't about a new test
		return
	}
	o.test = contract.NewTestResult(testNameReg[1])
	o.tests = append(o.tests, o.test)
	o.testMap[o.test.TestName] = o.test
}
func (o *outputParser) recordTestMetadata() {
	testName := strings.Split(o.line, " ")[2]
	if test, ok := o.testMap[testName]; ok {
		o.test = test
		o.test.Passed = !strings.HasPrefix(o.line, "--- FAIL: ")
		o.test.Skipped = strings.HasPrefix(o.line, "--- SKIP: ")
		o.test.Elapsed = parseTestFunctionDuration(o.line)
	}
}
func (o *outputParser) recordPackageMetadata() {
	if packageFailed(o.line) {
		o.recordTestingOutcome(contract.Failed)

	} else if packagePassed(o.line) {
		o.recordTestingOutcome(contract.Passed)

	} else if isCoverageSummary(o.line) {
		o.recordCoverageSummary(o.line)
	}
}
func (o *outputParser) recordTestingOutcome(outcome string) {
	o.result.Outcome = outcome
	fields := strings.Split(o.line, "\t")
	o.result.PackageName = strings.TrimSpace(fields[1])
	o.result.Elapsed = parseDurationInSeconds(fields[2], 3)
}
func (o *outputParser) recordCoverageSummary(summary string) {
	start := len("coverage: ")
	end := strings.Index(summary, "%")
	value := summary[start:end]
	parsed, err := strconv.ParseFloat(value, 64)
	if err != nil {
		o.result.Coverage = -1
	} else {
		o.result.Coverage = parsed
	}
}
func (o *outputParser) saveLineForParsingLater() {
	o.line = strings.TrimLeft(o.line, "\t")
	if o.test == nil {
		fmt.Println("Potential error parsing output of", o.result.PackageName, "; couldn't handle this stray line:", o.line)
		return
	}
	o.test.RawLines = append(o.test.RawLines, o.line)
}

// TestResults is a collection of TestResults that implements sort.Interface.
type TestResults []contract.TestResult

func (r TestResults) Len() int {
	return len(r)
}

// Less compares TestResults on TestName
func (r TestResults) Less(i, j int) bool {
	return r[i].TestName < r[j].TestName
}

func (r TestResults) Swap(i, j int) {
	r[i], r[j] = r[j], r[i]
}

func (o *outputParser) parseEachTestFunction() {
	for _, o.test = range o.tests {
		o.test = parseTestOutput(o.test)
		if o.test.Error != "" {
			o.result.Outcome = contract.Panicked
		}
		o.test.RawLines = []string{}
		o.result.TestResults = append(o.result.TestResults, *o.test)
	}
	sort.Sort(TestResults(o.result.TestResults))
}
