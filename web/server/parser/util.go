package parser

import (
	"math"
	"strings"
	"time"
)

// parseTestFunctionDuration parses the duration in seconds as a float64
// from a line of go test output that looks something like this:
// --- PASS: TestOldSchool_PassesWithMessage (0.03 seconds)
func parseTestFunctionDuration(line string) float64 {
	line = strings.Replace(line, "(", "", 1)
	line = strings.Replace(line, ")", "", 1)
	fields := strings.Split(line, " ")
	return parseDurationInSeconds(fields[3], 2)
}

func parseDurationInSeconds(raw string, precision int) float64 {
	elapsed, err := time.ParseDuration(raw)
	if err != nil {
		elapsed, _ = time.ParseDuration(raw + "s")
	}
	return round(elapsed.Seconds(), precision)
}

// round returns the rounded version of x with precision.
//
// Special cases are:
//  round(±0) = ±0
//  round(±Inf) = ±Inf
//  round(NaN) = NaN
//
// Why, oh why doesn't the math package come with a round function?
// Inspiration: http://play.golang.org/p/ZmFfr07oHp
func round(x float64, precision int) float64 {
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

// applyCarriageReturn process line just like terminals do (CR moves
// cursor to the start of line and next output overwrites previous one)
// and return final line without CR.
func applyCarriageReturn(line string) string {
	output := []rune(line)[:0]
	cursor := 0
	for _, r := range line {
		if r == '\r' {
			cursor = 0
		} else if cursor < len(output) {
			output[cursor] = r
			cursor++
		} else {
			output = append(output, r)
			cursor++
		}
	}
	return string(output)
}
