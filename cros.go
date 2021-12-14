package middleware

import (
	"net/http"
	"strconv"
	"strings"

	"github.com/ahmetb/go-linq/v3"
)

func (r *Router) UseCros() {
	r.UseCrosWith(nil, true, nil, nil, nil)
}

func (r *Router) UseCrosWith(allowOrigins *[]string, allowCredentials bool, exposeHeaders *[]string, allowHeaders *[]string, allowMethods *[]string) {
	r.use(crosMiddleware(allowOrigins, allowCredentials, exposeHeaders, allowHeaders, allowMethods))
}

func crosMiddleware(allowOrigins *[]string, allowCredentials bool, exposeHeaders *[]string, allowHeaders *[]string, allowMethods *[]string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if origin := r.Header.Get("Origin"); origin != "" {
				allowAccess := allowOrigins == nil || linq.From(*allowOrigins).
					AnyWith(func(i interface{}) bool {
						return strings.EqualFold(i.(string), origin)
					})

				if allowAccess {
					w.Header().Set("Access-Control-Allow-Origin", origin)
					w.Header().Set("Access-Control-Allow-Credentials", strconv.FormatBool(allowCredentials))

					if exposeHeaders == nil {
						w.Header().Set("Access-Control-Expose-Headers", "*")
					} else {
						w.Header().Set("Access-Control-Expose-Headers", strings.Join(*exposeHeaders, ", "))
					}

					if method := r.Header.Get("Access-Control-Request-Method"); method != "" {
						allowAccess := allowMethods == nil || linq.From(*allowMethods).
							AnyWith(func(i interface{}) bool {
								return strings.EqualFold(i.(string), method)
							})

						if allowAccess {
							if allowMethods == nil {
								w.Header().Set("Access-Control-Allow-Methods", "*")
							} else {
								w.Header().Set("Access-Control-Allow-Methods", strings.Join(*allowMethods, ", "))
							}
						} else {
							w.WriteHeader(http.StatusForbidden)

							return
						}
					}

					if headers := r.Header.Get("Access-Control-Request-Headers"); headers != "" {
						allowAccess := allowHeaders == nil || linq.From(strings.Split(headers, ",")).
							Select(func(i interface{}) interface{} {
								return strings.TrimSpace(i.(string))
							}).All(func(i interface{}) bool {
							return linq.From(*allowHeaders).
								AnyWith(func(j interface{}) bool {
									return strings.EqualFold(i.(string), j.(string))
								})
						})

						if allowAccess {
							if allowHeaders == nil {
								w.Header().Set("Access-Control-Allow-Headers", "*")
							} else {
								w.Header().Set("Access-Control-Allow-Headers", strings.Join(*allowHeaders, ", "))
							}
						} else {
							w.WriteHeader(http.StatusForbidden)

							return
						}
					}

					if r.Method == http.MethodOptions {
						w.WriteHeader(http.StatusOK)

						return
					}
				} else {
					w.WriteHeader(http.StatusForbidden)

					return
				}
			}

			next.ServeHTTP(w, r)
		})
	}
}
