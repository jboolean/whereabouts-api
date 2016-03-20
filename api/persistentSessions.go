package api

import (
	"github.com/jboolean/whereabouts-api/biz"
	"github.com/jboolean/whereabouts-api/router"
	"net/http"
)

var persistentSessionsRoutes = router.Routes{
	router.Route{
		Name:    "Create",
		Method:  "POST",
		Pattern: "/persistent-sessions",
		Handler: jsonResource(login),
	},
}

func login(r *http.Request) (interface{}, *resourceError) {
	newSessionRequest := new(persistentSessionResponse)

	var err error
	err = decodeJsonBody(newSessionRequest, r)
	if err != nil {
		return nil, &resourceError{err, "Error reading request body.", http.StatusBadRequest}
	}

	resultingSession, creationErr := biz.Login(newSessionRequest.Username, r.Header["X-Password"][0])
	if creationErr == biz.InvalidPasswordError {
		return nil, &resourceError{creationErr, "Incorrect password.", http.StatusUnauthorized}
	}

	return makePersistentSessionResponse(resultingSession), nil
}
