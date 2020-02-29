package gobuilder

import "plugin"

// !! TEMP !!
type PluginListener struct {
	onStartBuildFn   func()
	onBuildSuccessFn func()
	onBuildFailedFn  func()
}

func LoadPluginListener(soFile string) (*PluginListener, error) {
	pgin, err := plugin.Open(soFile)
	if err != nil {
		return nil, err
	}

	if onLoadFn := loadFn(pgin, "OnLoad"); onLoadFn != nil {
		onLoadFn()
	}

	return &PluginListener{
		onStartBuildFn: loadFn(pgin, "OnBuildStart"),
		onBuildSuccessFn: loadFn(pgin, "OnBuildSuccessful"),
		onBuildFailedFn: loadFn(pgin, "OnBuildError"),
	}, nil
}

func (pl *PluginListener) OnBuildTriggered() {
	if fn := pl.onStartBuildFn; fn != nil {
		fn()
	}
}

func (pl *PluginListener) OnBuildSuccess() {
	if fn := pl.onBuildSuccessFn; fn != nil {
		fn()
	}
}

func (pl *PluginListener) OnBuildFailed() {
	if fn := pl.onBuildFailedFn; fn != nil {
		fn()
	}
}

func loadFn(pgin *plugin.Plugin, name string) func() {
	symbol, err := pgin.Lookup(name)
	if err != nil {
		return nil
	}

	fn, isFn := symbol.(func())
	if !isFn {
		return nil
	}

	return fn
}
