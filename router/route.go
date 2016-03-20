package router

import (
	"net/http"
)

// A Route represents a web resource and a route to it
type Route struct {
	Name    string
	Method  HttpMethod
	Pattern string
	// Create your own stucts to provide addition information
	// That this router package need not be privy to
	// But your interceptors do
	Meta interface{}

	Handler http.Handler
}

// Routes is a slice of Routes
type Routes []Route
