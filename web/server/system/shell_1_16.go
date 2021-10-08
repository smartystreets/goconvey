// +build go1.16

package system

func compile(directory, gobin, tagsArg string) Command {
	return NewCommand(directory, gobin, "test", tagsArg)
}
