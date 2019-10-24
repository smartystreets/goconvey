// TODO: under unit test

package reporting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type JsonReporter struct {
	out        *Printer
	currentKey []string
	current    *ScopeResult
	index      map[string]*ScopeResult
	scopes     []*ScopeResult
}

func (j *JsonReporter) depth() int { return len(j.currentKey) }

func (j *JsonReporter) BeginStory(story *StoryReport) {}

func (j *JsonReporter) Enter(scope *ScopeReport) {
	j.currentKey = append(j.currentKey, scope.Title)
	ID := strings.Join(j.currentKey, "|")
	if _, found := j.index[ID]; !found {
		next := newScopeResult(scope.Title, j.depth(), scope.File, scope.Line)
		j.scopes = append(j.scopes, next)
		j.index[ID] = next
	}
	j.current = j.index[ID]
}

func (j *JsonReporter) Report(report *AssertionResult) {
	j.current.Assertions = append(j.current.Assertions, report)
}

func (j *JsonReporter) Exit() {
	j.currentKey = j.currentKey[:len(j.currentKey)-1]
}

func (j *JsonReporter) EndStory() {
	j.report()
	j.reset()
}
func (j *JsonReporter) report() {
	scopes := []string{}
	for _, scope := range j.scopes {
		serialized, err := json.Marshal(scope)
		if err != nil {
			j.out.Println(jsonMarshalFailure)
			panic(err)
		}
		var buffer bytes.Buffer
		json.Indent(&buffer, serialized, "", "  ")
		scopes = append(scopes, buffer.String())
	}
	j.out.Print(fmt.Sprintf("%s\n%s,\n%s\n", OpenJson, strings.Join(scopes, ","), CloseJson))
}
func (j *JsonReporter) reset() {
	j.scopes = []*ScopeResult{}
	j.index = map[string]*ScopeResult{}
	j.currentKey = nil
}

func (j *JsonReporter) Write(content []byte) (written int, err error) {
	j.current.Output += string(content)
	return len(content), nil
}

func NewJsonReporter(out *Printer) *JsonReporter {
	self := new(JsonReporter)
	self.out = out
	self.reset()
	return self
}

const OpenJson = ">->->OPEN-JSON->->->"   // "⌦"
const CloseJson = "<-<-<-CLOSE-JSON<-<-<" // "⌫"
const jsonMarshalFailure = `

GOCONVEY_JSON_MARSHALL_FAILURE: There was an error when attempting to convert test results to JSON.
Please file a bug report and reference the code that caused this failure if possible.

Here's the panic:

`
