package filewatcher

type WatchHandler interface {
	// OnFileModified is called when a file of interest has been modified
	OnFileModified(file string)
}
