package printing

import (
	"fmt"
	"io"
	"strings"
)

func (self *Printer) Println(message string, values ...interface{}) {
	formatted := self.format(message, values...) + newline
	self.out.Write([]byte(formatted))
}

func (self *Printer) Print(message string, values ...interface{}) {
	formatted := self.format(message, values...)
	self.out.Write([]byte(formatted))
}

func (self *Printer) Insert(text string) {
	self.out.Write([]byte(text))
}

func (self *Printer) format(message string, values ...interface{}) string {
	formatted := self.prefix + fmt.Sprintf(message, values...)
	indented := strings.Replace(formatted, newline, newline+self.prefix, -1)
	return strings.TrimRight(indented, space)
}

func (self *Printer) Indent() {
	self.prefix += pad
}

func (self *Printer) Dedent() {
	if len(self.prefix) >= padLength {
		self.prefix = self.prefix[:len(self.prefix)-padLength]
	}
}

const newline = "\n"
const space = " "
const pad = space + space
const padLength = len(pad)

type Printer struct {
	out    io.Writer
	prefix string
}

func NewPrinter(out io.Writer) *Printer {
	self := Printer{}
	self.out = out
	return &self
}
