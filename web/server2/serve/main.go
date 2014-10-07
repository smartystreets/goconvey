package main

import (
	"fmt"
	"time"

	"github.com/smartystreets/goconvey/web/server2/watch"
)

func main() {
	started := time.Now()

	items := watch.YieldFileSystemItems("/Users/mike/code/src/github.com/smartystreets/goconvey")
	folderItems, profileItems, goFileItems := watch.Categorize(items)

	for _, item := range profileItems {
		contents := watch.ReadContents(item.Path)
		item.ProfileDisabled, item.ProfileArguments = watch.ParseProfile(contents)
	}

	folders := watch.CreateFolders(folderItems)
	watch.LimitDepth(folders, -1)
	watch.AttachProfiles(folders, profileItems)
	watch.MarkIgnored(folders, nil)

	checksum := int64(len(watch.ActiveFolders(folders)))
	checksum += watch.Sum(folders, profileItems)
	checksum += watch.Sum(folders, goFileItems)

	// defer func() { this.fileSystemState = checksum }()
	// this.output <- messaging.WatcherCommand{Folders: folders}

	fmt.Println("ELAPSED:", time.Since(started))
	fmt.Println("CHECKSUM:", checksum)
	fmt.Println(len(folderItems), len(goFileItems), len(profileItems))
	fmt.Println(profileItems[0].ProfileArguments)
}
