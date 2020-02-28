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

var flagNoDevServer = flag.Bool("noserve", false, "Do not setup a dev server")
var flagOut = flag.String("o", "", "Target WASM output file (default = main.wasm)")

func main() {
	flag.Parse()

	conf, err := config.FromWasmDevFile()
	if err != nil {
		log.Println(err)
	}

	wasmExec := filepath.Join(runtime.GOROOT(), "misc/wasm/wasm_exec.js")
	if _, err := os.Stat(wasmExec); err != nil {
		log.Fatalf("cannot stat '%v': %v", wasmExec, err)
	}
	log.Printf("Using wasm_exec.js at %v", wasmExec)

	// Dev server
	if !*flagNoDevServer && conf.GetBool("devserver.enabled", true) {
		devServer := devserver.New(devserver.Config{
			Scripts: []devserver.Resource{
				// Extra script files
				{Href: "http://esironal.github.io/cmtouch/lib/codemirror.js"},
				{Href: "http://esironal.github.io/cmtouch/addon/hint/show-hint.js"},
				{Href: "http://esironal.github.io/cmtouch/addon/hint/xml-hint.js"},
				{Href: "http://esironal.github.io/cmtouch/addon/hint/html-hint.js"},
				{Href: "http://esironal.github.io/cmtouch/mode/xml/xml.js"},
				{Href: "http://esironal.github.io/cmtouch/mode/javascript/javascript.js"},
				{Href: "http://esironal.github.io/cmtouch/mode/css/css.js"},
				{Href: "http://esironal.github.io/cmtouch/mode/htmlmixed/htmlmixed.js"},
				{Href: "http://esironal.github.io/cmtouch/addon/selection/active-line.js"},
				{Href: "http://esironal.github.io/cmtouch/addon/edit/matchbrackets.js"},

				// The WASM execution script file
				{Href: "/wasm_exec.js", Source: wasmExec},
			},

			Stylesheets: []devserver.Resource{
				{Href: "http://esironal.github.io/cmtouch/lib/codemirror.css"},
				{Href: "http://esironal.github.io/cmtouch/addon/hint/show-hint.css"},
				{Href: "http://esironal.github.io/cmtouch/theme/neonsyntax.css"},

				// Local stylesheet
				{Href: "/styles.css", Source: "styles.css"},
			},

			TargetWasm: conf.GetString("build.target", "main.wasm", config.WithStringOverride(*flagOut)),
		})

		go func() {
			log.Println("Started dev server on :8080")
			http.ListenAndServe(conf.GetString("devserver.listen", ":8080"), devServer)
		}()
	}

	buildListener(conf)
}

// Build listener

func buildListener(conf *config.Config) {
	targetWasm := conf.GetString("build.target", "main.wasm", config.WithStringOverride(*flagOut))

	builder := goBuilder{
		builder: &gobuilder.GoBuilder{ TargetWasm: targetWasm },
	}

	// Run the build first
	builder.builder.Build()

	fw := filewatcher.New()
	fw.Handler = builder
	fw.ExcludeDirs = []string{
		".*",
	}

	if err := fw.Watch(); err != nil {
		log.Fatal(err)
	}
}

type goBuilder struct{
	builder	*gobuilder.GoBuilder
}

func (gb goBuilder) OnFileModified(file string) {
	gb.builder.Build()
}