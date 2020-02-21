package devserver

import (
	"fmt"
	"html"
	"html/template"
	"io"
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
		w.Header().Set("Content-type", "text/html; charset=utf-8")
		w.WriteHeader(200)

		if err := tmpl.Execute(w, data); err != nil {
			// TODO: shouldn't be a 200 here
			fmt.Printf("<html><body><h1>Template Error</h1><pre>%v</pre></html>", html.EscapeString(err.Error()))
		}
	})
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/html; charset=utf-8")
	w.WriteHeader(200)

	io.WriteString(w, `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8"/>

	<!-- TODO: make style injection configurable -->
	<link rel="stylesheet" href="/styles.css">
	<link rel="stylesheet" href="http://esironal.github.io/cmtouch/lib/codemirror.css">
	<link rel="stylesheet" href="http://esironal.github.io/cmtouch/addon/hint/show-hint.css">
	<link rel="stylesheet" href="http://esironal.github.io/cmtouch/theme/neonsyntax.css">

	<!-- TODO: and these -->
	<script src="http://esironal.github.io/cmtouch/lib/codemirror.js"></script>
	<script src="http://esironal.github.io/cmtouch/addon/hint/show-hint.js"></script>
	<script src="http://esironal.github.io/cmtouch/addon/hint/xml-hint.js"></script>
	<script src="http://esironal.github.io/cmtouch/addon/hint/html-hint.js"></script>
	<script src="http://esironal.github.io/cmtouch/mode/xml/xml.js"></script>
	<script src="http://esironal.github.io/cmtouch/mode/javascript/javascript.js"></script>
	<script src="http://esironal.github.io/cmtouch/mode/css/css.js"></script>
	<script src="http://esironal.github.io/cmtouch/mode/htmlmixed/htmlmixed.js"></script>
	<script src="http://esironal.github.io/cmtouch/addon/selection/active-line.js"></script>
	<script src="http://esironal.github.io/cmtouch/addon/edit/matchbrackets.js"></script>
	

	<script src="/wasm_exec.js"></script>
    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("/main.wasm"), go.importObject).then((result) => {
            go.run(result.instance);
        });
    </script>
</head>
<body>
</body>
</html>`)
}
