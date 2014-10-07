package messaging

///////////////////////////////////////////////////////////////////////////////

type ServerToWatcherCommand struct {
	Instruction int
	Details     string
}

const (
	WatcherPause = iota
	WatcherResume
	WatcherIgnore
	WatcherReinstate
	WatcherAdjustRoot
)

///////////////////////////////////////////////////////////////////////////////

type Folders map[string]*Folder

type Folder struct {
	Path          string // key
	Root          string
	Ignored       bool
	Disabled      bool
	TestArguments []string
}

///////////////////////////////////////////////////////////////////////////////
