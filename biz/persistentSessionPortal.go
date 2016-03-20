package biz

import (
	"fmt"
	"github.com/jboolean/whereabouts-api/dao"
	"github.com/jboolean/whereabouts-api/model"
)

var InvalidPasswordError = fmt.Errorf("Invalid password")

// Login and logout and previous sessions
func Login(username string, password string) (*model.PersistentSession, error) {
	if !CheckOrSetUserPassword(username, password) {
		return nil, InvalidPasswordError
	}

	persistentSession := &model.PersistentSession{
		Username: username,
		Key:      makeToken(),
	}

	dao.PersistentSessionDAO.Upsert(persistentSession)

	return persistentSession, nil
}

// Returns a session or nil
func FindSessionByKey(key string) *model.PersistentSession {
	return dao.PersistentSessionDAO.FindByKey(key)
}
