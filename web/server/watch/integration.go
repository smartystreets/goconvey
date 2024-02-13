package watch

import (
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/smartystreets/goconvey/web/server/messaging"
)

type Watcher struct {
	nap             time.Duration
	rootFolder      string
	folderDepth     int
	ignoredFolders  map[string]struct{}
	fileSystemState int64
	paused          bool
	stopped         bool
	watchSuffixes   []string
	excludedDirs    []string

	input  chan messaging.WatcherCommand
	output chan messaging.Folders

	lock sync.RWMutex
}

func NewWatcher(rootFolder string, folderDepth int, nap time.Duration,
	input chan messaging.WatcherCommand, output chan messaging.Folders, watchSuffixes string, excludedDirs []string) *Watcher {

	return &Watcher{
		nap:           nap,
		rootFolder:    rootFolder,
		folderDepth:   folderDepth,
		input:         input,
		output:        output,
		watchSuffixes: strings.Split(watchSuffixes, ","),
		excludedDirs:  excludedDirs,

		ignoredFolders: make(map[string]struct{}),
	}
}

func (w *Watcher) Listen() {
	for {
		if !w.paused {
			w.scan()
		}
		select {
		case command := <-w.input:
			if stopped := w.respond(command); stopped {
				return
			}
		case <-time.After(w.nap):
			break
		}
	}
}

// respond handles the command and returns whether the watcher is permanently stopped
func (w *Watcher) respond(command messaging.WatcherCommand) bool {
	switch command.Instruction {

	case messaging.WatcherAdjustRoot:
		log.Println("Adjusting root...")
		w.rootFolder = command.Details
		w.execute()

	case messaging.WatcherIgnore:
		log.Println("Ignoring specified folders")
		w.ignore(command.Details)
		// Prevent a filesystem change due to the number of active folders changing
		_, checksum := w.gather()
		w.set(checksum)

	case messaging.WatcherReinstate:
		log.Println("Reinstating specified folders")
		w.reinstate(command.Details)
		// Prevent a filesystem change due to the number of active folders changing
		_, checksum := w.gather()
		w.set(checksum)

	case messaging.WatcherPause:
		log.Println("Pausing watcher...")
		w.paused = true

	case messaging.WatcherResume:
		log.Println("Resuming watcher...")
		w.paused = false

	case messaging.WatcherExecute:
		log.Println("Gathering folders for immediate execution...")
		w.execute()

	case messaging.WatcherStop:
		log.Println("Stopping the watcher...")
		close(w.output)
		w.stopped = true
		return true

	default:
		log.Println("Unrecognized command from server:", command.Instruction)
	}
	return false
}

func (w *Watcher) execute() {
	folders, _ := w.gather()
	w.sendToExecutor(folders)
}

func (w *Watcher) scan() {
	folders, checksum := w.gather()

	if checksum == w.fileSystemState {
		return
	}

	log.Println("File system state modified, publishing current folders...", w.fileSystemState, checksum)

	defer w.set(checksum)
	w.sendToExecutor(folders)
}

func (w *Watcher) gather() (folders messaging.Folders, checksum int64) {
	items := YieldFileSystemItems(w.rootFolder, w.excludedDirs)
	folderItems, profileItems, goFileItems := Categorize(items, w.rootFolder, w.watchSuffixes)

	for _, item := range profileItems {
		// TODO: don't even bother if the item's size is over a few hundred bytes...
		contents := ReadContents(item.Path)
		item.ProfileDisabled, item.ProfileTags, item.ProfileArguments = ParseProfile(contents)
	}

	folders = CreateFolders(folderItems)
	LimitDepth(folders, w.folderDepth)
	AttachProfiles(folders, profileItems)
	w.protectedRead(func() { MarkIgnored(folders, w.ignoredFolders) })

	active := ActiveFolders(folders)
	checksum = int64(len(active))
	checksum += Sum(active, profileItems)
	checksum += Sum(active, goFileItems)

	return folders, checksum
}

func (w *Watcher) set(state int64) {
	w.fileSystemState = state
}

func (w *Watcher) sendToExecutor(folders messaging.Folders) {
	w.output <- folders
}

func (w *Watcher) ignore(paths string) {
	w.protectedWrite(func() {
		for _, folder := range strings.Split(paths, string(os.PathListSeparator)) {
			w.ignoredFolders[folder] = struct{}{}
			log.Println("Currently ignored folders:", w.ignoredFolders)
		}
	})
}
func (w *Watcher) reinstate(paths string) {
	w.protectedWrite(func() {
		for _, folder := range strings.Split(paths, string(os.PathListSeparator)) {
			delete(w.ignoredFolders, folder)
		}
	})
}
func (w *Watcher) protectedWrite(protected func()) {
	w.lock.Lock()
	defer w.lock.Unlock()
	protected()
}
func (w *Watcher) protectedRead(protected func()) {
	w.lock.RLock()
	defer w.lock.RUnlock()
	protected()
}
