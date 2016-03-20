package api

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"net/http"
)

// Reads json from request body into target, or returns any error encountered
func decodeJsonBody(target interface{}, r *http.Request) error {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		return err
	}

	if err := r.Body.Close(); err != nil {
		return err
	}
	if err := json.Unmarshal(body, target); err != nil {
		return err
	}

	return nil
}

func readBodyString(r *http.Request) string {
	body, err := ioutil.ReadAll(io.LimitReader(r.Body, 1048576))
	if err != nil {
		log.Panic(err)
	}
	if string(body) == "{}" {
		return ""
	}

	return string(body)
}
