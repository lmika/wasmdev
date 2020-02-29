package filewatcher

type WatchHandler interface {
	OnStartWatching()
	OnStopWatching()

	// OnFileModified is called when a file of interest has been modified
	OnFileModified(file string)
}
