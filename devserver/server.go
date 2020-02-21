package devserver

import "net/http"

type Resource struct {
	// Href is
	Href string

	// Source is the file that this stylesheet should be fetched from.  Can be empty, in
	// which case this stylesheet will not be served.
	Source string

	// Cache manages the cache control of the resource.
	Cache bool
}

// Options for the dev server
type Config struct {
	// Stylesheets hold the extract stylesheets to add to the page
	Stylesheets []Resource
	Scripts     []Resource

	// Target WASM file
	TargetWasm string
}

// DevServer manages the development server
type DevServer struct {
	config	Config
	mux *http.ServeMux
}

func New(config Config) *DevServer {
	ds := &DevServer{config: config, mux: http.NewServeMux()}
	ds.setupMux()

	return ds
}

func (ds *DevServer) setupMux() {
	ds.setupHandlersForResources(ds.config.Scripts)
	ds.setupHandlersForResources(ds.config.Stylesheets)

	ds.mux.Handle("/main.wasm", serveStaticFile(ds.config.TargetWasm, false))
	ds.mux.Handle("/", exactPathMux("/", http.HandlerFunc(serveIndex), http.FileServer(http.Dir(`.`))))
}

func (ds *DevServer) setupHandlersForResources(resources []Resource) {
	for _, resource := range resources {
		// TODO: parse the URL and only setup handlers for resources of the form '/bla'
		if resource.Source != "" {
			ds.mux.Handle(resource.Href, serveStaticFile(resource.Source, resource.Cache))
		}
	}
}

func (ds *DevServer) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	ds.mux.ServeHTTP(w, r)
}
