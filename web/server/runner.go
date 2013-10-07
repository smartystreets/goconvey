package main

import (
	"github.com/smartystreets/goconvey/web/server/results"
	"hash"
)

type WrappedOS interface {
	Exists(directory string) bool
	Chdir(directory string) error
	Execute(command string, args ...string) string
}

type WrappedWatcher interface {
	WatchedPackages() []string
}

type TestRunner struct {
	done     chan (bool)
	os       WrappedOS
	watcher  WrappedWatcher
	revision hash.Hash
	results  []*results.PackageResult
}

func newTestRunner(done chan (bool), os WrappedOS, watcher WrappedWatcher) *TestRunner {
	self := &TestRunner{}
	return self
}

func (self *TestRunner) RunAll() results.CompleteOutput {
	// foreach package in watcher
	//   chdir to watch
	//   if chdir fails
	//     skip? panic?
	//   execute tests in package
	//   update revision with raw package output
	//   parse package output
	// return complete output
	return results.CompleteOutput{}
}
