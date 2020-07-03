package middlemux

import "net/http"

// MiddleMux is a wrapper for http.ServeMux that supports registration of
// one to many Middleware objects which will wrap calls to ServeHTTP.
type MiddleMux struct {
	*http.ServeMux
	handlerFn http.HandlerFunc
}

// NewMiddleMux creates a new MiddleMux instance, initialized with a fresh
// http.ServeMux instance.
func NewMiddleMux() *MiddleMux {
	mux := http.NewServeMux()

	return &MiddleMux{
		ServeMux:  mux,
		handlerFn: mux.ServeHTTP,
	}
}

// Use causes the given Middleware function to wrap calls to ServeHTTP on this
// MiddleMux. The last piece of middleware registered is the first to be invoked
// on calls to ServeHTTP.
func (mm *MiddleMux) Use(m Middleware) {
	// If the provided func is nil, ignore this.
	if m == nil {
		return
	}

	mm.handlerFn = m(mm.handlerFn)
}

// ServeHTTP overrides the default behavior of the embedded http.ServeMux and
// invokes the MiddleMux's own handler function which will use the middleware.
func (mm *MiddleMux) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	mm.handlerFn(w, r)
}
