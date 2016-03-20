package router

import (
	"net/http"

	"github.com/gorilla/mux"
)

// RouteInterceptor is a function that decorates inner and has access to the route
type RouteInterceptor func(inner http.Handler, r Route) http.Handler

type routerBuilder struct {
	interceptors []RouteInterceptor
	routes       []Route
}

func NewRouterBuilder() *routerBuilder {
	return new(routerBuilder)
}

func (rb *routerBuilder) AddRoute(route Route) *routerBuilder {
	rb.routes = append(rb.routes, route)
	return rb
}

func (rb *routerBuilder) AddRoutes(routes []Route) *routerBuilder {
	for _, route := range routes {
		rb.AddRoute(route)
	}
	return rb
}

func (rb *routerBuilder) AddInterceptor(interceptor RouteInterceptor) *routerBuilder {
	rb.interceptors = append(rb.interceptors, interceptor)
	return rb
}

func (rb *routerBuilder) Build() http.Handler {
	router := mux.NewRouter().StrictSlash(true)
	for _, route := range rb.routes {
		var handler http.Handler
		handler = route.Handler
		for i := len(rb.interceptors) - 1; i >= 0; i-- {
			interceptor := rb.interceptors[i]
			handler = interceptor(handler, route)
		}

		router.
			Methods(string(route.Method)).
			Path(route.Pattern).
			Name(route.Name).
			Handler(handler)

	}

	var asHTTPHandler = http.Handler(router)
	return asHTTPHandler
}
