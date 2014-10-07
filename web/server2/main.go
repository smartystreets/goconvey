package main

import (
	"fmt"

	"github.com/smartystreets/goconvey/web/server2/messaging"
	"github.com/smartystreets/goconvey/web/server2/watch"
)

func main() {
	root := "/Users/mike/code/src/github.com/smartystreets/goconvey"

	serverToWatcher := make(chan messaging.ServerToWatcherCommand)
	watcherToExecutor := make(chan messaging.Folders)

	watcher := watch.NewWatcher(root, 2, serverToWatcher, watcherToExecutor)

	go watcher.Listen()

	for folders := range watcherToExecutor {
		fmt.Println("FOLDERS:", len(folders))
	}
}
