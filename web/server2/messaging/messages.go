package messaging

///////////////////////////////////////////////////////////////////////////////

type ServerCommand struct {
	Instruction int
	Details     string
}

const (
	ServerPause = iota
	ServerResume
	ServerIgnore
	ServerReinstate
	ServerAdjustRoot
)

///////////////////////////////////////////////////////////////////////////////

// type ( // Commands the user will send via the HTTP server:

// 	AdjustRootFolderCommand struct {
// 		NewRootFolder string
// 	}

// 	IgnorePackagesCommand struct {
// 		Paths []string
// 	}

// 	ReinstatePackagesCommand struct {
// 		Paths []string
// 	}

// 	ClearExecutionStatusCommand struct{} // not sure why we need this...

// 	ExecuteTestsCommand struct{}
// )

// type ( // Events:

// 	RootFolderAdjustedEvent struct {
// 		NewRootFolder string
// 	}

// 	FolderFoundEvent struct {
// 		Path string
// 	}

// 	FolderDeletedEvent struct {
// 		Path      string
// 		ImportDir string
// 	}

// 	FolderIgnoredViaUIEvent struct {
// 		Path      string
// 		ImportDir string
// 	}

// 	FolderReinstatedViaUIEvent struct {
// 		Path      string
// 		ImportDir string
// 	}

// 	FolderIgnoredViaProfileEvent struct {
// 		Path      string
// 		ImportDir string
// 	}

// 	FolderReinstatedViaProfileEvent struct {
// 		Path      string
// 		ImportDir string
// 	}

// 	FolderProfileProvidedEvent struct {
// 		PackagePath string
// 		ImportDir   string
// 		Flags       []string
// 	}

// 	LatestTestResultsRenderedStaleEvent struct {
// 		Stamp time.Time
// 	}

// 	NewTestResultsPublishedEvent struct {
// 		Stamp time.Time
// 		// TODO: need a 'result' type here...
// 	}
// )
