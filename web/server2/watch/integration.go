package watch

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/smartystreets/goconvey/web/server2/messaging"
)

type Watcher struct {
	rootFolder      string
	folderDepth     int
	ignoredFolders  map[string]struct{}
	fileSystemState int64
	paused          bool

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
			this.respond(command)

		default:
			if !this.paused {
				this.scan()
			}
			time.Sleep(nap)
		}
	}
}

func (this *Watcher) respond(command messaging.ServerToWatcherCommand) {
	log.Println("Received command from server:", command)

	switch command.Instruction {

	case messaging.WatcherAdjustRoot:
		log.Println("Adjusting root...")
		this.rootFolder = command.Details
		this.set(0)

	case messaging.WatcherIgnore:
		this.protectedWrite(func() {
			log.Println("Ignoring specified folders")
			for _, folder := range strings.Split(command.Details, string(os.PathListSeparator)) {
				this.ignoredFolders[folder] = struct{}{}
			}
		})
		this.set(0)

	case messaging.WatcherReinstate:
		this.protectedWrite(func() {
			log.Println("Reinstating specified folders")
			for _, folder := range strings.Split(command.Details, string(os.PathListSeparator)) {
				delete(this.ignoredFolders, folder)
			}
		})
		this.set(0)

	case messaging.WatcherPause:
		log.Println("Pausing watcher...")
		this.paused = true
		this.set(0)

	case messaging.WatcherResume:
		log.Println("Resuming watcher...")
		this.paused = false

	case messaging.WatcherExecute:
		log.Println("Gathering folders for immediate execution...")
		folders, _ := this.gather()
		this.sendToExecutor(folders)

	default:
		log.Println("Unrecognized command from server:", command.Instruction)
	}
}

func (this *Watcher) scan() {
	folders, checksum := this.gather()

	if checksum == this.fileSystemState {
		return
	}

	log.Println("File system state modified, publishing current folders...", this.fileSystemState, checksum)

	defer this.set(checksum)
	this.sendToExecutor(folders)
}

func (this *Watcher) gather() (folders messaging.Folders, checksum int64) {
	items := YieldFileSystemItems(this.rootFolder)
	folderItems, profileItems, goFileItems := Categorize(items)

	for _, item := range profileItems {
		contents := ReadContents(item.Path)
		item.ProfileDisabled, item.ProfileArguments = ParseProfile(contents)
	}

	folders = CreateFolders(folderItems)
	LimitDepth(folders, this.folderDepth)
	AttachProfiles(folders, profileItems)                                    // TODO: test drive
	this.protectedRead(func() { MarkIgnored(folders, this.ignoredFolders) }) // TODO: test drive

	checksum = int64(len(ActiveFolders(folders)))
	checksum += Sum(folders, profileItems)
	checksum += Sum(folders, goFileItems)

	return folders, checksum
}

func (this *Watcher) sendToExecutor(folders messaging.Folders) {
	this.output <- folders
}

func (this *Watcher) set(state int64) {
	this.fileSystemState = state
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

const nap = time.Millisecond * 250
