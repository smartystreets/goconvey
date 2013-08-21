package goconvey

import (
	"testing"
)

func Convey(situation string, test *testing.T, action func(c C)) {
	scene := NewStep(situation, action, test)

	for !scene.complete {
		scene.visit()
	}
}

type C interface {
	Convey(situation string, action func(C))
}