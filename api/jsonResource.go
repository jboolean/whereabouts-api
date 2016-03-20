package api

import (
	"encoding/json"
	"log"
	"net/http"
)

type resourceError struct {
	Source  error  `json:"-"` //nillable
	Message string `json:"error"`
	Code    int    `json:"-"`
}

func (err resourceError) Error() string {
	return err.Message
}

func (err *resourceError) WriteToResponseAsJson(w http.ResponseWriter) {
	if err != nil {
		if err.Source != nil {
			log.Printf("Resource error: %s\n", err.Source)
		}

		w.WriteHeader(err.Code)
		if errEncodeErr := json.NewEncoder(w).Encode(err); errEncodeErr != nil {
			log.Printf("Encoder error:%v:", errEncodeErr)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		return
	}
}

type jsonResource func(*http.Request) (interface{}, *resourceError)

// implement http.Handler
func (fn jsonResource) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	response, err := fn(r)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")

	if err != nil {
		err.WriteToResponseAsJson(w)
		return
	}

	if response == nil {
		w.WriteHeader(http.StatusNoContent)
		return
	}

	if err := json.NewEncoder(w).Encode(response); err != nil {
		log.Printf("Encoder error:%v:", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}
