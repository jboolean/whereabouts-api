package dao

import (
	"database/sql"
	"github.com/jboolean/whereabouts-api/model"
	"log"
	"time"
)

type processedLocationsDAO struct {
	db *sql.DB
}

func (dao processedLocationsDAO) Store(processedLocation *model.ProcessedLocation) {
	stmt,
		err := dao.db.Prepare("insert into processed_locations (username, time_to_home, distance_from_home) " +
		"values ($1, $2, $3)")

	if err != nil {
		log.Panic(err)
	}

	_, err = stmt.Exec(
		processedLocation.Username,
		processedLocation.TimeToHome,
		processedLocation.DistanceFromHome)

	if err != nil {
		// I hate go error handling
		log.Panic(err)
	}
}

func (dao processedLocationsDAO) FindLatest(username string) *model.ProcessedLocation {
	stmt, err := dao.db.Prepare("select username, time_to_home, distance_from_home, created from " +
		"processed_locations where username = $1 order by created desc limit 1")

	if err != nil {
		log.Panic(err)
	}

	var timeToHome int64

	var result = new(model.ProcessedLocation)
	err = stmt.QueryRow(username).Scan(
		&result.Username,
		&timeToHome,
		&result.DistanceFromHome,
		&result.CreatedOn)

	result.TimeToHome = time.Duration(timeToHome)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Panic(err)
	}

	return result
}

func (dao processedLocationsDAO) FindInTimeRange(username string, start time.Time, end time.Time) []model.ProcessedLocation {
	stmt, err := dao.db.Prepare("select username, time_to_home, distance_from_home, created from " +
		"processed_locations where username = $1 and created <= $2 and created >= $3")

	if err != nil {
		log.Panic(err)
	}

	var rows *sql.Rows
	rows, err = stmt.Query(username, start, end)

	defer rows.Close()

	records := make([]model.ProcessedLocation, 0)
	for rows.Next() {
		var record = new(model.ProcessedLocation)
		var timeToHome int64
		err = rows.Scan(
			&record.Username,
			&timeToHome,
			&record.DistanceFromHome,
			&record.CreatedOn)
		record.TimeToHome = time.Duration(timeToHome)
		records = append(records, *record)
	}

	err = rows.Err()

	if err != nil {
		log.Panic(err)
	}

	return records
}
