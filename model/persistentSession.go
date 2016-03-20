package model

import (
	"time"
)

type PersistentSession struct {
	Username string
	Key      string
	Updated  time.Time
}
