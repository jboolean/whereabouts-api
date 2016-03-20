package main

import (
	"github.com/jboolean/whereabouts-api/api"
	"log"
	"net/http"
)

func main() {

	log.Fatal(http.ListenAndServe(":8080", api.WhereaboutsHttpHandler))
}
