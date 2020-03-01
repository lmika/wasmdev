package main

import (
	"flag"
	"github.com/lmika/wasmdev/config"
	"github.com/lmika/wasmdev/devserver"
	"github.com/lmika/wasmdev/filewatcher"
	"github.com/lmika/wasmdev/gobuilder"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
)

var flagServe = flag.String("serve", ":8080", "Host:Port to bind the dev server")
var flagNoDevServer = flag.Bool("noserve", false, "Do not setup a dev server")
var flagOut = flag.String("o", "", "Target WASM output file (default = main.wasm)")

func main() {
	flag.Parse()

	conf, err := config.FromWasmDevFile()
	if err != nil {
		log.Println(err)
	}

	// TODO: this needs to be moved somewhere else
	wasmExec := filepath.Join(runtime.GOROOT(), "misc/wasm/wasm_exec.js")
	if _, err := os.Stat(wasmExec); err != nil {
		log.Fatalf("cannot stat '%v': %v", wasmExec, err)
	}
	log.Printf("Using wasm_exec.js at %v", wasmExec)

	// Dev server
	if !*flagNoDevServer && conf.GetBool("devserver.enabled", true) {
		devServer := devserver.New(devserver.Config{
			Scripts: []devserver.Resource{
				// The WASM execution script file
				{Href: "/wasm_exec.js", Source: wasmExec},
			},

			Stylesheets: []devserver.Resource{
				// Local stylesheet
				{Href: "/styles.css", Source: "styles.css"},
			},

			TargetWasm: conf.GetString("build.target", "main.wasm", config.WithStringOverride(*flagOut)),
		})

		go func() {
			bindAddr := conf.GetString("devserver.listen", ":8080", config.WithStringOverride(*flagServe))

			log.Printf("Started dev server on %v", bindAddr)
			http.ListenAndServe(bindAddr, devServer)
		}()
	}

	buildListener(conf)
}

// Build listener

func buildListener(conf *config.Config) {
	targetWasm := conf.GetString("build.target", "main.wasm", config.WithStringOverride(*flagOut))

	builder := &gobuilder.GoBuilder{
		TargetWasm: targetWasm,
	}
	continuousBuilder := &gobuilder.ContinuousGoBuilder{
		Builder: builder,
	}

	// Do the initial build first
	builder.Build()

	// Start watching for updates
	fw := filewatcher.New()
	fw.Handler = continuousBuilder
	fw.ExcludeDirs = []string{
		".*",
	}

	if err := fw.Watch(); err != nil {
		log.Fatal(err)
	}
}