package api

import (
	"github.com/jboolean/whereabouts-api/router"
	"log"
	"net/http"
	"time"
)

func logger(inner http.Handler, route router.Route) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		inner.ServeHTTP(w, r)

		log.Printf(
			"method=%s uri=%s route=%s time=%s",
			r.Method,
			r.RequestURI,
			route.Name,
			time.Since(start),
		)
	})
}
