// TODO: truncate to given line width (so we can format the end of the line with '.', 'X', or 'E')?

package execution

import (
	"fmt"
	"io"
	"strings"
)

func (self *printer) println(message string, values ...interface{}) {
	formatted := self.format(message, values...) + "\n"
	self.out.Write([]byte(formatted))
}

func (self *printer) print(message string, values ...interface{}) {
	formatted := self.format(message, values...)
	self.out.Write([]byte(formatted))
}

func (self *printer) format(message string, values ...interface{}) string {
	formatted := self.prefix + fmt.Sprintf(message, values...)
	indented := strings.Replace(formatted, "\n", "\n"+self.prefix, -1)
	return strings.TrimRight(indented, "\t")
}

func (self *printer) indent() {
	self.prefix += "\t"
}

func (self *printer) dedent() {
	self.prefix = self.prefix[:len(self.prefix)-1]
}

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
