package main

import (
	"fmt"
	"time"

	"github.com/smartystreets/goconvey/web/server2/watch"
)

func main() {
	started := time.Now()

	items := make(map[string]*watch.FileSystemItem)
	items["hi"] = &watch.FileSystemItem{}
	alter(items)

	fmt.Println("SIZE:", items["hi"].Size)
	fmt.Println(time.Since(started))
}

func alter(items map[string]*watch.FileSystemItem) {
	item := items["hi"]
	item.Size = 42
}
