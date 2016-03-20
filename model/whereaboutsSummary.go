package model

import "time"

type WhereaboutsSummary struct {
	Username   string
	UpdatedOn  time.Time
	TimeToHome time.Duration
	Velocity   float64
}
