package watch

import (
	"log"
	"sync"
	"time"

	"github.com/smartystreets/goconvey/web/server2/messaging"
)

type Watcher struct {
	rootFolder      string
	folderDepth     int
	ignoredFolders  map[string]struct{}
	fileSystemState int64

	input  chan messaging.ServerToWatcherCommand
	output chan messaging.Folders

	lock sync.RWMutex
}

func NewWatcher(rootFolder string, folderDepth int,
	input chan messaging.ServerToWatcherCommand, output chan messaging.Folders) *Watcher {

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

func (this *Watcher) execute(command messaging.ServerToWatcherCommand) bool {
	log.Println("Received command from server:", command)

	switch command.Instruction {

	case messaging.WatcherAdjustRoot:
		log.Println("Adjusting root...")
		this.rootFolder = command.Details

	case messaging.WatcherIgnore:
		log.Println("Ignoring specified folders") // TODO: protectedWrite(...)

	case messaging.WatcherReinstate:
		log.Println("Reinstating specified folders") // TODO: protectedWrite(...)

	default:
		log.Println("Unrecognized command from server:", command.Instruction)
		return false
	}

	return true
}

func (this *Watcher) scan() {
	items := YieldFileSystemItems(this.rootFolder)
	folderItems, profileItems, goFileItems := Categorize(items)

	for _, item := range profileItems {
		contents := ReadContents(item.Path)
		item.ProfileDisabled, item.ProfileArguments = ParseProfile(contents)
	}

	folders := CreateFolders(folderItems)
	// LimitDepth(folders, this.folderDepth)
	AttachProfiles(folders, profileItems)
	this.protectedRead(func() { MarkIgnored(folders, this.ignoredFolders) })

	checksum := int64(len(ActiveFolders(folders)))
	checksum += Sum(folders, profileItems)
	checksum += Sum(folders, goFileItems)

	if checksum == this.fileSystemState {
		return
	}

	log.Println("File system state modified, publishing current folders...", this.fileSystemState, checksum)

	defer func() { this.fileSystemState = checksum }()
	this.output <- folders
}

func (this *Watcher) protectedRead(protected func()) {
	this.lock.RLock()
	defer this.lock.RUnlock()
	protected()
}
func (this *Watcher) protectedWrite(protected func()) {
	this.lock.Lock()
	defer this.lock.Unlock()
	protected()
}
