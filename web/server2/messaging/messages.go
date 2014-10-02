package messaging

import "time"

type ( // Commands the user will send via the HTTP server:

	AdjustRootFolderCommand struct {
		NewRootFolder string
	}

	IgnorePackagesCommand struct {
		Paths []string
	}

	ReinstatePackagesCommand struct {
		Paths []string
	}

	ClearExecutionStatusCommand struct{} // not sure why we need this...

	ExecuteTestsCommand struct{}
)

type ( // Events:

	FileSystemItemFoundEvent struct {
		Root     string
		Path     string
		Name     string
		Size     int64
		Modified int64
		IsFolder bool
	}

	RootFolderAdjustedEvent struct {
		NewRootFolder string
	}

	FolderFoundEvent struct {
		Path string
	}

	GoFileFoundEvent struct {
		Path string
	}

	NewFolderToWatchEvent struct {
		Path      string
		ImportDir string
		// TODO: include other info
	}

	FolderDeletedEvent struct {
		Path      string
		ImportDir string
	}

	FolderIgnoredViaUIEvent struct {
		Path      string
		ImportDir string
	}

	FolderReinstatedViaUIEvent struct {
		Path      string
		ImportDir string
	}

	FolderIgnoredViaProfileEvent struct {
		Path      string
		ImportDir string
	}

	FolderReinstatedViaProfileEvent struct {
		Path      string
		ImportDir string
	}

	FolderProfileProvidedEvent struct {
		PackagePath string
		ImportDir   string
		Flags       []string
	}

	LatestTestResultsRenderedStaleEvent struct {
		Stamp time.Time
	}

	NewTestResultsPublishedEvent struct {
		Stamp time.Time
		// TODO: need a 'result' type here...
	}
)
