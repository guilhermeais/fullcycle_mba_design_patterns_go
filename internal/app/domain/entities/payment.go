package domain

import "time"

type Payment struct {
	Id     string
	Amount int
	Date   time.Time
}
