package parser

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	"github.com/smartystreets/goconvey/convey/reporting"
	"github.com/smartystreets/goconvey/web/server/contract"
)

type testParser struct {
	test       *contract.TestResult
	line       string
	index      int
	inJson     bool
	jsonLines  []string
	otherLines []string
}

func parseTestOutput(test *contract.TestResult) *contract.TestResult {
	parser := newTestParser(test)
	parser.parseTestFunctionOutput()
	return test
}

func newTestParser(test *contract.TestResult) *testParser {
	self := new(testParser)
	self.test = test
	return self
}

func (t *testParser) parseTestFunctionOutput() {
	if len(t.test.RawLines) > 0 {
		t.processLines()
		t.deserializeJson()
		t.composeCapturedOutput()
	}
}

func (t *testParser) processLines() {
	for t.index, t.line = range t.test.RawLines {
		if !t.processLine() {
			break
		}
	}
}

func (t *testParser) processLine() bool {
	if strings.HasSuffix(t.line, reporting.OpenJson) {
		t.inJson = true
		t.accountForOutputWithoutNewline()

	} else if t.line == reporting.CloseJson {
		t.inJson = false

	} else if t.inJson {
		t.jsonLines = append(t.jsonLines, t.line)

	} else if isPanic(t.line) {
		t.parsePanicOutput()
		return false

	} else if isGoTestLogOutput(t.line) {
		t.parseLogLocation()

	} else {
		t.otherLines = append(t.otherLines, t.line)
	}
	return true
}

// If fmt.Print(f) produces output with no \n and that output
// is that last output before the framework spits out json
// (which starts with ''>>>>>'') then without this code
// all of the json is counted as output, not as json to be
// parsed and displayed by the web UI.
func (t *testParser) accountForOutputWithoutNewline() {
	prefix := strings.Split(t.line, reporting.OpenJson)[0]
	if prefix != "" {
		t.otherLines = append(t.otherLines, prefix)
	}
}

func (t *testParser) deserializeJson() {
	formatted := createArrayForJsonItems(t.jsonLines)
	var scopes []reporting.ScopeResult
	err := json.Unmarshal(formatted, &scopes)
	if err != nil {
		panic(fmt.Sprintf(bugReportRequest, err, formatted))
	}
	t.test.Stories = scopes
}
func (t *testParser) parsePanicOutput() {
	for index, line := range t.test.RawLines[t.index:] {
		t.parsePanicLocation(index, line)
		t.preserveStackTraceIndentation(index, line)
	}
	t.test.Error = strings.Join(t.test.RawLines, "\n")
}
func (t *testParser) parsePanicLocation(index int, line string) {
	if !panicLineHasMetadata(line) {
		return
	}
	metaLine := t.test.RawLines[index+4]
	fields := strings.Split(metaLine, " ")
	fileAndLine := strings.Split(fields[0], ":")
	t.test.File = fileAndLine[0]
	if len(fileAndLine) >= 2 {
		t.test.Line, _ = strconv.Atoi(fileAndLine[1])
	}
}
func (t *testParser) preserveStackTraceIndentation(index int, line string) {
	if panicLineShouldBeIndented(index, line) {
		t.test.RawLines[index] = "\t" + line
	}
}
func (t *testParser) parseLogLocation() {
	t.otherLines = append(t.otherLines, t.line)
	lineFields := strings.TrimSpace(t.line)
	if strings.HasPrefix(lineFields, "Error Trace:") {
		lineFields = strings.TrimPrefix(lineFields, "Error Trace:")
	}
	fields := strings.Split(lineFields, ":")
	t.test.File = strings.TrimSpace(fields[0])
	t.test.Line, _ = strconv.Atoi(fields[1])
}

func (t *testParser) composeCapturedOutput() {
	t.test.Message = strings.Join(t.otherLines, "\n")
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

func isGoTestLogOutput(line string) bool {
	return strings.Count(line, ":") == 2
}

func isPanic(line string) bool {
	return strings.HasPrefix(line, "panic: ")
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
