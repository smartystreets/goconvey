package main

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"
)

func main() {
	starting := time.Now()
	newState := checksum("/Users/mike/work/dev/liveaddress-zipapi")
	fmt.Println(time.Since(starting), newState)
}

func checksum(root string) (sum int64) {
	goFiles := 0
	paths := 0
	filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		paths++
		if strings.HasSuffix(info.Name(), ".go") {
			sum += info.Size() + info.ModTime().Unix()
			goFiles++
		}
		return nil
	})
	fmt.Println(paths, goFiles)
	return
}
