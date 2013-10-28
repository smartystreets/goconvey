package system

import (
	"os"
	"os/exec"
)

type Shell struct{}

func (self *Shell) Execute(name string, args ...string) (output string, err error) {
	rawOutput, err := exec.Command(name, args...).CombinedOutput()
	output = string(rawOutput)
	return
}

func (self *Shell) Getenv(key string) string {
	return os.Getenv(key)
}

func (self *Shell) Setenv(key, value string) error {
	return os.Setenv(key, value)
}
