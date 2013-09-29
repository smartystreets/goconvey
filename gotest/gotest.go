// Package gotest contains internal functionality. Although this package
// contains one or more exported names it is not intended for public
// consumption. See the examples package for how to use this project.
package gotest

// This interface allows us to pass the *testing.T struct
// throughout the internals of this tool without ever
// having to import the "testing" package.
type T interface {
	Fail()
}
