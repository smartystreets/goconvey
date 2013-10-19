package system

type Shell struct{}

func (self *Shell) Execute(name string, args ...string) (output string, err error) {
	return "", nil
}
