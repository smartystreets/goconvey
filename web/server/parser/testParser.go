package parser

import (
	"encoding/json"
	"fmt"
	"github.com/smartystreets/goconvey/reporting"
	"github.com/smartystreets/goconvey/web/server/results"
	"strconv"
	"strings"
)

func parseTestOutput(test *results.TestResult) *results.TestResult {
	parser := newTestParser(test)
	parser.parse()
	return test
}

func newTestParser(test *results.TestResult) *testParser {
	self := &testParser{}
	self.test = test
	return self
}

func (self *testParser) parse() {
	if len(self.test.RawLines) == 0 {
		return
	} else if isJson(self.test.RawLines[0]) {
		self.deserializeScopes()
	} else {
		self.parseAdditionalGoTestOutput()
	}
}

func (self *testParser) deserializeScopes() {
	formatted := createArrayForJsonItems(self.test.RawLines)
	var scopes []reporting.ScopeResult
	err := json.Unmarshal(formatted, &scopes)
	if err != nil {
		panic(fmt.Sprintf(bugReportRequest, err, formatted))
	}
	self.test.Stories = scopes
}
func (self *testParser) parseAdditionalGoTestOutput() {
	if strings.HasPrefix(self.test.RawLines[0], "panic: ") {
		self.parsePanicOutput()
	} else {
		self.parseLoggedOutput()
		self.compileCompleteMessage()
	}
}
func (self *testParser) parsePanicOutput() {
	for index, line := range self.test.RawLines {
		self.parsePanicMetadata(index, line)
		self.preserveStackTraceIndentation(index, line)
	}
	self.test.Error = strings.Join(self.test.RawLines, "\n")
}
func (self *testParser) parsePanicMetadata(index int, line string) {
	if !panicLineHasMetadata(line) {
		return
	}
	metaLine := self.test.RawLines[index+4]
	fields := strings.Split(metaLine, " ")
	fileAndLine := strings.Split(fields[0], ":")
	self.test.File = fileAndLine[0]
	self.test.Line, _ = strconv.Atoi(fileAndLine[1])
}
func (self *testParser) preserveStackTraceIndentation(index int, line string) {
	if panicLineShouldBeIndented(index, line) {
		self.test.RawLines[index] = "\t" + line
	}
}
func (self *testParser) parseLoggedOutput() {
	lineFields := self.test.RawLines[0]
	fields := strings.Split(lineFields, ":")
	self.test.File = strings.TrimSpace(fields[0])
	self.test.Line, _ = strconv.Atoi(fields[1])
	self.test.Message = strings.TrimSpace(fields[2])
}
func (self *testParser) compileCompleteMessage() {
	if len(self.test.RawLines) > 1 {
		additionalLines := strings.Join(self.test.RawLines[1:], "\n")
		self.test.Message = self.test.Message + "\n" + additionalLines
	}
}

type testParser struct {
	test *results.TestResult
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
