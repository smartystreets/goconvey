// TODO: under unit test

package reporting

import (
	"bytes"
	"encoding/json"
	"fmt"
	"strings"
)

type JsonReporter struct {
	nestedReporter
	out  *Printer
	seed int64
}

func (s *JsonReporter) Close() {
	flattened := []*ScopeResult{}
	stack := []*ScopeResult{}

	top := func() *ScopeResult {
		return stack[len(stack)-1]
	}

	if s.seed != 0 {
		flattened = append(flattened, &ScopeResult{
			Title:      "Random Seed",
			Depth:      1,
			Assertions: []*AssertionResult{},
		}, &ScopeResult{
			Title:      fmt.Sprint(s.seed),
			Depth:      2,
			Assertions: []*AssertionResult{},
		})
	}

	s.Walk(func(obj interface{}) {
		switch obj := obj.(type) {
		case *NestedScopeResult:
			ent := &ScopeResult{
				Title:      obj.Title,
				File:       obj.File,
				Line:       obj.Line,
				Depth:      len(stack) + 1,
				Assertions: []*AssertionResult{},
			}
			stack = append(stack, ent)
			flattened = append(flattened, ent)

		case ScopeExit:
			stack = stack[:len(stack)-1]

		case string:
			top().Output += obj

		case *AssertionResult:
			top().Assertions = append(top().Assertions, obj)
		}
	})

	scopes := []string{}
	for _, scope := range flattened {
		serialized, err := json.Marshal(scope)
		if err != nil {
			s.out.Statement(jsonMarshalFailure)
			panic(err)
		}
		var buffer bytes.Buffer
		json.Indent(&buffer, serialized, "", "  ")
		scopes = append(scopes, buffer.String())
	}
	s.out.Insert(fmt.Sprintf("%s\n%s,\n%s\n", OpenJson, strings.Join(scopes, ","), CloseJson))
}

func NewJsonReporter(out *Printer, seed int64) Reporter {
	return &JsonReporter{out: out, seed: seed}
}

const OpenJson = ">->->OPEN-JSON->->->"   // "⌦"
const CloseJson = "<-<-<-CLOSE-JSON<-<-<" // "⌫"
const jsonMarshalFailure = `

GOCONVEY_JSON_MARSHALL_FAILURE: There was an error when attempting to convert test results to JSON.
Please file a bug report and reference the code that caused this failure if possible.

Here's the panic:

`
