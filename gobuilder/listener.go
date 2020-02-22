package gobuilder

// BuildHook can be used to listen to go build events
type BuildHook interface {
	// OnBuildTriggered is called when the Go build has started
	OnBuildTriggered()

	// OnBuildSuccess is called when the build was successful
	OnBuildSuccess()

	// OnBuildFailed is called when the build failed
	OnBuildFailed()
}


type ContinuousBuildHooks interface {
	// OnStartListening is called when the continuous build has started
	OnStartListening()

	// OnStopListening is called when the continuous build is complete
	OnStopListening()
}
