package printing

import (
	"fmt"
	"io"
)

type console struct{}

func (self *console) Write(p []byte) (n int, err error) {
	return fmt.Print(string(p))
}

func NewConsole() io.Writer {
	return &console{}
}
