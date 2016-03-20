package api

import (
	"github.com/jboolean/whereabouts-api/router"
	"log"
	"net/http"
)

// docs route
var docsRoute = router.Route{
	Name:    "Docs",
	Method:  "GET",
	Pattern: "/{rest:docs.*}",
	Handler: http.StripPrefix("/docs", http.FileServer(http.Dir("./docs"))),
}

func recoveryInterceptor(inner http.Handler, route router.Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		defer func() {
			if e := recover(); e != nil {
				log.Printf("Fatal error! %v", e)
				jsonErr := &resourceError{nil, "Something went wrong", http.StatusInternalServerError}
				jsonErr.WriteToResponseAsJson(w)
			}
		}()

		inner.ServeHTTP(w, r)

	})
}

var WhereaboutsHttpHandler = router.NewRouterBuilder().
	AddRoutes(userRoutes).
	AddRoutes(persistentSessionsRoutes).
	AddRoutes(whereaboutsRoutes).
	AddRoute(docsRoute).
	AddInterceptor(logger).
	AddInterceptor(recoveryInterceptor).
	Build()
