package main

import (
	"math"
	"time"
)

func parseDuration(raw string, precision int) float64 {
	elapsed, _ := time.ParseDuration(raw)
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
