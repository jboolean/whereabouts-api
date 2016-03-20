package dao

import (
	"database/sql"
	_ "github.com/lib/pq"
	"log"
)

var UserDAO userDAO
var PersistentSessionDAO persistentSessionDAO
var ProcessedLocationsDAO processedLocationsDAO
var WhereaboutsSummaryDAO whereaboutsSummaryDAO

/*
What's going on here?
An instance of each dao is initialized with the database connection and exposed.
This way the database isn't exposed from the package and each dao can get its own scope

Usage example:
`dao.UserDAO.FindUser(username)`


As a consequence, all daos must be registered here.
*/
func init() {
	db, err := sql.Open("postgres", "user=julianboilen dbname=whereabouts sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	UserDAO = userDAO{db: db}
	PersistentSessionDAO = persistentSessionDAO{db: db}
	ProcessedLocationsDAO = processedLocationsDAO{db: db}
	WhereaboutsSummaryDAO = whereaboutsSummaryDAO{db: db}
}
