package dao

import (
	"database/sql"
	"github.com/jboolean/whereabouts-api/model"
	"log"
	"time"
)

type whereaboutsSummaryDAO struct {
	db *sql.DB
}

func (dao whereaboutsSummaryDAO) Store(summary *model.WhereaboutsSummary) {
	tx, err := dao.db.Begin()

	defer tx.Rollback()

	var stmt *sql.Stmt

	tx.Exec("delete from whereabouts_summaries where username = $1", summary.Username)

	stmt, err = tx.Prepare("insert into whereabouts_summaries (username, updated, time_to_home, velocity) " +
		"values ($1, $2, $3, $4)")

	_, err = stmt.Exec(
		summary.Username,
		summary.UpdatedOn,
		summary.TimeToHome,
		summary.Velocity)

	err = tx.Commit()

	if err != nil {
		log.Panic(err)
	}
}

func (dao whereaboutsSummaryDAO) FindAll() []model.WhereaboutsSummary {
	stmt, err := dao.db.Prepare("select username, updated, time_to_home, velocity from " +
		"whereabouts_summaries")

	if err != nil {
		log.Panic(err)
	}

	var rows *sql.Rows
	rows, err = stmt.Query()

	defer rows.Close()

	records := make([]model.WhereaboutsSummary, 0)
	for rows.Next() {
		var record = new(model.WhereaboutsSummary)
		var timeToHome int64
		err = rows.Scan(
			&record.Username,
			&record.UpdatedOn,
			&timeToHome,
			&record.Velocity)
		record.TimeToHome = time.Duration(timeToHome)
		records = append(records, *record)
	}

	err = rows.Err()

	if err != nil {
		log.Panic(err)
	}

	return records
}

func (dao whereaboutsSummaryDAO) FindByUsername(username string) *model.WhereaboutsSummary {
	stmt, err := dao.db.Prepare("select username, updated, time_to_home, velocity from " +
		"whereabouts_summaries where username = $1")

	if err != nil {
		log.Panic(err)
	}

	var record = new(model.WhereaboutsSummary)
	var timeToHome int64
	err = stmt.QueryRow(username).Scan(
		&record.Username,
		&record.UpdatedOn,
		&timeToHome,
		&record.Velocity)

	if err != nil {
		if err == sql.ErrNoRows {
			return nil
		}
		log.Panic(err)
	}

	record.TimeToHome = time.Duration(timeToHome)

	return record
}
