package main

import (
	"github.com/mdwhatcott/goconvey"
)

var runner goconvey.runner

func main() {
	runner.Run()
}

func convey(situation string, action func()) {
	runner.Convey(situation, action)
}
