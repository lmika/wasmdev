package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/lmika/wasmdev/devserver"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

func main() {
	wasmExec := filepath.Join(runtime.GOROOT(), "misc/wasm/wasm_exec.js")
	if _, err := os.Stat(wasmExec); err != nil {
		log.Fatalf("cannot stat '%v': %v", wasmExec, err)
	}
	log.Printf("Using wasm_exec.js at %v", wasmExec)

	// Dev server
	devServer := devserver.New(devserver.Config{
		Scripts: []devserver.Resource{
			// Extra script files
			{ Href: "http://esironal.github.io/cmtouch/lib/codemirror.js" },
			{ Href: "http://esironal.github.io/cmtouch/addon/hint/show-hint.js" },
			{ Href: "http://esironal.github.io/cmtouch/addon/hint/xml-hint.js" },
			{ Href: "http://esironal.github.io/cmtouch/addon/hint/html-hint.js" },
			{ Href: "http://esironal.github.io/cmtouch/mode/xml/xml.js" },
			{ Href: "http://esironal.github.io/cmtouch/mode/javascript/javascript.js" },
			{ Href: "http://esironal.github.io/cmtouch/mode/css/css.js" },
			{ Href: "http://esironal.github.io/cmtouch/mode/htmlmixed/htmlmixed.js" },
			{ Href: "http://esironal.github.io/cmtouch/addon/selection/active-line.js" },
			{ Href: "http://esironal.github.io/cmtouch/addon/edit/matchbrackets.js" },

			// The WASM execution script file
			{ Href: "/wasm_exec.js", Source: wasmExec },
		},

		Stylesheets: []devserver.Resource{
			{ Href: "http://esironal.github.io/cmtouch/lib/codemirror.css" },
			{ Href: "http://esironal.github.io/cmtouch/addon/hint/show-hint.css" },
			{ Href: "http://esironal.github.io/cmtouch/theme/neonsyntax.css" },

			// Local stylesheet
			{ Href: "/styles.css", Source: "styles.css" },
		},

		TargetWasm: "main.wasm",
	})

	go func() {
		log.Println("Started dev server on :8080")
		http.ListenAndServe(`:8080`, devServer)
	}()

	buildListener()
}

// Build listener

func buildListener() {
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	watcher.Add(".")

	log.Println("Scanning for files")

	for {
		select {
		case event := <-watcher.Events:
			//log.Println("event:", event)
			if event.Op&fsnotify.Write == fsnotify.Write {
				if filepath.Ext(event.Name) == ".go" {
					log.Println("modified file:", event.Name,  " Rebuilding")
					runGoBuild()
				}
			}
		case err := <-watcher.Errors:
			log.Println("error:", err)
		}
	}
}

func runGoBuild() {
	cmd := exec.Command("go", "build", "-o", "main.wasm", ".")
	cmd.Env = append(os.Environ(),
		"GOOS=js",
		"GOARCH=wasm",
	)

	outErr, err := cmd.CombinedOutput()
	if err != nil {
		log.Printf("go build error:\n%v", string(outErr))
	} else {
		log.Println("rebuilt")
	}
}
