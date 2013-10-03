/*
- all go test output is compiled as an array of PackageResult structs
- a PackageResult is an object with a name, an elapsed time, and an array of StoryResults
- a StoryResult is an array of ScopeResults
- a ScopeResult is an object with a depth field and an array of AssertionResults
- a AssertionResults is a failure, pass, error, skip w/ stacktrace
*/

package main

import (
	"encoding/json"
	"github.com/smartystreets/goconvey/reporting"
	"math"
	"strings"
	"time"
)

type PackageResult struct {
	PackageName string
	Elapsed     float64
	Passed      bool
	Stories     []StoryResult
}

type StoryResult []reporting.ScopeResult

func parsePackageResult(rawResult string) PackageResult {
	passed, packageName, duration := parseMetadata(rawResult)

	var storyEnd int
	storyEnd = strings.Index(rawResult, ",PASS\nok")
	if storyEnd < 0 {
		storyEnd = strings.Index(rawResult, ",--- FAIL")
	}
	rawStories := "[" + rawResult[:storyEnd] + "]"

	var stories []StoryResult
	json.Unmarshal([]byte(rawStories), &stories) // TODO: returns err...

	return PackageResult{
		PackageName: packageName,
		Elapsed:     duration,
		Passed:      passed,
		Stories:     stories,
	}
}
func parseMetadata(rawResult string) (passed bool, packageName string, elapsed float64) {
	lines := strings.Split(strings.TrimSpace(rawResult), "\n")
	lastLine := lines[len(lines)-1]
	fields := strings.Split(lastLine, "\t")
	// fmt.Println(lastLine, fields)

	if len(fields) != 3 {
		// fmt.Println(lastLine)
	}

	passed = fields[0] == "ok  "
	packageName = fields[1]
	rawDuration := fields[2]

	duration, _ := time.ParseDuration(rawDuration)
	elapsed = Round(duration.Seconds(), 3)
	return
}

// Round returns the rounded version of x with precision.
//
// Special cases are:
//	Round(±0) = ±0
//	Round(±Inf) = ±Inf
//	Round(NaN) = NaN
//
// Why, oh why doesn't the math package come with a Round function?
// Inspiration: http://play.golang.org/p/ZmFfr07oHp
func Round(x float64, precision int) float64 {
	var rounder float64
	pow := math.Pow(10, float64(precision))
	intermediate := x * pow

	if intermediate < 0.0 {
		intermediate -= 0.5
	} else {
		intermediate += 0.5
	}
	rounder = float64(int64(intermediate))

	return rounder / float64(pow)
}
