package main

import (
	"github.com/mdwhatcott/goconvey"
)

var runner goconvey.runner

func main() {
	runner.Run()
}

func Convey(situation string, action func()) {
	runner.Convey(situation, action)
}

func Reset(action func()) {
	runner.Reset(action)
}
