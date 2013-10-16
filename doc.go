// package goconvey doesn't do anything except list external dependencies for
// all packages not at the root level so that when you 'go get' this project
// everything arrives at once and everything is built/installed successfully.
package goconvey

import (
	_ "github.com/howeyc/fsnotify"
	_ "github.com/jacobsa/oglematchers"
	_ "github.com/jacobsa/ogletest"
)
