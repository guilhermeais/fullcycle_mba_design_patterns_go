package domain

import "time"

type Payment struct {
	Id     string
	Amount float64
	Date   time.Time
}
