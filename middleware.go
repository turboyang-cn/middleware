package middleware

import (
	"net/http"
)

type middleware func(http.Handler) http.Handler

type Router struct {
	chain    []middleware
	serveMux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		serveMux: http.NewServeMux(),
	}
}

func (r *Router) use(m middleware) {
	r.chain = append(r.chain, m)
}

func (r *Router) Add(route string, f func(w http.ResponseWriter, r *http.Request)) http.Handler {
	var mergedHandler http.Handler = http.HandlerFunc(f)

	for i := len(r.chain) - 1; i >= 0; i-- {
		mergedHandler = r.chain[i](mergedHandler)
	}

	r.serveMux.Handle(route, mergedHandler)

	return mergedHandler
}

func (r *Router) ListenAndServe(server *http.Server) error {
	server.Handler = r.serveMux

	return server.ListenAndServe()
}

func (r *Router) ListenAndServeTLS(server *http.Server, certFile string, keyFile string) error {
	server.Handler = r.serveMux

	return server.ListenAndServeTLS(certFile, keyFile)
}
