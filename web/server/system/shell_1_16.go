// +build go1.16

package system

func compile(directory, gobin, tagsArg string) Command {
	// Don't pre-compile on 1.16
	// TODO(iannucci): rework the whole compilation feature to compile all
	// test binaries up front in a single `go` command.
	return Command{}
}
