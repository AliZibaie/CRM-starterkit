package domain

import "time"

type Setting struct {
	Key       string
	Value     string
	UpdatedBy string
	UpdatedAt time.Time
}
