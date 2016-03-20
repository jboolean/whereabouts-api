package api

import "github.com/jboolean/whereabouts-api/model"

type persistentSessionResponse struct {
	Username string `json:"username"`
	Key      string `json:"key"`
}

func makePersistentSessionResponse(persistentSession *model.PersistentSession) *persistentSessionResponse {
	return &persistentSessionResponse{
		Username: persistentSession.Username,
		Key:      persistentSession.Key,
	}
}
