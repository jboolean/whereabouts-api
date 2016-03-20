package main

import (
	"github.com/jboolean/whereabouts-api/api"
	"log"
	"net/http"
	"os"
)

func main() {
	port := os.Getenv("PORT")

	if port == "" {
		log.Fatal("$PORT must be set")
	}

	log.Fatal(http.ListenAndServe(":"+port, api.WhereaboutsHttpHandler))
}
