package domain

import (
	"time"
)

type Contract struct {
	Id, Description string
	Amount          float64
	Periods         int
	Date            time.Time
}
