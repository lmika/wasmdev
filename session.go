package main

import (
	"fmt"
	"github.com/lmika/wasmdev/config"
	"github.com/lmika/wasmdev/devserver"
	"github.com/lmika/wasmdev/filewatcher"
	"github.com/lmika/wasmdev/gobuilder"
	"github.com/lmika/wasmdev/profiles"
	"log"
	"os"
	"path/filepath"
	"runtime"
)

type Session struct {
	devServer   *devserver.DevServer
	builder     *gobuilder.GoBuilder
	fileWatcher *filewatcher.FileWatcher
}

func NewFromConfig(conf *config.Config) (*Session, error) {
	profileName := conf.GetString("profile", "go")
	profile, hasProfile := profiles.StandardProfiles[profileName]
	if !hasProfile {
		return nil, fmt.Errorf("unrecognised profile: %v", profileName)
	}

	devServer, err := configureDevServer(conf, profile)
	if err != nil {
		return nil, err
	}

	builder, err := configureBuilder(conf, profile)
	if err != nil {
		return nil, err
	}

	fileWatcher, err := configureFileWatcher(builder)
	if err != nil {
		return nil, err
	}

	return &Session{
		devServer:   devServer,
		builder:     builder,
		fileWatcher: fileWatcher,
	}, nil
}

func (session *Session) Watch() error {

}

func configureDevServer(conf *config.Config, profile profiles.Profile) (*devserver.DevServer, error) {
	if !conf.GetBool("devserver.enabled", true) {
		return nil, nil
	}

	// TODO: based on pipeline, may need to be different
	wasmExec := filepath.Join(runtime.GOROOT(), "misc/wasm/wasm_exec.js")
	if _, err := os.Stat(wasmExec); err != nil {
		log.Fatalf("cannot stat '%v': %v", wasmExec, err)
	}
	log.Printf("Using wasm_exec.js at %v", wasmExec)

	return devserver.New(devserver.Config{
		Scripts: []devserver.Resource{
			// The WASM execution script file
			{Href: "/wasm_exec.js", Source: wasmExec},
		},

		Stylesheets: []devserver.Resource{
			// Local stylesheet
			{Href: "/styles.css", Source: "styles.css"},
		},

		BindAddress: conf.GetString("devserver.listen", ":8080"),
		TargetWasm:  conf.GetString("build.target", "main.wasm"),
	}), nil
}

func configureBuilder(conf *config.Config, profile profiles.Profile) (*gobuilder.GoBuilder, error) {
	targetWasm := conf.GetString("build.target", "main.wasm")

	return &gobuilder.GoBuilder{
		TargetWasm: targetWasm,
		Pipeline:   profile.Pipeline,
	}, nil
}

func configureFileWatcher(builder *gobuilder.GoBuilder) (*filewatcher.FileWatcher, error) {
	fw := filewatcher.New()
	fw.Handler = &gobuilder.ContinuousGoBuilder{Builder: builder}
	fw.ExcludeDirs = []string{
		".*",
	}
	return fw, nil
}
