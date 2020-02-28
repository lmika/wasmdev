package devserver

import (
	"bytes"
	"fmt"
	"html"
	"html/template"
	"net/http"
)

func serveStaticFile(wasmExec string, cachable bool) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if !cachable {
			w.Header().Set("Pragma", "no-cache")
			w.Header().Set("Cache-Control", "no-store")
		}
		http.ServeFile(w, r, wasmExec)
	})
}

func exactPathMux(path string, exactHandler http.Handler, otherHandler http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path == path {
			exactHandler.ServeHTTP(w, r)
		} else {
			otherHandler.ServeHTTP(w, r)
		}
	})
}

func serveTemplate(templateSrc string, data interface{}) http.Handler {
	tmpl := template.Must(template.New("index.html").Parse(templateSrc))

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		respContent := new(bytes.Buffer)
		if err := tmpl.Execute(respContent, data); err != nil {
			w.Header().Set("Content-type", "text/plain; charset=utf-8")
			w.WriteHeader(500)
			fmt.Print(err.Error())
			return
		}

		w.Header().Set("Content-type", "text/html; charset=utf-8")
		w.WriteHeader(200)

		if err := tmpl.Execute(w, data); err != nil {
			fmt.Printf("<html><body><h1>Template Error</h1><pre>%v</pre></html>", html.EscapeString(err.Error()))
		}
	})
}