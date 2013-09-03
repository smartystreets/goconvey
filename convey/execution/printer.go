package execution

import (
	"fmt"
	"io"
	"strings"
)

func (self *printer) println(message string, values ...interface{}) {
	formatted := self.format(message, values...) + newline
	self.out.Write([]byte(formatted))
}

func (self *printer) print(message string, values ...interface{}) {
	formatted := self.format(message, values...)
	self.out.Write([]byte(formatted))
}

func (self *printer) insert(text string) {
	self.out.Write([]byte(text))
}

func (self *printer) format(message string, values ...interface{}) string {
	formatted := self.prefix + fmt.Sprintf(message, values...)
	indented := strings.Replace(formatted, newline, newline+self.prefix, -1)
	return strings.TrimRight(indented, space)
}

func (self *printer) indent() {
	self.prefix += pad
}

func (self *printer) dedent() {
	if len(self.prefix) >= padLength {
		self.prefix = self.prefix[:len(self.prefix)-padLength]
	}
}

const newline = "\n"
const space = " "
const pad = space + space
const padLength = len(pad)

type printer struct {
	out    io.Writer
	prefix string
}

func newPrinter(out io.Writer) *printer {
	self := printer{}
	self.out = out
	return &self
}

type console struct{}

func (self *console) Write(p []byte) (n int, err error) {
	return fmt.Print(string(p))
}

func newConsole() io.Writer {
	return &console{}
}
