package devserver

const indexTemplate = `<!DOCTYPE html>
<html>
<head>
    <meta charset="utf-8"/>

	{{range .Stylesheets}}
		<link rel="stylesheet" href="{{.Href}}">
	{{end}}

	{{range .Scripts}}
		<script src="{{.Href}}"></script>
	{{end}}

    <script>
        const go = new Go();
        WebAssembly.instantiateStreaming(fetch("/main.wasm"), go.importObject).then((result) => {
            go.run(result.instance);
        });
    </script>
</head>
<body>

</body>
</html>
`
