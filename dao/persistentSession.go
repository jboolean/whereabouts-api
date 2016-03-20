package dao

import (
	"database/sql"
	"github.com/jboolean/whereabouts-api/model"
	"log"
)

type persistentSessionDAO struct {
	db *sql.DB
}

func (dao persistentSessionDAO) Upsert(persistentSession *model.PersistentSession) {
	// manual upsert, too much of a pain in postgres

	username := persistentSession.Username
	key := persistentSession.Key

	stmt, err := dao.db.Prepare("select count(key) from persistent_sessions where username = $1")
	var existingCount int
	err = stmt.QueryRow(username).Scan(&existingCount)

	if existingCount > 0 {
		stmt, err = dao.db.Prepare("update persistent_sessions set key = $1 where username = $2")
	} else {
		stmt, err = dao.db.Prepare("insert into persistent_sessions (key, username) values ($1, $2)")
	}

	stmt.Exec(key, username)
	if err != nil {
		log.Panic(err)
	}
}

func (dao persistentSessionDAO) FindByKey(key string) *model.PersistentSession {
	stmt, err := dao.db.Prepare("select key, username, updated from persistent_sessions " +
		"where key = $1")

	var result = new(model.PersistentSession)
	err = stmt.QueryRow(key).
		Scan(&result.Key, &result.Username, &result.Updated)

	if err == sql.ErrNoRows {
		return nil
	}

	if err != nil {
		log.Panic(err)
	}

	return result
}
