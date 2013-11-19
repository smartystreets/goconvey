package system

import (
	"os"
	"os/exec"
)

type Shell struct{}

func (self *Shell) GoTest(directory string) (output string, err error) {
	output, err = self.execute(directory, "go", "test", "-i")
	if err == nil {
		output, err = self.execute(directory, "go", "test", "-v", "-timeout=-42s")
	}
	return
}

func (self *Shell) execute(directory, name string, args ...string) (output string, err error) {
	command := exec.Command(name, args...)
	command.Dir = directory
	rawOutput, err := command.CombinedOutput()
	output = string(rawOutput)
	return
}

func (self *Shell) Getenv(key string) string {
	return os.Getenv(key)
}

func (self *Shell) Setenv(key, value string) error {
	if self.Getenv(key) != value {
		return os.Setenv(key, value)
	}
	return nil
}

func NewShell() *Shell {
	return &Shell{}
}
