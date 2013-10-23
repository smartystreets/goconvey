package system

import (
	"os/exec"
)

type Shell struct{}

func (self *Shell) Execute(name string, args ...string) (output string, err error) {
	rawOutput, err := exec.Command(name, args...).CombinedOutput()
	output = string(rawOutput)
	return
}
