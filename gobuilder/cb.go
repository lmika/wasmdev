package gobuilder

type ContinuousGoBuilder struct {
	Builder *GoBuilder
	Hooks   ContinuousBuildHooks
}

func (cb *ContinuousGoBuilder) OnStartWatching() {
	if cb.Hooks != nil {
		cb.Hooks.OnStartListening()
	}
}

func (cb *ContinuousGoBuilder) OnStopWatching() {
	if cb.Hooks != nil {
		cb.Hooks.OnStopListening()
	}
}

func (cb *ContinuousGoBuilder) OnFileModified(file string) {
	cb.Builder.Build()
}