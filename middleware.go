package middlemux

import "net/http"

// Middleware is a type of function which wraps the given HandlerFunc and
// returns another. The returned function may perform actions prior to or after
// a particular call to the handler it invokes.
type Middleware func(h http.HandlerFunc) http.HandlerFunc
