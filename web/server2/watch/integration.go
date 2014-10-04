package watch

import (
	"log"
	"time"

	"github.com/smartystreets/goconvey/web/server2/messaging"
)

type Watcher struct {
	rootFolder      string
	folderDepth     int
	ignoredFolders  map[string]struct{}
	fileSystemState int64

	input  chan messaging.ServerCommand
	output chan messaging.WatcherCommand
}

func NewWatcher(rootFolder string, folderDepth int, input chan messaging.ServerCommand, output chan messaging.WatcherCommand) *Watcher {
	return &Watcher{
		rootFolder:  rootFolder,
		folderDepth: folderDepth,
		input:       input,
		output:      output,

		ignoredFolders: make(map[string]struct{}),
	}
}

func (this *Watcher) Listen() {
	for {
		select {

		case command := <-this.input:
			if this.execute(command) {
				this.fileSystemState = 0
			}

		default:
			this.scan()
			time.Sleep(time.Millisecond * 250)
		}
	}
}

func (this *Watcher) execute(command messaging.ServerCommand) bool {
	log.Println("Received command from server:", command)

	switch command.Instruction {

	case messaging.ServerAdjustRoot:
		log.Println("Adjusting root...")
		this.rootFolder = command.Details

	case messaging.ServerIgnore:
		log.Println("Ignoring specified folders")

	case messaging.ServerReinstate:
		log.Println("Reinstating specified folders")

	default:
		log.Println("Unrecognized command from server:", command.Instruction)
		return false
	}

	return true
}

func (this *Watcher) scan() {
	items := YieldFileSystemItems(this.rootFolder)
	folderItems, profileItems, goFileItems := Categorize(items)
	rawProfiles := ReadProfiles(profileItems)
	profiles := ParseProfiles(rawProfiles)
	folders := CreateFolders(folderItems, goFileItems, profiles)
	folders = FilterDepth(folders, this.folderDepth)
	folders = FlagIgnored(folders, this.ignoredFolders)
	checksum := Checksum(folders)

	if checksum == this.fileSystemState {
		return
	}

	defer func() { this.fileSystemState = checksum }()
	this.output <- messaging.WatcherCommand{Folders: folders}
}
