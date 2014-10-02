package main

import (
	"log"

	"github.com/smartystreets/goconvey/web/server2/functional_core"
	"github.com/smartystreets/goconvey/web/server2/imperative_shell"
)

func main() {
	root := "/Users/mike/code/src/github.com/smartystreets/goconvey"
	items := imperative_shell.ScanFileSystem(root)
	items = functional_core.LimitDepth(items, 4)
	sum := functional_core.Checksum(items)
	log.Println(sum)
}
