# wasmdev

A simple dev server for developing Go WASM projects.

The goal of this is to provide that rapid development approach commonly seen in
JavaScript for frontend Go development.

## Getting Started

The simplest use for a Go WASM project is with no arguments:

```
$ wasmdev
```

Doing this will:

- Build the current project using the Go compiler with `GOOS=js` and `GOARCH=wasm` 
- Start a development on port 8080 with a HTML shell that will load the compiled WASM file with the "wasm_exec.js" bootstrap file 
- Start a file listener that will rebuild the project when any Go files within the project change.

## Usage

A couple of arguments are supported at the moment (more likely to come):

- `-o`: set the target WASM file.  The default is "main.wasm"
- `-noserve`: do not start the dev server.  The file listener will still be used.

## Roadmap

This project is still in it's early stages but here are the things I'm thinking of adding down the line:

- Loading external JS and CSS resources in the generated HTML shell.
- Proxying requests to another service, e.g. an API.
- Defining configuration in a build file.
- Adding to or customising the build steps.
- Adding hooks to the build process.
- Maybe adding automatic refreshing to the page.
- Possibly a simple plugin architecture to customise the tool in other ways.

There are also a list of known bugs and limitations:

- Can't configure the port.
- Adding new directories to the project do not get picked up by the file listener. 

## License

License under the [MIT License](LICENSE.md).
