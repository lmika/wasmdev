package gobuilder

import "log"

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
	log.Println("Starting build")

	if err := cb.Builder.Build(); err != nil {
		log.Printf("Build failed: %v", err)
	} else {
		log.Println("Build complete")
	}
}