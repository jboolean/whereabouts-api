package model

import "time"

type ProcessedLocation struct {
	Username         string
	TimeToHome       time.Duration
	DistanceFromHome int
	CreatedOn        time.Time
}
