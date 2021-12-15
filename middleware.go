package middleware

import "net/http"

type middleware func(http.Handler) http.Handler

type Router struct {
	chain    []middleware
	ServeMux *http.ServeMux
}

func NewRouter() *Router {
	return &Router{
		ServeMux: http.NewServeMux(),
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

	r.ServeMux.Handle(route, mergedHandler)

	return mergedHandler
}
